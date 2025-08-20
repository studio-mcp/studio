package blueprint

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBlueprint_BuildCommandArgs(t *testing.T) {
	t.Run("builds simple command without templates", func(t *testing.T) {
		bp, err := FromArgs([]string{"echo", "hello", "world"})
		require.NoError(t, err)
		args, err := bp.BuildCommandArgs(map[string]interface{}{})

		assert.NoError(t, err)
		assert.Equal(t, []string{"echo", "hello", "world"}, args)
	})

	t.Run("builds command with required template", func(t *testing.T) {
		bp, err := FromArgs([]string{"echo", "{{message}}"})
		require.NoError(t, err)
		args, err := bp.BuildCommandArgs(map[string]interface{}{
			"message": "Hello World",
		})

		assert.NoError(t, err)
		assert.Equal(t, []string{"echo", "Hello World"}, args)
	})

	t.Run("builds command with template in middle of arg", func(t *testing.T) {
		bp, err := FromArgs([]string{"curl", "https://api.example.com/{{endpoint}}"})
		require.NoError(t, err)
		args, err := bp.BuildCommandArgs(map[string]interface{}{
			"endpoint": "users/123",
		})

		assert.NoError(t, err)
		assert.Equal(t, []string{"curl", "https://api.example.com/users/123"}, args)
	})

	t.Run("builds command with multiple templates in one arg", func(t *testing.T) {
		bp, err := FromArgs([]string{"echo", "{{greeting}} {{name}}!"})
		require.NoError(t, err)
		args, err := bp.BuildCommandArgs(map[string]interface{}{
			"greeting": "Hello",
			"name":     "World",
		})

		assert.NoError(t, err)
		assert.Equal(t, []string{"echo", "Hello World!"}, args)
	})

	t.Run("builds command with array argument", func(t *testing.T) {
		bp, err := FromArgs([]string{"echo", "[files...]"})
		require.NoError(t, err)
		args, err := bp.BuildCommandArgs(map[string]interface{}{
			"files": []string{"file1.txt", "file2.txt", "file3.txt"},
		})

		assert.NoError(t, err)
		assert.Equal(t, []string{"echo", "file1.txt", "file2.txt", "file3.txt"}, args)
	})

	t.Run("builds command with empty array argument", func(t *testing.T) {
		bp, err := FromArgs([]string{"echo", "prefix", "[files...]"})
		require.NoError(t, err)
		args, err := bp.BuildCommandArgs(map[string]interface{}{
			"files": []string{},
		})

		assert.NoError(t, err)
		assert.Equal(t, []string{"echo", "prefix"}, args)
	})

	t.Run("builds command with optional string argument provided", func(t *testing.T) {
		bp, err := FromArgs([]string{"echo", "hello", "[name]"})
		require.NoError(t, err)
		args, err := bp.BuildCommandArgs(map[string]interface{}{
			"name": "world",
		})

		assert.NoError(t, err)
		assert.Equal(t, []string{"echo", "hello", "world"}, args)
	})

	t.Run("builds command with optional string argument omitted", func(t *testing.T) {
		bp, err := FromArgs([]string{"echo", "hello", "[name]"})
		require.NoError(t, err)
		args, err := bp.BuildCommandArgs(map[string]interface{}{})

		assert.NoError(t, err)
		assert.Equal(t, []string{"echo", "hello"}, args)
	})

	t.Run("builds command with blueprint arguments containing spaces", func(t *testing.T) {
		bp, err := FromArgs([]string{"echo", "{{text#text to echo}}"})
		require.NoError(t, err)
		args, err := bp.BuildCommandArgs(map[string]interface{}{
			"text": "Hello World",
		})

		assert.NoError(t, err)
		assert.Equal(t, []string{"echo", "Hello World"}, args)
	})

	t.Run("builds command with mixed blueprint with and without descriptions", func(t *testing.T) {
		bp, err := FromArgs([]string{"echo", "{{greeting#The greeting}}", "{{name}}"})
		require.NoError(t, err)
		args, err := bp.BuildCommandArgs(map[string]interface{}{
			"greeting": "Hello",
			"name":     "World",
		})

		assert.NoError(t, err)
		assert.Equal(t, []string{"echo", "Hello", "World"}, args)
	})

	t.Run("builds command with blueprint arguments in mixed content", func(t *testing.T) {
		bp, err := FromArgs([]string{"echo", "simon says {{text#text for simon to say}}"})
		require.NoError(t, err)
		args, err := bp.BuildCommandArgs(map[string]interface{}{
			"text": "Hello World",
		})

		assert.NoError(t, err)
		assert.Equal(t, []string{"echo", "simon says Hello World"}, args)
	})

	t.Run("builds command with blueprint arguments containing special shell characters", func(t *testing.T) {
		bp, err := FromArgs([]string{"echo", "{{text#text to echo}}"})
		require.NoError(t, err)
		args, err := bp.BuildCommandArgs(map[string]interface{}{
			"text": "Hello & World; echo pwned",
		})

		assert.NoError(t, err)
		assert.Equal(t, []string{"echo", "Hello & World; echo pwned"}, args)
	})

	t.Run("builds command with blueprint in middle of argument", func(t *testing.T) {
		bp, err := FromArgs([]string{"echo", "--message={{text#message content}}"})
		require.NoError(t, err)
		args, err := bp.BuildCommandArgs(map[string]interface{}{
			"text": "Hello World",
		})

		assert.NoError(t, err)
		assert.Equal(t, []string{"echo", "--message=Hello World"}, args)
	})

	t.Run("builds command with blueprint with prefix and suffix", func(t *testing.T) {
		bp, err := FromArgs([]string{"echo", "prefix-{{text#middle part}}-suffix"})
		require.NoError(t, err)
		args, err := bp.BuildCommandArgs(map[string]interface{}{
			"text": "Hello World",
		})

		assert.NoError(t, err)
		assert.Equal(t, []string{"echo", "prefix-Hello World-suffix"}, args)
	})

	t.Run("builds command with mixed blueprint and non-blueprint arguments", func(t *testing.T) {
		bp, err := FromArgs([]string{"echo", "static", "{{dynamic#dynamic content}}", "more-static"})
		require.NoError(t, err)
		args, err := bp.BuildCommandArgs(map[string]interface{}{
			"dynamic": "Hello World",
		})

		assert.NoError(t, err)
		assert.Equal(t, []string{"echo", "static", "Hello World", "more-static"}, args)
	})

	t.Run("preserves shell safety with complex blueprint values", func(t *testing.T) {
		bp, err := FromArgs([]string{"echo", "Result: {{text#text content}}"})
		require.NoError(t, err)
		args, err := bp.BuildCommandArgs(map[string]interface{}{
			"text": "$(echo 'dangerous'); echo 'safe'",
		})

		assert.NoError(t, err)
		assert.Equal(t, []string{"echo", "Result: $(echo 'dangerous'); echo 'safe'"}, args)
	})

	t.Run("builds command with mixed string and array arguments", func(t *testing.T) {
		bp, err := FromArgs([]string{"echo", "{{prefix#Prefix text}}", "[files...]"})
		require.NoError(t, err)
		args, err := bp.BuildCommandArgs(map[string]interface{}{
			"prefix": "Files:",
			"files":  []string{"a.txt", "b.txt"},
		})

		assert.NoError(t, err)
		assert.Equal(t, []string{"echo", "Files:", "a.txt", "b.txt"}, args)
	})

	t.Run("builds command with mixed required and optional arguments", func(t *testing.T) {
		bp, err := FromArgs([]string{"echo", "{{required#Required text}}", "[optional]"})
		require.NoError(t, err)

		// With both provided
		args, err := bp.BuildCommandArgs(map[string]interface{}{
			"required": "hello",
			"optional": "world",
		})
		assert.NoError(t, err)
		assert.Equal(t, []string{"echo", "hello", "world"}, args)

		// With only required provided
		args, err = bp.BuildCommandArgs(map[string]interface{}{
			"required": "hello",
		})
		assert.NoError(t, err)
		assert.Equal(t, []string{"echo", "hello"}, args)
	})
}

