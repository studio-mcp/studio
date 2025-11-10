import { describe, it, expect } from 'vitest';
import { fromArgs } from './parse.js';
import { buildCommandArgs } from './render.js';

describe('Blueprint.BuildCommandArgs', () => {
  it('builds simple command without templates', () => {
    const bp = fromArgs(['echo', 'hello', 'world']);
    const args = buildCommandArgs(bp, {});
    expect(args).toEqual(['echo', 'hello', 'world']);
  });

  it('builds command with required template', () => {
    const bp = fromArgs(['echo', '{{message}}']);
    const args = buildCommandArgs(bp, { message: 'Hello World' });
    expect(args).toEqual(['echo', 'Hello World']);
  });

  it('builds command with template in middle of arg', () => {
    const bp = fromArgs(['curl', 'https://api.example.com/{{endpoint}}']);
    const args = buildCommandArgs(bp, { endpoint: 'users/123' });
    expect(args).toEqual(['curl', 'https://api.example.com/users/123']);
  });

  it('builds command with multiple templates in one arg', () => {
    const bp = fromArgs(['echo', '{{greeting}} {{name}}!']);
    const args = buildCommandArgs(bp, { greeting: 'Hello', name: 'World' });
    expect(args).toEqual(['echo', 'Hello World!']);
  });

  it('builds command with array argument', () => {
    const bp = fromArgs(['echo', '[files...]']);
    const args = buildCommandArgs(bp, {
      files: ['file1.txt', 'file2.txt', 'file3.txt'],
    });
    expect(args).toEqual(['echo', 'file1.txt', 'file2.txt', 'file3.txt']);
  });

  it('builds command with empty array argument', () => {
    const bp = fromArgs(['echo', 'prefix', '[files...]']);
    const args = buildCommandArgs(bp, { files: [] });
    expect(args).toEqual(['echo', 'prefix']);
  });

  it('builds command with optional string argument provided', () => {
    const bp = fromArgs(['echo', 'hello', '[name]']);
    const args = buildCommandArgs(bp, { name: 'world' });
    expect(args).toEqual(['echo', 'hello', 'world']);
  });

  it('builds command with optional string argument omitted', () => {
    const bp = fromArgs(['echo', 'hello', '[name]']);
    const args = buildCommandArgs(bp, {});
    expect(args).toEqual(['echo', 'hello']);
  });

  it('builds command with blueprint arguments containing spaces', () => {
    const bp = fromArgs(['echo', '{{text#text to echo}}']);
    const args = buildCommandArgs(bp, { text: 'Hello World' });
    expect(args).toEqual(['echo', 'Hello World']);
  });

  it('builds command with mixed blueprint with and without descriptions', () => {
    const bp = fromArgs(['echo', '{{greeting#The greeting}}', '{{name}}']);
    const args = buildCommandArgs(bp, { greeting: 'Hello', name: 'World' });
    expect(args).toEqual(['echo', 'Hello', 'World']);
  });

  it('builds command with blueprint arguments in mixed content', () => {
    const bp = fromArgs(['echo', 'simon says {{text#text for simon to say}}']);
    const args = buildCommandArgs(bp, { text: 'Hello World' });
    expect(args).toEqual(['echo', 'simon says Hello World']);
  });

  it('builds command with blueprint arguments containing special shell characters', () => {
    const bp = fromArgs(['echo', '{{text#text to echo}}']);
    const args = buildCommandArgs(bp, { text: 'Hello & World; echo pwned' });
    expect(args).toEqual(['echo', 'Hello & World; echo pwned']);
  });

  it('builds command with blueprint in middle of argument', () => {
    const bp = fromArgs(['echo', '--message={{text#message content}}']);
    const args = buildCommandArgs(bp, { text: 'Hello World' });
    expect(args).toEqual(['echo', '--message=Hello World']);
  });

  it('builds command with blueprint with prefix and suffix', () => {
    const bp = fromArgs(['echo', 'prefix-{{text#middle part}}-suffix']);
    const args = buildCommandArgs(bp, { text: 'Hello World' });
    expect(args).toEqual(['echo', 'prefix-Hello World-suffix']);
  });

  it('builds command with mixed blueprint and non-blueprint arguments', () => {
    const bp = fromArgs([
      'echo',
      'static',
      '{{dynamic#dynamic content}}',
      'more-static',
    ]);
    const args = buildCommandArgs(bp, { dynamic: 'Hello World' });
    expect(args).toEqual(['echo', 'static', 'Hello World', 'more-static']);
  });

  it('preserves shell safety with complex blueprint values', () => {
    const bp = fromArgs(['echo', 'Result: {{text#text content}}']);
    const args = buildCommandArgs(bp, {
      text: "$(echo 'dangerous'); echo 'safe'",
    });
    expect(args).toEqual(['echo', "Result: $(echo 'dangerous'); echo 'safe'"]);
  });

  it('builds command with mixed string and array arguments', () => {
    const bp = fromArgs(['echo', '{{prefix#Prefix text}}', '[files...]']);
    const args = buildCommandArgs(bp, {
      prefix: 'Files:',
      files: ['a.txt', 'b.txt'],
    });
    expect(args).toEqual(['echo', 'Files:', 'a.txt', 'b.txt']);
  });

  it('builds command with mixed required and optional arguments - both provided', () => {
    const bp = fromArgs(['echo', '{{required#Required text}}', '[optional]']);
    const args = buildCommandArgs(bp, { required: 'hello', optional: 'world' });
    expect(args).toEqual(['echo', 'hello', 'world']);
  });

  it('builds command with mixed required and optional arguments - only required', () => {
    const bp = fromArgs(['echo', '{{required#Required text}}', '[optional]']);
    const args = buildCommandArgs(bp, { required: 'hello' });
    expect(args).toEqual(['echo', 'hello']);
  });
});

