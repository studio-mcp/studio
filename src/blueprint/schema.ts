import { Blueprint, FieldToken } from './types.js';

/**
 * JSON Schema interface
 */
export interface Schema {
  type: string;
  properties?: Record<string, SchemaProperty>;
  required?: string[];
  items?: SchemaProperty;
  description?: string;
}

export interface SchemaProperty {
  type: string;
  description?: string;
  items?: SchemaProperty;
}

/**
 * Generates a JSON schema from the tokenized shell words
 */
export function generateInputSchema(bp: Blueprint): Schema {
  const properties: Record<string, SchemaProperty> = {};
  const required: string[] = [];

  // Iterate through all shell words and their tokens
  for (const tokens of bp.shellWords) {
    for (const token of tokens) {
      if (token instanceof FieldToken) {
        // Use normalized name for schema properties (dashes to underscores)
        const normalizedName = token.name.replace(/-/g, '_');

        // Skip if we already have this property and it has a description
        const existingProp = properties[normalizedName];
        if (existingProp) {
          if (token.description && !existingProp.description) {
            // Update existing property with description
            existingProp.description = token.description;
          }
          // Handle required status - if any instance is required, make it required
          if (token.required && !required.includes(normalizedName)) {
            required.push(normalizedName);
          }
          continue;
        }

        // Create new property based on token type
        let prop: SchemaProperty;

        if (token.originalFlag) {
          // Boolean flag
          let description = token.description;
          if (!description) {
            description = `Enable ${token.originalFlag} flag`;
          }
          prop = {
            type: 'boolean',
            description,
          };
        } else if (token.isArray) {
          // Array field
          let description = token.description;
          if (!description) {
            description = 'Additional command line arguments';
          }
          prop = {
            type: 'array',
            items: { type: 'string' },
            description,
          };
          // Array fields follow the same required logic as other fields
          if (token.required && !required.includes(normalizedName)) {
            required.push(normalizedName);
          }
        } else {
          // String field
          prop = { type: 'string' };
          if (token.description) {
            prop.description = token.description;
          }

          // Add to required if the token is marked as required
          if (token.required && !required.includes(normalizedName)) {
            required.push(normalizedName);
          }
        }

        properties[normalizedName] = prop;
      }
    }
  }

  return {
    type: 'object',
    properties,
    required,
  };
}
