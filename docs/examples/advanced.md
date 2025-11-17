# Advanced Examples

This guide covers advanced use cases and complex tool configurations.

## Multiple Templates in URLs

Combine multiple parameters in a single URL.

### GitHub API with Version and Endpoint

```json
{
  "mcpServers": {
    "github-api": {
      "command": "npx",
      "args": [
        "-y",
        "@studio-mcp/studio",
        "curl",
        "-s",
        "-H",
        "Accept: application/vnd.github.{api_version}+json",
        "https://api.github.com/{endpoint}"
      ]
    }
  }
}
```

**Usage:** "Fetch users/octocat with api_version: v3"

### RESTful API with Dynamic Paths

```json
{
  "mcpServers": {
    "rest-api": {
      "command": "npx",
      "args": [
        "-y",
        "@studio-mcp/studio",
        "curl",
        "-s",
        "-X",
        "{method}",
        "https://api.example.com/{version}/{resource}/{id}"
      ]
    }
  }
}
```

**Usage:** "Call GET on v1/users/123"

## Environment-Based Configuration

Use environment variables for sensitive data.

### API with Token

```json
{
  "mcpServers": {
    "secure-api": {
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

**Security Note:** Store tokens in environment files, not in the config itself.

### Database Query with Credentials

```json
{
  "mcpServers": {
    "db-query": {
      "command": "npx",
      "args": [
        "-y",
        "@studio-mcp/studio",
        "psql",
        "-h",
        "${DB_HOST}",
        "-U",
        "${DB_USER}",
        "-d",
        "${DB_NAME}",
        "-c",
        "{query}"
      ],
      "env": {
        "DB_HOST": "localhost",
        "DB_USER": "admin",
        "DB_NAME": "mydb",
        "PGPASSWORD": "secret"
      }
    }
  }
}
```

## Complex Git Operations

Advanced git workflows.

### Git Commit with Message

```json
{
  "mcpServers": {
    "git-commit": {
      "command": "npx",
      "args": [
        "-y",
        "@studio-mcp/studio",
        "git",
        "-C",
        "{repo}",
        "commit",
        "[--all]",
        "-m",
        "{message}"
      ]
    }
  }
}
```

**Usage:** "Commit changes in /path/to/repo with message: 'Fix bug'"

### Git Push with Branch

```json
{
  "mcpServers": {
    "git-push": {
      "command": "npx",
      "args": [
        "-y",
        "@studio-mcp/studio",
        "git",
        "-C",
        "{repo}",
        "push",
        "[--force]",
        "{remote}",
        "{branch}"
      ]
    }
  }
}
```

**Usage:** "Push repo /path/to/repo to origin main"

### Git Cherry-pick

```json
{
  "mcpServers": {
    "git-cherry-pick": {
      "command": "npx",
      "args": [
        "-y",
        "@studio-mcp/studio",
        "git",
        "-C",
        "{repo}",
        "cherry-pick",
        "{commits...}"
      ]
    }
  }
}
```

**Usage:** "Cherry-pick commits: abc123, def456 in /path/to/repo"

## Docker Compose

Manage multi-container applications.

### Docker Compose Up

```json
{
  "mcpServers": {
    "compose-up": {
      "command": "npx",
      "args": [
        "-y",
        "@studio-mcp/studio",
        "docker-compose",
        "-f",
        "{compose_file}",
        "up",
        "[--detach]",
        "[services...]"
      ]
    }
  }
}
```

**Usage:** "Start services web, db from docker-compose.yml with detach: true"

### Docker Compose Logs

```json
{
  "mcpServers": {
    "compose-logs": {
      "command": "npx",
      "args": [
        "-y",
        "@studio-mcp/studio",
        "docker-compose",
        "-f",
        "{compose_file}",
        "logs",
        "[--tail {lines}]",
        "[--follow]",
        "[service]"
      ]
    }
  }
}
```

**Usage:** "Show logs for service web with lines: 100"

## Data Processing Pipelines

Complex command chains.

### JSON Processing

```json
{
  "mcpServers": {
    "process-json": {
      "command": "npx",
      "args": [
        "-y",
        "@studio-mcp/studio",
        "sh",
        "-c",
        "curl -s {url} | jq '{filter}'"
      ]
    }
  }
}
```

**Usage:** "Process JSON from https://api.example.com/data with filter: '.results[]'"

### Log Analysis

```json
{
  "mcpServers": {
    "analyze-logs": {
      "command": "npx",
      "args": [
        "-y",
        "@studio-mcp/studio",
        "sh",
        "-c",
        "grep '{pattern}' {log_file} | wc -l"
      ]
    }
  }
}
```

**Usage:** "Count errors in /var/log/app.log"

## File Synchronization

### Rsync with Options

```json
{
  "mcpServers": {
    "rsync": {
      "command": "npx",
      "args": [
        "-y",
        "@studio-mcp/studio",
        "rsync",
        "-avz",
        "[--delete]",
        "[--dry-run]",
        "{source}",
        "{destination}"
      ]
    }
  }
}
```

**Usage:** "Sync /source/ to /dest/ with delete: true, dry_run: true"

### SCP File Transfer

```json
{
  "mcpServers": {
    "scp": {
      "command": "npx",
      "args": [
        "-y",
        "@studio-mcp/studio",
        "scp",
        "[--recursive]",
        "{source}",
        "{user}@{host}:{destination}"
      ]
    }
  }
}
```

**Usage:** "Copy file.txt to user@server.com:/tmp/"

## System Monitoring

### Disk Usage

```json
{
  "mcpServers": {
    "disk-usage": {
      "command": "npx",
      "args": [
        "-y",
        "@studio-mcp/studio",
        "du",
        "-sh",
        "{path}"
      ]
    }
  }
}
```

**Usage:** "Check disk usage of /var/log"

### Process Monitoring

```json
{
  "mcpServers": {
    "ps-grep": {
      "command": "npx",
      "args": [
        "-y",
        "@studio-mcp/studio",
        "sh",
        "-c",
        "ps aux | grep '{process}'"
      ]
    }
  }
}
```

**Usage:** "Find processes matching 'node'"

### Network Connections

```json
{
  "mcpServers": {
    "netstat": {
      "command": "npx",
      "args": [
        "-y",
        "@studio-mcp/studio",
        "netstat",
        "-an",
        "|",
        "grep",
        "{port}"
      ]
    }
  }
}
```

**Usage:** "Check connections on port 8080"

## Build and Deploy

### Maven Build

```json
{
  "mcpServers": {
    "maven": {
      "command": "npx",
      "args": [
        "-y",
        "@studio-mcp/studio",
        "mvn",
        "{goals...}",
        "[--file {pom}]"
      ]
    }
  }
}
```

**Usage:** "Run maven with goals: clean, install"

### Gradle Build

```json
{
  "mcpServers": {
    "gradle": {
      "command": "npx",
      "args": [
        "-y",
        "@studio-mcp/studio",
        "gradle",
        "{tasks...}"
      ]
    }
  }
}
```

**Usage:** "Run gradle with tasks: clean, build"

### Make

```json
{
  "mcpServers": {
    "make": {
      "command": "npx",
      "args": [
        "-y",
        "@studio-mcp/studio",
        "make",
        "[--directory {dir}]",
        "{target}"
      ]
    }
  }
}
```

**Usage:** "Make target all in directory /path/to/project"

## Testing Tools

### Jest with Coverage

```json
{
  "mcpServers": {
    "jest": {
      "command": "npx",
      "args": [
        "-y",
        "@studio-mcp/studio",
        "jest",
        "[--coverage]",
        "[--watch]",
        "[pattern]"
      ]
    }
  }
}
```

**Usage:** "Run jest with coverage: true, pattern: user.test.ts"

### Pytest

```json
{
  "mcpServers": {
    "pytest": {
      "command": "npx",
      "args": [
        "-y",
        "@studio-mcp/studio",
        "pytest",
        "[--verbose]",
        "[--coverage]",
        "{tests...}"
      ]
    }
  }
}
```

**Usage:** "Run pytest on tests: test_user.py, test_auth.py with verbose: true"

## Code Quality

### ESLint

```json
{
  "mcpServers": {
    "eslint": {
      "command": "npx",
      "args": [
        "-y",
        "@studio-mcp/studio",
        "eslint",
        "[--fix]",
        "{files...}"
      ]
    }
  }
}
```

**Usage:** "Lint src/main.ts with fix: true"

### Prettier

```json
{
  "mcpServers": {
    "prettier": {
      "command": "npx",
      "args": [
        "-y",
        "@studio-mcp/studio",
        "prettier",
        "[--write]",
        "{files...}"
      ]
    }
  }
}
```

**Usage:** "Format files: src/*.ts with write: true"

### TypeScript Compiler

```json
{
  "mcpServers": {
    "tsc": {
      "command": "npx",
      "args": [
        "-y",
        "@studio-mcp/studio",
        "tsc",
        "[--noEmit]",
        "[--project {tsconfig}]"
      ]
    }
  }
}
```

**Usage:** "Type-check with noEmit: true"

## Database Operations

### PostgreSQL Query

```json
{
  "mcpServers": {
    "psql": {
      "command": "npx",
      "args": [
        "-y",
        "@studio-mcp/studio",
        "psql",
        "-h",
        "{host}",
        "-d",
        "{database}",
        "-c",
        "{query}"
      ]
    }
  }
}
```

**Usage:** "Query database mydb on localhost: SELECT * FROM users"

### MySQL Query

```json
{
  "mcpServers": {
    "mysql": {
      "command": "npx",
      "args": [
        "-y",
        "@studio-mcp/studio",
        "mysql",
        "-h",
        "{host}",
        "-D",
        "{database}",
        "-e",
        "{query}"
      ]
    }
  }
}
```

**Usage:** "Query database mydb: SELECT COUNT(*) FROM users"

### Redis Commands

```json
{
  "mcpServers": {
    "redis": {
      "command": "npx",
      "args": [
        "-y",
        "@studio-mcp/studio",
        "redis-cli",
        "-h",
        "{host}",
        "{command...}"
      ]
    }
  }
}
```

**Usage:** "Run redis command GET user:123 on localhost"

## Cloud CLIs

### AWS CLI

```json
{
  "mcpServers": {
    "aws": {
      "command": "npx",
      "args": [
        "-y",
        "@studio-mcp/studio",
        "aws",
        "{service}",
        "{command}",
        "[args...]"
      ]
    }
  }
}
```

**Usage:** "Run AWS s3 ls with args: s3://my-bucket"

### Google Cloud

```json
{
  "mcpServers": {
    "gcloud": {
      "command": "npx",
      "args": [
        "-y",
        "@studio-mcp/studio",
        "gcloud",
        "{service}",
        "{command}",
        "[args...]"
      ]
    }
  }
}
```

**Usage:** "Run gcloud compute instances list"

### Kubernetes

```json
{
  "mcpServers": {
    "kubectl": {
      "command": "npx",
      "args": [
        "-y",
        "@studio-mcp/studio",
        "kubectl",
        "{command}",
        "[args...]"
      ]
    }
  }
}
```

**Usage:** "Run kubectl get pods with args: -n, production"

## Custom Scripts

### Shell Script Wrapper

```json
{
  "mcpServers": {
    "deploy": {
      "command": "npx",
      "args": [
        "-y",
        "@studio-mcp/studio",
        "/path/to/deploy.sh",
        "{environment}",
        "[--dry-run]"
      ]
    }
  }
}
```

**Usage:** "Deploy to production with dry_run: true"

### Python Script

```json
{
  "mcpServers": {
    "analyze": {
      "command": "npx",
      "args": [
        "-y",
        "@studio-mcp/studio",
        "python3",
        "/path/to/analyze.py",
        "--input",
        "{input_file}",
        "--output",
        "{output_file}",
        "[--verbose]"
      ]
    }
  }
}
```

**Usage:** "Analyze data.csv output to results.json with verbose: true"

## Working Directory Context

### Execute in Specific Directory

```json
{
  "mcpServers": {
    "build-project": {
      "command": "npx",
      "args": [
        "-y",
        "@studio-mcp/studio",
        "sh",
        "-c",
        "cd {project_dir} && npm run build"
      ]
    }
  }
}
```

**Usage:** "Build project in /path/to/project"

## Debugging Complex Tools

### Full Debug Configuration

```json
{
  "mcpServers": {
    "debug-tool": {
      "command": "npx",
      "args": [
        "-y",
        "@studio-mcp/studio",
        "--debug",
        "--log",
        "/tmp/studio-debug.log",
        "your-command",
        "{args...}"
      ]
    }
  }
}
```

Check `/tmp/studio-debug.log` for detailed execution logs.

## Best Practices

### 1. Parameter Validation

Let commands handle validation:

```json
{
  "command": "npx",
  "args": [
    "-y",
    "@studio-mcp/studio",
    "curl",
    "--fail",
    "{url}"
  ]
}
```

The `--fail` flag makes curl return errors for HTTP errors.

### 2. Error Handling

Use command flags for better errors:

```json
{
  "command": "npx",
  "args": [
    "-y",
    "@studio-mcp/studio",
    "git",
    "-C",
    "{repo}",
    "status",
    "--porcelain"
  ]
}
```

`--porcelain` gives machine-readable output.

### 3. Idempotent Operations

Prefer commands that can run multiple times:

```json
{
  "command": "npx",
  "args": [
    "-y",
    "@studio-mcp/studio",
    "mkdir",
    "-p",
    "{directory}"
  ]
}
```

`-p` won't error if directory exists.

### 4. Descriptive Tool Names

Use clear, descriptive names:

```json
// Good
"git-commit-all"
"docker-compose-up"
"npm-build-production"

// Less clear
"gc"
"dup"
"nb"
```

### 5. Document Parameters

Always add descriptions:

```json
{
  "args": [
    "-y",
    "@studio-mcp/studio",
    "command",
    "{param # What this parameter does}"
  ]
}
```

## Troubleshooting Advanced Tools

### Shell Escaping Issues

When using `sh -c`, escape carefully:

```json
{
  "args": [
    "-y",
    "@studio-mcp/studio",
    "sh",
    "-c",
    "echo '{message}' | tr a-z A-Z"
  ]
}
```

### Complex Command Chains

Test each part separately first:

```bash
# Test individually
curl -s https://api.example.com/data
jq '.results[]' data.json

# Then combine
sh -c "curl -s https://api.example.com/data | jq '.results[]'"
```

### Debugging Pipeline Failures

Add debug at each stage:

```json
{
  "args": [
    "-y",
    "@studio-mcp/studio",
    "sh",
    "-c",
    "set -x && curl -s {url} | jq '{filter}'"
  ]
}
```

`set -x` shows each command as it executes.

## Next Steps

- [Template Syntax](/reference/template-syntax) - Master all syntax features
- [Debugging](/guide/debugging) - Troubleshoot complex issues
- [Architecture](/reference/architecture) - Understand internals
