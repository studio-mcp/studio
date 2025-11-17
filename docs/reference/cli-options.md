# Command Line Options

Studio MCP can be run directly from the command line with various options. This reference covers all available flags and arguments.

## Basic Usage

```bash
studio [options] <command> [template-args...]
```

## Synopsis

```bash
studio [--help] [--version] [--debug] [--log <file>] [--] <command> [args...]
```

## Options

### `--help`, `-h`

Display help information and exit.

```bash
studio --help
studio -h
```

**Output:**
```
Studio MCP - Turn any CLI command into an MCP tool

Usage: studio [options] <command> [template-args...]

Options:
  -h, --help           Show this help message
  --version            Show version information
  --debug              Print debug logs to stderr
  --log <filename>     Write debug logs to file
  --                   Stop parsing flags

Examples:
  studio echo "{message}"
  studio --debug git status
  studio --log /tmp/studio.log curl "{url}"
```

### `--version`

Display version information and exit.

```bash
studio --version
```

**Output:**
```
0.0.0
```

### `--debug`

Enable debug logging to stderr.

```bash
studio --debug echo "{message}"
```

**Debug output includes:**
- Template parsing details
- Generated JSON schema
- Parameter validation
- Command execution details
- Exit codes and errors

**Example output:**
```
[DEBUG] Parsing template: echo "{message}"
[DEBUG] Generated schema: {"type":"object","properties":...}
[DEBUG] Executing command: ["echo", "Hello"]
[DEBUG] Exit code: 0
```

**Note:** Debug output goes to stderr, MCP protocol goes to stdout.

### `--log <filename>`

Write debug logs to a specified file instead of stderr.

```bash
studio --log /tmp/studio.log echo "{message}"
```

**Log file format:**
Same as `--debug`, but written to file:
```
2025-01-15T12:34:56.789Z [DEBUG] Parsing template...
2025-01-15T12:34:56.790Z [DEBUG] Generated schema...
```

**Use cases:**
- Debugging MCP client integrations
- Persistent logging
- Log analysis
- Troubleshooting production issues

### `--`

Stop parsing flags and treat remaining arguments as command.

```bash
studio -- command --with-its-own-flags "{arg}"
```

**Useful when:**
- Command has flags that conflict with Studio's flags
- Command starts with `--`
- You want to be explicit about where Studio options end

**Example:**
```bash
# Without --, Studio tries to parse --verbose
studio echo --verbose "{message}"  # Error

# With --, --verbose is passed to echo
studio -- echo --verbose "{message}"  # Works
```

## Combining Options

Options can be combined:

```bash
# Debug to both stderr and file
studio --debug --log /tmp/studio.log echo "{message}"

# Multiple options before command
studio --debug --log /tmp/test.log -- git "{args...}"
```

## Command Arguments

Everything after options is the command template:

```bash
studio [options] <command> [template-args...]
```

### Command

The actual command to execute:

```bash
studio echo "{message}"
#      ^^^^ command
```

Can be:
- System commands: `echo`, `curl`, `git`
- Scripts: `/path/to/script.sh`
- Executables: `/usr/bin/python3`

### Template Arguments

Arguments with template syntax:

```bash
studio curl "https://api.example.com/{endpoint}"
#           ^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^ template arg
```

See [Template Syntax](/reference/template-syntax) for details.

## Examples

### Basic Echo

```bash
studio echo "{message}"
```

### Debug Mode

```bash
studio --debug curl "https://api.github.com/{endpoint}"
```

### Log to File

```bash
studio --log /tmp/studio.log git status
```

### Complex Command with Flags

```bash
studio -- git -C "{repo_path}" log --oneline "{branch}"
```

### Multiple Templates

```bash
studio curl -X "{method}" "https://api.example.com/{endpoint}"
```

## MCP Client Usage

In MCP client configurations, options go in the `args` array:

### Claude Desktop

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

### Cursor

