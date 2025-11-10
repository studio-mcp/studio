import { Template, Token, TextToken, FieldToken } from "./types.js";
import { generateInputSchema } from "./schema.js";

/**
 * Normalizes field names by converting dashes to underscores
 */
function normalizeFieldName(name: string): string {
  return name.replace(/-/g, "_");
}

/**
 * Finds a parameter value by name, handling dash-underscore equivalence
 */
function findParamValue(
  params: Record<string, any>,
  fieldName: string,
): [any, boolean] {
  // Try exact match first
  if (fieldName in params) {
    return [params[fieldName], true];
  }

  // Try normalized version (dashes to underscores)
  const normalized = normalizeFieldName(fieldName);
  if (normalized in params) {
    return [params[normalized], true];
  }

  // Try reverse (underscores to dashes) if original had underscores
  if (fieldName.includes("_")) {
    const dashed = fieldName.replace(/_/g, "-");
    if (dashed in params) {
      return [params[dashed], true];
    }
  }

  return [undefined, false];
}

/**
 * Checks if a value is meaningful (not empty)
 */
function hasValue(value: any): boolean {
  if (typeof value === "string") {
    return value !== "";
  }
  if (typeof value === "boolean") {
    return value;
  }
  if (Array.isArray(value)) {
    return value.length > 0;
  }
  return value != null;
}

/**
 * Converts a value to its string representation
 */
function valueToString(value: any): string {
  if (typeof value === "string") {
    return value;
  }
  if (typeof value === "boolean") {
    return value ? "true" : "false";
  }
  return String(value);
}

/**
 * Renders a single optional field token
 */
function renderSingleOptionalField(
  bp: Template,
  fieldToken: FieldToken,
  params: Record<string, any>,
): [boolean, string[]] {
  const [value, exists] = findParamValue(params, fieldToken.name);
  if (!exists) {
    return [false, []];
  }

  const inputSchema = generateInputSchema(bp);
  const normalizedName = normalizeFieldName(fieldToken.name);

  // Check if this is a boolean flag
  if (inputSchema.properties?.[normalizedName]?.type === "boolean") {
    if (typeof value === "boolean") {
      if (value) {
        // Use the original flag format if available
        if (fieldToken.originalFlag) {
          return [true, [fieldToken.originalFlag]];
        } else {
          return [true, [`-${fieldToken.name}`]];
        }
      }
    }
    return [false, []];
  }

  // Regular optional field
  const strValue = valueToString(value);
  if (strValue !== "") {
    return [true, [strValue]];
  }

  return [false, []];
}

/**
 * Renders an array field
 */
function renderArrayField(
  fieldToken: FieldToken,
  params: Record<string, any>,
): [boolean, string[]] {
  const [value, exists] = findParamValue(params, fieldToken.name);
  if (!exists) {
    return [false, []];
  }

  if (Array.isArray(value)) {
    const result: string[] = value.map((item) =>
      typeof item === "string" ? item : String(item),
    );
    return [result.length > 0, result];
  }

  return [false, []];
}

/**
 * Renders a single shell word from its tokens
 */
function renderShellWord(
  bp: Template,
  tokens: Token[],
  params: Record<string, any>,
): [boolean, string[]] {
  // Check if this word contains only optional fields that are not provided
  let hasRequiredContent = false;
  let allOptionalFieldsEmpty = true;

  // First pass: check if we should include this word at all
  for (const token of tokens) {
    if (token instanceof TextToken) {
      hasRequiredContent = true;
    } else if (token instanceof FieldToken) {
      const [value, exists] = findParamValue(params, token.name);
      if (exists) {
        if (token.required) {
          hasRequiredContent = true;
          allOptionalFieldsEmpty = false;
        } else {
          // Optional field - check if it has a meaningful value
          if (hasValue(value)) {
            allOptionalFieldsEmpty = false;
          }
        }
      } else if (token.required) {
        hasRequiredContent = true;
      }
    }
  }

  // If this word has only optional fields and they're all empty, skip it
  if (!hasRequiredContent && allOptionalFieldsEmpty) {
    return [false, []];
  }

  // Handle special cases for single field tokens
  if (tokens.length === 1) {
    const token = tokens[0];
    if (token instanceof FieldToken) {
      const inputSchema = generateInputSchema(bp);
      const normalizedName = normalizeFieldName(token.name);

      // Check if this is an array field first (arrays take precedence)
      if (inputSchema.properties?.[normalizedName]?.type === "array") {
        return renderArrayField(token, params);
      }

      // Then check if it's an optional field
      if (!token.required) {
        return renderSingleOptionalField(bp, token, params);
      }
    }
  }

  // Render as a single concatenated string
  const parts: string[] = [];
  for (const token of tokens) {
    if (token instanceof TextToken) {
      parts.push(token.value);
    } else if (token instanceof FieldToken) {
      const [value, exists] = findParamValue(params, token.name);
      if (exists) {
        const strValue = valueToString(value);
        if (strValue !== "") {
          parts.push(strValue);
        }
      }
    }
  }

  if (parts.length > 0) {
    return [true, [parts.join("")]];
  }

  return [false, []];
}

/**
 * Builds command arguments from a template and parameters
 */
export function buildCommandArgs(
  bp: Template,
  params: Record<string, any>,
): string[] {
  const inputSchema = generateInputSchema(bp);

  // Validate required parameters
  for (const required of inputSchema.required || []) {
    const [, exists] = findParamValue(params, required);
    if (!exists) {
      throw new Error(`missing required parameter: ${required}`);
    }
  }

  // Validate parameter types
  for (const [name, param] of Object.entries(params)) {
    const normalizedName = normalizeFieldName(name);
    const schema = inputSchema.properties?.[normalizedName];
    if (schema && schema.type === "array") {
      if (!Array.isArray(param)) {
        throw new Error(
          `parameter '${name}' must be an array, got ${typeof param}`,
        );
      }
    }
  }

  const result: string[] = [];

  for (const shellWord of bp.shellWords) {
    const [shouldInclude, wordResult] = renderShellWord(bp, shellWord, params);
    if (shouldInclude) {
      if (wordResult.length === 0) {
        continue;
      }
      result.push(...wordResult);
    }
  }

  return result;
}
