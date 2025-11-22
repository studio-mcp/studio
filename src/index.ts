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
import {
  loadConfig,
  filterCommands,
  getDefaultConfigPath,
  type StudioConfig,
} from "./studio/config.js";

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
  configPath: string;
  commandsFilter: string[];
  commandArgs: string[];
} {
  let i = 0;
  let debug = false;
  let version = false;
  let logFile = "";
  let configPath = "";
  const commandsFilter: string[] = [];
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
      case "--config":
        // Check if we have a next argument for the config path
        if (i + 1 >= args.length) {
          throw new Error("--config requires a path argument");
        }
        i++;
        // Check if the next argument is another flag
        if (args[i].startsWith("-")) {
          throw new Error("--config requires a path argument");
        }
        configPath = args[i];
        break;
      case "--commands":
        // Check if we have a next argument for the commands list
        if (i + 1 >= args.length) {
          throw new Error("--commands requires a comma-separated list");
        }
        i++;
        // Check if the next argument is another flag
        if (args[i].startsWith("-")) {
          throw new Error("--commands requires a comma-separated list");
        }
        commandsFilter.push(
          ...args[i].split(",").map((s) => s.trim()).filter((s) => s.length > 0),
        );
        break;
      case "-h":
      case "--help":
        console.log(`studio - One word MCP for any CLI command

Usage:
  studio [--debug] [--log filename] [--] <command> [args...]
  studio [--debug] [--log filename] [--config <path>] [--commands <list>]

Options:
  -h, --help            Show this help message and exit
  --version             Show version information and exit
  --debug               Print debug logs to stderr
  --log <filename>      Write debug logs to specified file
  --config <path>       Load commands from config file (default: ~/.config/studio/config.jsonc)
  --commands <list>     Comma-separated list of commands to load from config
  --                    End flag parsing, treat rest as command

Template Syntax:
  {{name}}              Required argument
  {name}                Required argument (single brace)
  [name]                Optional argument
  [args...]             Optional array argument
  {args...}             Required array argument
  {name#description}    Argument with description
  [-f]                  Boolean flag (optional)

Config File Format (.config/studio/config.jsonc):
  {
    "echo": {
      "description": "Run the echo command",  // optional
      "command": ["echo", "Hello World"]      // required
    }
  }

Examples:
  studio echo "{text#message to echo}"
  studio git "[args...]"
  studio --config myconfig.jsonc
  studio --commands "echo,git"
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

  return { debug, version, logFile, configPath, commandsFilter, commandArgs };
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
async function execute(command: string, args: string[]): Promise<string> {
  return new Promise((resolve, reject) => {
    const proc = spawn(command, args);
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
 * Register a config-based command as an MCP tool
 */
function registerConfigCommand(
  server: McpServer,
  name: string,
  config: { description?: string; command: string[] },
  debug: boolean,
) {
  const toolName = name.replace(/-/g, "_");
  const toolDescription =
    config.description || `Run the command: ${config.command.join(" ")}`;

  server.registerTool(
    toolName,
    {
      title: name,
      description: toolDescription,
      inputSchema: {}, // No parameters for config commands
    },
    async () => {
      try {
        if (debug) {
          console.error(`[Studio MCP] Executing: ${config.command.join(" ")}`);
        }

        const output = await execute(config.command[0], config.command.slice(1));

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
}

/**
 * Main function
 */
async function main() {
  const args = process.argv.slice(2);

  try {
    const { debug, version, logFile, configPath, commandsFilter, commandArgs } =
      parseArgs(args);

    // Handle version flag
    if (version) {
      console.log(`studio ${VERSION}`);
      console.log(`commit: dev`);
      console.log(`built: ${new Date().toISOString()}`);
      return;
    }

    // Create MCP server
    const server = new McpServer({
      name: "studio",
      version: VERSION,
    });

    // Determine if we're in config mode
    const useConfigMode = configPath || commandsFilter.length > 0 || commandArgs.length === 0;

    if (useConfigMode) {
      // Config mode: load commands from config file
      const actualConfigPath = configPath || getDefaultConfigPath();

      if (debug) {
        console.error(`[Studio MCP] Loading config from: ${actualConfigPath}`);
      }

      let config: StudioConfig;
      try {
        config = loadConfig(actualConfigPath);
      } catch (err: any) {
        if (!configPath && commandsFilter.length === 0 && commandArgs.length === 0) {
          // No config file found and no explicit config requested
          console.error(
            'usage: studio <command> --example "{{req # required arg}}" "[args... # array of args]"',
          );
          console.error(`\nOr create a config file at: ${actualConfigPath}`);
          process.exit(1);
        }
        throw err;
      }

      // Filter commands if specified
      if (commandsFilter.length > 0) {
        config = filterCommands(config, commandsFilter);
      }

      // Register all commands from config
      for (const [name, cmdConfig] of Object.entries(config)) {
        registerConfigCommand(server, name, cmdConfig, debug);
      }

      if (debug) {
        console.error(
          `[Studio MCP] Registered ${Object.keys(config).length} command(s) from config`,
        );
      }
    } else {
      // Template mode: single command from args
      const template = fromArgs(commandArgs);
      const inputSchema = generateInputSchema(template);
      const zodSchema = schemaToZod(inputSchema);

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
            const output = await execute(fullCommand[0], fullCommand.slice(1));

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
    }

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
