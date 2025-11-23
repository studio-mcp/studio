# Configuring VSCode

Visual Studio Code supports MCP through extensions and configuration. This guide shows you how to add Studio MCP tools to VSCode.

## Prerequisites

VSCode MCP support requires:
- VSCode version 1.85 or later
- An MCP-compatible extension (varies by implementation)

Check with your specific MCP extension documentation for exact requirements.

## Configuration File Location

VSCode MCP configuration is typically in:

**All Platforms:**
```
.vscode/settings.json (workspace)
```

Or global settings:
```
~/Library/Application Support/Code/User/settings.json (macOS)
%APPDATA%\Code\User\settings.json (Windows)
~/.config/Code/User/settings.json (Linux)
```

## Basic Configuration

The exact configuration format depends on your MCP extension. Here's a common pattern:

```json
{
  "mcp.servers": {
    "tool-name": {
      "type": "stdio",
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
  "mcp.servers": {
    "echo": {
      "type": "stdio",
      "command": "npx",
      "args": [
        "-y",
        "@studio-mcp/studio",
        "echo",
        "{message # What to echo}"
      ]
    }
  }
}
```

## Common Examples

### Echo Tool

```json
{
  "mcp.servers": {
    "echo": {
      "type": "stdio",
      "command": "npx",
      "args": [
        "-y",
        "@studio-mcp/studio",
        "echo",
        "{text # Text to echo}"
      ]
    }
  }
}
```

### Git Commands

```json
{
  "mcp.servers": {
    "git": {
      "type": "stdio",
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

### File Operations

```json
{
  "mcp.servers": {
    "cat": {
      "type": "stdio",
      "command": "npx",
      "args": [
        "-y",
        "@studio-mcp/studio",
        "cat",
        "{file_path # Path to file}"
      ]
    }
  }
}
```

### Curl for API Calls

```json
{
  "mcp.servers": {
    "api-request": {
      "type": "stdio",
      "command": "npx",
      "args": [
        "-y",
        "@studio-mcp/studio",
        "curl",
        "-s",
        "https://api.example.com/{endpoint}"
      ]
    }
  }
}
```

### NPM Package Search

```json
{
  "mcp.servers": {
    "npm-search": {
      "type": "stdio",
      "command": "npx",
      "args": [
        "-y",
        "@studio-mcp/studio",
        "npm",
        "search",
        "{package}"
      ]
    }
  }
}
```

## Workspace Configuration

For project-specific tools, create `.vscode/settings.json` in your project:

```json
{
  "mcp.servers": {
    "project-build": {
      "type": "stdio",
      "command": "npx",
      "args": [
        "-y",
        "@studio-mcp/studio",
        "npm",
        "run",
        "{script # npm script name}"
      ]
    },
    "test": {
      "type": "stdio",
      "command": "npx",
      "args": [
        "-y",
        "@studio-mcp/studio",
        "npm",
        "test",
        "[file # test file path]"
      ]
    }
  }
}
```

## Using Global Installation

If Studio is installed globally:

```json
{
  "mcp.servers": {
    "echo": {
      "type": "stdio",
      "command": "studio",
      "args": ["echo", "{message}"]
    }
  }
}
```

Requires: `npm install -g @studio-mcp/studio`

## Environment Variables

Set environment variables for your tools:

```json
{
  "mcp.servers": {
    "api-call": {
      "type": "stdio",
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

## Debugging Configuration

Enable debug logging:

```json
{
  "mcp.servers": {
    "echo": {
      "type": "stdio",
      "command": "npx",
      "args": [
        "-y",
        "@studio-mcp/studio",
        "--debug",
        "--log",
        "${workspaceFolder}/studio.log",
        "echo",
        "{message}"
      ]
    }
  }
}
```

Use VSCode variables:
- `${workspaceFolder}` - Project root
- `${file}` - Current file path
- `${fileBasename}` - Current file name
- `${fileDirname}` - Current file directory

## Multiple Tools Configuration

A comprehensive VSCode setup:

```json
{
  "mcp.servers": {
    "echo": {
      "type": "stdio",
      "command": "npx",
      "args": ["-y", "@studio-mcp/studio", "echo", "{text}"]
    },
    "git": {
      "type": "stdio",
      "command": "npx",
      "args": ["-y", "@studio-mcp/studio", "git", "[args...]"]
    },
    "curl": {
      "type": "stdio",
      "command": "npx",
      "args": [
        "-y",
        "@studio-mcp/studio",
        "curl",
        "-s",
        "{url}"
      ]
    },
    "npm-run": {
      "type": "stdio",
      "command": "npx",
      "args": [
        "-y",
        "@studio-mcp/studio",
        "npm",
        "run",
        "{script}"
      ]
    }
  }
}
```

## Reloading Configuration

After editing settings:

1. **Save settings.json**
2. **Reload window** - `Cmd+Shift+P` (Mac) or `Ctrl+Shift+P` (Windows/Linux)
3. **Search "Developer: Reload Window"**
4. **Test the tool** - Try using your MCP tool

## Troubleshooting

### Tool not appearing

1. Check the MCP extension is installed and enabled
2. Verify JSON syntax (no trailing commas)
3. Reload VSCode window
4. Check Output panel for errors

### Command fails

1. Test command in integrated terminal
2. Check file paths are correct
3. Verify npx is available: `npx --version`
4. Use absolute paths if needed

### Permission errors

1. Ensure execute permissions on scripts
2. On macOS/Linux, check file permissions
3. Use full paths to executables

## VSCode-Specific Tips

### Using Workspace Variables

```json
{
  "mcp.servers": {
    "lint-current-file": {
      "type": "stdio",
      "command": "npx",
      "args": [
        "-y",
        "@studio-mcp/studio",
        "eslint",
        "${file}"
      ]
    }
  }
}
```

### Project-Specific Tools

```json
{
  "mcp.servers": {
    "build": {
      "type": "stdio",
      "command": "npx",
      "args": [
        "-y",
        "@studio-mcp/studio",
        "npm",
        "run",
        "build"
      ],
      "cwd": "${workspaceFolder}"
    }
  }
}
```

### Conditional Configuration

Use VSCode's settings precedence:
1. Workspace settings (`.vscode/settings.json`)
2. User settings (global)

This allows project-specific overrides.

## Complete Example

A full VSCode workspace configuration:

```json
{
  "mcp.servers": {
    "echo": {
      "type": "stdio",
      "command": "npx",
      "args": ["-y", "@studio-mcp/studio", "echo", "{text}"]
    },
    "git-status": {
      "type": "stdio",
      "command": "npx",
      "args": [
        "-y",
        "@studio-mcp/studio",
        "git",
        "-C",
        "${workspaceFolder}",
        "status"
      ]
    },
    "build": {
      "type": "stdio",
      "command": "npx",
      "args": [
        "-y",
        "@studio-mcp/studio",
        "npm",
        "run",
        "build"
      ]
    },
    "test": {
      "type": "stdio",
      "command": "npx",
      "args": [
        "-y",
        "@studio-mcp/studio",
        "npm",
        "test",
        "[file]"
      ]
    },
    "curl": {
      "type": "stdio",
      "command": "npx",
      "args": [
        "-y",
        "@studio-mcp/studio",
        "curl",
        "-s",
        "{url}"
      ]
    }
  }
}
```

## Next Steps

- [Template Syntax](/reference/template-syntax) - Learn advanced template features
- [Examples](/examples/basic) - More tool examples
- [Debugging](/guide/debugging) - Troubleshooting guide