func TestBlueprint_TemplateValidation(t *testing.T) {
	t.Run("validates missing required parameters", func(t *testing.T) {
		bp, err := FromArgs([]string{"echo", "{{required}}"})
		require.NoError(t, err)
		_, err = bp.BuildCommandArgs(map[string]interface{}{})
		assert.Error(t, err) // Should error on missing required param
	})

	t.Run("validates parameter type mismatches", func(t *testing.T) {
		bp, err := FromArgs([]string{"echo", "[files...]"})
		require.NoError(t, err)
		_, err = bp.BuildCommandArgs(map[string]interface{}{
			"files": "not-an-array", // Should be []string
		})
		assert.Error(t, err)
	})
}

func TestBlueprint_EnhancedTemplateProcessing(t *testing.T) {
	t.Run("handles malformed template syntax", func(t *testing.T) {
		bp, err := FromArgs([]string{"echo", "{{incomplete"})
		require.NoError(t, err)
		// Should ignore any bad template syntax and print as normal text (always fallback on parse error to literal text)
		args, err := bp.BuildCommandArgs(map[string]interface{}{})
		assert.NoError(t, err)
		assert.Equal(t, []string{"echo", "{{incomplete"}, args)
	})

	t.Run("handles malformed template with only opening braces", func(t *testing.T) {
		bp, err := FromArgs([]string{"echo", "{{no_closing_braces"})
		require.NoError(t, err)
		args, err := bp.BuildCommandArgs(map[string]interface{}{})
		assert.NoError(t, err)
		assert.Equal(t, []string{"echo", "{{no_closing_braces"}, args)
	})

	t.Run("handles malformed template with only closing braces", func(t *testing.T) {
		bp, err := FromArgs([]string{"echo", "no_opening_braces}}"})
		require.NoError(t, err)
		args, err := bp.BuildCommandArgs(map[string]interface{}{})
		assert.NoError(t, err)
		assert.Equal(t, []string{"echo", "no_opening_braces}}"}, args)
	})

	t.Run("handles mixed valid and malformed templates", func(t *testing.T) {
		bp, err := FromArgs([]string{"echo", "{{valid}}", "{{incomplete", "}}"})
		require.NoError(t, err)
		args, err := bp.BuildCommandArgs(map[string]interface{}{
			"valid": "works",
		})
		assert.NoError(t, err)
		assert.Equal(t, []string{"echo", "works", "{{incomplete", "}}"}, args)
	})

	t.Run("handles empty template braces", func(t *testing.T) {
		bp, err := FromArgs([]string{"echo", "{{}}"})
		require.NoError(t, err)
		args, err := bp.BuildCommandArgs(map[string]interface{}{})
		assert.NoError(t, err)
		assert.Equal(t, []string{"echo", "{{}}"}, args)
	})
}

