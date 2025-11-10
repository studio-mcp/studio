import type { Schema } from './schema.js';

/**
 * Token represents a part of a shell word after parsing
 */
export interface Token {
  toString(): string;
}

/**
 * TextToken represents literal text in a shell word
 */
export class TextToken implements Token {
  constructor(public value: string) {}

  toString(): string {
    return this.value;
  }
}

/**
 * FieldToken represents a template field in a shell word
 */
export class FieldToken implements Token {
  constructor(
    public name: string,
    public description: string,
    public required: boolean,
    public isArray: boolean = false,
    public originalFlag: string = ''
  ) {}

  toString(): string {
    if (this.required) {
      return `{{${this.name}}}`;
    }
    return `[${this.name}]`;
  }
}

/**
 * Blueprint represents a parsed command template
 */
export class Blueprint {
  constructor(
    public baseCommand: string,
    public shellWords: Token[][]
  ) {}

  /**
   * Returns the base command
   */
  getBaseCommand(): string {
    return this.baseCommand;
  }

  /**
   * Returns the command format for display
   */
  getCommandFormat(): string {
    const parts = this.shellWords.map((tokens) =>
      this.renderTokensForDisplay(tokens)
    );
    return parts.join(' ');
  }

  /**
   * Renders tokens for display purposes
   */
  private renderTokensForDisplay(tokens: Token[]): string {
    if (tokens.length === 1) {
      if (tokens[0] instanceof FieldToken) {
        return this.renderFieldTokenForDisplay(tokens[0]);
      }
      return tokens[0].toString();
    }

    let result = '';
    for (const token of tokens) {
      if (token instanceof TextToken) {
        result += token.value;
      } else if (token instanceof FieldToken) {
        result += this.renderFieldTokenForDisplay(token);
      }
    }
    return result;
  }

  /**
   * Renders a single field token for display
   */
  private renderFieldTokenForDisplay(token: FieldToken): string {
    let name = token.name;

    // For boolean flags, use the original flag format
    if (token.originalFlag !== '') {
      name = token.originalFlag;
    }

    if (token.isArray) {
      name = name + '...';
    }

    // For required fields, use template format
    if (token.required) {
      return `{{${name}}}`;
    }

    // Regular optional field
    return `[${name}]`;
  }

  /**
   * Returns the input schema for this blueprint
   */
  getInputSchema(): Schema {
    return this.generateInputSchema();
  }

  /**
   * Generates the input schema from tokenized shell words
   */
  generateInputSchema(): Schema {
    // Import dynamically to avoid circular dependency
    return (require('./schema.js') as typeof import('./schema.js')).generateInputSchema(this);
  }
}
