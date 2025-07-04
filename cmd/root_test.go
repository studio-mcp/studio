package cmd

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseArgs(t *testing.T) {
	tests := []struct {
		name            string
		args            []string
		expectedDebug   bool
		expectedVersion bool
		expectedLogFile string
		expectedCommand []string
		expectedError   string
	}{
		{
			name:            "no flags, simple command",
			args:            []string{"echo", "hello"},
			expectedDebug:   false,
			expectedVersion: false,
			expectedLogFile: "",
			expectedCommand: []string{"echo", "hello"},
		},
		{
			name:            "debug flag before command",
			args:            []string{"--debug", "echo", "hello"},
			expectedDebug:   true,
			expectedVersion: false,
			expectedLogFile: "",
			expectedCommand: []string{"echo", "hello"},
		},
		{
			name:            "version flag only",
			args:            []string{"--version"},
			expectedDebug:   false,
			expectedVersion: true,
			expectedLogFile: "",
			expectedCommand: []string{},
		},
		{
			name:            "log flag with filename",
			args:            []string{"--log", "debug.log", "echo", "hello"},
			expectedDebug:   false,
			expectedVersion: false,
			expectedLogFile: "debug.log",
			expectedCommand: []string{"echo", "hello"},
		},
		{
			name:            "debug and log flags together",
			args:            []string{"--debug", "--log", "app.log", "echo", "hello"},
			expectedDebug:   true,
			expectedVersion: false,
			expectedLogFile: "app.log",
			expectedCommand: []string{"echo", "hello"},
		},
		{
			name:            "log flag with debug flag in different order",
			args:            []string{"--log", "test.log", "--debug", "echo", "hello"},
			expectedDebug:   true,
			expectedVersion: false,
			expectedLogFile: "test.log",
			expectedCommand: []string{"echo", "hello"},
		},
		{
			name:          "log flag followed by another flag",
			args:          []string{"--log", "--debug", "echo", "hello"},
			expectedError: "--log requires a filename argument",
		},
		{
			name:          "log flag at end without filename",
			args:          []string{"--log"},
			expectedError: "--log requires a filename argument",
		},
		{
			name:            "double dash terminates flag parsing",
			args:            []string{"--debug", "--", "echo", "--version"},
			expectedDebug:   true,
			expectedVersion: false,
			expectedLogFile: "",
			expectedCommand: []string{"echo", "--version"},
		},
		{
			name:            "double dash with log flag",
			args:            []string{"--log", "test.log", "--", "echo", "--debug"},
			expectedDebug:   false,
			expectedVersion: false,
			expectedLogFile: "test.log",
			expectedCommand: []string{"echo", "--debug"},
		},
		{
			name:            "double dash at start",
			args:            []string{"--", "--debug", "echo", "hello"},
			expectedDebug:   false,
			expectedVersion: false,
			expectedLogFile: "",
			expectedCommand: []string{"--debug", "echo", "hello"},
		},
		{
			name:            "double dash only",
			args:            []string{"--"},
			expectedDebug:   false,
			expectedVersion: false,
			expectedLogFile: "",
			expectedCommand: []string{},
		},
		{
			name:            "double dash prevents log flag parsing",
			args:            []string{"--", "--log", "test.log", "echo", "hello"},
			expectedDebug:   false,
			expectedVersion: false,
			expectedLogFile: "",
			expectedCommand: []string{"--log", "test.log", "echo", "hello"},
		},
		{
			name:            "command with flags that should not be parsed by studio",
			args:            []string{"say", "-v", "siri", "{{speech#message}}"},
			expectedDebug:   false,
			expectedVersion: false,
			expectedLogFile: "",
			expectedCommand: []string{"say", "-v", "siri", "{{speech#message}}"},
		},
		{
			name:            "debug flag before command with flags",
			args:            []string{"--debug", "say", "-v", "siri", "{{speech#message}}"},
			expectedDebug:   true,
			expectedVersion: false,
			expectedLogFile: "",
			expectedCommand: []string{"say", "-v", "siri", "{{speech#message}}"},
		},
		{
			name:            "command with multiple flags",
			args:            []string{"curl", "-X", "POST", "-H", "Content-Type: application/json", "{{url}}"},
			expectedDebug:   false,
			expectedVersion: false,
			expectedLogFile: "",
			expectedCommand: []string{"curl", "-X", "POST", "-H", "Content-Type: application/json", "{{url}}"},
		},
		{
			name:          "unknown studio flag",
			args:          []string{"--unknown", "echo", "hello"},
			expectedError: "unknown flag: --unknown",
		},
		{
			name:          "help flag",
			args:          []string{"--help"},
			expectedError: "help requested",
		},
		{
			name:          "help flag short",
			args:          []string{"-h"},
			expectedError: "help requested",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			debug, version, logFile, command, err := parseArgs(tt.args)

			if tt.expectedError != "" {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError)
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, tt.expectedDebug, debug)
			assert.Equal(t, tt.expectedVersion, version)
			assert.Equal(t, tt.expectedLogFile, logFile)
			assert.Equal(t, tt.expectedCommand, command)
		})
	}
}

func TestVersionFlagParsing(t *testing.T) {
	t.Run("identifies version flag correctly", func(t *testing.T) {
		debug, version, logFile, command, err := parseArgs([]string{"--version"})
		assert.NoError(t, err)
		assert.False(t, debug)
		assert.True(t, version)
		assert.Empty(t, logFile)
		assert.Empty(t, command)
	})
}

func TestEmptyArgs(t *testing.T) {
	t.Run("handles empty args", func(t *testing.T) {
		debug, version, logFile, command, err := parseArgs([]string{})
		assert.NoError(t, err)
		assert.False(t, debug)
		assert.False(t, version)
		assert.Empty(t, logFile)
		assert.Empty(t, command)
	})
}

// Test the specific regression case that was fixed
func TestSayCommandWithVoiceFlag(t *testing.T) {
	t.Run("say command with -v flag should not be parsed as studio flag", func(t *testing.T) {
		args := []string{"say", "-v", "siri", "{{speech#A very concise message to say out loud to the user}}"}

		debug, version, logFile, command, err := parseArgs(args)

		assert.NoError(t, err)
		assert.False(t, debug)
		assert.False(t, version)
		assert.Empty(t, logFile)
		assert.Equal(t, args, command)
	})

	t.Run("debug flag followed by say command with -v flag", func(t *testing.T) {
		args := []string{"--debug", "say", "-v", "siri", "{{speech#message}}"}

		debug, version, logFile, command, err := parseArgs(args)

		assert.NoError(t, err)
		assert.True(t, debug)
		assert.False(t, version)
		assert.Empty(t, logFile)
		assert.Equal(t, []string{"say", "-v", "siri", "{{speech#message}}"}, command)
	})
}
