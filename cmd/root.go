/*
Copyright © 2025 Martin Emde me@martinemde.com

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/studio-mcp/studio/internal/studio"

	"github.com/spf13/cobra"
)

var (
	// Version information set by main
	Version string
	Commit  string
	Date    string
)

// parseArgs parses arguments manually, stopping flag parsing at first non-flag
func parseArgs(args []string) (debugFlag bool, versionFlag bool, logFile string, commandArgs []string, err error) {
	i := 0

	// Parse studio flags until we hit a non-flag or --
	for i < len(args) {
		arg := args[i]

		// If it's --, stop parsing flags and treat everything after as command args
		if arg == "--" {
			i++
			break
		}

		// If it doesn't start with -, we're done with studio flags
		if !strings.HasPrefix(arg, "-") {
			break
		}

		switch arg {
		case "--debug":
			debugFlag = true
		case "--version":
			versionFlag = true
		case "--log":
			// Check if we have a next argument for the filename
			if i+1 >= len(args) {
				return false, false, "", nil, fmt.Errorf("--log requires a filename argument")
			}
			i++
			// Check if the next argument is another flag
			if strings.HasPrefix(args[i], "-") {
				return false, false, "", nil, fmt.Errorf("--log requires a filename argument")
			}
			logFile = args[i]
		case "-h", "--help":
			// Let cobra handle help
			return false, false, "", nil, fmt.Errorf("help requested")
		default:
			return false, false, "", nil, fmt.Errorf("unknown flag: %s", arg)
		}

		i++
	}

	// Everything from i onwards goes to blueprint parsing
	commandArgs = args[i:]

	return debugFlag, versionFlag, logFile, commandArgs, nil
}

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "studio [--debug] [--log filename] [--] <command> --example \"{req # required arg}\" \"[args... # array of args]\"",
	Short: "A tool for running a single command MCP server",
	Long: `studio is a tool for running a single command MCP server.

  -h, --help - Show this help message and exit.
  --version - Show version information and exit.
  --debug - Print debug logs to stderr to diagnose MCP server issues.
  --log <filename> - Write debug logs to the specified file instead of stderr.
  -- - End of flag parsing. Everything after this is treated as command arguments.

the command starts at the first non-flag argument:

  <command> - the shell command to run.
  --example - in the example, this includes the literal argument '--example' is part of the command.
              literal args any shell word.

arguments can be templated as their own shellword or as part of a shellword:

  "{req # required arg}" - tell the LLM about a required arg named 'req'.
  "[args... # array of args]" - tell the LLM about an optional array of args named 'args'.
  "[opt # optional string]" - a optional string arg named 'opt' (not in example).
  "https://en.wikipedia.org/wiki/{wiki_page_name}" - an example partially templated words.

Example:
  studio say -v siri "{speech # a concise phrase to say outloud to the user}"`,
	DisableFlagParsing: true, // Disable cobra's flag parsing so we can do custom parsing
	Args: func(cmd *cobra.Command, args []string) error {
		// Custom argument parsing
		_, versionFlag, _, commandArgs, err := parseArgs(args)
		if err != nil {
			if err.Error() == "help requested" {
				return nil // Let cobra handle help
			}
			return err
		}

		// If version flag is set, don't validate command args
		if versionFlag {
			return nil
		}

		if len(commandArgs) == 0 {
			return fmt.Errorf("usage: studio <command> --example \"{req # required arg}\" \"[args... # array of args]\"")
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		// Parse arguments manually
		debugFlag, versionFlag, logFile, commandArgs, err := parseArgs(args)
		if err != nil {
			if err.Error() == "help requested" {
				return cmd.Help()
			}
			return err
		}

		// Debug logging - log the raw arguments received
		// Write to log file if specified, otherwise stderr
		var logWriter *os.File
		if logFile != "" {
			logWriter, _ = os.OpenFile(logFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
		}

		writeDebug := func(format string, args ...interface{}) {
			msg := fmt.Sprintf("[Studio MCP Root] "+format+"\n", args...)
			if logWriter != nil {
				logWriter.WriteString(msg)
			} else {
				fmt.Fprint(os.Stderr, msg)
			}
		}

		writeDebug("Raw args received: %d arguments", len(args))
		for i, arg := range args {
			writeDebug("  raw[%d]: %q", i, arg)
		}
		writeDebug("Parsed command args: %d arguments", len(commandArgs))
		for i, arg := range commandArgs {
			writeDebug("  cmd[%d]: %q", i, arg)
		}

		if logWriter != nil {
			logWriter.Close()
		}

		// Handle version flag
		if versionFlag {
			cmd.Printf("studio %s\n", Version)
			cmd.Printf("commit: %s\n", Commit)
			cmd.Printf("built: %s\n", Date)
			return nil
		}

		// Create a new Studio instance with the command args
		s, err := studio.New(commandArgs, debugFlag, logFile, Version)
		if err != nil {
			return err
		}

		// Start the MCP server
		return s.Serve()
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute(version, commit, date string) {
	Version = version
	Commit = commit
	Date = date

	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// No flags defined here since we're doing custom parsing
}