describe('Blueprint.TemplateValidation', () => {
  it('validates missing required parameters', () => {
    const bp = fromArgs(['echo', '{{required}}']);
    expect(() => buildCommandArgs(bp, {})).toThrow('missing required parameter');
  });

  it('validates parameter type mismatches', () => {
    const bp = fromArgs(['echo', '[files...]']);
    expect(() => buildCommandArgs(bp, { files: 'not-an-array' })).toThrow(
      "parameter 'files' must be an array"
    );
  });
});

describe('Blueprint.EnhancedTemplateProcessing', () => {
  it('handles malformed template syntax', () => {
    const bp = fromArgs(['echo', '{{incomplete']);
    const args = buildCommandArgs(bp, {});
    expect(args).toEqual(['echo', '{{incomplete']);
  });

  it('handles malformed template with only opening braces', () => {
    const bp = fromArgs(['echo', '{{no_closing_braces']);
    const args = buildCommandArgs(bp, {});
    expect(args).toEqual(['echo', '{{no_closing_braces']);
  });

  it('handles malformed template with only closing braces', () => {
    const bp = fromArgs(['echo', 'no_opening_braces}}']);
    const args = buildCommandArgs(bp, {});
    expect(args).toEqual(['echo', 'no_opening_braces}}']);
  });

  it('handles mixed valid and malformed templates', () => {
    const bp = fromArgs(['echo', '{{valid}}', '{{incomplete', '}}']);
    const args = buildCommandArgs(bp, { valid: 'works' });
    expect(args).toEqual(['echo', 'works', '{{incomplete', '}}']);
  });

  it('handles empty template braces', () => {
    const bp = fromArgs(['echo', '{{}}']);
    const args = buildCommandArgs(bp, {});
    expect(args).toEqual(['echo', '{{}}']);
  });
});

describe('Blueprint.BuildCommandArgsWithBooleanFlags', () => {
  it('builds command with boolean flag enabled', () => {
    const bp = fromArgs(['ls', '[-f]']);
    const args = buildCommandArgs(bp, { f: true });
    expect(args).toEqual(['ls', '-f']);
  });

  it('builds command with boolean flag disabled', () => {
    const bp = fromArgs(['ls', '[-f]']);
    const args = buildCommandArgs(bp, { f: false });
    expect(args).toEqual(['ls']);
  });

  it('builds command with boolean flag omitted', () => {
    const bp = fromArgs(['ls', '[-f]']);
    const args = buildCommandArgs(bp, {});
    expect(args).toEqual(['ls']);
  });

  it('builds command with long boolean flag enabled', () => {
    const bp = fromArgs(['ls', '[--force]']);
    const args = buildCommandArgs(bp, { force: true });
    expect(args).toEqual(['ls', '--force']);
  });

  it('builds command with mixed boolean and string arguments', () => {
    const bp = fromArgs(['cp', '[-r]', '{{source}}', '{{dest}}']);
    const args = buildCommandArgs(bp, {
      r: true,
      source: 'file1.txt',
      dest: 'file2.txt',
    });
    expect(args).toEqual(['cp', '-r', 'file1.txt', 'file2.txt']);
  });

  it('builds command with mixed boolean flags some enabled some disabled', () => {
    const bp = fromArgs(['ls', '[-l]', '[-a]', '[--human-readable]']);
    const args = buildCommandArgs(bp, {
      l: true,
      a: false,
      human_readable: true,
    });
    expect(args).toEqual(['ls', '-l', '--human-readable']);
  });
});

