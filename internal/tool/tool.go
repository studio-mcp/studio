package tool

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/modelcontextprotocol/go-sdk/jsonschema"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// Blueprint interface defines what we need from a blueprint
type Blueprint interface {
	BuildCommandArgs(args map[string]interface{}) ([]string, error)
	GetBaseCommand() string
	GetCommandFormat() string
	GetInputSchema() interface{}
}

var debugMode bool
var logFile *os.File
var logger *log.Logger

// SetDebugMode enables or disables debug mode
func SetDebugMode(enabled bool) {
	debugMode = enabled
}

// SetLogFile sets the log file for debug output
func SetLogFile(filename string) error {
	// Close existing log file if open
	if logFile != nil {
		logFile.Close()
	}

	// Open new log file
	file, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return fmt.Errorf("failed to open log file %s: %w", filename, err)
	}

	logFile = file
	logger = log.New(file, "[Studio MCP] ", log.LstdFlags)
	return nil
}

// IsDebugMode returns whether debug mode is enabled
func IsDebugMode() bool {
	return debugMode && os.Getenv("NODE_ENV") != "test"
}

// debug logs a message to stderr if debug mode is enabled
func debug(format string, args ...interface{}) {
	if IsDebugMode() {
		if logger != nil {
			logger.Printf(format, args...)
		} else {
			fmt.Fprintf(os.Stderr, "[Studio MCP] "+format+"\n", args...)
		}
	}
}

// Execute runs a command and returns trimmed combined stdout+stderr or an error
func Execute(command string, args ...string) (string, error) {
	debug("Executing command: %s %s", command, strings.Join(args, " "))

	cmd := exec.Command(command, args...)

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()

	// Always combine outputs for visibility
	output := strings.TrimSpace(stdout.String() + "\n" + stderr.String())

	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			debug("Command completed with non-zero exit code: %d", exitErr.ExitCode())
			debug("Final output length: %d chars", len(output))
			return output, fmt.Errorf("command failed with exit code %d", exitErr.ExitCode())
		}
		debug("Spawn error: %s", err.Error())
		return output, fmt.Errorf("Studio error: %w", err)
	}

	debug("Command completed successfully with exit code 0")
	debug("Final output length: %d chars", len(output))

	return output, nil
}

// CreateToolFunction creates a tool handler for the given blueprint
func CreateToolFunction(blueprint Blueprint) mcp.ToolHandlerFor[map[string]any, map[string]any] {
	return func(ctx context.Context, session *mcp.ServerSession, params *mcp.CallToolParamsFor[map[string]any]) (*mcp.CallToolResultFor[map[string]any], error) {
		debug("Tool called with args: %v", params.Arguments)

		fullCommand, err := blueprint.BuildCommandArgs(params.Arguments)
		if err != nil {
			return createToolResult(fmt.Sprintf("Validation error: %s", err.Error()), true), nil
		}

		debug("Built command: %s", strings.Join(fullCommand, " "))

		output, err := Execute(fullCommand[0], fullCommand[1:]...)
		isError := err != nil

		if isError {
			debug("Execution error: %s", err)
		}

		return createToolResult(output, isError), nil
	}
}

// GenerateToolName generates a tool name from a base command by replacing dashes with underscores
func GenerateToolName(baseCommand string) string {
	return strings.ReplaceAll(baseCommand, "-", "_")
}

// CreateServerTool creates a complete MCP server tool from a blueprint
func CreateServerTool(blueprint Blueprint) *mcp.ServerTool {
	schema, ok := blueprint.GetInputSchema().(*jsonschema.Schema)
	if !ok {
		// This should never happen if the Blueprint interface is implemented correctly
		panic("blueprint.GetInputSchema() must return *jsonschema.Schema")
	}

	// Debug logging
	debug("CreateServerTool called")
	debug("  Schema type: %s", schema.Type)
	debug("  Schema properties count: %d", len(schema.Properties))
	for name, prop := range schema.Properties {
		debug("    property %s: type=%s, description=%s", name, prop.Type, prop.Description)
	}
	debug("  Schema required count: %d", len(schema.Required))
	for _, req := range schema.Required {
		debug("    required: %s", req)
	}

	return mcp.NewServerTool(
		GenerateToolName(blueprint.GetBaseCommand()),
		GetToolDescription(blueprint),
		CreateToolFunction(blueprint),
		mcp.Input(mcp.Schema(schema)),
	)
}

func createToolResult(output string, isError bool) *mcp.CallToolResultFor[map[string]any] {
	return &mcp.CallToolResultFor[map[string]any]{
		Content: []mcp.Content{
			&mcp.TextContent{Text: output},
		},
		IsError: isError,
	}
}

// GetToolDescription generates the tool description from a blueprint
func GetToolDescription(blueprint Blueprint) string {
	return "Run the shell command `" + blueprint.GetCommandFormat() + "`"
}
