#!/usr/bin/env node
import { McpServer } from "@modelcontextprotocol/sdk/server/mcp.js";
import { StdioServerTransport } from "@modelcontextprotocol/sdk/server/stdio.js";
import { z } from "zod";
import { spawn } from "child_process";
import * as fs from "fs";
import * as path from "path";
import { fileURLToPath } from "url";
import { fromArgs } from "./studio/parse.js";
import { generateInputSchema, type Schema } from "./studio/schema.js";
import { buildCommandArgs } from "./studio/render.js";

// Read version from package.json
const __filename = fileURLToPath(import.meta.url);
const __dirname = path.dirname(__filename);
const packageJsonPath = path.join(__dirname, "..", "package.json");
const packageJson = JSON.parse(fs.readFileSync(packageJsonPath, "utf-8"));
const VERSION = packageJson.version;

/**
 * Parses command line arguments
 */
function parseArgs(args: string[]): {
  debug: boolean;
  version: boolean;
  logFile: string;
  envVars: Record<string, string>;
  commandArgs: string[];
} {
  let i = 0;
  let debug = false;
  let version = false;
  let logFile = "";
  const envVars: Record<string, string> = {};
  const commandArgs: string[] = [];

  // Parse studio flags until we hit a non-flag or --
  while (i < args.length) {
    const arg = args[i];

    // If it's --, stop parsing flags and treat everything after as command args
    if (arg === "--") {
      i++;
      break;
    }

    // If it doesn't start with -, we're done with studio flags
    if (!arg.startsWith("-")) {
      break;
    }

    switch (arg) {
      case "--debug":
        debug = true;
        break;
      case "--version":
        version = true;
        break;
      case "--log":
        // Check if we have a next argument for the filename
        if (i + 1 >= args.length) {
          throw new Error("--log requires a filename argument");
        }
        i++;
        // Check if the next argument is another flag
        if (args[i].startsWith("-")) {
          throw new Error("--log requires a filename argument");
        }
        logFile = args[i];
        break;
      case "-e":
      case "--env":
        // Check if we have a next argument for the env var
        if (i + 1 >= args.length) {
          throw new Error(`${arg} requires an environment variable in KEY=VALUE format`);
        }
        i++;
        const envArg = args[i];
        // Parse KEY=VALUE format
        const eqIndex = envArg.indexOf("=");
        if (eqIndex === -1) {
          throw new Error(`${arg} requires an environment variable in KEY=VALUE format, got: ${envArg}`);
        }
        const key = envArg.substring(0, eqIndex);
        const value = envArg.substring(eqIndex + 1);
        if (!key) {
          throw new Error(`${arg} requires a non-empty key in KEY=VALUE format`);
        }
        envVars[key] = value;
        break;
      case "-h":
      case "--help":
        console.log(`studio - One word MCP for any CLI command

Usage: studio [--debug] [--log filename] [-e KEY=VALUE] [--] <command> [args...]

Options:
  -h, --help            Show this help message and exit
  --version             Show version information and exit
  --debug               Print debug logs to stderr
  --log <filename>      Write debug logs to specified file
  -e, --env KEY=VALUE   Set environment variable for command (can be used multiple times)
  --                    End flag parsing, treat rest as command

Template Syntax:
  {{name}}              Required argument
  {name}                Required argument (single brace)
  [name]                Optional argument
  [args...]             Optional array argument
  {args...}             Required array argument
  {name#description}    Argument with description
  [-f]                  Boolean flag (optional)

Example:
  studio echo "{text#message to echo}"
  studio git "[args...]"
`);
        process.exit(0);
        break;
      default:
        throw new Error(`unknown flag: ${arg}`);
    }

    i++;
  }

  // Everything from i onwards goes to command template parsing
  commandArgs.push(...args.slice(i));

  return { debug, version, logFile, envVars, commandArgs };
}

/**
 * Converts a JSON schema to Zod schema
 */
function schemaToZod(schema: Schema) {
  const shape: Record<string, z.ZodTypeAny> = {};

  for (const [key, prop] of Object.entries(schema.properties || {})) {
    let fieldSchema: z.ZodTypeAny;

    if (prop.type === "array") {
      fieldSchema = z.array(z.string());
    } else if (prop.type === "boolean") {
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
async function execute(
  command: string,
  args: string[],
  envVars: Record<string, string> = {},
): Promise<string> {
  return new Promise((resolve, reject) => {
    // Prepare environment variables
    const env = { ...process.env, ...envVars };

    // Handle PWD special case
    let cwd: string | undefined = undefined;
    if (envVars.PWD) {
      cwd = envVars.PWD;
      // Remove PWD from env since we're setting it via cwd option
      delete env.PWD;
    }

    const proc = spawn(command, args, { env, cwd });
    let stdout = "";
    let stderr = "";

    proc.stdout.on("data", (data) => {
      stdout += data.toString();
    });

    proc.stderr.on("data", (data) => {
      stderr += data.toString();
    });

    proc.on("close", (code) => {
      const output = (stdout + "\n" + stderr).trim();
      if (code !== 0) {
        reject(new Error(`command failed with exit code ${code}`));
      } else {
        resolve(output);
      }
    });

    proc.on("error", (err) => {
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
    const { debug, version, logFile, envVars, commandArgs } = parseArgs(args);

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
        'usage: studio <command> --example "{{req # required arg}}" "[args... # array of args]"',
      );
      process.exit(1);
    }

    // Create template
    const template = fromArgs(commandArgs);
    const inputSchema = generateInputSchema(template);
    const zodSchema = schemaToZod(inputSchema);

    // Create MCP server
    const server = new McpServer({
      name: "studio",
      version: VERSION,
    });

    // Generate tool name from base command
    const toolName = template.baseCommand.replace(/-/g, "_");
    const toolDescription = `Run the shell command \`${template.getCommandFormat()}\``;

    // Register tool
    server.registerTool(
      toolName,
      {
        title: template.baseCommand,
        description: toolDescription,
        inputSchema: zodSchema,
      },
      async (params: any) => {
        try {
          // Build command args
          const fullCommand = buildCommandArgs(template, params);

          if (debug) {
            console.error(`[Studio MCP] Executing: ${fullCommand.join(" ")}`);
          }

          // Execute command
          const output = await execute(fullCommand[0], fullCommand.slice(1), envVars);

          return {
            content: [
              {
                type: "text" as const,
                text: output,
              },
            ],
          };
        } catch (err: any) {
          const errorMsg = err.message || String(err);
          return {
            content: [
              {
                type: "text" as const,
                text: errorMsg,
              },
            ],
            isError: true,
          };
        }
      },
    );

    // Create transport
    const transport = new StdioServerTransport();

    // Connect server to transport
    await server.connect(transport);

    if (debug) {
      console.error("[Studio MCP] Server started and listening");
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
