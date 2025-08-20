package blueprint

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBlueprint_BaseCommand(t *testing.T) {
	tests := []struct {
		name     string
		args     []string
		expected string
		wantErr  bool
		errMsg   string
	}{
		{
			name:     "simple command name",
			args:     []string{"git", "status"},
			expected: "git",
		},
		{
			name:     "command name with dashes",
			args:     []string{"git-status"},
			expected: "git-status",
		},
		{
			name:     "command name with underscores",
			args:     []string{"do_thing", "[paths...]"},
			expected: "do_thing",
		},
		{
			name:    "empty args",
			args:    []string{},
			wantErr: true,
			errMsg:  "no command provided",
		},
		{
			name:    "empty command",
			args:    []string{""},
			wantErr: true,
			errMsg:  "empty command provided",
		},
		{
			name:    "whitespace-only command",
			args:    []string{"   "},
			wantErr: true,
			errMsg:  "empty command provided",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bp, err := FromArgs(tt.args)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
				return
			}

			require.NoError(t, err)
			assert.Equal(t, tt.expected, bp.BaseCommand)
		})
	}
}

func TestBlueprint_GetCommandFormat(t *testing.T) {
	tests := []struct {
		name     string
		args     []string
		expected string
	}{
		{
			name:     "command without args",
			args:     []string{"git"},
			expected: "git",
		},
		{
			name:     "command with literal args",
			args:     []string{"git", "status"},
			expected: "git status",
		},
		{
			name:     "command with array args",
			args:     []string{"git", "[args...]"},
			expected: "git [args...]",
		},
		{
			name:     "command with array args",
			args:     []string{"echo", "{{args...}}"},
			expected: "echo {{args...}}",
		},
		{
			name:     "command with literal and array args",
			args:     []string{"git", "status", "[args...]"},
			expected: "git status [args...]",
		},
		{
			name:     "command with array with description",
			args:     []string{"git", "[args... #Additional command line arguments]"},
			expected: "git [args...]",
		},
		{
			name:     "command with required field",
			args:     []string{"echo", "{{text}}"},
			expected: "echo {{text}}",
		},
		{
			name:     "command with required field with description",
			args:     []string{"echo", "{{text#A required field}}"},
			expected: "echo {{text}}",
		},
		{
			name:     "command with optional field",
			args:     []string{"echo", "[text]"},
			expected: "echo [text]",
		},
		{
			name:     "command with optional field with description",
			args:     []string{"echo", "[text#description]"},
			expected: "echo [text]",
		},
		{
			name:     "command with suffix literal",
			args:     []string{"echo", "[text]suffix"},
			expected: "echo [text]suffix",
		},
		{
			name:     "command with prefix literal",
			args:     []string{"echo", "prefix[text]"},
			expected: "echo prefix[text]",
		},
		{
			name:     "command with prefix and suffix literal",
			args:     []string{"echo", "prefix[text]suffix"},
			expected: "echo prefix[text]suffix",
		},
		{
			name:     "command with embedded template",
			args:     []string{"curl", "https://en.m.wikipedia.org/wiki/{{page#A valid wikipedia page}}"},
			expected: "curl https://en.m.wikipedia.org/wiki/{{page}}",
		},
		{
			name:     "command with spaces in field description",
			args:     []string{"curl", "https://en.m.wikipedia.org/wiki/{{page # A valid wikipedia page}}"},
			expected: "curl https://en.m.wikipedia.org/wiki/{{page}}",
		},
		{
			name:     "mixed blueprints with required and optional arguments",
			args:     []string{"command", "{{arg1#Custom description}}", "[arg2]"},
			expected: "command {{arg1}} [arg2]",
		},
		{
			name:     "prioritizes explicit description over default",
			args:     []string{"echo", "{{text#Explicit description}}", "{{text}}"},
			expected: "echo {{text}} {{text}}",
		},
		{
			name:     "preserves underscores in field names for display",
			args:     []string{"echo", "[has_underscores]"},
			expected: "echo [has_underscores]",
		},
		{
			name:     "preserves dashes in field names for display",
			args:     []string{"echo", "[has-dashes]"},
			expected: "echo [has-dashes]",
		},
		{
			name:     "short boolean flag without description",
			args:     []string{"ls", "[-f]"},
			expected: "ls [-f]",
		},
		{
			name:     "long boolean flag without description",
			args:     []string{"ls", "[--force]"},
			expected: "ls [--force]",
		},
		{
			name:     "boolean flag with description",
			args:     []string{"rm", "[-f#force removal]"},
			expected: "rm [-f]",
		},
		{
			name:     "long boolean flag with description",
			args:     []string{"ls", "[--force#force removal]"},
			expected: "ls [--force]",
		},
		{
			name:     "required flag with description",
			args:     []string{"cp", "{{-r#recursive}}"},
			expected: "cp {{-r}}",
		},
		{
			name:     "complicated template with mixed text and fields",
			args:     []string{"curl", "http[s # use https]://api.com/{{endpoint#API endpoint}}", "[--verbose]"},
			expected: "curl http[s]://api.com/{{endpoint}} [--verbose]",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bp, err := FromArgs(tt.args)
			require.NoError(t, err)
			assert.Equal(t, tt.expected, bp.GetCommandFormat())
		})
	}
}

