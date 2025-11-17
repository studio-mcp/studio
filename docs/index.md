---
layout: home

hero:
  name: Studio MCP
  text: Turn any CLI into an MCP tool
  tagline: The simplest way to add command-line tools to AI assistants
  actions:
    - theme: brand
      text: Get Started
      link: /guide/getting-started
    - theme: alt
      text: View on GitHub
      link: https://github.com/studio-mcp/studio

features:
  - icon: 🚀
    title: Simple & Lightweight
    details: Convert any CLI command into an MCP tool with a single command. No complex setup or configuration required.

  - icon: 🎯
    title: Template-Based
    details: Use intuitive Mustache-like syntax to define command arguments. Studio automatically generates the schema.

  - icon: 🔌
    title: Universal Compatibility
    details: Works with Claude Desktop, Cursor, VSCode, and any MCP-compatible client.

  - icon: ⚡
    title: Zero Dependencies
    details: Use npx to run studio without installation. No global dependencies or package management hassles.

  - icon: 🎨
    title: Flexible Arguments
    details: Support for required/optional parameters, arrays, boolean flags, and complex templates.

  - icon: 🔍
    title: Built-in Debugging
    details: Debug mode with stderr or file logging helps troubleshoot tool execution.
---

## Quick Start

Install globally:

```bash
npm install -g @studio-mcp/studio
```

Or use with npx (no installation):

```bash
npx -y @studio-mcp/studio echo "{message}"
```

## Example Usage

Add a simple echo tool to Claude Desktop:

```json
{
  "mcpServers": {
    "echo": {
      "command": "npx",
      "args": ["-y", "@studio-mcp/studio", "echo", "{text}"]
    }
  }
}
```

Now Claude can use the echo tool with any text input!

## Why Studio MCP?

Traditional MCP servers require:
- Writing TypeScript/JavaScript code
- Setting up build tools
- Managing dependencies
- Publishing to npm

With Studio MCP, you just:
1. Choose a CLI command
2. Add template syntax for arguments
3. Add to your MCP client config

That's it! Studio handles schema generation, validation, and execution automatically.