func TestBlueprint_BuildCommandArgsWithBooleanFlags(t *testing.T) {
	t.Run("builds command with boolean flag enabled", func(t *testing.T) {
		bp, err := FromArgs([]string{"ls", "[-f]"})
		require.NoError(t, err)
		args, err := bp.BuildCommandArgs(map[string]interface{}{
			"f": true,
		})

		assert.NoError(t, err)
		assert.Equal(t, []string{"ls", "-f"}, args)
	})

	t.Run("builds command with boolean flag disabled", func(t *testing.T) {
		bp, err := FromArgs([]string{"ls", "[-f]"})
		require.NoError(t, err)
		args, err := bp.BuildCommandArgs(map[string]interface{}{
			"f": false,
		})

		assert.NoError(t, err)
		assert.Equal(t, []string{"ls"}, args)
	})

	t.Run("builds command with boolean flag omitted", func(t *testing.T) {
		bp, err := FromArgs([]string{"ls", "[-f]"})
		require.NoError(t, err)
		args, err := bp.BuildCommandArgs(map[string]interface{}{})

		assert.NoError(t, err)
		assert.Equal(t, []string{"ls"}, args)
	})

	t.Run("builds command with long boolean flag enabled", func(t *testing.T) {
		bp, err := FromArgs([]string{"ls", "[--force]"})
		require.NoError(t, err)
		args, err := bp.BuildCommandArgs(map[string]interface{}{
			"force": true,
		})

		assert.NoError(t, err)
		assert.Equal(t, []string{"ls", "--force"}, args)
	})

	t.Run("builds command with mixed boolean and string arguments", func(t *testing.T) {
		bp, err := FromArgs([]string{"cp", "[-r]", "{{source}}", "{{dest}}"})
		require.NoError(t, err)
		args, err := bp.BuildCommandArgs(map[string]interface{}{
			"r":      true,
			"source": "file1.txt",
			"dest":   "file2.txt",
		})

		assert.NoError(t, err)
		assert.Equal(t, []string{"cp", "-r", "file1.txt", "file2.txt"}, args)
	})

	t.Run("builds command with mixed boolean flags some enabled some disabled", func(t *testing.T) {
		bp, err := FromArgs([]string{"ls", "[-l]", "[-a]", "[--human-readable]"})
		require.NoError(t, err)
		args, err := bp.BuildCommandArgs(map[string]interface{}{
			"l":              true,
			"a":              false,
			"human_readable": true,
		})

		assert.NoError(t, err)
		assert.Equal(t, []string{"ls", "-l", "--human-readable"}, args)
	})
}

