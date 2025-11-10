#!/usr/bin/env node
import { McpServer } from '@modelcontextprotocol/sdk/server/mcp.js';
import { StdioServerTransport } from '@modelcontextprotocol/sdk/server/stdio.js';
import { z } from 'zod';
import { spawn } from 'child_process';
import * as fs from 'fs';
import * as path from 'path';
import { fromArgs } from './blueprint/parse.js';
import { generateInputSchema, type Schema } from './blueprint/schema.js';
import { buildCommandArgs } from './blueprint/render.js';

const VERSION = '0.0.0';

/**
 * Parses command line arguments
 */
function parseArgs(args: string[]): {
  debug: boolean;
  version: boolean;
  logFile: string;
  commandArgs: string[];
} {
  let i = 0;
  let debug = false;
  let version = false;
  let logFile = '';
  const commandArgs: string[] = [];

  // Parse studio flags until we hit a non-flag or --
  while (i < args.length) {
    const arg = args[i];

    // If it's --, stop parsing flags and treat everything after as command args
    if (arg === '--') {
      i++;
      break;
    }

    // If it doesn't start with -, we're done with studio flags
    if (!arg.startsWith('-')) {
      break;
    }

    switch (arg) {
      case '--debug':
        debug = true;
        break;
      case '--version':
        version = true;
        break;
      case '--log':
        // Check if we have a next argument for the filename
        if (i + 1 >= args.length) {
          throw new Error('--log requires a filename argument');
        }
        i++;
        // Check if the next argument is another flag
        if (args[i].startsWith('-')) {
          throw new Error('--log requires a filename argument');
        }
        logFile = args[i];
        break;
      case '-h':
      case '--help':
        console.log(`studio - One word MCP for any CLI command

Usage: studio [--debug] [--log filename] [--] <command> [args...]

Options:
  -h, --help            Show this help message and exit
  --version             Show version information and exit
  --debug               Print debug logs to stderr
  --log <filename>      Write debug logs to specified file
  --                    End flag parsing, treat rest as command

Template Syntax:
  {{name}}              Required argument
  {name}                Required argument (single brace)
  [name]                Optional argument
  [args...]             Optional array argument
  {{args...}}           Required array argument
  {{name#description}}  Argument with description
  [-f]                  Boolean flag (optional)

Example:
  studio echo "{{text#message to echo}}"
  studio git "[args...]"
`);
        process.exit(0);
        break;
      default:
        throw new Error(`unknown flag: ${arg}`);
    }

    i++;
  }

  // Everything from i onwards goes to blueprint parsing
  commandArgs.push(...args.slice(i));

  return { debug, version, logFile, commandArgs };
}

/**
 * Converts a JSON schema to Zod schema
 */
function schemaToZod(schema: Schema) {
  const shape: Record<string, z.ZodTypeAny> = {};

  for (const [key, prop] of Object.entries(schema.properties || {})) {
    let fieldSchema: z.ZodTypeAny;

    if (prop.type === 'array') {
      fieldSchema = z.array(z.string());
    } else if (prop.type === 'boolean') {
      fieldSchema = z.boolean();
    } else {
      fieldSchema = z.string();
    }

    // Add description if available
    if (prop.description) {
      fieldSchema = fieldSchema.describe(prop.description);
    }

    // Make optional if not in required array
    if (!schema.required?.includes(key)) {
      fieldSchema = fieldSchema.optional();
    }

    shape[key] = fieldSchema;
  }

  return shape;
}

/**
 * Executes a command
 */
async function execute(command: string, args: string[]): Promise<string> {
  return new Promise((resolve, reject) => {
    const proc = spawn(command, args);
    let stdout = '';
    let stderr = '';

    proc.stdout.on('data', (data) => {
      stdout += data.toString();
    });

    proc.stderr.on('data', (data) => {
      stderr += data.toString();
    });

    proc.on('close', (code) => {
      const output = (stdout + '\n' + stderr).trim();
      if (code !== 0) {
        reject(new Error(`command failed with exit code ${code}`));
      } else {
        resolve(output);
      }
    });

    proc.on('error', (err) => {
      reject(new Error(`Studio error: ${err.message}`));
    });
  });
}

/**
 * Main function
 */
async function main() {
  const args = process.argv.slice(2);

  try {
    const { debug, version, logFile, commandArgs } = parseArgs(args);

    // Handle version flag
    if (version) {
      console.log(`studio ${VERSION}`);
      console.log(`commit: dev`);
      console.log(`built: ${new Date().toISOString()}`);
      return;
    }

    // Validate command args
    if (commandArgs.length === 0) {
      console.error(
        'usage: studio <command> --example "{{req # required arg}}" "[args... # array of args]"'
      );
      process.exit(1);
    }

    // Create blueprint
    const blueprint = fromArgs(commandArgs);
    const inputSchema = generateInputSchema(blueprint);
    const zodSchema = schemaToZod(inputSchema);

    // Create MCP server
    const server = new McpServer({
      name: 'studio',
      version: VERSION,
    });

    // Generate tool name from base command
    const toolName = blueprint.baseCommand.replace(/-/g, '_');
    const toolDescription = `Run the shell command \`${blueprint.getCommandFormat()}\``;

    // Register tool
    server.registerTool(
      toolName,
      {
        title: blueprint.baseCommand,
        description: toolDescription,
        inputSchema: zodSchema,
      },
      async (params: any) => {
        try {
          // Build command args
          const fullCommand = buildCommandArgs(blueprint, params);

          if (debug) {
            console.error(`[Studio MCP] Executing: ${fullCommand.join(' ')}`);
          }

          // Execute command
          const output = await execute(fullCommand[0], fullCommand.slice(1));

          return {
            content: [
              {
                type: 'text' as const,
                text: output,
              },
            ],
          };
        } catch (err: any) {
          const errorMsg = err.message || String(err);
          return {
            content: [
              {
                type: 'text' as const,
                text: errorMsg,
              },
            ],
            isError: true,
          };
        }
      }
    );

    // Create transport
    const transport = new StdioServerTransport();

    // Connect server to transport
    await server.connect(transport);

    if (debug) {
      console.error('[Studio MCP] Server started and listening');
    }
  } catch (err: any) {
    console.error(`Error: ${err.message}`);
    process.exit(1);
  }
}

// Run main function
main().catch((err) => {
  console.error(`Fatal error: ${err.message}`);
  process.exit(1);
});
