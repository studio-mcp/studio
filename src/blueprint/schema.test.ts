import { describe, it, expect } from 'vitest';
import { fromArgs } from './parse.js';
import { generateInputSchema, type Schema } from './schema.js';

describe('Blueprint.GenerateInputSchema', () => {
  const testCases: Array<{
    name: string;
    args: string[];
    expectedSchema: Schema;
  }> = [
    {
      name: 'simple command without args',
      args: ['git', 'status'],
      expectedSchema: {
        type: 'object',
        properties: {},
        required: [],
      },
    },
    {
      name: 'simple command with array args',
      args: ['git', 'status', '[args...]'],
      expectedSchema: {
        type: 'object',
        properties: {
          args: {
            type: 'array',
            items: { type: 'string' },
            description: 'Additional command line arguments',
          },
        },
        required: [],
      },
    },
    {
      name: 'template command with description',
      args: ['curl', 'https://en.m.wikipedia.org/wiki/{{page#A valid wikipedia page}}'],
      expectedSchema: {
        type: 'object',
        properties: {
          page: {
            type: 'string',
            description: 'A valid wikipedia page',
          },
        },
        required: ['page'],
      },
    },
    {
      name: 'template command without description',
      args: ['echo', '{{text}}'],
      expectedSchema: {
        type: 'object',
        properties: {
          text: {
            type: 'string',
          },
        },
        required: ['text'],
      },
    },
    {
      name: 'template with spaces in description',
      args: ['curl', 'https://en.m.wikipedia.org/wiki/{{page # A valid wikipedia page}}'],
      expectedSchema: {
        type: 'object',
        properties: {
          page: {
            type: 'string',
            description: 'A valid wikipedia page',
          },
        },
        required: ['page'],
      },
    },
    {
      name: 'mixed required and optional arguments',
      args: ['command', '{{arg1#Custom description}}', '[arg2]'],
      expectedSchema: {
        type: 'object',
        properties: {
          arg1: {
            type: 'string',
            description: 'Custom description',
          },
          arg2: {
            type: 'string',
          },
        },
        required: ['arg1'],
      },
    },
    {
      name: 'duplicate template with explicit description priority',
      args: ['echo', '{{text#Explicit description}}', '{{text}}'],
      expectedSchema: {
        type: 'object',
        properties: {
          text: {
            type: 'string',
            description: 'Explicit description',
          },
        },
        required: ['text'],
      },
    },
    {
      name: 'array arguments with default description',
      args: ['echo', '[files...]'],
      expectedSchema: {
        type: 'object',
        properties: {
          files: {
            type: 'array',
            items: { type: 'string' },
            description: 'Additional command line arguments',
          },
        },
        required: [],
      },
    },
    {
      name: 'array arguments with custom description',
      args: ['ls', '[files...#Files to list]'],
      expectedSchema: {
        type: 'object',
        properties: {
          files: {
            type: 'array',
            items: { type: 'string' },
            description: 'Files to list',
          },
        },
        required: [],
      },
    },
    {
      name: 'optional string field',
      args: ['echo', '[optional]'],
      expectedSchema: {
        type: 'object',
        properties: {
          optional: {
            type: 'string',
          },
        },
        required: [],
      },
    },
    {
      name: 'optional field with description',
      args: ['echo', "[name#Person's name]"],
      expectedSchema: {
        type: 'object',
        properties: {
          name: {
            type: 'string',
            description: "Person's name",
          },
        },
        required: [],
      },
    },
    {
      name: 'dashes converted to underscores',
      args: ['echo', '[has-dashes]'],
      expectedSchema: {
        type: 'object',
        properties: {
          has_dashes: {
            type: 'string',
          },
        },
        required: [],
      },
    },
    {
      name: 'boolean flag short form',
      args: ['ls', '[-f]'],
      expectedSchema: {
        type: 'object',
        properties: {
          f: {
            type: 'boolean',
            description: 'Enable -f flag',
          },
        },
        required: [],
      },
    },
    {
      name: 'boolean flag long form',
      args: ['ls', '[--force]'],
      expectedSchema: {
        type: 'object',
        properties: {
          force: {
            type: 'boolean',
            description: 'Enable --force flag',
          },
        },
        required: [],
      },
    },
    {
      name: 'boolean flag with custom description',
      args: ['rm', '[-f#force removal]'],
      expectedSchema: {
        type: 'object',
        properties: {
          f: {
            type: 'boolean',
            description: 'force removal',
          },
        },
        required: [],
      },
    },
    {
      name: 'mixed boolean flags and required arguments',
      args: ['cp', '[-r#recursive]', '{{source}}', '{{dest}}'],
      expectedSchema: {
        type: 'object',
        properties: {
          r: {
            type: 'boolean',
            description: 'recursive',
          },
          source: {
            type: 'string',
          },
          dest: {
            type: 'string',
          },
        },
        required: ['source', 'dest'],
      },
    },
    {
      name: 'mixed string and array arguments',
      args: ['command', '{{flag#Command flag}}', '[files...]'],
      expectedSchema: {
        type: 'object',
        properties: {
          flag: {
            type: 'string',
            description: 'Command flag',
          },
          files: {
            type: 'array',
            items: { type: 'string' },
            description: 'Additional command line arguments',
          },
        },
        required: ['flag'],
      },
    },
  ];

  testCases.forEach((tc) => {
    it(tc.name, () => {
      const bp = fromArgs(tc.args);
      const actualSchema = generateInputSchema(bp);
      expect(actualSchema).toEqual(tc.expectedSchema);
    });
  });
});

describe('Blueprint.GenerateInputSchema EdgeCases', () => {
  it('mixed optional arguments with and without descriptions', () => {
    const bp = fromArgs(['cmd', '[arg1]', '[arg2#Custom desc]']);
    const schema = generateInputSchema(bp);

    // Check properties individually since the order might vary
    expect(schema.properties?.arg1?.type).toBe('string');
    expect(schema.properties?.arg1?.description).toBeUndefined();

    expect(schema.properties?.arg2?.type).toBe('string');
    expect(schema.properties?.arg2?.description).toBe('Custom desc');

    // Neither should be required
    expect(schema.required).toEqual([]);
  });
});
