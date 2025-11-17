# What is Studio MCP?

Studio MCP is a lightweight tool that converts any command-line interface (CLI) command into a Model Context Protocol (MCP) tool that can be used by AI assistants like Claude, Cursor, and other MCP-compatible clients.

## The Problem

Building MCP servers traditionally requires:

- Writing TypeScript or JavaScript code
- Setting up build systems and dependencies
- Creating JSON schemas manually
- Handling parameter validation
- Managing process execution
- Publishing and maintaining packages

This complexity creates a high barrier for adding simple CLI tools to AI assistants.

## The Solution

Studio MCP eliminates all this complexity. Instead of writing code, you simply:

1. **Choose a CLI command** - Any shell command you want to expose
2. **Add template syntax** - Use simple `{argument}` syntax for parameters
3. **Configure your client** - Add one JSON entry to your MCP config

Studio automatically:
- Parses your template
- Generates the JSON schema
- Validates parameters
- Executes the command
- Returns the output

## Example

Instead of writing a full MCP server to expose the `echo` command, you just write:

```bash
studio echo "{message}"
```

Studio converts this into an MCP tool with:
- Tool name: `echo`
- Parameter: `message` (string, required)
- Execution: Runs `echo` with the provided message

## Use Cases

Studio MCP is perfect for:

- **Quick prototyping** - Test MCP integration ideas without building servers
- **System utilities** - Expose `curl`, `git`, `grep`, and other CLI tools
- **Custom scripts** - Make your shell scripts AI-accessible
- **Text-to-speech** - Let AI assistants speak (`say` command)
- **API calls** - Wrap HTTP requests with templated URLs
- **File operations** - Safe, parameterized file manipulation

## How It Works

```
┌─────────────────┐
│  MCP Client     │ (Claude, Cursor, etc.)
│  (AI Assistant) │
└────────┬────────┘
         │ MCP Protocol
         │
┌────────▼────────┐
│  Studio MCP     │
│  - Parse template
│  - Generate schema
│  - Validate args
│  - Execute command
└────────┬────────┘
         │ spawn()
         │
┌────────▼────────┐
│  CLI Command    │ (echo, curl, git, etc.)
└─────────────────┘
```

## What Studio Is Not

Studio MCP is intentionally simple and focused. It is **not**:

- A replacement for complex MCP servers with custom logic
- A security sandbox (it runs commands with your permissions)
- A package manager or dependency resolver
- A programming language or scripting environment

For complex use cases requiring custom logic, state management, or advanced features, you should build a traditional MCP server. Studio is for when you just want to expose a CLI command quickly and easily.

## Next Steps

- [Getting Started](/guide/getting-started) - Install and run your first command
- [Template Syntax](/reference/template-syntax) - Learn the template language
- [Examples](/examples/basic) - See common usage patterns
