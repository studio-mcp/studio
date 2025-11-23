# Common Tools

This guide shows how to integrate popular development tools with Studio MCP.

## Git

### Git Status

Check repository status.

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
        "status",
        "[--short]"
      ]
    }
  }
}
```

**Usage:** "Check git status of /path/to/repo with short: true"

### Git Log

View commit history.

```json
{
  "mcpServers": {
    "git-log": {
      "command": "npx",
      "args": [
        "-y",
        "@studio-mcp/studio",
        "git",
        "-C",
        "{repo_path}",
        "log",
        "--oneline",
        "-n",
        "{count # Number of commits}"
      ]
    }
  }
}
```

**Usage:** "Show last 10 commits from /path/to/repo"

### Git Diff

Show changes in repository.

```json
{
  "mcpServers": {
    "git-diff": {
      "command": "npx",
      "args": [
        "-y",
        "@studio-mcp/studio",
        "git",
        "-C",
        "{repo_path}",
        "diff",
        "[file]"
      ]
    }
  }
}
```

**Usage:** "Show diff for file src/main.ts in /path/to/repo"

### Git Clone

Clone a repository.

```json
{
  "mcpServers": {
    "git-clone": {
      "command": "npx",
      "args": [
        "-y",
        "@studio-mcp/studio",
        "git",
        "clone",
        "{url}",
        "[directory]"
      ]
    }
  }
}
```

**Usage:** "Clone https://github.com/user/repo.git"

### Git Branch

List or manage branches.

```json
{
  "mcpServers": {
    "git-branch": {
      "command": "npx",
      "args": [
        "-y",
        "@studio-mcp/studio",
        "git",
        "-C",
        "{repo_path}",
        "branch",
        "[--all]"
      ]
    }
  }
}
```

**Usage:** "List all branches in /path/to/repo"

### Flexible Git

General-purpose git tool with variable arguments.

```json
{
  "mcpServers": {
    "git": {
      "command": "npx",
      "args": [
        "-y",
        "@studio-mcp/studio",
        "git",
        "[args... # Git command and arguments]"
      ]
    }
  }
}
```

**Usage:** "Run git with args: status, --short"

## npm

### npm Install

Install packages.

```json
{
  "mcpServers": {
    "npm-install": {
      "command": "npx",
      "args": [
        "-y",
        "@studio-mcp/studio",
        "npm",
        "install",
        "{package}"
      ]
    }
  }
}
```

**Usage:** "Install package lodash"

### npm Search

Search for packages.

```json
{
  "mcpServers": {
    "npm-search": {
      "command": "npx",
      "args": [
        "-y",
        "@studio-mcp/studio",
        "npm",
        "search",
        "{query}"
      ]
    }
  }
}
```

**Usage:** "Search npm for react packages"

### npm Run Script

Run package.json scripts.

```json
{
  "mcpServers": {
    "npm-run": {
      "command": "npx",
      "args": [
        "-y",
        "@studio-mcp/studio",
        "npm",
        "run",
        "{script # Script name from package.json}"
      ]
    }
  }
}
```

**Usage:** "Run script: build"

### npm Info

Get package information.

```json
{
  "mcpServers": {
    "npm-info": {
      "command": "npx",
      "args": [
        "-y",
        "@studio-mcp/studio",
        "npm",
        "info",
        "{package}",
        "[field]"
      ]
    }
  }
}
```

**Usage:** "Get info for package react, field: version"

### npm List

List installed packages.

```json
{
  "mcpServers": {
    "npm-list": {
      "command": "npx",
      "args": [
        "-y",
        "@studio-mcp/studio",
        "npm",
        "list",
        "[--depth {depth}]"
      ]
    }
  }
}
```

**Usage:** "List packages with depth: 0"

## Docker

### Docker PS

List containers.

```json
{
  "mcpServers": {
    "docker-ps": {
      "command": "npx",
      "args": [
        "-y",
        "@studio-mcp/studio",
        "docker",
        "ps",
        "[--all]"
      ]
    }
  }
}
```

**Usage:** "List all docker containers"

### Docker Images

List images.

```json
{
  "mcpServers": {
    "docker-images": {
      "command": "npx",
      "args": [
        "-y",
        "@studio-mcp/studio",
        "docker",
        "images"
      ]
    }
  }
}
```

**Usage:** "Show docker images"

### Docker Logs

View container logs.

```json
{
  "mcpServers": {
    "docker-logs": {
      "command": "npx",
      "args": [
        "-y",
        "@studio-mcp/studio",
        "docker",
        "logs",
        "[--tail {lines}]",
        "{container}"
      ]
    }
  }
}
```

**Usage:** "Show logs for container my-app with lines: 100"

### Docker Exec

Execute command in container.

```json
{
  "mcpServers": {
    "docker-exec": {
      "command": "npx",
      "args": [
        "-y",
        "@studio-mcp/studio",
        "docker",
        "exec",
        "{container}",
        "{command...}"
      ]
    }
  }
}
```

**Usage:** "Execute in container my-app: ls, -la, /app"

## curl

### Basic GET Request

```json
{
  "mcpServers": {
    "curl-get": {
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

**Usage:** "Fetch https://api.github.com/users/octocat"

### POST Request

```json
{
  "mcpServers": {
    "curl-post": {
      "command": "npx",
      "args": [
        "-y",
        "@studio-mcp/studio",
        "curl",
        "-X",
        "POST",
        "-H",
        "Content-Type: application/json",
        "-d",
        "{data}",
        "{url}"
      ]
    }
  }
}
```

**Usage:** "Post to https://api.example.com/data with data: {\"key\":\"value\"}"

### Custom Headers

```json
{
  "mcpServers": {
    "curl-auth": {
      "command": "npx",
      "args": [
        "-y",
        "@studio-mcp/studio",
        "curl",
        "-H",
        "Authorization: Bearer {token}",
        "{url}"
      ]
    }
  }
}
```

**Usage:** "Fetch https://api.example.com/protected with token: abc123"

### Download File

```json
{
  "mcpServers": {
    "curl-download": {
      "command": "npx",
      "args": [
        "-y",
        "@studio-mcp/studio",
        "curl",
        "-o",
        "{output_file}",
        "{url}"
      ]
    }
  }
}
```

**Usage:** "Download https://example.com/file.zip to /tmp/file.zip"

## Python

### Run Python Script

```json
{
  "mcpServers": {
    "python-run": {
      "command": "npx",
      "args": [
        "-y",
        "@studio-mcp/studio",
        "python3",
        "{script_path}",
        "[args...]"
      ]
    }
  }
}
```

**Usage:** "Run script.py with args: --verbose, input.txt"

### Python One-Liner

```json
{
  "mcpServers": {
    "python-eval": {
      "command": "npx",
      "args": [
        "-y",
        "@studio-mcp/studio",
        "python3",
        "-c",
        "{code}"
      ]
    }
  }
}
```

**Usage:** "Execute python code: print(2 + 2)"

### pip Install

```json
{
  "mcpServers": {
    "pip-install": {
      "command": "npx",
      "args": [
        "-y",
        "@studio-mcp/studio",
        "pip3",
        "install",
        "{package}"
      ]
    }
  }
}
```

**Usage:** "Install package requests"

## Grep

### Search Files

```json
{
  "mcpServers": {
    "grep": {
      "command": "npx",
      "args": [
        "-y",
        "@studio-mcp/studio",
        "grep",
        "[--ignore-case]",
        "{pattern}",
        "{files...}"
      ]
    }
  }
}
```

**Usage:** "Search for 'TODO' in files: src/main.ts, src/utils.ts"

### Recursive Search

```json
{
  "mcpServers": {
    "grep-recursive": {
      "command": "npx",
      "args": [
        "-y",
        "@studio-mcp/studio",
        "grep",
        "-r",
        "[--ignore-case]",
        "{pattern}",
        "{directory}"
      ]
    }
  }
}
```

**Usage:** "Recursively search for 'function' in directory: src/"

## Tar

### Create Archive

```json
{
  "mcpServers": {
    "tar-create": {
      "command": "npx",
      "args": [
        "-y",
        "@studio-mcp/studio",
        "tar",
        "-czf",
        "{archive_name}",
        "{files...}"
      ]
    }
  }
}
```

**Usage:** "Create archive backup.tar.gz with files: file1.txt, file2.txt"

### Extract Archive

```json
{
  "mcpServers": {
    "tar-extract": {
      "command": "npx",
      "args": [
        "-y",
        "@studio-mcp/studio",
        "tar",
        "-xzf",
        "{archive}",
        "[--directory {dir}]"
      ]
    }
  }
}
```

**Usage:** "Extract backup.tar.gz to directory: /tmp/restore"

### List Archive Contents

```json
{
  "mcpServers": {
    "tar-list": {
      "command": "npx",
      "args": [
        "-y",
        "@studio-mcp/studio",
        "tar",
        "-tzf",
        "{archive}"
      ]
    }
  }
}
```

**Usage:** "List contents of backup.tar.gz"

## jq

### Parse JSON

```json
{
  "mcpServers": {
    "jq-parse": {
      "command": "npx",
      "args": [
        "-y",
        "@studio-mcp/studio",
        "jq",
        "{filter}",
        "{file}"
      ]
    }
  }
}
```

**Usage:** "Parse data.json with filter: '.users[0].name'"

### Format JSON

```json
{
  "mcpServers": {
    "jq-format": {
      "command": "npx",
      "args": [
        "-y",
        "@studio-mcp/studio",
        "jq",
        ".",
        "{file}"
      ]
    }
  }
}
```

**Usage:** "Format file data.json"

## sed

### Replace Text

```json
{
  "mcpServers": {
    "sed-replace": {
      "command": "npx",
      "args": [
        "-y",
        "@studio-mcp/studio",
        "sed",
        "s/{pattern}/{replacement}/g",
        "{file}"
      ]
    }
  }
}
```

**Usage:** "Replace 'foo' with 'bar' in file.txt"

## awk

### Extract Columns

```json
{
  "mcpServers": {
    "awk-column": {
      "command": "npx",
      "args": [
        "-y",
        "@studio-mcp/studio",
        "awk",
        "{'{print $'}{column}{'}'}",
        "{file}"
      ]
    }
  }
}
```

**Usage:** "Extract column 2 from data.txt"

## Complete Development Toolkit

A comprehensive configuration with common tools:

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
        "{repo}",
        "status",
        "[--short]"
      ]
    },
    "npm-run": {
      "command": "npx",
      "args": [
        "-y",
        "@studio-mcp/studio",
        "npm",
        "run",
        "{script}"
      ]
    },
    "docker-ps": {
      "command": "npx",
      "args": [
        "-y",
        "@studio-mcp/studio",
        "docker",
        "ps",
        "[--all]"
      ]
    },
    "curl-api": {
      "command": "npx",
      "args": [
        "-y",
        "@studio-mcp/studio",
        "curl",
        "-s",
        "{url}"
      ]
    },
    "grep-search": {
      "command": "npx",
      "args": [
        "-y",
        "@studio-mcp/studio",
        "grep",
        "-r",
        "{pattern}",
        "{directory}"
      ]
    }
  }
}
```

## Next Steps

- [Advanced Examples](/examples/advanced) - Complex use cases
- [Template Syntax](/reference/template-syntax) - Full syntax reference
- [Debugging](/guide/debugging) - Troubleshooting tools
