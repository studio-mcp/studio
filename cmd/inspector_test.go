package cmd

import (
	"context"
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// InspectorToolsListResponse represents the JSON response from tools/list
type InspectorToolsListResponse struct {
	Tools []InspectorTool `json:"tools"`
}

// InspectorTool represents a tool in the inspector response
type InspectorTool struct {
	Name        string               `json:"name"`
	Description string               `json:"description"`
	InputSchema InspectorInputSchema `json:"inputSchema"`
}

// InspectorInputSchema represents the input schema for a tool
type InspectorInputSchema struct {
	Type       string                       `json:"type"`
	Properties map[string]InspectorProperty `json:"properties"`
	Required   []string                     `json:"required"`
}

// InspectorProperty represents a property in the input schema
type InspectorProperty struct {
	Type        string                 `json:"type"`
	Description string                 `json:"description,omitempty"`
	Items       *InspectorPropertyItem `json:"items,omitempty"`
}

// InspectorPropertyItem represents items in an array property
type InspectorPropertyItem struct {
	Type string `json:"type"`
}

// InspectorToolCallResponse represents the JSON response from tools/call
type InspectorToolCallResponse struct {
	Content []InspectorContent `json:"content"`
	IsError bool               `json:"isError,omitempty"`
}

// InspectorContent represents content in a tool call response
type InspectorContent struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

// TestMCPInspectorSmokeTest tests studio-mcp using the MCP inspector
// This is a smoke test that verifies the MCP server works correctly with the
// official MCP inspector tool.
func TestMCPInspectorSmokeTest(t *testing.T) {
	timeout := 10 * time.Second

	// Build the project first
	projectRoot, err := filepath.Abs("..")
	require.NoError(t, err)

	// Ensure bin directory exists
	binDir := filepath.Join(projectRoot, "bin")
	err = os.MkdirAll(binDir, 0755)
	require.NoError(t, err)

	// Build to bin/studio-mcp
	binaryPath := filepath.Join(binDir, "studio-mcp")
	buildCmd := exec.Command("go", "build", "-o", binaryPath, ".")
	buildCmd.Dir = projectRoot
	err = buildCmd.Run()
	require.NoError(t, err, "Failed to build project")

	// Helper function to run MCP inspector commands
	runInspectorCmd := func(t *testing.T, args []string) (string, string, error) {
		ctx, cancel := context.WithTimeout(context.Background(), timeout)
		defer cancel()

		cmd := exec.CommandContext(ctx, "npx", args...)
		cmd.Dir = projectRoot

		var stdout, stderr strings.Builder
		cmd.Stdout = &stdout
		cmd.Stderr = &stderr

		err := cmd.Run()
		return stdout.String(), stderr.String(), err
	}

	t.Run("BasicInspectorConnection", func(t *testing.T) {
		t.Run("can list simple echo tools", func(t *testing.T) {
			args := []string{
				"@modelcontextprotocol/inspector",
				"--cli",
				binaryPath,
				"echo",
				"hello",
				"--method",
				"tools/list",
			}

			stdout, stderr, err := runInspectorCmd(t, args)

			// The inspector should be able to connect without errors
			assert.NoError(t, err, "Inspector should connect successfully. Stdout: %s, Stderr: %s", stdout, stderr)

			// Parse JSON response
			var response InspectorToolsListResponse
			err = json.Unmarshal([]byte(stdout), &response)
			require.NoError(t, err, "Should parse JSON response. Stdout: %s", stdout)

			// Validate response structure
			require.Len(t, response.Tools, 1, "Should have exactly one tool")

			tool := response.Tools[0]
			assert.Equal(t, "echo", tool.Name)
			assert.Equal(t, "Run the shell command `echo hello`", tool.Description)
			assert.Equal(t, "object", tool.InputSchema.Type)

			// Simple commands without [args...] don't have properties
			assert.Empty(t, tool.InputSchema.Properties, "Simple commands should not have properties")
		})

		t.Run("can list echo tools with args", func(t *testing.T) {
			args := []string{
				"@modelcontextprotocol/inspector",
				"--cli",
				binaryPath,
				"echo",
				"hello",
				"[args...]",
				"--method",
				"tools/list",
			}

			stdout, stderr, err := runInspectorCmd(t, args)

			// The inspector should be able to connect without errors
			assert.NoError(t, err, "Inspector should connect successfully. Stdout: %s, Stderr: %s", stdout, stderr)

			// Parse JSON response
			var response InspectorToolsListResponse
			err = json.Unmarshal([]byte(stdout), &response)
			require.NoError(t, err, "Should parse JSON response. Stdout: %s", stdout)

			// Validate response structure
			require.Len(t, response.Tools, 1, "Should have exactly one tool")

			tool := response.Tools[0]
			assert.Equal(t, "echo", tool.Name)
			assert.Equal(t, "Run the shell command `echo hello [args...]`", tool.Description)
			assert.Equal(t, "object", tool.InputSchema.Type)

			// Should have args property for commands with [args...]
			assert.Contains(t, tool.InputSchema.Properties, "args")
			argsProperty := tool.InputSchema.Properties["args"]
			assert.Equal(t, "array", argsProperty.Type)
			assert.Equal(t, "Additional command line arguments", argsProperty.Description)
		})

		t.Run("can call simple echo tool", func(t *testing.T) {
			args := []string{
				"@modelcontextprotocol/inspector",
				"--cli",
				binaryPath,
				"echo",
				"hello",
				"--method",
				"tools/call",
				"--tool-name",
				"echo",
			}

			stdout, stderr, err := runInspectorCmd(t, args)

			// Should succeed
			assert.NoError(t, err, "Inspector should call echo tool successfully. Stdout: %s, Stderr: %s", stdout, stderr)

			// Parse JSON response
			var response InspectorToolCallResponse
			err = json.Unmarshal([]byte(stdout), &response)
			require.NoError(t, err, "Should parse JSON response. Stdout: %s", stdout)

			// Validate response structure
			require.Len(t, response.Content, 1, "Should have exactly one content item")

			content := response.Content[0]
			assert.Equal(t, "text", content.Type)
			assert.Equal(t, "hello", content.Text)
			assert.False(t, response.IsError, "Should not be an error")
		})
	})

	t.Run("BlueprintedCommand", func(t *testing.T) {
		t.Run("can list blueprinted tools with proper schema", func(t *testing.T) {
			args := []string{
				"@modelcontextprotocol/inspector",
				"--cli",
				binaryPath,
				"echo",
				"{{text#message to echo}}",
				"--method",
				"tools/list",
			}

			stdout, stderr, err := runInspectorCmd(t, args)

			// Should succeed
			assert.NoError(t, err, "Inspector should list blueprinted tools successfully. Stdout: %s, Stderr: %s", stdout, stderr)

			// Parse JSON response
			var response InspectorToolsListResponse
			err = json.Unmarshal([]byte(stdout), &response)
			require.NoError(t, err, "Should parse JSON response. Stdout: %s", stdout)

			// Validate response structure
			require.Len(t, response.Tools, 1, "Should have exactly one tool")

			tool := response.Tools[0]
			assert.Equal(t, "echo", tool.Name)
			assert.Equal(t, "Run the shell command `echo {{text}}`", tool.Description)
			assert.Equal(t, "object", tool.InputSchema.Type)

			// Should have text property for blueprinted commands
			assert.Contains(t, tool.InputSchema.Properties, "text")
			textProperty := tool.InputSchema.Properties["text"]
			assert.Equal(t, "string", textProperty.Type)
			assert.Equal(t, "message to echo", textProperty.Description)

			// Should have text as required parameter
			assert.Contains(t, tool.InputSchema.Required, "text")

			// Should NOT have args property for blueprinted commands
			assert.NotContains(t, tool.InputSchema.Properties, "args")
		})

		t.Run("can call blueprinted tool with arguments", func(t *testing.T) {
			args := []string{
				"@modelcontextprotocol/inspector",
				"--cli",
				binaryPath,
				"echo",
				"{{text#message to echo}}",
				"--method",
				"tools/call",
				"--tool-name",
				"echo",
				"--tool-arg",
				"text=Hello from inspector!",
			}

			stdout, stderr, err := runInspectorCmd(t, args)

			// Should succeed
			assert.NoError(t, err, "Inspector should call blueprinted tool successfully. Stdout: %s, Stderr: %s", stdout, stderr)

			// Parse JSON response
			var response InspectorToolCallResponse
			err = json.Unmarshal([]byte(stdout), &response)
			require.NoError(t, err, "Should parse JSON response. Stdout: %s", stdout)

			// Validate response structure
			require.Len(t, response.Content, 1, "Should have exactly one content item")

			content := response.Content[0]
			assert.Equal(t, "text", content.Type)
			assert.Equal(t, "Hello from inspector!", content.Text)
			assert.False(t, response.IsError, "Should not be an error")
		})

		t.Run("can handle mixed blueprint with additional args", func(t *testing.T) {
			args := []string{
				"@modelcontextprotocol/inspector",
				"--cli",
				binaryPath,
				"echo",
				"{{text#message to echo}}",
				"[args...]",
				"--method",
				"tools/list",
			}

			stdout, stderr, err := runInspectorCmd(t, args)

			// Should succeed
			assert.NoError(t, err, "Inspector should handle mixed blueprint successfully. Stdout: %s, Stderr: %s", stdout, stderr)

			// Parse JSON response
			var response InspectorToolsListResponse
			err = json.Unmarshal([]byte(stdout), &response)
			require.NoError(t, err, "Should parse JSON response. Stdout: %s", stdout)

			// Validate response structure
			require.Len(t, response.Tools, 1, "Should have exactly one tool")

			tool := response.Tools[0]
			assert.Equal(t, "echo", tool.Name)
			assert.Equal(t, "Run the shell command `echo {{text}} [args...]`", tool.Description)

			// Should have both text and args properties
			assert.Contains(t, tool.InputSchema.Properties, "text")
			assert.Contains(t, tool.InputSchema.Properties, "args")

			textProperty := tool.InputSchema.Properties["text"]
			assert.Equal(t, "string", textProperty.Type)
			assert.Equal(t, "message to echo", textProperty.Description)

			argsProperty := tool.InputSchema.Properties["args"]
			assert.Equal(t, "array", argsProperty.Type)

			// Only text should be required
			assert.Contains(t, tool.InputSchema.Required, "text")
			assert.NotContains(t, tool.InputSchema.Required, "args")
		})
	})

	t.Run("ComplexCommand", func(t *testing.T) {
		t.Run("can handle git status command", func(t *testing.T) {
			args := []string{
				"@modelcontextprotocol/inspector",
				"--cli",
				binaryPath,
				"git",
				"status",
				"--method",
				"tools/list",
			}

			stdout, stderr, err := runInspectorCmd(t, args)

			// Should succeed
			assert.NoError(t, err, "Inspector should handle git command successfully. Stdout: %s, Stderr: %s", stdout, stderr)

			// Parse JSON response
			var response InspectorToolsListResponse
			err = json.Unmarshal([]byte(stdout), &response)
			require.NoError(t, err, "Should parse JSON response. Stdout: %s", stdout)

			// Validate response structure
			require.Len(t, response.Tools, 1, "Should have exactly one tool")

			tool := response.Tools[0]
			assert.Equal(t, "git", tool.Name)
			assert.Equal(t, "Run the shell command `git status`", tool.Description)
			assert.Equal(t, "object", tool.InputSchema.Type)

			// Simple git status command without [args...] doesn't have properties
			assert.Empty(t, tool.InputSchema.Properties, "Simple git status should not have properties")
		})

		t.Run("can handle git status command with args", func(t *testing.T) {
			args := []string{
				"@modelcontextprotocol/inspector",
				"--cli",
				binaryPath,
				"git",
				"status",
				"[args...]",
				"--method",
				"tools/list",
			}

			stdout, stderr, err := runInspectorCmd(t, args)

			// Should succeed
			assert.NoError(t, err, "Inspector should handle git command with args successfully. Stdout: %s, Stderr: %s", stdout, stderr)

			// Parse JSON response
			var response InspectorToolsListResponse
			err = json.Unmarshal([]byte(stdout), &response)
			require.NoError(t, err, "Should parse JSON response. Stdout: %s", stdout)

			// Validate response structure
			require.Len(t, response.Tools, 1, "Should have exactly one tool")

			tool := response.Tools[0]
			assert.Equal(t, "git", tool.Name)
			assert.Equal(t, "Run the shell command `git status [args...]`", tool.Description)
			assert.Equal(t, "object", tool.InputSchema.Type)

			// Should have args property
			assert.Contains(t, tool.InputSchema.Properties, "args")
			argsProperty := tool.InputSchema.Properties["args"]
			assert.Equal(t, "array", argsProperty.Type)
			assert.Equal(t, "Additional command line arguments", argsProperty.Description)
		})

		t.Run("can execute git status", func(t *testing.T) {
			args := []string{
				"@modelcontextprotocol/inspector",
				"--cli",
				binaryPath,
				"git",
				"status",
				"--method",
				"tools/call",
				"--tool-name",
				"git",
			}

			stdout, stderr, err := runInspectorCmd(t, args)

			// Should succeed (git status should work in any directory)
			assert.NoError(t, err, "Inspector should execute git status successfully. Stdout: %s, Stderr: %s", stdout, stderr)

			// Parse JSON response
			var response InspectorToolCallResponse
			err = json.Unmarshal([]byte(stdout), &response)
			require.NoError(t, err, "Should parse JSON response. Stdout: %s", stdout)

			// Validate response structure
			require.Len(t, response.Content, 1, "Should have exactly one content item")

			content := response.Content[0]
			assert.Equal(t, "text", content.Type)
			assert.Contains(t, content.Text, "On branch", "Should contain git status output")
			assert.False(t, response.IsError, "Should not be an error")
		})
	})

	t.Run("SayCommandRegression", func(t *testing.T) {
		t.Run("can handle say command with -v flag", func(t *testing.T) {
			args := []string{
				"@modelcontextprotocol/inspector",
				"--cli",
				binaryPath,
				"say",
				"-v",
				"siri",
				"{{speech#message to speak}}",
				"--method",
				"tools/list",
			}

			stdout, stderr, err := runInspectorCmd(t, args)

			// Should succeed (this was previously failing with flag parsing errors)
			assert.NoError(t, err, "Inspector should handle say command with -v flag successfully. Stdout: %s, Stderr: %s", stdout, stderr)

			// Parse JSON response
			var response InspectorToolsListResponse
			err = json.Unmarshal([]byte(stdout), &response)
			require.NoError(t, err, "Should parse JSON response. Stdout: %s", stdout)

			// Validate response structure
			require.Len(t, response.Tools, 1, "Should have exactly one tool")

			tool := response.Tools[0]
			assert.Equal(t, "say", tool.Name)
			assert.Equal(t, "Run the shell command `say -v siri {{speech}}`", tool.Description)
			assert.Equal(t, "object", tool.InputSchema.Type)

			// Should have speech property
			assert.Contains(t, tool.InputSchema.Properties, "speech")
			speechProperty := tool.InputSchema.Properties["speech"]
			assert.Equal(t, "string", speechProperty.Type)
			assert.Equal(t, "message to speak", speechProperty.Description)

			// Should have speech as required parameter
			assert.Contains(t, tool.InputSchema.Required, "speech")
		})
	})

	t.Run("MultipleBlueprints", func(t *testing.T) {
		t.Run("can handle multiple blueprint parameters", func(t *testing.T) {
			args := []string{
				"@modelcontextprotocol/inspector",
				"--cli",
				binaryPath,
				"rails",
				"generate",
				"{{generator#Rails generator name}}",
				"{{name#Resource name}}",
				"--method",
				"tools/list",
			}

			stdout, stderr, err := runInspectorCmd(t, args)

			// Should succeed
			assert.NoError(t, err, "Inspector should handle multiple blueprints successfully. Stdout: %s, Stderr: %s", stdout, stderr)

			// Parse JSON response
			var response InspectorToolsListResponse
			err = json.Unmarshal([]byte(stdout), &response)
			require.NoError(t, err, "Should parse JSON response. Stdout: %s", stdout)

			// Validate response structure
			require.Len(t, response.Tools, 1, "Should have exactly one tool")

			tool := response.Tools[0]
			assert.Equal(t, "rails", tool.Name)
			assert.Equal(t, "Run the shell command `rails generate {{generator}} {{name}}`", tool.Description)
			assert.Equal(t, "object", tool.InputSchema.Type)

			// Should have both generator and name properties
			assert.Contains(t, tool.InputSchema.Properties, "generator")
			assert.Contains(t, tool.InputSchema.Properties, "name")

			generatorProperty := tool.InputSchema.Properties["generator"]
			assert.Equal(t, "string", generatorProperty.Type)
			assert.Equal(t, "Rails generator name", generatorProperty.Description)

			nameProperty := tool.InputSchema.Properties["name"]
			assert.Equal(t, "string", nameProperty.Type)
			assert.Equal(t, "Resource name", nameProperty.Description)

			// Both should be required
			assert.Contains(t, tool.InputSchema.Required, "generator")
			assert.Contains(t, tool.InputSchema.Required, "name")

			// Should NOT have args property
			assert.NotContains(t, tool.InputSchema.Properties, "args")
		})
	})
}
