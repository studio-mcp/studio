# Basic Examples

This guide provides simple, practical examples to get you started with Studio MCP.

## Echo Tool

The simplest possible tool - just echo text back.

### Configuration

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

### Usage

"Use the echo tool to say 'Hello, World!'"

### Output

```
Hello, World!
```

## Date and Time

Get the current date and time.

### Configuration

```json
{
  "mcpServers": {
    "date": {
      "command": "npx",
      "args": ["-y", "@studio-mcp/studio", "date"]
    }
  }
}
```

### Usage

"What's the current date and time?"

### Output

```
Mon Jan 15 14:23:45 PST 2025
```

## Weather Check

Check the weather using wttr.in.

### Configuration

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

### Usage

"Check the weather in London"

### Output

```
London: ☁️  +7°C
```

## Text-to-Speech (macOS)

Make your computer speak text aloud.

### Configuration

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

### Usage

"Use the speak tool to say 'Hello from Studio MCP'"

### Result

Your computer speaks the text using the Samantha voice.

## File Reading

Read the contents of a file.

### Configuration

```json
{
  "mcpServers": {
    "read-file": {
      "command": "npx",
      "args": [
        "-y",
        "@studio-mcp/studio",
        "cat",
        "{file_path # Path to the file to read}"
      ]
    }
  }
}
```

### Usage

"Read the file at /tmp/test.txt"

### Output

```
Contents of the file...
```

## Directory Listing

List files in a directory.

### Configuration

```json
{
  "mcpServers": {
    "list-files": {
      "command": "npx",
      "args": [
        "-y",
        "@studio-mcp/studio",
        "ls",
        "-la",
        "{directory # Directory path}"
      ]
    }
  }
}
```

### Usage

"List files in /tmp"

### Output

```
total 8
drwxr-xr-x   3 user  wheel   96 Jan 15 14:23 .
drwxr-xr-x  12 user  wheel  384 Jan 15 14:23 ..
-rw-r--r--   1 user  wheel   42 Jan 15 14:23 test.txt
```

## Command with Optional Arguments

Echo with optional name parameter.

### Configuration

```json
{
  "mcpServers": {
    "greet": {
      "command": "npx",
      "args": [
        "-y",
        "@studio-mcp/studio",
        "echo",
        "Hello",
        "[name]"
      ]
    }
  }
}
```

### Usage 1: With name

"Use greet with name: Alice"

**Output:**
```
Hello Alice
```

### Usage 2: Without name

"Use greet without any name"

**Output:**
```
Hello
```

## Command with Boolean Flags

List files with optional flags.

### Configuration

```json
{
  "mcpServers": {
    "ls": {
      "command": "npx",
      "args": [
        "-y",
        "@studio-mcp/studio",
        "ls",
        "[--all]",
        "[--long]",
        "{path}"
      ]
    }
  }
}
```

### Usage 1: With all flag

"Use ls with path /tmp and all: true"

**Output:**
```
. .. .hidden-file test.txt
```

### Usage 2: With both flags

"Use ls with path /tmp, all: true, and long: true"

**Output:**
```
total 8
drwxr-xr-x   3 user  wheel   96 Jan 15 14:23 .
drwxr-xr-x  12 user  wheel  384 Jan 15 14:23 ..
-rw-r--r--   1 user  wheel    0 Jan 15 14:23 .hidden-file
-rw-r--r--   1 user  wheel   42 Jan 15 14:23 test.txt
```

## Command with Array Arguments

Concatenate multiple files.

### Configuration

```json
{
  "mcpServers": {
    "concat": {
      "command": "npx",
      "args": [
        "-y",
        "@studio-mcp/studio",
        "cat",
        "{files... # List of files to concatenate}"
      ]
    }
  }
}
```

### Usage

"Use concat with files: /tmp/file1.txt, /tmp/file2.txt"

### Output

```
Contents of file1.txt
Contents of file2.txt
```

## GitHub API Access

Fetch data from GitHub's API.

