package tool

import (
	"context"
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/modelcontextprotocol/go-sdk/jsonschema"
	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/stretchr/testify/assert"
)

func TestTool_Execute(t *testing.T) {
	tests := []struct {
		name           string
		args           []string
		expectError    bool
		expectOutput   string
		containsOutput string
	}{
		{
			name:         "executes simple command successfully",
			args:         []string{"echo", "hello"},
			expectOutput: "hello",
		},
		{
			name:         "executes command with multiple arguments",
			args:         []string{"echo", "hello", "world"},
			expectOutput: "hello world",
		},
		{
			name:           "captures stderr output",
			args:           []string{"sh", "-c", "echo 'error message' >&2"},
			containsOutput: "error message",
		},
		{
			name:        "handles command failure",
			args:        []string{"false"},
			expectError: true,
		},
		{
			name:        "handles non-existent command",
			args:        []string{"this-command-does-not-exist-12345"},
			expectError: true,
		},
		{
			name:         "handles command with spaces",
			args:         []string{"echo", "hello world with spaces"},
			expectOutput: "hello world with spaces",
		},
		{
			name:         "handles command with no output",
			args:         []string{"true"},
			expectOutput: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			output, err := Execute(tt.args[0], tt.args[1:]...)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			if tt.expectOutput != "" {
				assert.Equal(t, tt.expectOutput, strings.TrimSpace(output))
			}

			if tt.containsOutput != "" {
				assert.Contains(t, output, tt.containsOutput)
			}
		})
	}
}

func TestTool_DebugMode(t *testing.T) {
	t.Run("debug mode is off by default", func(t *testing.T) {
		assert.False(t, IsDebugMode())
	})

	t.Run("can enable debug mode", func(t *testing.T) {
		SetDebugMode(true)
		assert.True(t, IsDebugMode())

		// Reset for other tests
		SetDebugMode(false)
		assert.False(t, IsDebugMode())
	})
}

func TestTool_LogFile(t *testing.T) {
	// Create a temporary log file
	tmpFile, err := os.CreateTemp("", "test_log_*.log")
	assert.NoError(t, err)
	tmpFile.Close()
	defer os.Remove(tmpFile.Name())

	t.Run("can set log file", func(t *testing.T) {
		err := SetLogFile(tmpFile.Name())
		assert.NoError(t, err)
	})

	t.Run("debug messages are written to log file", func(t *testing.T) {
		// Enable debug mode and set log file
		SetDebugMode(true)
		err := SetLogFile(tmpFile.Name())
		assert.NoError(t, err)

		// Execute a command that will generate debug output
		_, err = Execute("echo", "test log message")
		assert.NoError(t, err)

		// Read the log file
		content, err := os.ReadFile(tmpFile.Name())
		assert.NoError(t, err)

		// Check that debug messages are in the log file
		logContent := string(content)
		assert.Contains(t, logContent, "Executing command: echo test log message")
		assert.Contains(t, logContent, "Command completed successfully with exit code 0")

		// Reset for other tests
		SetDebugMode(false)
		if logFile != nil {
			logFile.Close()
			logFile = nil
			logger = nil
		}
	})

	t.Run("handles invalid log file path", func(t *testing.T) {
		err := SetLogFile("/invalid/path/that/does/not/exist.log")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to open log file")
	})
}

func TestTool_CreateToolFunction(t *testing.T) {
	tests := []struct {
		name           string
		blueprint      Blueprint
		args           map[string]any
		expectText     string
		expectIsError  bool
		expectContains string
	}{
		{
			name:       "creates function for simple command",
			blueprint:  &MockBlueprint{commandArgs: []string{"echo", "hello"}},
			args:       map[string]any{},
			expectText: "hello",
		},
		{
			name:       "creates function with template arguments",
			blueprint:  &MockBlueprint{commandArgs: []string{"echo", "Hello World"}},
			args:       map[string]any{"message": "Hello World"},
			expectText: "Hello World",
		},
		{
			name:          "handles command errors",
			blueprint:     &MockBlueprint{commandArgs: []string{"false"}},
			args:          map[string]any{},
			expectText:    "",
			expectIsError: true,
		},
		{
			name:           "handles command errors with stderr output",
			blueprint:      &MockBlueprint{commandArgs: []string{"sh", "-c", "echo 'error message' >&2; exit 1"}},
			args:           map[string]any{},
			expectContains: "error message",
			expectIsError:  true,
		},
		{
			name:           "handles blueprint validation errors",
			blueprint:      &MockBlueprintWithError{err: fmt.Errorf("missing required parameter: name")},
			args:           map[string]any{},
			expectIsError:  true,
			expectContains: "Validation error: missing required parameter: name",
		},
		{
			name:          "handles command with no output",
			blueprint:     &MockBlueprint{commandArgs: []string{"true"}},
			args:          map[string]any{},
			expectText:    "",
			expectIsError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := CreateToolFunction(tt.blueprint)

			// Create MCP parameters
			params := &mcp.CallToolParamsFor[map[string]any]{
				Arguments: tt.args,
			}

			result, err := handler(context.Background(), nil, params)

			assert.NoError(t, err)
			assert.Len(t, result.Content, 1)

			// Cast to TextContent to access the text
			textContent, ok := result.Content[0].(*mcp.TextContent)
			assert.True(t, ok, "Expected content to be TextContent")

			text := textContent.Text

			if tt.expectContains != "" {
				assert.Contains(t, text, tt.expectContains)
			} else {
				assert.Equal(t, tt.expectText, text)
			}

			assert.Equal(t, tt.expectIsError, result.IsError)
		})
	}
}

func TestTool_GenerateToolName(t *testing.T) {
	tests := []struct {
		name        string
		baseCommand string
		expected    string
	}{
		{
			name:        "simple command without dashes",
			baseCommand: "git",
			expected:    "git",
		},
		{
			name:        "command with single dash",
			baseCommand: "git-flow",
			expected:    "git_flow",
		},
		{
			name:        "command with multiple dashes",
			baseCommand: "my-long-command-name",
			expected:    "my_long_command_name",
		},
		{
			name:        "command with no changes needed",
			baseCommand: "simple_command",
			expected:    "simple_command",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := GenerateToolName(tt.baseCommand)
			assert.Equal(t, tt.expected, result)
		})
	}
}

// MockBlueprint is a test helper that implements the Blueprint interface
type MockBlueprint struct {
	commandArgs []string
}

func (m *MockBlueprint) BuildCommandArgs(args map[string]interface{}) ([]string, error) {
	return m.commandArgs, nil
}

func (m *MockBlueprint) GetBaseCommand() string {
	return "mock-tool"
}

func (m *MockBlueprint) GetCommandFormat() string {
	return "mock-tool"
}

func (m *MockBlueprint) GetInputSchema() interface{} {
	return &jsonschema.Schema{
		Type:       "object",
		Properties: make(map[string]*jsonschema.Schema),
	}
}

// MockBlueprintWithError is a test helper that returns an error
type MockBlueprintWithError struct {
	err error
}

func (m *MockBlueprintWithError) BuildCommandArgs(args map[string]interface{}) ([]string, error) {
	return nil, m.err
}

func (m *MockBlueprintWithError) GetBaseCommand() string {
	return "mock-error-tool"
}

func (m *MockBlueprintWithError) GetCommandFormat() string {
	return "mock-error-tool"
}

func (m *MockBlueprintWithError) GetInputSchema() interface{} {
	return &jsonschema.Schema{
		Type:       "object",
		Properties: make(map[string]*jsonschema.Schema),
	}
}
