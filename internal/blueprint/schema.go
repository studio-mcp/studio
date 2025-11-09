package blueprint

import (
	"fmt"
	"strings"

	"github.com/google/jsonschema-go/jsonschema"
)

// GenerateInputSchema creates a JSON schema from the tokenized shell words
func (bp *Blueprint) GenerateInputSchema() *jsonschema.Schema {
	properties := make(map[string]*jsonschema.Schema)
	required := []string{}

	// Iterate through all shell words and their tokens
	for _, tokens := range bp.ShellWords {
		for _, token := range tokens {
			if fieldToken, ok := token.(FieldToken); ok {
				// Use normalized name for schema properties (dashes to underscores)
				normalizedName := strings.ReplaceAll(fieldToken.Name, "-", "_")

				// Skip if we already have this property and it has a description
				if existingProp, exists := properties[normalizedName]; exists {
					if fieldToken.Description != "" && existingProp.Description == "" {
						// Update existing property with description
						existingProp.Description = fieldToken.Description
					}
					// Handle required status - if any instance is required, make it required
					if fieldToken.Required && !contains(required, normalizedName) {
						required = append(required, normalizedName)
					}
					continue
				}

				// Create new property based on token type
				var prop *jsonschema.Schema

				if fieldToken.OriginalFlag != "" {
					// Boolean flag
					description := fieldToken.Description
					if description == "" {
						description = fmt.Sprintf("Enable %s flag", fieldToken.OriginalFlag)
					}
					prop = &jsonschema.Schema{
						Type:        "boolean",
						Description: description,
					}
				} else if fieldToken.IsArray {
					// Array field
					description := fieldToken.Description
					if description == "" {
						description = "Additional command line arguments"
					}
					prop = &jsonschema.Schema{
						Type:        "array",
						Items:       &jsonschema.Schema{Type: "string"},
						Description: description,
					}
					// Array fields follow the same required logic as other fields
					if fieldToken.Required && !contains(required, normalizedName) {
						required = append(required, normalizedName)
					}
				} else {
					// String field
					prop = &jsonschema.Schema{Type: "string"}
					if fieldToken.Description != "" {
						prop.Description = fieldToken.Description
					}

					// Add to required if the token is marked as required
					if fieldToken.Required && !contains(required, normalizedName) {
						required = append(required, normalizedName)
					}
				}

				properties[normalizedName] = prop
			}
		}
	}

	schema := &jsonschema.Schema{
		Type:       "object",
		Properties: properties,
		Required:   required, // Always set, even if empty
	}

	// Debug logging
	debug("GenerateInputSchema created schema with %d properties, %d required", len(properties), len(required))
	for name, prop := range properties {
		debug("  property %s: type=%s, description=%s", name, prop.Type, prop.Description)
	}
	for _, req := range required {
		debug("  required: %s", req)
	}

	return schema
}
