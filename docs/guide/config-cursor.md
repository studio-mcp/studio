# Configuring Cursor

Cursor is an AI-powered code editor with MCP support. This guide shows you how to add Studio MCP tools to Cursor.

## Configuration File Location

Cursor uses a different configuration file than Claude Desktop:

**All Platforms:**
```
~/.cursor/mcp.json
```

**Windows:**
```
%USERPROFILE%\.cursor\mcp.json
```

## Important: Cursor Limitations

Cursor currently has a bug that prevents spaces in argument descriptions. You must:

1. **Use underscores instead of spaces** in descriptions
2. **Avoid spaces in the template syntax**

This is a known issue with Cursor's MCP implementation.

## Basic Configuration

Create or edit `~/.cursor/mcp.json`:

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

Add an echo tool (note: no spaces in description):

```json
{
  "mcpServers": {
    "echo": {
      "command": "npx",
      "args": [
        "-y",
        "@studio-mcp/studio",
        "echo",
        "{message#text_to_echo}"
      ]
    }
  }
}
```

**Important:** Notice `{message#text_to_echo}` uses underscores, not spaces.

## Correct vs Incorrect

### ❌ Incorrect (will fail in Cursor):

```json
{
  "mcpServers": {
    "speak": {
      "command": "npx",
      "args": [
        "-y",
        "@studio-mcp/studio",
        "say",
        "{speech # A message to say}"
      ]
    }
  }
}
```

The spaces in `"# A message to say"` will cause errors.

### ✅ Correct (works in Cursor):

```json
{
  "mcpServers": {
    "speak": {
      "command": "npx",
      "args": [
        "-y",
        "@studio-mcp/studio",
        "say",
        "{speech#message_to_say}"
      ]
    }
  }
}
```

Use underscores or omit the description entirely.

## Common Examples

### Echo Tool

```json
{
  "mcpServers": {
    "echo": {
      "command": "npx",
      "args": [
        "-y",
        "@studio-mcp/studio",
        "echo",
        "{text#message_to_echo}"
      ]
    }
  }
}
```

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
        "{speech#text_to_speak}"
      ]
    }
  }
}
```

### Weather Check

```json
{
  "mcpServers": {
    "weather": {
      "command": "npx",
      "args": [
        "-y",
        "@studio-mcp/studio",
        "curl",
        "-s",
        "wttr.in/{city}?format=3"
      ]
    }
  }
}
```

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

### API Requests

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
        "https://api.github.com/{endpoint#api_path}"
      ]
    }
  }
}
```

## Multiple Tools

Configure multiple tools in one file:

```json
{
  "mcpServers": {
    "echo": {
      "command": "npx",
      "args": ["-y", "@studio-mcp/studio", "echo", "{text}"]
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
    "git": {
      "command": "npx",
      "args": ["-y", "@studio-mcp/studio", "git", "[args...]"]
    }
  }
}
```

## Using Global Installation

If you've installed Studio globally:

```json
{
  "mcpServers": {
    "echo": {
      "command": "studio",
      "args": ["echo", "{message}"]
    }
  }
}
```

Requires: `npm install -g @studio-mcp/studio`

## Debugging in Cursor

Add debug flags to troubleshoot:

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
        "/tmp/cursor-studio.log",
        "echo",
        "{message}"
      ]
    }
  }
}
```

Check `/tmp/cursor-studio.log` for debug output.

## Reloading Configuration

After editing `~/.cursor/mcp.json`:

1. **Save the file**
2. **Reload Cursor** - Use `Cmd+Shift+P` (Mac) or `Ctrl+Shift+P` (Windows/Linux)
3. **Search for "Reload Window"** and execute
4. **Test the tool** - Ask Cursor's AI to use your tool

## Troubleshooting

### Tool not appearing

1. Verify file location: `~/.cursor/mcp.json`
2. Check JSON syntax (no trailing commas)
3. Reload Cursor window
4. Check Cursor's developer console for errors

### "Invalid argument" errors

1. Remove spaces from descriptions
2. Use underscores: `{arg#like_this}` not `{arg # like this}`
3. Test the command in terminal first

### Command not found

1. Verify npx is available: `npx --version`
2. Test Studio: `npx -y @studio-mcp/studio --version`
3. Use full path to npx if needed

## Complete Example

A comprehensive Cursor configuration:

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
        "{speech#text_to_speak}"
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
    "git": {
      "command": "npx",
      "args": ["-y", "@studio-mcp/studio", "git", "[args...]"]
    },
    "npm-search": {
      "command": "npx",
      "args": [
        "-y",
        "@studio-mcp/studio",
        "npm",
        "search",
        "{package_name}"
      ]
    }
  }
}
```

## Cursor-Specific Tips

1. **Keep descriptions short** - Use single words or underscores
2. **Test commands first** - Verify in terminal before adding to config
3. **Use simple tools** - Complex commands may not work well in Cursor
4. **Check logs** - Use `--log` flag for debugging
5. **Reload frequently** - Cursor sometimes needs a full reload

## Next Steps

- [Template Syntax](/reference/template-syntax) - Learn template features
- [Examples](/examples/basic) - More tool examples
- [Debugging](/guide/debugging) - Troubleshooting guide