func TestBlueprint_BuildCommandArgsTokenized(t *testing.T) {
	t.Run("builds simple command without templates", func(t *testing.T) {
		bp, err := FromArgs([]string{"echo", "hello", "world"})
		require.NoError(t, err)

		args, err := bp.BuildCommandArgs(map[string]interface{}{})
		assert.NoError(t, err)
		assert.Equal(t, []string{"echo", "hello", "world"}, args)
	})

	t.Run("builds command with required template", func(t *testing.T) {
		bp, err := FromArgs([]string{"echo", "{{message}}"})
		require.NoError(t, err)

		args, err := bp.BuildCommandArgs(map[string]interface{}{
			"message": "Hello World",
		})
		assert.NoError(t, err)
		assert.Equal(t, []string{"echo", "Hello World"}, args)
	})

	t.Run("builds command with template in middle of arg", func(t *testing.T) {
		bp, err := FromArgs([]string{"curl", "https://api.example.com/{{endpoint}}"})
		require.NoError(t, err)

		args, err := bp.BuildCommandArgs(map[string]interface{}{
			"endpoint": "users/123",
		})
		assert.NoError(t, err)
		assert.Equal(t, []string{"curl", "https://api.example.com/users/123"}, args)
	})

	t.Run("builds command with multiple templates in one arg", func(t *testing.T) {
		bp, err := FromArgs([]string{"echo", "{{greeting}} {{name}}!"})
		require.NoError(t, err)

		args, err := bp.BuildCommandArgs(map[string]interface{}{
			"greeting": "Hello",
			"name":     "World",
		})
		assert.NoError(t, err)
		assert.Equal(t, []string{"echo", "Hello World!"}, args)
	})

	t.Run("builds command with optional field provided", func(t *testing.T) {
		bp, err := FromArgs([]string{"echo", "hello", "[name]"})
		require.NoError(t, err)

		args, err := bp.BuildCommandArgs(map[string]interface{}{
			"name": "world",
		})
		assert.NoError(t, err)
		assert.Equal(t, []string{"echo", "hello", "world"}, args)
	})

	t.Run("builds command with optional field omitted", func(t *testing.T) {
		bp, err := FromArgs([]string{"echo", "hello", "[name]"})
		require.NoError(t, err)

		args, err := bp.BuildCommandArgs(map[string]interface{}{})
		assert.NoError(t, err)
		assert.Equal(t, []string{"echo", "hello"}, args)
	})

	t.Run("builds command with boolean flag enabled", func(t *testing.T) {
		bp, err := FromArgs([]string{"ls", "[--verbose]"})
		require.NoError(t, err)

		args, err := bp.BuildCommandArgs(map[string]interface{}{
			"verbose": true,
		})
		assert.NoError(t, err)
		assert.Equal(t, []string{"ls", "--verbose"}, args)
	})

	t.Run("builds command with boolean flag disabled", func(t *testing.T) {
		bp, err := FromArgs([]string{"ls", "[--verbose]"})
		require.NoError(t, err)

		args, err := bp.BuildCommandArgs(map[string]interface{}{
			"verbose": false,
		})
		assert.NoError(t, err)
		assert.Equal(t, []string{"ls"}, args)
	})

	t.Run("builds command with array argument", func(t *testing.T) {
		bp, err := FromArgs([]string{"echo", "[files...]"})
		require.NoError(t, err)

		args, err := bp.BuildCommandArgs(map[string]interface{}{
			"files": []string{"file1.txt", "file2.txt", "file3.txt"},
		})
		assert.NoError(t, err)
		assert.Equal(t, []string{"echo", "file1.txt", "file2.txt", "file3.txt"}, args)
	})

	t.Run("handles dash-underscore equivalence", func(t *testing.T) {
		bp, err := FromArgs([]string{"echo", "{{my-var}}", "[--my-flag]"})
		require.NoError(t, err)

		// Should accept both dash and underscore versions
		args, err := bp.BuildCommandArgs(map[string]interface{}{
			"my_var":  "hello",
			"my_flag": true,
		})
		assert.NoError(t, err)
		assert.Equal(t, []string{"echo", "hello", "--my-flag"}, args)
	})

	t.Run("returns error for missing required parameter", func(t *testing.T) {
		bp, err := FromArgs([]string{"echo", "{{message}}"})
		require.NoError(t, err)

		_, err = bp.BuildCommandArgs(map[string]interface{}{})
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "missing required parameter")
	})
}

