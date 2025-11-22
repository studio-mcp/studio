import * as fs from "fs";
import * as path from "path";
import * as os from "os";

/**
 * Config file schema
 */
export interface CommandConfig {
  description?: string;
  command: string[];
}

export interface StudioConfig {
  [commandName: string]: CommandConfig;
}

/**
 * Get XDG config directory
 */
function getXdgConfigHome(): string {
  return process.env.XDG_CONFIG_HOME || path.join(os.homedir(), ".config");
}

/**
 * Get default config file path
 */
export function getDefaultConfigPath(): string {
  return path.join(getXdgConfigHome(), "studio", "config.jsonc");
}

/**
 * Strip comments from JSONC content
 */
function stripJsonComments(content: string): string {
  // Remove single-line comments
  content = content.replace(/\/\/.*$/gm, "");

  // Remove multi-line comments
  content = content.replace(/\/\*[\s\S]*?\*\//g, "");

  return content;
}

/**
 * Load and parse a config file
 */
export function loadConfig(configPath: string): StudioConfig {
  if (!fs.existsSync(configPath)) {
    throw new Error(`Config file not found: ${configPath}`);
  }

  const content = fs.readFileSync(configPath, "utf-8");
  const jsonContent = stripJsonComments(content);

  let config: any;
  try {
    config = JSON.parse(jsonContent);
  } catch (err: any) {
    throw new Error(`Failed to parse config file: ${err.message}`);
  }

  // Validate config structure
  if (typeof config !== "object" || config === null || Array.isArray(config)) {
    throw new Error("Config must be an object");
  }

  for (const [name, cmd] of Object.entries(config)) {
    if (typeof cmd !== "object" || cmd === null) {
      throw new Error(`Command "${name}" must be an object`);
    }

    const cmdConfig = cmd as any;

    if (!cmdConfig.command) {
      throw new Error(`Command "${name}" is missing required "command" field`);
    }

    if (!Array.isArray(cmdConfig.command)) {
      throw new Error(`Command "${name}": "command" field must be an array`);
    }

    if (cmdConfig.command.length === 0) {
      throw new Error(`Command "${name}": "command" array cannot be empty`);
    }

    for (const arg of cmdConfig.command) {
      if (typeof arg !== "string") {
        throw new Error(`Command "${name}": all command elements must be strings`);
      }
    }

    if (cmdConfig.description !== undefined && typeof cmdConfig.description !== "string") {
      throw new Error(`Command "${name}": "description" must be a string`);
    }
  }

  return config as StudioConfig;
}

/**
 * Filter config to only include specified commands
 */
export function filterCommands(config: StudioConfig, commandNames: string[]): StudioConfig {
  const filtered: StudioConfig = {};

  for (const name of commandNames) {
    if (!(name in config)) {
      throw new Error(`Command "${name}" not found in config`);
    }
    filtered[name] = config[name];
  }

  return filtered;
}
