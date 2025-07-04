package studio

import (
	"context"
	"fmt"
	"io"
	"os"

	"github.com/studio-mcp/studio/internal/blueprint"
	"github.com/studio-mcp/studio/internal/tool"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// Studio represents the main application logic
type Studio struct {
	Blueprint *blueprint.Blueprint
	DebugMode bool
	LogFile   string
	Version   string
}

// New creates a new Studio instance from command arguments
func New(args []string, debugMode bool, logFile string, version string) (*Studio, error) {
	if len(args) == 0 {
		return nil, fmt.Errorf("no command provided")
	}

	bp, err := blueprint.FromArgs(args)
	if err != nil {
		return nil, fmt.Errorf("failed to create blueprint: %w", err)
	}

	// Set debug mode and log file on tool
	tool.SetDebugMode(debugMode)
	if logFile != "" {
		err = tool.SetLogFile(logFile)
		if err != nil {
			return nil, fmt.Errorf("failed to set log file: %w", err)
		}
	}

	return &Studio{
		Blueprint: bp,
		DebugMode: debugMode,
		LogFile:   logFile,
		Version:   version,
	}, nil
}

// Serve starts the MCP server over stdio
func (s *Studio) Serve() error {
	return s.ServeWithContext(context.Background())
}

// ServeWithContext starts the MCP server over stdio with a context
func (s *Studio) ServeWithContext(ctx context.Context) error {
	// Create server with version from build
	server := mcp.NewServer("studio", s.Version, nil)

	// Add the tool to the server using CreateServerTool from tool package
	serverTool := tool.CreateServerTool(s.Blueprint)

	server.AddTools(serverTool)

	// Create base transport
	var transport mcp.Transport = mcp.NewStdioTransport()

	// Wrap with logging transport if debug mode is enabled or log file is specified
	if s.DebugMode || s.LogFile != "" {
		var logWriter io.Writer
		if s.LogFile != "" {
			// Open log file for transport logging
			file, err := os.OpenFile(s.LogFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
			if err != nil {
				return fmt.Errorf("failed to open log file for transport logging: %w", err)
			}
			// Don't close the file immediately - let the transport use it
			// The file will be closed when the process exits
			logWriter = file
		} else {
			// Use stderr for transport logging when debug mode is enabled
			logWriter = os.Stderr
		}
		transport = mcp.NewLoggingTransport(transport, logWriter)
	}

	// Run the server with the configured transport
	return server.Run(ctx, transport)
}