### Configuration

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
        "https://api.github.com/{endpoint}"
      ]
    }
  }
}
```

### Usage

"Use github-api with endpoint: users/octocat"

### Output

```json
{
  "login": "octocat",
  "id": 583231,
  "name": "The Octocat",
  "bio": null,
  ...
}
```

## URL Shortener Check

Check where a short URL redirects.

### Configuration

```json
{
  "mcpServers": {
    "check-redirect": {
      "command": "npx",
      "args": [
        "-y",
        "@studio-mcp/studio",
        "curl",
        "-I",
        "-s",
        "{url}"
      ]
    }
  }
}
```

### Usage

"Check where https://bit.ly/example redirects"

### Output

```
HTTP/2 301
location: https://example.com
...
```

## Calculator

Perform basic arithmetic.

### Configuration

```json
{
  "mcpServers": {
    "calc": {
      "command": "npx",
      "args": [
        "-y",
        "@studio-mcp/studio",
        "bc",
        "-l"
      ]
    }
  }
}
```

### Note

This requires piping input, which is more complex. See [Advanced Examples](/examples/advanced) for calculator implementations.

## System Information

Get system information.

### Configuration (macOS/Linux)

```json
{
  "mcpServers": {
    "sysinfo": {
      "command": "npx",
      "args": [
        "-y",
        "@studio-mcp/studio",
        "uname",
        "-a"
      ]
    }
  }
}
```

### Usage

"What's the system information?"

### Output

```
Darwin MacBook.local 21.6.0 Darwin Kernel Version 21.6.0 x86_64
```

## Find Files

Search for files by name.

### Configuration

```json
{
  "mcpServers": {
    "find-files": {
      "command": "npx",
      "args": [
        "-y",
        "@studio-mcp/studio",
        "find",
        "{directory}",
        "-name",
        "{pattern}"
      ]
    }
  }
}
```

### Usage

"Find files in /tmp matching *.txt"

### Output

```
/tmp/file1.txt
/tmp/subdir/file2.txt
/tmp/test.txt
```

## Word Count

Count words in a file.

### Configuration

```json
{
  "mcpServers": {
    "word-count": {
      "command": "npx",
      "args": [
        "-y",
        "@studio-mcp/studio",
        "wc",
        "-w",
        "{file}"
      ]
    }
  }
}
```

### Usage

"Count words in /tmp/document.txt"

### Output

```
     1234 /tmp/document.txt
```

## Common Patterns

### Pattern 1: Simple Command

No parameters needed:

```json
{
  "command": "npx",
  "args": ["-y", "@studio-mcp/studio", "date"]
}
```

### Pattern 2: Required Parameter

One required string:

```json
{
  "command": "npx",
  "args": ["-y", "@studio-mcp/studio", "echo", "{message}"]
}
```

### Pattern 3: Optional Parameter

One optional string:

```json
{
  "command": "npx",
  "args": ["-y", "@studio-mcp/studio", "echo", "Hello", "[name]"]
}
```

### Pattern 4: Multiple Parameters

Mix of required and optional:

```json
{
  "command": "npx",
  "args": [
    "-y",
    "@studio-mcp/studio",
    "curl",
    "[--verbose]",
    "{url}"
  ]
}
```

### Pattern 5: Array Parameter

Variable number of arguments:

```json
{
  "command": "npx",
  "args": ["-y", "@studio-mcp/studio", "cat", "{files...}"]
}
```

## Testing Your Tools

### Command Line Testing

Test tools before adding to MCP clients:

```bash
# Test basic execution
npx -y @studio-mcp/studio echo "{message}"

# Test with MCP Inspector
npm install -g @modelcontextprotocol/inspector
mcp-inspector npx -y @studio-mcp/studio echo "{message}"
```

### Debug Mode

Add debug logging to troubleshoot:

```json
{
  "command": "npx",
  "args": [
    "-y",
    "@studio-mcp/studio",
    "--debug",
    "--log",
    "/tmp/tool.log",
    "echo",
    "{message}"
  ]
}
```

## Next Steps

- [Common Tools](/examples/common-tools) - Git, npm, docker examples
- [Advanced Examples](/examples/advanced) - Complex use cases
- [Template Syntax](/reference/template-syntax) - Full syntax reference
