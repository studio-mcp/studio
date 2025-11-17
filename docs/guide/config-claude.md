# Configuring Claude Desktop

Claude Desktop supports MCP servers through a JSON configuration file. This guide shows you how to add Studio MCP tools to Claude Desktop.

## Configuration File Location

The configuration file location varies by platform:

**macOS:**
```
~/Library/Application Support/Claude/claude_desktop_config.json
```

**Windows:**
```
%APPDATA%\Claude\claude_desktop_config.json
```

## Basic Configuration

Create or edit the configuration file with this structure:

```json
{
  "mcpServers": {
    "tool-name": {
      "command": "npx",
      "args": [
        "-y",
        "@studio-mcp/studio",
        "your-command",
        "{parameters}"
      ]
    }
  }
}
```

## Simple Example

Add an echo tool:

```json
{
  "mcpServers": {
    "echo": {
      "command": "npx",
      "args": [
        "-y",
        "@studio-mcp/studio",
        "echo",
        "{message # The text to echo}"
      ]
    }
  }
}
```

After saving, restart Claude Desktop to load the new configuration.

## Multiple Tools

Add multiple tools to the same configuration:

```json
{
  "mcpServers": {
    "echo": {
      "command": "npx",
      "args": [
        "-y",
        "@studio-mcp/studio",
        "echo",
        "{text}"
      ]
    },
    "speak": {
      "command": "npx",
      "args": [
        "-y",
        "@studio-mcp/studio",
        "say",
        "-v",
        "Samantha",
        "{speech # A concise message to say out loud}"
      ]
    },
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

## Common Examples

### Text-to-Speech (macOS)

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

**Usage:** "Use the speak tool to say hello"

### Curl with API Endpoints

```json
{
  "mcpServers": {
    "api-get": {
      "command": "npx",
      "args": [
        "-y",
        "@studio-mcp/studio",
        "curl",
        "-s",
        "https://api.github.com/{endpoint}"
      ]
    }
  }
}
```

**Usage:** "Use api-get to fetch users/octocat"

### Git Commands

```json
{
  "mcpServers": {
    "git": {
      "command": "npx",
      "args": [
        "-y",
        "@studio-mcp/studio",
        "git",
        "[args...]"
      ]
    }
  }
}
```

**Usage:** "Use git tool with args: status, --short"

### Date and Time

```json
{
  "mcpServers": {
    "date": {
      "command": "npx",
      "args": [
        "-y",
        "@studio-mcp/studio",
        "date",
        "[format]"
      ]
    }
  }
}
```

**Usage:** "What's the current date?" or "Use date with format: +%Y-%m-%d"

## Using Global Installation

If you've installed Studio globally, you can use it directly:

```json
{
  "mcpServers": {
    "echo": {
      "command": "studio",
      "args": [
        "echo",
        "{message}"
      ]
    }
  }
}
```

**Note:** This requires `npm install -g @studio-mcp/studio` first.

## Debugging Configuration

Add debug flags to troubleshoot issues:

```json
{
  "mcpServers": {
    "echo": {
      "command": "npx",
      "args": [
        "-y",
        "@studio-mcp/studio",
        "--debug",
        "--log",
        "/tmp/studio-echo.log",
        "echo",
        "{message}"
      ]
    }
  }
}
```

Check the log file at `/tmp/studio-echo.log` for debug output.

## Environment Variables

Set environment variables for your tools:

```json
{
  "mcpServers": {
    "api-call": {
      "command": "npx",
      "args": [
        "-y",
        "@studio-mcp/studio",
        "curl",
        "-H",
        "Authorization: Bearer ${API_TOKEN}",
        "https://api.example.com/{endpoint}"
      ],
      "env": {
        "API_TOKEN": "your-token-here"
      }
    }
  }
}
```

**Security Note:** Never commit API tokens to version control. Consider using a separate file or environment variable manager.

## Testing Configuration

After editing the configuration:

1. **Save the file** - Ensure JSON is valid (no trailing commas!)
2. **Restart Claude Desktop** - Completely quit and reopen
3. **Check the Tools** - Look for the hammer icon in Claude Desktop
4. **Test a tool** - Ask Claude to use your tool

## Troubleshooting

### Tool not appearing

1. Check JSON syntax - Use a JSON validator
2. Ensure Claude Desktop was fully restarted
3. Check the configuration file path is correct
4. Verify npx works: `npx -y @studio-mcp/studio --version`

### Tool fails to execute

1. Add `--debug` and `--log` flags to see errors
2. Test the command directly in terminal
3. Check file permissions and paths
4. Verify the command is available on your system

### Invalid JSON errors

Common mistakes:
- Trailing commas (not allowed in JSON)
- Missing quotes around strings
- Unescaped backslashes in Windows paths

Use a JSON validator or editor with JSON support.

## Complete Example

Here's a comprehensive configuration with multiple useful tools:

```json
{
  "mcpServers": {
    "echo": {
      "command": "npx",
      "args": ["-y", "@studio-mcp/studio", "echo", "{text}"]
    },
    "speak": {
      "command": "npx",
      "args": [
        "-y",
        "@studio-mcp/studio",
        "say",
        "-v",
        "Samantha",
        "{speech # Short phrase to speak}"
      ]
    },
    "weather": {
      "command": "npx",
      "args": [
        "-y",
        "@studio-mcp/studio",
        "curl",
        "-s",
        "wttr.in/{city}?format=3"
      ]
    },
    "git-status": {
      "command": "npx",
      "args": [
        "-y",
        "@studio-mcp/studio",
        "git",
        "-C",
        "{repo_path}",
        "status",
        "[--short]"
      ]
    }
  }
}
```

## Next Steps

- [Template Syntax](/reference/template-syntax) - Learn advanced template features
- [Examples](/examples/basic) - More tool examples
- [Debugging](/guide/debugging) - Troubleshooting guide
