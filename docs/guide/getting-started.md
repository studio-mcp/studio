# Getting Started

This guide will walk you through installing Studio MCP and creating your first MCP tool.

## Prerequisites

Studio MCP requires Node.js 18 or later. Check your version:

```bash
node --version
```

## Installation

You have two options for using Studio MCP:

### Option 1: Use npx (Recommended)

No installation required! Use `npx` to run Studio on demand:

```bash
npx -y @studio-mcp/studio echo "{message}"
```

The `-y` flag automatically accepts the package installation.

### Option 2: Global Installation

Install Studio globally for faster execution:

```bash
npm install -g @studio-mcp/studio
```

Then use it directly:

```bash
studio echo "{message}"
```

## Your First Tool

Let's create a simple echo tool that AI assistants can use.

### Step 1: Test the Command

First, test Studio from the command line to verify it works:

```bash
npx -y @studio-mcp/studio echo "{message}"
```

This starts an MCP server that exposes the `echo` command with a `message` parameter.

### Step 2: Add to MCP Client

Now add it to your MCP client configuration. For Claude Desktop, edit:

**macOS**: `~/Library/Application Support/Claude/claude_desktop_config.json`
**Windows**: `%APPDATA%\Claude\claude_desktop_config.json`

Add this configuration:

```json
{
  "mcpServers": {
    "echo": {
      "command": "npx",
      "args": [
        "-y",
        "@studio-mcp/studio",
        "echo",
        "{message}"
      ]
    }
  }
}
```

### Step 3: Restart Claude Desktop

Close and reopen Claude Desktop to load the new configuration.

### Step 4: Use the Tool

In Claude Desktop, you can now say:

> "Use the echo tool to say 'Hello, World!'"

Claude will invoke the tool and show you the output.

## Understanding the Configuration

Let's break down what each part does:

```json
{
  "mcpServers": {
    "echo": {                    // Server name (arbitrary)
      "command": "npx",          // How to run Studio
      "args": [
        "-y",                    // Auto-accept npx install
        "@studio-mcp/studio",    // Package name
        "echo",                  // CLI command to wrap
        "{message}"              // Template with parameters
      ]
    }
  }
}
```

The template `"{message}"` tells Studio:
- Create a tool parameter called `message`
- The parameter is required (using `{}` not `[]`)
- The parameter is a string
- Pass it as an argument to `echo`

## More Examples

### Example 1: Text-to-Speech

Make Claude speak out loud (macOS):

```json
{
  "mcpServers": {
    "speak": {
      "command": "npx",
      "args": [
        "-y",
        "@studio-mcp/studio",
        "say",
        "-v",
        "Samantha",
        "{speech # A short phrase to speak out loud}"
      ]
    }
  }
}
```

Usage: "Use the speak tool to say hello"

### Example 2: Weather Check

Wrap curl to check weather:

```json
{
  "mcpServers": {
    "weather": {
      "command": "npx",
      "args": [
        "-y",
        "@studio-mcp/studio",
        "curl",
        "wttr.in/{city}?format=3"
      ]
    }
  }
}
```

Usage: "Check the weather in London"

### Example 3: Git Status

Check git repository status:

```json
{
  "mcpServers": {
    "git-status": {
      "command": "npx",
      "args": [
        "-y",
        "@studio-mcp/studio",
        "git",
        "-C",
        "{repo_path}",
        "status"
      ]
    }
  }
}
```

Usage: "Check the status of /path/to/repo"

## Command Line Usage

You can also run Studio directly for testing:

```bash
# Basic syntax
studio <command> [template args...]

# With debugging
studio --debug echo "{message}"

# Log to file
studio --log /tmp/studio.log curl "https://api.github.com/{endpoint}"

# Stop flag parsing with --
studio -- command --with-flags "{arg}"
```

## Next Steps

Now that you have Studio working, explore:

- [Template Syntax](/reference/template-syntax) - Learn all template features
- [Configuration Guides](/guide/config-claude) - Client-specific setup
- [Examples](/examples/basic) - More real-world examples
- [Debugging](/guide/debugging) - Troubleshoot issues
