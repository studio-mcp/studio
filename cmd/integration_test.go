package cmd

import (
	"bufio"
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

// MCPRequest represents a JSON-RPC request
type MCPRequest struct {
	JSONRPC string      `json:"jsonrpc"`
	ID      string      `json:"id"`
	Method  string      `json:"method"`
	Params  interface{} `json:"params,omitempty"`
}

// MCPResponse represents a JSON-RPC response
type MCPResponse struct {
	JSONRPC string      `json:"jsonrpc"`
	ID      string      `json:"id"`
	Result  interface{} `json:"result,omitempty"`
	Error   interface{} `json:"error,omitempty"`
}

// InitializeParams represents initialize request parameters
type InitializeParams struct {
	ProtocolVersion string                 `json:"protocolVersion"`
	Capabilities    map[string]interface{} `json:"capabilities"`
	ClientInfo      map[string]interface{} `json:"clientInfo"`
}

// sendMCPRequest spawns the Go binary and sends an MCP request over stdio
func sendMCPRequest(t *testing.T, commandArgs []string, request MCPRequest, timeout time.Duration) MCPResponse {
	// Build the project first
	projectRoot, err := filepath.Abs("..")
	require.NoError(t, err)

	// Ensure bin directory exists
	binDir := filepath.Join(projectRoot, "bin")
	err = os.MkdirAll(binDir, 0755)
	require.NoError(t, err)

	// Build to bin/studio
	binaryPath := filepath.Join(binDir, "studio")
	buildCmd := exec.Command("go", "build", "-o", binaryPath, ".")
	buildCmd.Dir = projectRoot
	err = buildCmd.Run()
	require.NoError(t, err, "Failed to build project")

	// Prepare the command
	args := append([]string{}, commandArgs...)
	cmd := exec.Command(binaryPath, args...)

	// Set up pipes
	stdin, err := cmd.StdinPipe()
	require.NoError(t, err)

	stdout, err := cmd.StdoutPipe()
	require.NoError(t, err)

	stderr, err := cmd.StderrPipe()
	require.NoError(t, err)

	// Start the process
	err = cmd.Start()
	require.NoError(t, err)

	defer func() {
		cmd.Process.Kill()
		cmd.Wait()
	}()

	// For tools/ methods, we need to initialize first
	needsInit := strings.HasPrefix(request.Method, "tools/")

	if needsInit {
		// Send initialize request first
		initRequest := MCPRequest{
			JSONRPC: "2.0",
			ID:      "init",
			Method:  "initialize",
			Params: InitializeParams{
				ProtocolVersion: "2024-11-05",
				Capabilities:    map[string]interface{}{},
				ClientInfo: map[string]interface{}{
					"name":    "test-client",
					"version": "1.0.0",
				},
			},
		}

		initJSON, err := json.Marshal(initRequest)
		require.NoError(t, err)

		_, err = stdin.Write(append(initJSON, '\n'))
		require.NoError(t, err)
	}

	// Send the actual request
	requestJSON, err := json.Marshal(request)
	require.NoError(t, err)

	_, err = stdin.Write(append(requestJSON, '\n'))
	require.NoError(t, err)
	stdin.Close()

	// Read response with timeout
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	responseData := make(chan []byte, 10) // Buffer multiple responses
	errorData := make(chan string, 1)

	go func() {
		scanner := bufio.NewScanner(stdout)
		for scanner.Scan() {
			line := scanner.Text()
			if strings.TrimSpace(line) != "" {
				responseData <- []byte(line)
			}
		}
		close(responseData)
	}()

	go func() {
		scanner := bufio.NewScanner(stderr)
		var errLines []string
		for scanner.Scan() {
			errLines = append(errLines, scanner.Text())
		}
		if len(errLines) > 0 {
			errorData <- strings.Join(errLines, "\n")
		}
	}()

	// Collect all responses and find the one matching our request ID
	var targetResponse MCPResponse
	found := false

	for {
		select {
		case data, ok := <-responseData:
			if !ok {
				// Channel closed, we're done
				if !found {
					t.Fatalf("Did not receive response for request ID %s", request.ID)
				}
				return targetResponse
			}

			var response MCPResponse
			err := json.Unmarshal(data, &response)
			require.NoError(t, err, "Failed to parse JSON response: %s", string(data))

			// Check if this is the response we're looking for
			if response.ID == request.ID {
				targetResponse = response
				found = true
			}

		case errMsg := <-errorData:
			t.Fatalf("Process error: %s", errMsg)
		case <-ctx.Done():
			t.Fatalf("Request timed out after %v", timeout)
		}
	}
}

func TestStudioMCPServerIntegration(t *testing.T) {
	timeout := 5 * time.Second

	t.Run("BasicMCPProtocol", func(t *testing.T) {
		t.Run("responds to ping requests", func(t *testing.T) {
			request := MCPRequest{
				JSONRPC: "2.0",
				ID:      "1",
				Method:  "ping",
			}

			response := sendMCPRequest(t, []string{"echo", "hello"}, request, timeout)

			assert.Equal(t, "2.0", response.JSONRPC)
			assert.Equal(t, "1", response.ID)
			assert.Equal(t, map[string]interface{}{}, response.Result)
		})

		t.Run("responds to initialize requests", func(t *testing.T) {
			request := MCPRequest{
				JSONRPC: "2.0",
				ID:      "1",
				Method:  "initialize",
				Params: InitializeParams{
					ProtocolVersion: "2024-11-05",
					Capabilities:    map[string]interface{}{},
					ClientInfo: map[string]interface{}{
						"name":    "test-client",
						"version": "1.0.0",
					},
				},
			}

			response := sendMCPRequest(t, []string{"echo", "hello"}, request, timeout)

			assert.Equal(t, "2.0", response.JSONRPC)
			assert.Equal(t, "1", response.ID)

			result, ok := response.Result.(map[string]interface{})
			require.True(t, ok, "Result should be an object")

			assert.Contains(t, result, "protocolVersion")
			assert.Contains(t, result, "capabilities")
			assert.Contains(t, result, "serverInfo")

			serverInfo, ok := result["serverInfo"].(map[string]interface{})
			require.True(t, ok, "serverInfo should be an object")
			assert.Equal(t, "studio", serverInfo["name"])
		})
	})

	t.Run("ToolsFunctionality", func(t *testing.T) {
		t.Run("with simple echo command", func(t *testing.T) {
			t.Run("lists available tools", func(t *testing.T) {
				request := MCPRequest{
					JSONRPC: "2.0",
					ID:      "2",
					Method:  "tools/list",
				}

				response := sendMCPRequest(t, []string{"echo", "hello", "[args...]"}, request, timeout)

				assert.Equal(t, "2.0", response.JSONRPC)
				assert.Equal(t, "2", response.ID)

				result, ok := response.Result.(map[string]interface{})
				require.True(t, ok)

				tools, ok := result["tools"].([]interface{})
				require.True(t, ok)
				assert.Len(t, tools, 1)

				tool, ok := tools[0].(map[string]interface{})
				require.True(t, ok)
				assert.Equal(t, "echo", tool["name"])
				assert.Equal(t, "Run the shell command `echo hello [args...]`", tool["description"])

				inputSchema, ok := tool["inputSchema"].(map[string]interface{})
				require.True(t, ok)
				assert.Equal(t, "object", inputSchema["type"])

				properties, ok := inputSchema["properties"].(map[string]interface{})
				require.True(t, ok)
				assert.Contains(t, properties, "args")
			})

			t.Run("executes simple echo command", func(t *testing.T) {
				request := MCPRequest{
					JSONRPC: "2.0",
					ID:      "3",
					Method:  "tools/call",
					Params: map[string]interface{}{
						"name": "echo",
						"arguments": map[string]interface{}{
							"args": []string{"hello", "world"},
						},
					},
				}

				response := sendMCPRequest(t, []string{"echo", "hello", "[args...]"}, request, timeout)

				assert.Equal(t, "2.0", response.JSONRPC)
				assert.Equal(t, "3", response.ID)

				result, ok := response.Result.(map[string]interface{})
				require.True(t, ok)

				content, ok := result["content"].([]interface{})
				require.True(t, ok)
				require.Len(t, content, 1)

				textContent, ok := content[0].(map[string]interface{})
				require.True(t, ok)
				assert.Equal(t, "text", textContent["type"])
				assert.Equal(t, "hello hello world", textContent["text"])

				// isError should be false or omitted (omitted means false)
				if isError, exists := result["isError"]; exists {
					assert.Equal(t, false, isError)
				}
			})
		})

		t.Run("with blueprinted echo command", func(t *testing.T) {
			t.Run("lists blueprinted tool with proper schema", func(t *testing.T) {
				request := MCPRequest{
					JSONRPC: "2.0",
					ID:      "4",
					Method:  "tools/list",
				}

				response := sendMCPRequest(t, []string{"echo", "{{text#the text to echo}}"}, request, timeout)

				assert.Equal(t, "2.0", response.JSONRPC)
				assert.Equal(t, "4", response.ID)

				result, ok := response.Result.(map[string]interface{})
				require.True(t, ok)

				tools, ok := result["tools"].([]interface{})
				require.True(t, ok)
				require.Len(t, tools, 1)

				tool, ok := tools[0].(map[string]interface{})
				require.True(t, ok)
				assert.Equal(t, "echo", tool["name"])
				assert.Equal(t, "Run the shell command `echo {{text}}`", tool["description"])

				inputSchema, ok := tool["inputSchema"].(map[string]interface{})
				require.True(t, ok)

				properties, ok := inputSchema["properties"].(map[string]interface{})
				require.True(t, ok)
				assert.Contains(t, properties, "text")
				assert.NotContains(t, properties, "args")

				textProp, ok := properties["text"].(map[string]interface{})
				require.True(t, ok)
				assert.Equal(t, "string", textProp["type"])
				assert.Equal(t, "the text to echo", textProp["description"])
			})

			t.Run("executes blueprinted command", func(t *testing.T) {
				request := MCPRequest{
					JSONRPC: "2.0",
					ID:      "5",
					Method:  "tools/call",
					Params: map[string]interface{}{
						"name": "echo",
						"arguments": map[string]interface{}{
							"text": "Hello Blueprint!",
						},
					},
				}

				response := sendMCPRequest(t, []string{"echo", "{{text#the text to echo}}"}, request, timeout)

				assert.Equal(t, "2.0", response.JSONRPC)
				assert.Equal(t, "5", response.ID)

				result, ok := response.Result.(map[string]interface{})
				require.True(t, ok)

				content, ok := result["content"].([]interface{})
				require.True(t, ok)
				require.Len(t, content, 1)

				textContent, ok := content[0].(map[string]interface{})
				require.True(t, ok)
				assert.Equal(t, "text", textContent["type"])
				assert.Equal(t, "Hello Blueprint!", textContent["text"])

				// isError should be false or omitted (omitted means false)
				if isError, exists := result["isError"]; exists {
					assert.Equal(t, false, isError)
				}
			})

			t.Run("executes blueprinted command with additional args", func(t *testing.T) {
				request := MCPRequest{
					JSONRPC: "2.0",
					ID:      "6",
					Method:  "tools/call",
					Params: map[string]interface{}{
						"name": "echo",
						"arguments": map[string]interface{}{
							"text": "Hello",
							"args": []string{"World", "from", "args"},
						},
					},
				}

				response := sendMCPRequest(t, []string{"echo", "{{text#the text to echo}}", "[args...]"}, request, timeout)

				assert.Equal(t, "2.0", response.JSONRPC)
				assert.Equal(t, "6", response.ID)

				result, ok := response.Result.(map[string]interface{})
				require.True(t, ok)

				content, ok := result["content"].([]interface{})
				require.True(t, ok)
				require.Len(t, content, 1)

				textContent, ok := content[0].(map[string]interface{})
				require.True(t, ok)
				assert.Equal(t, "text", textContent["type"])
				assert.Equal(t, "Hello World from args", textContent["text"])

				if isError, exists := result["isError"]; exists {
					assert.Equal(t, false, isError)
				}
			})

			t.Run("handles blueprint arguments with spaces", func(t *testing.T) {
				request := MCPRequest{
					JSONRPC: "2.0",
					ID:      "7",
					Method:  "tools/call",
					Params: map[string]interface{}{
						"name": "echo",
						"arguments": map[string]interface{}{
							"text": "Hello World with spaces",
						},
					},
				}

				response := sendMCPRequest(t, []string{"echo", "{{text#the text to echo}}"}, request, timeout)

				assert.Equal(t, "2.0", response.JSONRPC)
				assert.Equal(t, "7", response.ID)

				result, ok := response.Result.(map[string]interface{})
				require.True(t, ok)

				content, ok := result["content"].([]interface{})
				require.True(t, ok)
				require.Len(t, content, 1)

				textContent, ok := content[0].(map[string]interface{})
				require.True(t, ok)
				assert.Equal(t, "text", textContent["type"])
				assert.Equal(t, "Hello World with spaces", textContent["text"])

				if isError, exists := result["isError"]; exists {
					assert.Equal(t, false, isError)
				}
			})

			t.Run("handles mixed blueprint with spaces", func(t *testing.T) {
				request := MCPRequest{
					JSONRPC: "2.0",
					ID:      "8",
					Method:  "tools/call",
					Params: map[string]interface{}{
						"name": "echo",
						"arguments": map[string]interface{}{
							"text": "Hello World",
						},
					},
				}

				response := sendMCPRequest(t, []string{"echo", "simon says {{text#the text to echo}}"}, request, timeout)

				assert.Equal(t, "2.0", response.JSONRPC)
				assert.Equal(t, "8", response.ID)

				result, ok := response.Result.(map[string]interface{})
				require.True(t, ok)

				content, ok := result["content"].([]interface{})
				require.True(t, ok)
				require.Len(t, content, 1)

				textContent, ok := content[0].(map[string]interface{})
				require.True(t, ok)
				assert.Equal(t, "text", textContent["type"])
				assert.Equal(t, "simon says Hello World", textContent["text"])

				if isError, exists := result["isError"]; exists {
					assert.Equal(t, false, isError)
				}
			})

			t.Run("handles blueprint definitions with spaces around hash", func(t *testing.T) {
				request := MCPRequest{
					JSONRPC: "2.0",
					ID:      "9",
					Method:  "tools/call",
					Params: map[string]interface{}{
						"name": "echo",
						"arguments": map[string]interface{}{
							"text": "Hello World",
						},
					},
				}

				response := sendMCPRequest(t, []string{"echo", "{{text # the text to echo}}"}, request, timeout)

				assert.Equal(t, "2.0", response.JSONRPC)
				assert.Equal(t, "9", response.ID)

				result, ok := response.Result.(map[string]interface{})
				require.True(t, ok)

				content, ok := result["content"].([]interface{})
				require.True(t, ok)
				require.Len(t, content, 1)

				textContent, ok := content[0].(map[string]interface{})
				require.True(t, ok)
				assert.Equal(t, "text", textContent["type"])
				assert.Equal(t, "Hello World", textContent["text"])

				if isError, exists := result["isError"]; exists {
					assert.Equal(t, false, isError)
				}
			})

			t.Run("handles blueprints without descriptions", func(t *testing.T) {
				request := MCPRequest{
					JSONRPC: "2.0",
					ID:      "10",
					Method:  "tools/list",
				}

				response := sendMCPRequest(t, []string{"echo", "{{text}}"}, request, timeout)

				result, ok := response.Result.(map[string]interface{})
				require.True(t, ok)

				tools, ok := result["tools"].([]interface{})
				require.True(t, ok)
				require.Len(t, tools, 1)

				tool, ok := tools[0].(map[string]interface{})
				require.True(t, ok)
				assert.Equal(t, "echo", tool["name"])
				assert.Equal(t, "Run the shell command `echo {{text}}`", tool["description"])

				inputSchema, ok := tool["inputSchema"].(map[string]interface{})
				require.True(t, ok)

				properties, ok := inputSchema["properties"].(map[string]interface{})
				require.True(t, ok)
				assert.Contains(t, properties, "text")

				textProp, ok := properties["text"].(map[string]interface{})
				require.True(t, ok)
				assert.Equal(t, "string", textProp["type"])
				// Description should be undefined/omitted when not provided
				_, hasDescription := textProp["description"]
				assert.False(t, hasDescription)
			})

			t.Run("executes blueprints without descriptions", func(t *testing.T) {
				request := MCPRequest{
					JSONRPC: "2.0",
					ID:      "11",
					Method:  "tools/call",
					Params: map[string]interface{}{
						"name": "echo",
						"arguments": map[string]interface{}{
							"text": "Hello Blueprint!",
						},
					},
				}

				response := sendMCPRequest(t, []string{"echo", "{{text}}"}, request, timeout)

				assert.Equal(t, "2.0", response.JSONRPC)
				assert.Equal(t, "11", response.ID)

				result, ok := response.Result.(map[string]interface{})
				require.True(t, ok)

				content, ok := result["content"].([]interface{})
				require.True(t, ok)
				require.Len(t, content, 1)

				textContent, ok := content[0].(map[string]interface{})
				require.True(t, ok)
				assert.Equal(t, "text", textContent["type"])
				assert.Equal(t, "Hello Blueprint!", textContent["text"])

				if isError, exists := result["isError"]; exists {
					assert.Equal(t, false, isError)
				}
			})
		})
	})

	t.Run("ErrorHandling", func(t *testing.T) {
		t.Run("handles command errors gracefully", func(t *testing.T) {
			request := MCPRequest{
				JSONRPC: "2.0",
				ID:      "12",
				Method:  "tools/call",
				Params: map[string]interface{}{
					"name":      "false",
					"arguments": map[string]interface{}{},
				},
			}

			response := sendMCPRequest(t, []string{"false"}, request, timeout)

			assert.Equal(t, "2.0", response.JSONRPC)
			assert.Equal(t, "12", response.ID)

			result, ok := response.Result.(map[string]interface{})
			require.True(t, ok)
			assert.Equal(t, true, result["isError"])

			content, ok := result["content"].([]interface{})
			require.True(t, ok)
			require.Len(t, content, 1)
		})

		t.Run("handles commands with no output", func(t *testing.T) {
			request := MCPRequest{
				JSONRPC: "2.0",
				ID:      "12a",
				Method:  "tools/call",
				Params: map[string]interface{}{
					"name":      "true",
					"arguments": map[string]interface{}{},
				},
			}

			response := sendMCPRequest(t, []string{"true"}, request, timeout)

			assert.Equal(t, "2.0", response.JSONRPC)
			assert.Equal(t, "12a", response.ID)

			result, ok := response.Result.(map[string]interface{})
			require.True(t, ok)

			// Should not be an error
			if isError, exists := result["isError"]; exists {
				assert.Equal(t, false, isError)
			}

			content, ok := result["content"].([]interface{})
			require.True(t, ok)
			require.Len(t, content, 1)

			textContent, ok := content[0].(map[string]interface{})
			require.True(t, ok)
			assert.Equal(t, "text", textContent["type"])
			assert.Equal(t, "", textContent["text"]) // Empty string, but field must exist
		})

		t.Run("handles missing required arguments", func(t *testing.T) {
			request := MCPRequest{
				JSONRPC: "2.0",
				ID:      "12b",
				Method:  "tools/call",
				Params: map[string]interface{}{
					"name":      "echo",
					"arguments": map[string]interface{}{
						// Missing required "text" argument
					},
				},
			}

			response := sendMCPRequest(t, []string{"echo", "{{text#Required text to echo}}"}, request, timeout)

			assert.Equal(t, "2.0", response.JSONRPC)
			assert.Equal(t, "12b", response.ID)

			result, ok := response.Result.(map[string]interface{})
			require.True(t, ok)
			assert.Equal(t, true, result["isError"])

			content, ok := result["content"].([]interface{})
			require.True(t, ok)
			require.Len(t, content, 1)

			textContent, ok := content[0].(map[string]interface{})
			require.True(t, ok)
			assert.Equal(t, "text", textContent["type"])
			assert.Contains(t, textContent["text"], "Validation error: missing required parameter: text")
		})

		t.Run("handles nonexistent tools", func(t *testing.T) {
			request := MCPRequest{
				JSONRPC: "2.0",
				ID:      "13",
				Method:  "tools/call",
				Params: map[string]interface{}{
					"name":      "nonexistent",
					"arguments": map[string]interface{}{},
				},
			}

			response := sendMCPRequest(t, []string{"echo", "hello"}, request, timeout)

			assert.Equal(t, "2.0", response.JSONRPC)
			assert.Equal(t, "13", response.ID)
			assert.NotNil(t, response.Error)
		})
	})

	t.Run("ComplexCommands", func(t *testing.T) {
		t.Run("works with git commands", func(t *testing.T) {
			request := MCPRequest{
				JSONRPC: "2.0",
				ID:      "14",
				Method:  "tools/list",
			}

			response := sendMCPRequest(t, []string{"git", "status", "[args...]"}, request, timeout)

			result, ok := response.Result.(map[string]interface{})
			require.True(t, ok)

			tools, ok := result["tools"].([]interface{})
			require.True(t, ok)
			require.Len(t, tools, 1)

			tool, ok := tools[0].(map[string]interface{})
			require.True(t, ok)
			assert.Equal(t, "git", tool["name"])
			assert.Equal(t, "Run the shell command `git status [args...]`", tool["description"])
		})

		t.Run("works with multiple blueprint arguments", func(t *testing.T) {
			request := MCPRequest{
				JSONRPC: "2.0",
				ID:      "15",
				Method:  "tools/list",
			}

			response := sendMCPRequest(t, []string{
				"rails", "generate",
				"{{generator#Rails generator name}}",
				"{{name#Resource name}}",
			}, request, timeout)

			result, ok := response.Result.(map[string]interface{})
			require.True(t, ok)

			tools, ok := result["tools"].([]interface{})
			require.True(t, ok)
			require.Len(t, tools, 1)

			tool, ok := tools[0].(map[string]interface{})
			require.True(t, ok)
			assert.Equal(t, "rails", tool["name"])
			assert.Equal(t, "Run the shell command `rails generate {{generator}} {{name}}`", tool["description"])

			inputSchema, ok := tool["inputSchema"].(map[string]interface{})
			require.True(t, ok)

			properties, ok := inputSchema["properties"].(map[string]interface{})
			require.True(t, ok)
			assert.Contains(t, properties, "generator")
			assert.Contains(t, properties, "name")
			assert.NotContains(t, properties, "args")
		})
	})
}

// TestArgumentParsingRegression tests the specific issue where flags in command
// templates (like -v in "say -v siri") were incorrectly parsed as studio flags
func TestArgumentParsingRegression(t *testing.T) {
	timeout := 5 * time.Second

	t.Run("SayCommandWithVoiceFlag", func(t *testing.T) {
		t.Run("should not parse -v as studio flag", func(t *testing.T) {
			request := MCPRequest{
				JSONRPC: "2.0",
				ID:      "say-test-1",
				Method:  "tools/list",
			}

			// This should work without errors - previously this would fail with
			// "Error: unknown shorthand flag: 'v' in -v"
			response := sendMCPRequest(t, []string{
				"say", "-v", "siri", "{{speech#A very concise message to say out loud to the user}}",
			}, request, timeout)

			assert.Equal(t, "2.0", response.JSONRPC)
			assert.Equal(t, "say-test-1", response.ID)
			assert.Nil(t, response.Error, "Should not have parsing errors")

			result, ok := response.Result.(map[string]interface{})
			require.True(t, ok)

			tools, ok := result["tools"].([]interface{})
			require.True(t, ok)
			require.Len(t, tools, 1)

			tool, ok := tools[0].(map[string]interface{})
			require.True(t, ok)
			assert.Equal(t, "say", tool["name"])
			assert.Equal(t, "Run the shell command `say -v siri {{speech}}`", tool["description"])

			inputSchema, ok := tool["inputSchema"].(map[string]interface{})
			require.True(t, ok)

			properties, ok := inputSchema["properties"].(map[string]interface{})
			require.True(t, ok)
			assert.Contains(t, properties, "speech")

			speechProp, ok := properties["speech"].(map[string]interface{})
			require.True(t, ok)
			assert.Equal(t, "string", speechProp["type"])
			assert.Equal(t, "A very concise message to say out loud to the user", speechProp["description"])
		})

		t.Run("should work with debug flag before say command", func(t *testing.T) {
			request := MCPRequest{
				JSONRPC: "2.0",
				ID:      "say-test-2",
				Method:  "tools/list",
			}

			// Test that --debug flag is parsed correctly while -v is not
			response := sendMCPRequest(t, []string{
				"--debug", "say", "-v", "siri", "{{speech#message}}",
			}, request, timeout)

			assert.Equal(t, "2.0", response.JSONRPC)
			assert.Equal(t, "say-test-2", response.ID)
			assert.Nil(t, response.Error, "Should not have parsing errors")

			result, ok := response.Result.(map[string]interface{})
			require.True(t, ok)

			tools, ok := result["tools"].([]interface{})
			require.True(t, ok)
			require.Len(t, tools, 1)

			tool, ok := tools[0].(map[string]interface{})
			require.True(t, ok)
			assert.Equal(t, "say", tool["name"])
			assert.Equal(t, "Run the shell command `say -v siri {{speech}}`", tool["description"])
		})
	})

	t.Run("CurlCommandWithMultipleFlags", func(t *testing.T) {
		t.Run("should handle multiple flags in command template", func(t *testing.T) {
			request := MCPRequest{
				JSONRPC: "2.0",
				ID:      "curl-test-1",
				Method:  "tools/list",
			}

			// Test a more complex command with multiple flags
			response := sendMCPRequest(t, []string{
				"curl", "-X", "POST", "-H", "Content-Type: application/json", "{{url#The URL to request}}",
			}, request, timeout)

			assert.Equal(t, "2.0", response.JSONRPC)
			assert.Equal(t, "curl-test-1", response.ID)
			assert.Nil(t, response.Error, "Should not have parsing errors")

			result, ok := response.Result.(map[string]interface{})
			require.True(t, ok)

			tools, ok := result["tools"].([]interface{})
			require.True(t, ok)
			require.Len(t, tools, 1)

			tool, ok := tools[0].(map[string]interface{})
			require.True(t, ok)
			assert.Equal(t, "curl", tool["name"])
			assert.Equal(t, "Run the shell command `curl -X POST -H Content-Type: application/json {{url}}`", tool["description"])
		})
	})
}