func TestBlueprint_BuildCommandArgs_SingleBrace(t *testing.T) {
	t.Run("builds command with single-brace required template", func(t *testing.T) {
		bp, err := FromArgs([]string{"echo", "{message}"})
		require.NoError(t, err)
		args, err := bp.BuildCommandArgs(map[string]interface{}{
			"message": "Hello World",
		})
		assert.NoError(t, err)
		assert.Equal(t, []string{"echo", "Hello World"}, args)
	})

	t.Run("builds command with single-brace in middle of arg", func(t *testing.T) {
		bp, err := FromArgs([]string{"curl", "https://api.example.com/{endpoint}"})
		require.NoError(t, err)
		args, err := bp.BuildCommandArgs(map[string]interface{}{
			"endpoint": "users/123",
		})
		assert.NoError(t, err)
		assert.Equal(t, []string{"curl", "https://api.example.com/users/123"}, args)
	})

	t.Run("builds command with multiple single-brace templates in one arg", func(t *testing.T) {
		bp, err := FromArgs([]string{"echo", "{greeting} {name}!"})
		require.NoError(t, err)
		args, err := bp.BuildCommandArgs(map[string]interface{}{
			"greeting": "Hello",
			"name":     "World",
		})
		assert.NoError(t, err)
		assert.Equal(t, []string{"echo", "Hello World!"}, args)
	})

	t.Run("builds command with single-brace required array", func(t *testing.T) {
		bp, err := FromArgs([]string{"echo", "{files...}"})
		require.NoError(t, err)
		args, err := bp.BuildCommandArgs(map[string]interface{}{
			"files": []string{"a.txt", "b.txt"},
		})
		assert.NoError(t, err)
		assert.Equal(t, []string{"echo", "a.txt", "b.txt"}, args)
	})

	t.Run("errors when single-brace required arg missing", func(t *testing.T) {
		bp, err := FromArgs([]string{"echo", "{text}"})
		require.NoError(t, err)
		_, err = bp.BuildCommandArgs(map[string]interface{}{})
		assert.Error(t, err)
	})
}