```json
{
  "mcpServers": {
    "git": {
      "command": "npx",
      "args": [
        "-y",
        "@studio-mcp/studio",
        "--log",
        "/tmp/cursor-git.log",
        "git",
        "[args...]"
      ]
    }
  }
}
```

### VSCode

```json
{
  "mcp.servers": {
    "curl": {
      "type": "stdio",
      "command": "npx",
      "args": [
        "-y",
        "@studio-mcp/studio",
        "--debug",
        "curl",
        "{url}"
      ]
    }
  }
}
```

## Environment Variables

Studio doesn't use environment variables for configuration, but respects:

### `PATH`

Used to find commands:

```bash
# These use PATH
studio echo "{message}"
studio git "{args...}"
studio curl "{url}"
```

### `HOME`

Used for `~` expansion in file paths:

```bash
studio --log ~/studio.log echo "{message}"
```

## Exit Codes

Studio returns exit codes to indicate success or failure:

| Code | Meaning |
|------|---------|
| 0 | Success |
| 1 | General error |
| 2 | Invalid arguments |
| 126 | Command not executable |
| 127 | Command not found |
| 130 | Terminated by Ctrl+C |

**Examples:**

```bash
# Success
studio echo "Hello"
echo $?  # 0

# Command not found
studio nonexistent-command
echo $?  # 127

# Invalid template
studio echo "{"
echo $?  # 2
```

## Standard Streams

### stdin

MCP protocol input (JSON-RPC messages).

```bash
# MCP client sends requests to stdin
echo '{"jsonrpc":"2.0","method":"tools/list"}' | studio echo "{message}"
```

### stdout

MCP protocol output (JSON-RPC responses).

```bash
# Studio sends responses to stdout
studio echo "{message}" > mcp-output.json
```

### stderr

Debug logs and errors.

```bash
# Debug output goes to stderr
studio --debug echo "{message}" 2> debug.log
```

## Debugging Options

### Full Debug Output

```bash
studio --debug echo "{message}" 2>&1 | tee full-output.log
```

This captures:
- Debug logs (stderr)
- MCP protocol (stdout)

### Separate Streams

```bash
studio --debug echo "{message}" \
  1> mcp-protocol.log \
  2> debug.log
```

This separates:
- MCP protocol → `mcp-protocol.log`
- Debug output → `debug.log`

### Both Debug and Log File

```bash
studio --debug --log /tmp/studio.log echo "{message}"
```

This outputs debug to:
- stderr (terminal)
- `/tmp/studio.log` (file)

## Shell Integration

### Bash/Zsh

Add to `.bashrc` or `.zshrc`:

```bash
# Alias for common commands
alias studio-debug='studio --debug --log /tmp/studio.log'

# Function wrapper
studio-test() {
  studio --debug "$@" 2>&1 | less
}
```

### Fish

Add to `~/.config/fish/config.fish`:

```fish
# Alias
alias studio-debug='studio --debug --log /tmp/studio.log'

# Function
function studio-test
    studio --debug $argv 2>&1 | less
end
```

## Testing Commands

### Test Template Parsing

```bash
# See generated schema
studio --debug echo "{message}" 2>&1 | grep "Generated schema"
```

### Test Command Execution

```bash
# See what command runs
studio --debug curl "{url}" 2>&1 | grep "Executing"
```

### Test with MCP Inspector

```bash
mcp-inspector studio --debug --log /tmp/inspector.log echo "{message}"
```

## Common Patterns

### Development Testing

```bash
studio --debug --log ./studio-dev.log echo "{message}"
```

### Production Use

```bash
# No debug output
studio curl "https://api.example.com/{endpoint}"
```

### Troubleshooting

```bash
studio --debug --log /tmp/troubleshoot.log -- problematic-command "{arg}"
```

### Quick Prototyping

```bash
# Test different templates quickly
studio echo "{a} {b} {c}"
studio echo "[a] [b] [c]"
studio echo "[args...]"
```

## Next Steps

- [Template Syntax](/reference/template-syntax) - Learn template features
- [Debugging](/guide/debugging) - Troubleshooting guide
- [Examples](/examples/basic) - Real-world usage
