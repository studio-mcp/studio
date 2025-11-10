import { describe, it, expect, beforeAll } from "vitest";
import { execFile, exec } from "child_process";
import { promisify } from "util";
import * as path from "path";
import * as fs from "fs/promises";
import { fileURLToPath } from "url";

const execFileAsync = promisify(execFile);
const execAsync = promisify(exec);

const __filename = fileURLToPath(import.meta.url);
const __dirname = path.dirname(__filename);

// Type definitions for MCP Inspector responses
interface InspectorToolsListResponse {
  tools: InspectorTool[];
}

interface InspectorTool {
  name: string;
  description: string;
  inputSchema: InspectorInputSchema;
}

interface InspectorInputSchema {
  type: string;
  properties?: Record<string, InspectorProperty>;
  required?: string[];
}

interface InspectorProperty {
  type: string;
  description?: string;
  items?: InspectorPropertyItem;
}

interface InspectorPropertyItem {
  type: string;
}

interface InspectorToolCallResponse {
  content: InspectorContent[];
  isError?: boolean;
}

interface InspectorContent {
  type: string;
  text: string;
}

const TIMEOUT = 10000;

describe("MCP Inspector Smoke Test", () => {
  let binaryPath: string;

  beforeAll(async () => {
    const projectRoot = path.join(__dirname, "..");
    const distDir = path.join(projectRoot, "dist");
    binaryPath = path.join(distDir, "index.js");

    // Ensure dist directory exists
    await fs.mkdir(distDir, { recursive: true });

    // Build the project
    await execAsync("bun run build", { cwd: projectRoot });

    // Verify binary exists
    await fs.access(binaryPath);
  }, 30000);

  async function runInspectorCmd(
    args: string[],
  ): Promise<{ stdout: string; stderr: string }> {
    const fullArgs = ["--yes", ...args];
    const projectRoot = path.join(__dirname, "..");

    try {
      const result = await execFileAsync("npx", fullArgs, {
        cwd: projectRoot,
        timeout: TIMEOUT,
        maxBuffer: 1024 * 1024 * 10, // 10MB buffer
      });

      return {
        stdout: result.stdout || "",
        stderr: result.stderr || "",
      };
    } catch (error: any) {
      // execFile throws on non-zero exit, but we want to return the output
      return {
        stdout: error.stdout || "",
        stderr: error.stderr || "",
      };
    }
  }

  describe("BasicInspectorConnection", () => {
    it("can list simple echo tools", async () => {
      const args = [
        "@modelcontextprotocol/inspector",
        "--cli",
        binaryPath,
        "echo",
        "hello",
        "--method",
        "tools/list",
      ];

      const { stdout, stderr } = await runInspectorCmd(args);

      // Parse JSON response
      const response: InspectorToolsListResponse = JSON.parse(stdout);

      // Validate response structure
      expect(response.tools).toHaveLength(1);

      const tool = response.tools[0];
      expect(tool.name).toBe("echo");
      expect(tool.description).toBe("Run the shell command `echo hello`");
      expect(tool.inputSchema.type).toBe("object");

      // Simple commands without [args...] don't have properties
      // or have an empty properties object
      if (tool.inputSchema.properties) {
        expect(Object.keys(tool.inputSchema.properties)).toHaveLength(0);
      }
    }, 10000);

    it("can list echo tools with args", async () => {
      const args = [
        "@modelcontextprotocol/inspector",
        "--cli",
        binaryPath,
        "echo",
        "hello",
        "[args...]",
        "--method",
        "tools/list",
      ];

      const { stdout, stderr } = await runInspectorCmd(args);

      // Parse JSON response
      const response: InspectorToolsListResponse = JSON.parse(stdout);

      // Validate response structure
      expect(response.tools).toHaveLength(1);

      const tool = response.tools[0];
      expect(tool.name).toBe("echo");
      expect(tool.description).toBe(
        "Run the shell command `echo hello [args...]`",
      );
      expect(tool.inputSchema.type).toBe("object");

      // Should have args property for commands with [args...]
      expect(tool.inputSchema.properties).toBeDefined();
      expect(tool.inputSchema.properties!.args).toBeDefined();
      const argsProperty = tool.inputSchema.properties!.args;
      expect(argsProperty.type).toBe("array");
      expect(argsProperty.description).toBe(
        "Additional command line arguments",
      );
    });

    it("can call simple echo tool", async () => {
      const args = [
        "@modelcontextprotocol/inspector",
        "--cli",
        binaryPath,
        "echo",
        "hello",
        "--method",
        "tools/call",
        "--tool-name",
        "echo",
      ];

      const { stdout, stderr } = await runInspectorCmd(args);

      // Parse JSON response
      const response: InspectorToolCallResponse = JSON.parse(stdout);

      // Validate response structure
      expect(response.content).toHaveLength(1);

      const content = response.content[0];
      expect(content.type).toBe("text");
      expect(content.text).toBe("hello");
      expect(response.isError).toBeFalsy();
    });
  });

  describe("Command", () => {
    it("can list templated tools with proper schema", async () => {
      const args = [
        "@modelcontextprotocol/inspector",
        "--cli",
        binaryPath,
        "echo",
        "{{text#message to echo}}",
        "--method",
        "tools/list",
      ];

      const { stdout, stderr } = await runInspectorCmd(args);

      // Parse JSON response
      const response: InspectorToolsListResponse = JSON.parse(stdout);

      // Validate response structure
      expect(response.tools).toHaveLength(1);

      const tool = response.tools[0];
      expect(tool.name).toBe("echo");
      expect(tool.description).toBe("Run the shell command `echo {{text}}`");
      expect(tool.inputSchema.type).toBe("object");

      // Should have text property for tepmlated commands
      expect(tool.inputSchema.properties).toBeDefined();
      expect(tool.inputSchema.properties!.text).toBeDefined();
      const textProperty = tool.inputSchema.properties!.text;
      expect(textProperty.type).toBe("string");
      expect(textProperty.description).toBe("message to echo");

      // Should have text as required parameter
      expect(tool.inputSchema.required).toContain("text");

      // Should NOT have args property for templated commands
      expect(tool.inputSchema.properties!.args).toBeUndefined();
    });

    it("can call templated tool with arguments", async () => {
      const args = [
        "@modelcontextprotocol/inspector",
        "--cli",
        binaryPath,
        "echo",
        "{{text#message to echo}}",
        "--method",
        "tools/call",
        "--tool-name",
        "echo",
        "--tool-arg",
        "text=Hello from inspector!",
      ];

      const { stdout, stderr } = await runInspectorCmd(args);

      // Parse JSON response
      const response: InspectorToolCallResponse = JSON.parse(stdout);

      // Validate response structure
      expect(response.content).toHaveLength(1);

      const content = response.content[0];
      expect(content.type).toBe("text");
      expect(content.text).toBe("Hello from inspector!");
      expect(response.isError).toBeFalsy();
    });

    it("can handle mixed template with additional args", async () => {
      const args = [
        "@modelcontextprotocol/inspector",
        "--cli",
        binaryPath,
        "echo",
        "{{text#message to echo}}",
        "[args...]",
        "--method",
        "tools/list",
      ];

      const { stdout, stderr } = await runInspectorCmd(args);

      // Parse JSON response
      const response: InspectorToolsListResponse = JSON.parse(stdout);

      // Validate response structure
      expect(response.tools).toHaveLength(1);

      const tool = response.tools[0];
      expect(tool.name).toBe("echo");
      expect(tool.description).toBe(
        "Run the shell command `echo {{text}} [args...]`",
      );

      // Should have both text and args properties
      expect(tool.inputSchema.properties).toBeDefined();
      expect(tool.inputSchema.properties!.text).toBeDefined();
      expect(tool.inputSchema.properties!.args).toBeDefined();

      const textProperty = tool.inputSchema.properties!.text;
      expect(textProperty.type).toBe("string");
      expect(textProperty.description).toBe("message to echo");

      const argsProperty = tool.inputSchema.properties!.args;
      expect(argsProperty.type).toBe("array");

      // Only text should be required
      expect(tool.inputSchema.required).toContain("text");
      expect(tool.inputSchema.required).not.toContain("args");
    });
  });

  describe("ComplexCommand", () => {
    it("can handle git status command", async () => {
      const args = [
        "@modelcontextprotocol/inspector",
        "--cli",
        binaryPath,
        "git",
        "status",
        "--method",
        "tools/list",
      ];

      const { stdout, stderr } = await runInspectorCmd(args);

      // Parse JSON response
      const response: InspectorToolsListResponse = JSON.parse(stdout);

      // Validate response structure
      expect(response.tools).toHaveLength(1);

      const tool = response.tools[0];
      expect(tool.name).toBe("git");
      expect(tool.description).toBe("Run the shell command `git status`");
      expect(tool.inputSchema.type).toBe("object");

      // Simple git status command without [args...] doesn't have properties
      // or has an empty properties object
      if (tool.inputSchema.properties) {
        expect(Object.keys(tool.inputSchema.properties)).toHaveLength(0);
      }
    });

    it("can handle git status command with args", async () => {
      const args = [
        "@modelcontextprotocol/inspector",
        "--cli",
        binaryPath,
        "git",
        "status",
        "[args...]",
        "--method",
        "tools/list",
      ];

      const { stdout, stderr } = await runInspectorCmd(args);

      // Parse JSON response
      const response: InspectorToolsListResponse = JSON.parse(stdout);

      // Validate response structure
      expect(response.tools).toHaveLength(1);

      const tool = response.tools[0];
      expect(tool.name).toBe("git");
      expect(tool.description).toBe(
        "Run the shell command `git status [args...]`",
      );
      expect(tool.inputSchema.type).toBe("object");

      // Should have args property
      expect(tool.inputSchema.properties).toBeDefined();
      expect(tool.inputSchema.properties!.args).toBeDefined();
      const argsProperty = tool.inputSchema.properties!.args;
      expect(argsProperty.type).toBe("array");
      expect(argsProperty.description).toBe(
        "Additional command line arguments",
      );
    });

    it("can execute git status", async () => {
      const args = [
        "@modelcontextprotocol/inspector",
        "--cli",
        binaryPath,
        "git",
        "status",
        "--method",
        "tools/call",
        "--tool-name",
        "git",
      ];

      const { stdout, stderr } = await runInspectorCmd(args);

      // Parse JSON response
      const response: InspectorToolCallResponse = JSON.parse(stdout);

      // Validate response structure
      expect(response.content).toHaveLength(1);

      const content = response.content[0];
      expect(content.type).toBe("text");
      // Git status output varies, but should contain some recognizable text
      expect(
        content.text.includes("On branch") ||
          content.text.includes("HEAD detached") ||
          content.text.includes("nothing to commit"),
      ).toBe(true);
      expect(response.isError).toBeFalsy();
    });
  });

  describe("SayCommandRegression", () => {
    it("can handle say command with -v flag", async () => {
      const args = [
        "@modelcontextprotocol/inspector",
        "--cli",
        binaryPath,
        "say",
        "-v",
        "siri",
        "{{speech#message to speak}}",
        "--method",
        "tools/list",
      ];

      const { stdout, stderr } = await runInspectorCmd(args);

      // Parse JSON response
      const response: InspectorToolsListResponse = JSON.parse(stdout);

      // Validate response structure
      expect(response.tools).toHaveLength(1);

      const tool = response.tools[0];
      expect(tool.name).toBe("say");
      expect(tool.description).toBe(
        "Run the shell command `say -v siri {{speech}}`",
      );
      expect(tool.inputSchema.type).toBe("object");

      // Should have speech property
      expect(tool.inputSchema.properties).toBeDefined();
      expect(tool.inputSchema.properties!.speech).toBeDefined();
      const speechProperty = tool.inputSchema.properties!.speech;
      expect(speechProperty.type).toBe("string");
      expect(speechProperty.description).toBe("message to speak");

      // Should have speech as required parameter
      expect(tool.inputSchema.required).toContain("speech");
    });
  });

  describe("Multipletemplates", () => {
    it("can handle multiple template parameters", async () => {
      const args = [
        "@modelcontextprotocol/inspector",
        "--cli",
        binaryPath,
        "rails",
        "generate",
        "{{generator#Rails generator name}}",
        "{{name#Resource name}}",
        "--method",
        "tools/list",
      ];

      const { stdout, stderr } = await runInspectorCmd(args);

      // Parse JSON response
      const response: InspectorToolsListResponse = JSON.parse(stdout);

      // Validate response structure
      expect(response.tools).toHaveLength(1);

      const tool = response.tools[0];
      expect(tool.name).toBe("rails");
      expect(tool.description).toBe(
        "Run the shell command `rails generate {{generator}} {{name}}`",
      );
      expect(tool.inputSchema.type).toBe("object");

      // Should have both generator and name properties
      expect(tool.inputSchema.properties).toBeDefined();
      expect(tool.inputSchema.properties!.generator).toBeDefined();
      expect(tool.inputSchema.properties!.name).toBeDefined();

      const generatorProperty = tool.inputSchema.properties!.generator;
      expect(generatorProperty.type).toBe("string");
      expect(generatorProperty.description).toBe("Rails generator name");

      const nameProperty = tool.inputSchema.properties!.name;
      expect(nameProperty.type).toBe("string");
      expect(nameProperty.description).toBe("Resource name");

      // Both should be required
      expect(tool.inputSchema.required).toContain("generator");
      expect(tool.inputSchema.required).toContain("name");

      // Should NOT have args property
      expect(tool.inputSchema.properties!.args).toBeUndefined();
    });
  });
});
