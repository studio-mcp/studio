import { describe, expect, test } from "vitest";
import { loadConfig, filterCommands, getDefaultConfigPath } from "./config.js";
import * as fs from "fs";
import * as os from "os";
import * as path from "path";

describe("config", () => {
  describe("getDefaultConfigPath", () => {
    test("returns XDG_CONFIG_HOME path when set", () => {
      const originalXdg = process.env.XDG_CONFIG_HOME;
      process.env.XDG_CONFIG_HOME = "/custom/config";
      expect(getDefaultConfigPath()).toBe("/custom/config/studio/config.jsonc");
      if (originalXdg) {
        process.env.XDG_CONFIG_HOME = originalXdg;
      } else {
        delete process.env.XDG_CONFIG_HOME;
      }
    });

    test("returns default path when XDG_CONFIG_HOME not set", () => {
      const originalXdg = process.env.XDG_CONFIG_HOME;
      delete process.env.XDG_CONFIG_HOME;
      expect(getDefaultConfigPath()).toBe(
        path.join(os.homedir(), ".config", "studio", "config.jsonc"),
      );
      if (originalXdg) {
        process.env.XDG_CONFIG_HOME = originalXdg;
      }
    });
  });

  describe("loadConfig", () => {
    test("loads valid JSON config", () => {
      const tempFile = `/tmp/test-config-${Date.now()}.json`;
      fs.writeFileSync(
        tempFile,
        JSON.stringify({
          echo: {
            description: "Echo command",
            command: ["echo", "hello"],
          },
        }),
      );

      const config = loadConfig(tempFile);
      expect(config).toEqual({
        echo: {
          description: "Echo command",
          command: ["echo", "hello"],
        },
      });

      fs.unlinkSync(tempFile);
    });

    test("loads valid JSONC with comments", () => {
      const tempFile = `/tmp/test-config-${Date.now()}.jsonc`;
      fs.writeFileSync(
        tempFile,
        `{
          // This is a comment
          "echo": {
            "description": "Echo command", // inline comment
            "command": ["echo", "hello"]
          }
          /* Multi-line
             comment */
        }`,
      );

      const config = loadConfig(tempFile);
      expect(config).toEqual({
        echo: {
          description: "Echo command",
          command: ["echo", "hello"],
        },
      });

      fs.unlinkSync(tempFile);
    });

    test("throws error for non-existent file", () => {
      expect(() => loadConfig("/nonexistent/file.jsonc")).toThrow(
        "Config file not found",
      );
    });

    test("throws error for invalid JSON", () => {
      const tempFile = `/tmp/test-config-${Date.now()}.jsonc`;
      fs.writeFileSync(tempFile, "{ invalid json }");

      expect(() => loadConfig(tempFile)).toThrow("Failed to parse config file");

      fs.unlinkSync(tempFile);
    });

    test("throws error when config is not an object", () => {
      const tempFile = `/tmp/test-config-${Date.now()}.jsonc`;
      fs.writeFileSync(tempFile, "[]");

      expect(() => loadConfig(tempFile)).toThrow("Config must be an object");

      fs.unlinkSync(tempFile);
    });

    test("throws error when command missing command field", () => {
      const tempFile = `/tmp/test-config-${Date.now()}.jsonc`;
      fs.writeFileSync(
        tempFile,
        JSON.stringify({
          echo: {
            description: "Echo command",
          },
        }),
      );

      expect(() => loadConfig(tempFile)).toThrow(
        'Command "echo" is missing required "command" field',
      );

      fs.unlinkSync(tempFile);
    });

    test("throws error when command field is not an array", () => {
      const tempFile = `/tmp/test-config-${Date.now()}.jsonc`;
      fs.writeFileSync(
        tempFile,
        JSON.stringify({
          echo: {
            command: "echo hello",
          },
        }),
      );

      expect(() => loadConfig(tempFile)).toThrow(
        'Command "echo": "command" field must be an array',
      );

      fs.unlinkSync(tempFile);
    });

    test("throws error when command array is empty", () => {
      const tempFile = `/tmp/test-config-${Date.now()}.jsonc`;
      fs.writeFileSync(
        tempFile,
        JSON.stringify({
          echo: {
            command: [],
          },
        }),
      );

      expect(() => loadConfig(tempFile)).toThrow(
        'Command "echo": "command" array cannot be empty',
      );

      fs.unlinkSync(tempFile);
    });

    test("throws error when command array contains non-strings", () => {
      const tempFile = `/tmp/test-config-${Date.now()}.jsonc`;
      fs.writeFileSync(
        tempFile,
        JSON.stringify({
          echo: {
            command: ["echo", 123],
          },
        }),
      );

      expect(() => loadConfig(tempFile)).toThrow(
        'Command "echo": all command elements must be strings',
      );

      fs.unlinkSync(tempFile);
    });

    test("allows optional description", () => {
      const tempFile = `/tmp/test-config-${Date.now()}.jsonc`;
      fs.writeFileSync(
        tempFile,
        JSON.stringify({
          echo: {
            command: ["echo", "hello"],
          },
        }),
      );

      const config = loadConfig(tempFile);
      expect(config).toEqual({
        echo: {
          command: ["echo", "hello"],
        },
      });

      fs.unlinkSync(tempFile);
    });

    test("throws error when description is not a string", () => {
      const tempFile = `/tmp/test-config-${Date.now()}.jsonc`;
      fs.writeFileSync(
        tempFile,
        JSON.stringify({
          echo: {
            description: 123,
            command: ["echo", "hello"],
          },
        }),
      );

      expect(() => loadConfig(tempFile)).toThrow(
        'Command "echo": "description" must be a string',
      );

      fs.unlinkSync(tempFile);
    });
  });

  describe("filterCommands", () => {
    test("filters to specified commands", () => {
      const config = {
        echo: { command: ["echo"] },
        date: { command: ["date"] },
        pwd: { command: ["pwd"] },
      };

      const filtered = filterCommands(config, ["echo", "pwd"]);
      expect(filtered).toEqual({
        echo: { command: ["echo"] },
        pwd: { command: ["pwd"] },
      });
    });

    test("throws error when command not found", () => {
      const config = {
        echo: { command: ["echo"] },
      };

      expect(() => filterCommands(config, ["nonexistent"])).toThrow(
        'Command "nonexistent" not found in config',
      );
    });

    test("returns empty object when filtering to empty list", () => {
      const config = {
        echo: { command: ["echo"] },
      };

      const filtered = filterCommands(config, []);
      expect(filtered).toEqual({});
    });
  });
});
