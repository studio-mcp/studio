package blueprint

import (
	"testing"

	"github.com/google/jsonschema-go/jsonschema"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBlueprint_GenerateInputSchema(t *testing.T) {
	testCases := []struct {
		name           string
		args           []string
		expectedSchema *jsonschema.Schema
	}{
		{
			name: "simple command without args",
			args: []string{"git", "status"},
			expectedSchema: &jsonschema.Schema{
				Type:       "object",
				Properties: map[string]*jsonschema.Schema{},
				Required:   []string{},
			},
		},
		{
			name: "simple command with array args",
			args: []string{"git", "status", "[args...]"},
			expectedSchema: &jsonschema.Schema{
				Type: "object",
				Properties: map[string]*jsonschema.Schema{
					"args": {
						Type:        "array",
						Items:       &jsonschema.Schema{Type: "string"},
						Description: "Additional command line arguments",
					},
				},
				Required: []string{},
			},
		},
		{
			name: "template command with description",
			args: []string{"curl", "https://en.m.wikipedia.org/wiki/{{page#A valid wikipedia page}}"},
			expectedSchema: &jsonschema.Schema{
				Type: "object",
				Properties: map[string]*jsonschema.Schema{
					"page": {
						Type:        "string",
						Description: "A valid wikipedia page",
					},
				},
				Required: []string{"page"},
			},
		},
		{
			name: "template command without description",
			args: []string{"echo", "{{text}}"},
			expectedSchema: &jsonschema.Schema{
				Type: "object",
				Properties: map[string]*jsonschema.Schema{
					"text": {
						Type: "string",
					},
				},
				Required: []string{"text"},
			},
		},
		{
			name: "template with spaces in description",
			args: []string{"curl", "https://en.m.wikipedia.org/wiki/{{page # A valid wikipedia page}}"},
			expectedSchema: &jsonschema.Schema{
				Type: "object",
				Properties: map[string]*jsonschema.Schema{
					"page": {
						Type:        "string",
						Description: "A valid wikipedia page",
					},
				},
				Required: []string{"page"},
			},
		},
		{
			name: "mixed required and optional arguments",
			args: []string{"command", "{{arg1#Custom description}}", "[arg2]"},
			expectedSchema: &jsonschema.Schema{
				Type: "object",
				Properties: map[string]*jsonschema.Schema{
					"arg1": {
						Type:        "string",
						Description: "Custom description",
					},
					"arg2": {
						Type: "string",
					},
				},
				Required: []string{"arg1"},
			},
		},
		{
			name: "duplicate template with explicit description priority",
			args: []string{"echo", "{{text#Explicit description}}", "{{text}}"},
			expectedSchema: &jsonschema.Schema{
				Type: "object",
				Properties: map[string]*jsonschema.Schema{
					"text": {
						Type:        "string",
						Description: "Explicit description",
					},
				},
				Required: []string{"text"},
			},
		},
		{
			name: "array arguments with default description",
			args: []string{"echo", "[files...]"},
			expectedSchema: &jsonschema.Schema{
				Type: "object",
				Properties: map[string]*jsonschema.Schema{
					"files": {
						Type:        "array",
						Items:       &jsonschema.Schema{Type: "string"},
						Description: "Additional command line arguments",
					},
				},
				Required: []string{},
			},
		},
		{
			name: "array arguments with custom description",
			args: []string{"ls", "[files...#Files to list]"},
			expectedSchema: &jsonschema.Schema{
				Type: "object",
				Properties: map[string]*jsonschema.Schema{
					"files": {
						Type:        "array",
						Items:       &jsonschema.Schema{Type: "string"},
						Description: "Files to list",
					},
				},
				Required: []string{},
			},
		},
		{
			name: "optional string field",
			args: []string{"echo", "[optional]"},
			expectedSchema: &jsonschema.Schema{
				Type: "object",
				Properties: map[string]*jsonschema.Schema{
					"optional": {
						Type: "string",
					},
				},
				Required: []string{},
			},
		},
		{
			name: "optional field with description",
			args: []string{"echo", "[name#Person's name]"},
			expectedSchema: &jsonschema.Schema{
				Type: "object",
				Properties: map[string]*jsonschema.Schema{
					"name": {
						Type:        "string",
						Description: "Person's name",
					},
				},
				Required: []string{},
			},
		},
		{
			name: "dashes converted to underscores",
			args: []string{"echo", "[has-dashes]"},
			expectedSchema: &jsonschema.Schema{
				Type: "object",
				Properties: map[string]*jsonschema.Schema{
					"has_dashes": {
						Type: "string",
					},
				},
				Required: []string{},
			},
		},
		{
			name: "boolean flag short form",
			args: []string{"ls", "[-f]"},
			expectedSchema: &jsonschema.Schema{
				Type: "object",
				Properties: map[string]*jsonschema.Schema{
					"f": {
						Type:        "boolean",
						Description: "Enable -f flag",
					},
				},
				Required: []string{},
			},
		},
		{
			name: "boolean flag long form",
			args: []string{"ls", "[--force]"},
			expectedSchema: &jsonschema.Schema{
				Type: "object",
				Properties: map[string]*jsonschema.Schema{
					"force": {
						Type:        "boolean",
						Description: "Enable --force flag",
					},
				},
				Required: []string{},
			},
		},
		{
			name: "boolean flag with custom description",
			args: []string{"rm", "[-f#force removal]"},
			expectedSchema: &jsonschema.Schema{
				Type: "object",
				Properties: map[string]*jsonschema.Schema{
					"f": {
						Type:        "boolean",
						Description: "force removal",
					},
				},
				Required: []string{},
			},
		},
		{
			name: "mixed boolean flags and required arguments",
			args: []string{"cp", "[-r#recursive]", "{{source}}", "{{dest}}"},
			expectedSchema: &jsonschema.Schema{
				Type: "object",
				Properties: map[string]*jsonschema.Schema{
					"r": {
						Type:        "boolean",
						Description: "recursive",
					},
					"source": {
						Type: "string",
					},
					"dest": {
						Type: "string",
					},
				},
				Required: []string{"source", "dest"},
			},
		},
		{
			name: "mixed string and array arguments",
			args: []string{"command", "{{flag#Command flag}}", "[files...]"},
			expectedSchema: &jsonschema.Schema{
				Type: "object",
				Properties: map[string]*jsonschema.Schema{
					"flag": {
						Type:        "string",
						Description: "Command flag",
					},
					"files": {
						Type:        "array",
						Items:       &jsonschema.Schema{Type: "string"},
						Description: "Additional command line arguments",
					},
				},
				Required: []string{"flag"},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			bp, err := FromArgs(tc.args)
			require.NoError(t, err)

			actualSchema := bp.GenerateInputSchema()
			assert.Equal(t, tc.expectedSchema, actualSchema)
		})
	}
}

func TestBlueprint_GenerateInputSchema_EdgeCases(t *testing.T) {
	t.Run("mixed optional arguments with and without descriptions", func(t *testing.T) {
		bp, err := FromArgs([]string{"cmd", "[arg1]", "[arg2#Custom desc]"})
		require.NoError(t, err)

		schema := bp.GenerateInputSchema()

		// Check properties individually since the order might vary
		assert.Equal(t, "string", schema.Properties["arg1"].Type)
		assert.Equal(t, "", schema.Properties["arg1"].Description)

		assert.Equal(t, "string", schema.Properties["arg2"].Type)
		assert.Equal(t, "Custom desc", schema.Properties["arg2"].Description)

		// Neither should be required
		assert.Empty(t, schema.Required)
	})
}
