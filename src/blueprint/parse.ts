import { Blueprint, Token, TextToken, FieldToken } from './types.js';

/**
 * Template match represents a found template in the text
 */
interface TemplateMatch {
  start: number;
  end: number;
  type: string;
  endMarker: string;
  markerLength: number;
}

/**
 * Creates a new Blueprint from command arguments using tokenization
 */
export function fromArgs(args: string[]): Blueprint {
  if (args.length === 0) {
    throw new Error('cannot create blueprint: no command provided');
  }

  if (args[0].trim() === '') {
    throw new Error('cannot create blueprint: empty command provided');
  }

  const bp = new Blueprint(args[0], new Array(args.length));

  // Tokenize each shell word
  for (let i = 0; i < args.length; i++) {
    bp.shellWords[i] = tokenizeShellWord(args[i]);
  }

  return bp;
}

/**
 * Tokenizes a single shell word into tokens
 */
function tokenizeShellWord(word: string): Token[] {
  const tokens: Token[] = [];
  let pos = 0;

  while (pos < word.length) {
    const templateStart = findNextTemplate(word, pos);
    if (!templateStart) {
      // No more templates, add remaining text
      if (pos < word.length) {
        tokens.push(new TextToken(word.slice(pos)));
      }
      break;
    }

    // Add text before template
    if (templateStart.start > pos) {
      tokens.push(new TextToken(word.slice(pos, templateStart.start)));
    }

    // Parse template
    const templateText = word.slice(templateStart.start, templateStart.end);
    const token = parseField(templateText);
    if (token) {
      tokens.push(token);
    } else {
      tokens.push(new TextToken(templateText));
    }

    pos = templateStart.end;
  }

  // Ensure we always return at least one token
  if (tokens.length === 0) {
    tokens.push(new TextToken(word));
  }

  return tokens;
}

/**
 * Finds the next template starting from the given position
 */
function findNextTemplate(word: string, startPos: number): TemplateMatch | null {
  const remaining = word.slice(startPos);

  // Find possible starts
  let startDouble = remaining.indexOf('{{');
  let startSingle = remaining.indexOf('{');
  const startOptional = remaining.indexOf('[');

  // If single coincides with double, prefer double
  if (startSingle !== -1 && startDouble !== -1 && startSingle === startDouble) {
    startSingle = -1;
  }

  // Choose earliest non -1 among double, single, optional
  let nextStart = -1;
  let templateType = '';
  let endMarker = '';
  let markerLength = 0;

  const candidates = [
    { idx: startDouble, typ: 'required_double', end: '}}', len: 2 },
    { idx: startSingle, typ: 'required_single', end: '}', len: 1 },
    { idx: startOptional, typ: 'optional', end: ']', len: 1 },
  ];

  for (const c of candidates) {
    if (c.idx !== -1 && (nextStart === -1 || c.idx < nextStart)) {
      nextStart = c.idx;
      templateType = c.typ;
      endMarker = c.end;
      markerLength = c.len;
    }
  }

  if (nextStart === -1) {
    return null; // No templates found
  }

  const absoluteStart = startPos + nextStart;

  // Find template end
  const endIndex = remaining.slice(nextStart).indexOf(endMarker);
  if (endIndex === -1) {
    // Malformed template - treat rest as text by returning no match
    return null;
  }

  const absoluteEnd = absoluteStart + endIndex + markerLength;

  return {
    start: absoluteStart,
    end: absoluteEnd,
    type: templateType,
    endMarker,
    markerLength,
  };
}

/**
 * Parses a field enclosed in {{ }}, { } or [ ]
 */
function parseField(field: string): Token | null {
  let content: string;
  let required: boolean;

  // Determine field type and extract content
  if (field.startsWith('{{') && field.endsWith('}}')) {
    content = field.slice(2, -2); // Remove {{ }}
    required = true;
  } else if (field.startsWith('{') && field.endsWith('}')) {
    content = field.slice(1, -1); // Remove { }
    required = true;
  } else if (field.startsWith('[') && field.endsWith(']')) {
    content = field.slice(1, -1); // Remove [ ]
    required = false;
  } else {
    return null; // Not a valid field
  }

  // Parse content for name, description, and modifiers
  let name: string;
  let description = '';
  let isArray = false;
  let originalFlag = '';

  // Check for description (split on #)
  const parts = content.split('#');
  name = parts[0].trim();

  // If name is empty, this is not a valid field (e.g., {{}} or {})
  if (name === '') {
    return null;
  }

  if (parts.length > 1) {
    description = parts.slice(1).join('#').trim();
  }

  // Check for array notation (...)
  if (name.endsWith('...')) {
    isArray = true;
    name = name.slice(0, -3).trim();
  }

  // Check for boolean flag (starts with - or --)
  if (!required && (name.startsWith('-') || name.startsWith('--'))) {
    originalFlag = name;
    name = name.replace(/^-+/, '');
    if (description === '') {
      description = `Enable ${originalFlag} flag`;
    }
  }

  return new FieldToken(name, description, required, isArray, originalFlag);
}