func TestBlueprint_FromArgsTokenization(t *testing.T) {
	t.Run("tokenizes simple command without templates", func(t *testing.T) {
		bp, err := FromArgs([]string{"echo", "hello"})
		require.NoError(t, err)

		expected := [][]Token{
			{TextToken{Value: "echo"}},
			{TextToken{Value: "hello"}},
		}
		assert.Equal(t, expected, bp.ShellWords)
	})

	t.Run("tokenizes command with simple template", func(t *testing.T) {
		bp, err := FromArgs([]string{"echo", "{{text}}"})
		require.NoError(t, err)

		expected := [][]Token{
			{TextToken{Value: "echo"}},
			{FieldToken{Name: "text", Description: "", Required: true, OriginalFlag: ""}},
		}
		assert.Equal(t, expected, bp.ShellWords)
	})

	t.Run("tokenizes command with optional field", func(t *testing.T) {
		bp, err := FromArgs([]string{"echo", "[optional]"})
		require.NoError(t, err)

		expected := [][]Token{
			{TextToken{Value: "echo"}},
			{FieldToken{Name: "optional", Description: "", Required: false, OriginalFlag: ""}},
		}
		assert.Equal(t, expected, bp.ShellWords)
	})

	t.Run("tokenizes command with template and description", func(t *testing.T) {
		bp, err := FromArgs([]string{"echo", "{{text#message to echo}}"})
		require.NoError(t, err)

		expected := [][]Token{
			{TextToken{Value: "echo"}},
			{FieldToken{Name: "text", Description: "message to echo", Required: true, OriginalFlag: ""}},
		}
		assert.Equal(t, expected, bp.ShellWords)
	})

	t.Run("tokenizes command with prefix text and template", func(t *testing.T) {
		bp, err := FromArgs([]string{"echo", "prefix{{text#desc}}"})
		require.NoError(t, err)

		expected := [][]Token{
			{TextToken{Value: "echo"}},
			{
				TextToken{Value: "prefix"},
				FieldToken{Name: "text", Description: "desc", Required: true, OriginalFlag: ""},
			},
		}
		assert.Equal(t, expected, bp.ShellWords)
	})

	t.Run("tokenizes command with suffix text and template", func(t *testing.T) {
		bp, err := FromArgs([]string{"echo", "{{text#desc}}suffix"})
		require.NoError(t, err)

		expected := [][]Token{
			{TextToken{Value: "echo"}},
			{
				FieldToken{Name: "text", Description: "desc", Required: true, OriginalFlag: ""},
				TextToken{Value: "suffix"},
			},
		}
		assert.Equal(t, expected, bp.ShellWords)
	})

	t.Run("tokenizes command with mixed text and template", func(t *testing.T) {
		bp, err := FromArgs([]string{"echo", "prefix{{text#desc}}suffix"})
		require.NoError(t, err)

		expected := [][]Token{
			{TextToken{Value: "echo"}},
			{
				TextToken{Value: "prefix"},
				FieldToken{Name: "text", Description: "desc", Required: true, OriginalFlag: ""},
				TextToken{Value: "suffix"},
			},
		}
		assert.Equal(t, expected, bp.ShellWords)
	})

	t.Run("tokenizes command with prefix text and optional field", func(t *testing.T) {
		bp, err := FromArgs([]string{"echo", "prefix[text]"})
		require.NoError(t, err)

		expected := [][]Token{
			{TextToken{Value: "echo"}},
			{
				TextToken{Value: "prefix"},
				FieldToken{Name: "text", Description: "", Required: false, OriginalFlag: ""},
			},
		}
		assert.Equal(t, expected, bp.ShellWords)
	})

	t.Run("tokenizes command with suffix text and optional field", func(t *testing.T) {
		bp, err := FromArgs([]string{"echo", "[text#description]suffix"})
		require.NoError(t, err)

		expected := [][]Token{
			{TextToken{Value: "echo"}},
			{
				FieldToken{Name: "text", Description: "description", Required: false, OriginalFlag: ""},
				TextToken{Value: "suffix"},
			},
		}
		assert.Equal(t, expected, bp.ShellWords)
	})

	t.Run("tokenizes command with mixed text and optional field", func(t *testing.T) {
		bp, err := FromArgs([]string{"echo", "prefix[text]suffix"})
		require.NoError(t, err)

		expected := [][]Token{
			{TextToken{Value: "echo"}},
			{
				TextToken{Value: "prefix"},
				FieldToken{Name: "text", Description: "", Required: false, OriginalFlag: ""},
				TextToken{Value: "suffix"},
			},
		}
		assert.Equal(t, expected, bp.ShellWords)
	})

	t.Run("tokenizes command with optional field and description", func(t *testing.T) {
		bp, err := FromArgs([]string{"echo", "[optional#optional text]"})
		require.NoError(t, err)

		expected := [][]Token{
			{TextToken{Value: "echo"}},
			{FieldToken{Name: "optional", Description: "optional text", Required: false, OriginalFlag: ""}},
		}
		assert.Equal(t, expected, bp.ShellWords)
	})

	t.Run("tokenizes complex mixed command", func(t *testing.T) {
		bp, err := FromArgs([]string{"curl", "http[s # use https]://api.com/{{endpoint#API endpoint}}", "[--verbose]"})
		require.NoError(t, err)

		expected := [][]Token{
			{TextToken{Value: "curl"}},
			{
				TextToken{Value: "http"},
				FieldToken{Name: "s", Description: "use https", Required: false, OriginalFlag: ""},
				TextToken{Value: "://api.com/"},
				FieldToken{Name: "endpoint", Description: "API endpoint", Required: true, OriginalFlag: ""},
			},
			{FieldToken{Name: "verbose", Description: "Enable --verbose flag", Required: false, OriginalFlag: "--verbose"}},
		}
		assert.Equal(t, expected, bp.ShellWords)
	})
}

