import { describe, it, expect } from 'vitest';
import { fromArgs } from './parse.js';
import { TextToken, FieldToken } from './types.js';

describe('Blueprint', () => {
  describe('BaseCommand', () => {
    const tests = [
      {
        name: 'simple command name',
        args: ['git', 'status'],
        expected: 'git',
      },
      {
        name: 'command name with dashes',
        args: ['git-status'],
        expected: 'git-status',
      },
      {
        name: 'command name with underscores',
        args: ['do_thing', '[paths...]'],
        expected: 'do_thing',
      },
      {
        name: 'empty args',
        args: [],
        wantErr: true,
        errMsg: 'no command provided',
      },
      {
        name: 'empty command',
        args: [''],
        wantErr: true,
        errMsg: 'empty command provided',
      },
      {
        name: 'whitespace-only command',
        args: ['   '],
        wantErr: true,
        errMsg: 'empty command provided',
      },
    ];

    tests.forEach((tt) => {
      it(tt.name, () => {
        if (tt.wantErr) {
          expect(() => fromArgs(tt.args)).toThrow(tt.errMsg);
        } else {
          const bp = fromArgs(tt.args);
          expect(bp.baseCommand).toBe(tt.expected);
        }
      });
    });
  });

  describe('GetCommandFormat', () => {
    const tests = [
      {
        name: 'command without args',
        args: ['git'],
        expected: 'git',
      },
      {
        name: 'command with literal args',
        args: ['git', 'status'],
        expected: 'git status',
      },
      {
        name: 'command with array args',
        args: ['git', '[args...]'],
        expected: 'git [args...]',
      },
      {
        name: 'command with array args',
        args: ['echo', '{{args...}}'],
        expected: 'echo {{args...}}',
      },
      {
        name: 'command with literal and array args',
        args: ['git', 'status', '[args...]'],
        expected: 'git status [args...]',
      },
      {
        name: 'command with array with description',
        args: ['git', '[args... #Additional command line arguments]'],
        expected: 'git [args...]',
      },
      {
        name: 'command with required field',
        args: ['echo', '{{text}}'],
        expected: 'echo {{text}}',
      },
      {
        name: 'command with required field with description',
        args: ['echo', '{{text#A required field}}'],
        expected: 'echo {{text}}',
      },
      {
        name: 'command with optional field',
        args: ['echo', '[text]'],
        expected: 'echo [text]',
      },
      {
        name: 'command with optional field with description',
        args: ['echo', '[text#description]'],
        expected: 'echo [text]',
      },
      {
        name: 'command with suffix literal',
        args: ['echo', '[text]suffix'],
        expected: 'echo [text]suffix',
      },
      {
        name: 'command with prefix literal',
        args: ['echo', 'prefix[text]'],
        expected: 'echo prefix[text]',
      },
      {
        name: 'command with prefix and suffix literal',
        args: ['echo', 'prefix[text]suffix'],
        expected: 'echo prefix[text]suffix',
      },
      {
        name: 'command with embedded template',
        args: ['curl', 'https://en.m.wikipedia.org/wiki/{{page#A valid wikipedia page}}'],
        expected: 'curl https://en.m.wikipedia.org/wiki/{{page}}',
      },
      {
        name: 'command with spaces in field description',
        args: ['curl', 'https://en.m.wikipedia.org/wiki/{{page # A valid wikipedia page}}'],
        expected: 'curl https://en.m.wikipedia.org/wiki/{{page}}',
      },
      {
        name: 'mixed blueprints with required and optional arguments',
        args: ['command', '{{arg1#Custom description}}', '[arg2]'],
        expected: 'command {{arg1}} [arg2]',
      },
      {
        name: 'prioritizes explicit description over default',
        args: ['echo', '{{text#Explicit description}}', '{{text}}'],
        expected: 'echo {{text}} {{text}}',
      },
      {
        name: 'preserves underscores in field names for display',
        args: ['echo', '[has_underscores]'],
        expected: 'echo [has_underscores]',
      },
      {
        name: 'preserves dashes in field names for display',
        args: ['echo', '[has-dashes]'],
        expected: 'echo [has-dashes]',
      },
      {
        name: 'short boolean flag without description',
        args: ['ls', '[-f]'],
        expected: 'ls [-f]',
      },
      {
        name: 'long boolean flag without description',
        args: ['ls', '[--force]'],
        expected: 'ls [--force]',
      },
      {
        name: 'boolean flag with description',
        args: ['rm', '[-f#force removal]'],
        expected: 'rm [-f]',
      },
      {
        name: 'long boolean flag with description',
        args: ['ls', '[--force#force removal]'],
        expected: 'ls [--force]',
      },
      {
        name: 'required flag with description',
        args: ['cp', '{{-r#recursive}}'],
        expected: 'cp {{-r}}',
      },
      {
        name: 'complicated template with mixed text and fields',
        args: ['curl', 'http[s # use https]://api.com/{{endpoint#API endpoint}}', '[--verbose]'],
        expected: 'curl http[s]://api.com/{{endpoint}} [--verbose]',
      },
    ];

    tests.forEach((tt) => {
      it(tt.name, () => {
        const bp = fromArgs(tt.args);
        expect(bp.getCommandFormat()).toBe(tt.expected);
      });
    });
  });

  describe('FromArgsTokenization', () => {
    it('tokenizes simple command without templates', () => {
      const bp = fromArgs(['echo', 'hello']);

      const expected = [
        [new TextToken('echo')],
        [new TextToken('hello')],
      ];
      expect(bp.shellWords).toEqual(expected);
    });

    it('tokenizes command with simple template', () => {
      const bp = fromArgs(['echo', '{{text}}']);

      const expected = [
        [new TextToken('echo')],
        [new FieldToken('text', '', true, false, '')],
      ];
      expect(bp.shellWords).toEqual(expected);
    });

    it('tokenizes command with optional field', () => {
      const bp = fromArgs(['echo', '[optional]']);

      const expected = [
        [new TextToken('echo')],
        [new FieldToken('optional', '', false, false, '')],
      ];
      expect(bp.shellWords).toEqual(expected);
    });

    it('tokenizes command with template and description', () => {
      const bp = fromArgs(['echo', '{{text#message to echo}}']);

      const expected = [
        [new TextToken('echo')],
        [new FieldToken('text', 'message to echo', true, false, '')],
      ];
      expect(bp.shellWords).toEqual(expected);
    });

    it('tokenizes command with prefix text and template', () => {
      const bp = fromArgs(['echo', 'prefix{{text#desc}}']);

      const expected = [
        [new TextToken('echo')],
        [
          new TextToken('prefix'),
          new FieldToken('text', 'desc', true, false, ''),
        ],
      ];
      expect(bp.shellWords).toEqual(expected);
    });

    it('tokenizes command with suffix text and template', () => {
      const bp = fromArgs(['echo', '{{text#desc}}suffix']);

      const expected = [
        [new TextToken('echo')],
        [
          new FieldToken('text', 'desc', true, false, ''),
          new TextToken('suffix'),
        ],
      ];
      expect(bp.shellWords).toEqual(expected);
    });

    it('tokenizes command with mixed text and template', () => {
      const bp = fromArgs(['echo', 'prefix{{text#desc}}suffix']);

      const expected = [
        [new TextToken('echo')],
        [
          new TextToken('prefix'),
          new FieldToken('text', 'desc', true, false, ''),
          new TextToken('suffix'),
        ],
      ];
      expect(bp.shellWords).toEqual(expected);
    });

    it('tokenizes command with prefix text and optional field', () => {
      const bp = fromArgs(['echo', 'prefix[text]']);

      const expected = [
        [new TextToken('echo')],
        [
          new TextToken('prefix'),
          new FieldToken('text', '', false, false, ''),
        ],
      ];
      expect(bp.shellWords).toEqual(expected);
    });

    it('tokenizes command with suffix text and optional field', () => {
      const bp = fromArgs(['echo', '[text#description]suffix']);

      const expected = [
        [new TextToken('echo')],
        [
          new FieldToken('text', 'description', false, false, ''),
          new TextToken('suffix'),
        ],
      ];
      expect(bp.shellWords).toEqual(expected);
    });

    it('tokenizes command with mixed text and optional field', () => {
      const bp = fromArgs(['echo', 'prefix[text]suffix']);

      const expected = [
        [new TextToken('echo')],
        [
          new TextToken('prefix'),
          new FieldToken('text', '', false, false, ''),
          new TextToken('suffix'),
        ],
      ];
      expect(bp.shellWords).toEqual(expected);
    });

    it('tokenizes command with optional field and description', () => {
      const bp = fromArgs(['echo', '[optional#optional text]']);

      const expected = [
        [new TextToken('echo')],
        [new FieldToken('optional', 'optional text', false, false, '')],
      ];
      expect(bp.shellWords).toEqual(expected);
    });

    it('tokenizes complex mixed command', () => {
      const bp = fromArgs([
        'curl',
        'http[s # use https]://api.com/{{endpoint#API endpoint}}',
        '[--verbose]',
      ]);

      const expected = [
        [new TextToken('curl')],
        [
          new TextToken('http'),
          new FieldToken('s', 'use https', false, false, ''),
          new TextToken('://api.com/'),
          new FieldToken('endpoint', 'API endpoint', true, false, ''),
        ],
        [new FieldToken('verbose', 'Enable --verbose flag', false, false, '--verbose')],
      ];
      expect(bp.shellWords).toEqual(expected);
    });
  });

  describe('SingleBraceTokenization', () => {
    it('tokenizes command with single-brace required field', () => {
      const bp = fromArgs(['echo', '{text}']);

      const expected = [
        [new TextToken('echo')],
        [new FieldToken('text', '', true, false, '')],
      ];
      expect(bp.shellWords).toEqual(expected);
      expect(bp.getCommandFormat()).toBe('echo {{text}}');
    });

    it('tokenizes command with single-brace and description', () => {
      const bp = fromArgs(['echo', '{text#message to echo}']);

      const expected = [
        [new TextToken('echo')],
        [new FieldToken('text', 'message to echo', true, false, '')],
      ];
      expect(bp.shellWords).toEqual(expected);
      expect(bp.getCommandFormat()).toBe('echo {{text}}');
    });

    it('tokenizes command with prefix and single-brace template', () => {
      const bp = fromArgs(['echo', 'prefix{text#desc}']);

      const expected = [
        [new TextToken('echo')],
        [
          new TextToken('prefix'),
          new FieldToken('text', 'desc', true, false, ''),
        ],
      ];
      expect(bp.shellWords).toEqual(expected);
      expect(bp.getCommandFormat()).toBe('echo prefix{{text}}');
    });

    it('tokenizes command with suffix and single-brace template', () => {
      const bp = fromArgs(['echo', '{text#desc}suffix']);

      const expected = [
        [new TextToken('echo')],
        [
          new FieldToken('text', 'desc', true, false, ''),
          new TextToken('suffix'),
        ],
      ];
      expect(bp.shellWords).toEqual(expected);
      expect(bp.getCommandFormat()).toBe('echo {{text}}suffix');
    });

    it('tokenizes command with mixed text and single-brace template', () => {
      const bp = fromArgs(['echo', 'prefix{text#desc}suffix']);

      const expected = [
        [new TextToken('echo')],
        [
          new TextToken('prefix'),
          new FieldToken('text', 'desc', true, false, ''),
          new TextToken('suffix'),
        ],
      ];
      expect(bp.shellWords).toEqual(expected);
      expect(bp.getCommandFormat()).toBe('echo prefix{{text}}suffix');
    });
  });

  describe('SingleBraceEnhancedTemplateProcessing', () => {
    it('handles malformed single-brace template syntax', () => {
      const bp = fromArgs(['echo', '{incomplete']);
      // Implementation will be tested when render is implemented
      // For now just verify it parses
      expect(bp.baseCommand).toBe('echo');
    });

    it('handles malformed single-brace with only closing brace', () => {
      const bp = fromArgs(['echo', 'no_opening_brace}']);
      expect(bp.baseCommand).toBe('echo');
    });

    it('handles empty single-brace template braces', () => {
      const bp = fromArgs(['echo', '{}']);
      expect(bp.baseCommand).toBe('echo');
    });
  });
});