describe('Blueprint.BuildCommandArgsTokenized', () => {
  it('builds simple command without templates', () => {
    const bp = fromArgs(['echo', 'hello', 'world']);
    const args = buildCommandArgs(bp, {});
    expect(args).toEqual(['echo', 'hello', 'world']);
  });

  it('builds command with required template', () => {
    const bp = fromArgs(['echo', '{{message}}']);
    const args = buildCommandArgs(bp, { message: 'Hello World' });
    expect(args).toEqual(['echo', 'Hello World']);
  });

  it('builds command with template in middle of arg', () => {
    const bp = fromArgs(['curl', 'https://api.example.com/{{endpoint}}']);
    const args = buildCommandArgs(bp, { endpoint: 'users/123' });
    expect(args).toEqual(['curl', 'https://api.example.com/users/123']);
  });

  it('builds command with multiple templates in one arg', () => {
    const bp = fromArgs(['echo', '{{greeting}} {{name}}!']);
    const args = buildCommandArgs(bp, { greeting: 'Hello', name: 'World' });
    expect(args).toEqual(['echo', 'Hello World!']);
  });

  it('builds command with optional field provided', () => {
    const bp = fromArgs(['echo', 'hello', '[name]']);
    const args = buildCommandArgs(bp, { name: 'world' });
    expect(args).toEqual(['echo', 'hello', 'world']);
  });

  it('builds command with optional field omitted', () => {
    const bp = fromArgs(['echo', 'hello', '[name]']);
    const args = buildCommandArgs(bp, {});
    expect(args).toEqual(['echo', 'hello']);
  });

  it('builds command with boolean flag enabled', () => {
    const bp = fromArgs(['ls', '[--verbose]']);
    const args = buildCommandArgs(bp, { verbose: true });
    expect(args).toEqual(['ls', '--verbose']);
  });

  it('builds command with boolean flag disabled', () => {
    const bp = fromArgs(['ls', '[--verbose]']);
    const args = buildCommandArgs(bp, { verbose: false });
    expect(args).toEqual(['ls']);
  });

  it('builds command with array argument', () => {
    const bp = fromArgs(['echo', '[files...]']);
    const args = buildCommandArgs(bp, {
      files: ['file1.txt', 'file2.txt', 'file3.txt'],
    });
    expect(args).toEqual(['echo', 'file1.txt', 'file2.txt', 'file3.txt']);
  });

  it('handles dash-underscore equivalence', () => {
    const bp = fromArgs(['echo', '{{my-var}}', '[--my-flag]']);
    const args = buildCommandArgs(bp, { my_var: 'hello', my_flag: true });
    expect(args).toEqual(['echo', 'hello', '--my-flag']);
  });

  it('returns error for missing required parameter', () => {
    const bp = fromArgs(['echo', '{{message}}']);
    expect(() => buildCommandArgs(bp, {})).toThrow('missing required parameter');
  });
});

describe('Blueprint.BuildCommandArgs_SingleBrace', () => {
  it('builds command with single-brace required template', () => {
    const bp = fromArgs(['echo', '{message}']);
    const args = buildCommandArgs(bp, { message: 'Hello World' });
    expect(args).toEqual(['echo', 'Hello World']);
  });

  it('builds command with single-brace in middle of arg', () => {
    const bp = fromArgs(['curl', 'https://api.example.com/{endpoint}']);
    const args = buildCommandArgs(bp, { endpoint: 'users/123' });
    expect(args).toEqual(['curl', 'https://api.example.com/users/123']);
  });

  it('builds command with multiple single-brace templates in one arg', () => {
    const bp = fromArgs(['echo', '{greeting} {name}!']);
    const args = buildCommandArgs(bp, { greeting: 'Hello', name: 'World' });
    expect(args).toEqual(['echo', 'Hello World!']);
  });

  it('builds command with single-brace required array', () => {
    const bp = fromArgs(['echo', '{files...}']);
    const args = buildCommandArgs(bp, { files: ['a.txt', 'b.txt'] });
    expect(args).toEqual(['echo', 'a.txt', 'b.txt']);
  });

  it('errors when single-brace required arg missing', () => {
    const bp = fromArgs(['echo', '{text}']);
    expect(() => buildCommandArgs(bp, {})).toThrow();
  });
});