func TestBlueprint_SingleBraceTokenization(t *testing.T) {
	t.Run("tokenizes command with single-brace required field", func(t *testing.T) {
		bp, err := FromArgs([]string{"echo", "{text}"})
		require.NoError(t, err)

		expected := [][]Token{
			{TextToken{Value: "echo"}},
			{FieldToken{Name: "text", Description: "", Required: true, OriginalFlag: ""}},
		}
		assert.Equal(t, expected, bp.ShellWords)
		assert.Equal(t, "echo {{text}}", bp.GetCommandFormat())
	})

	t.Run("tokenizes command with single-brace and description", func(t *testing.T) {
		bp, err := FromArgs([]string{"echo", "{text#message to echo}"})
		require.NoError(t, err)

		expected := [][]Token{
			{TextToken{Value: "echo"}},
			{FieldToken{Name: "text", Description: "message to echo", Required: true, OriginalFlag: ""}},
		}
		assert.Equal(t, expected, bp.ShellWords)
		assert.Equal(t, "echo {{text}}", bp.GetCommandFormat())
	})

	t.Run("tokenizes command with prefix and single-brace template", func(t *testing.T) {
		bp, err := FromArgs([]string{"echo", "prefix{text#desc}"})
		require.NoError(t, err)

		expected := [][]Token{
			{TextToken{Value: "echo"}},
			{
				TextToken{Value: "prefix"},
				FieldToken{Name: "text", Description: "desc", Required: true, OriginalFlag: ""},
			},
		}
		assert.Equal(t, expected, bp.ShellWords)
		assert.Equal(t, "echo prefix{{text}}", bp.GetCommandFormat())
	})

	t.Run("tokenizes command with suffix and single-brace template", func(t *testing.T) {
		bp, err := FromArgs([]string{"echo", "{text#desc}suffix"})
		require.NoError(t, err)

		expected := [][]Token{
			{TextToken{Value: "echo"}},
			{
				FieldToken{Name: "text", Description: "desc", Required: true, OriginalFlag: ""},
				TextToken{Value: "suffix"},
			},
		}
		assert.Equal(t, expected, bp.ShellWords)
		assert.Equal(t, "echo {{text}}suffix", bp.GetCommandFormat())
	})

	t.Run("tokenizes command with mixed text and single-brace template", func(t *testing.T) {
		bp, err := FromArgs([]string{"echo", "prefix{text#desc}suffix"})
		require.NoError(t, err)

		expected := [][]Token{
			{TextToken{Value: "echo"}},
			{
				TextToken{Value: "prefix"},
				FieldToken{Name: "text", Description: "desc", Required: true, OriginalFlag: ""},
				TextToken{Value: "suffix"},
			},
		}
		assert.Equal(t, expected, bp.ShellWords)
		assert.Equal(t, "echo prefix{{text}}suffix", bp.GetCommandFormat())
	})
}

func TestBlueprint_SingleBraceEnhancedTemplateProcessing(t *testing.T) {
	t.Run("handles malformed single-brace template syntax", func(t *testing.T) {
		bp, err := FromArgs([]string{"echo", "{incomplete"})
		require.NoError(t, err)
		args, err := bp.BuildCommandArgs(map[string]interface{}{})
		assert.NoError(t, err)
		assert.Equal(t, []string{"echo", "{incomplete"}, args)
	})

	t.Run("handles malformed single-brace with only closing brace", func(t *testing.T) {
		bp, err := FromArgs([]string{"echo", "no_opening_brace}"})
		require.NoError(t, err)
		args, err := bp.BuildCommandArgs(map[string]interface{}{})
		assert.NoError(t, err)
		assert.Equal(t, []string{"echo", "no_opening_brace}"}, args)
	})

	t.Run("handles empty single-brace template braces", func(t *testing.T) {
		bp, err := FromArgs([]string{"echo", "{}"})
		require.NoError(t, err)
		args, err := bp.BuildCommandArgs(map[string]interface{}{})
		assert.NoError(t, err)
		assert.Equal(t, []string{"echo", "{}"}, args)
	})
}
