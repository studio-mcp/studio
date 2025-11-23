# Template Syntax

Studio MCP uses a simple Mustache-like template syntax to define command arguments. This reference covers all syntax features and rules.

## Basic Syntax

Templates use curly braces `{}` and square brackets `[]` to define parameters:

| Syntax | Type | Required | Example |
|--------|------|----------|---------|
| `{name}` | String | Yes | `{message}` |
| `{{name}}` | String | Yes | `{{message}}` (legacy) |
| `[name]` | String | No | `[optional]` |
| `{name...}` | Array | Yes (1+) | `{files...}` |
| `[name...]` | Array | No | `[args...]` |
| `[--flag]` | Boolean | No | `[--verbose]` |

## Required Arguments

Use curly braces to define required parameters.

### Single Brace (Recommended)

```bash
studio echo "{message}"
```

Creates a required `message` parameter (string).

### Double Brace (Legacy)

```bash
studio echo "{{message}}"
```

Also creates a required parameter. Both syntaxes are equivalent.

### Example

```bash
# Command
studio curl "https://api.github.com/{endpoint}"

# Schema
{
  "endpoint": {
    "type": "string",
    "required": true
  }
}

# Usage
endpoint: "users/octocat"
# Executes: curl https://api.github.com/users/octocat
```

## Optional Arguments

Use square brackets for optional parameters.

```bash
studio echo "Hello" "[name]"
```

Creates an optional `name` parameter.

### Example

```bash
# Command
studio git status "[--short]"

# When provided
git status --short

# When omitted
git status
```

## Array Arguments

Arrays allow multiple values for a single parameter.

### Required Array

```bash
studio tar -czf archive.tar.gz "{files...}"
```

Requires at least one file.

```bash
# Usage
files: ["file1.txt", "file2.txt", "file3.txt"]
# Executes: tar -czf archive.tar.gz file1.txt file2.txt file3.txt
```

### Optional Array

```bash
studio git log "[paths...]"
```

Accepts zero or more paths.

```bash
# With values
paths: ["src/", "tests/"]
# Executes: git log src/ tests/

# Without values
paths: []
# Executes: git log
```

## Boolean Flags

Boolean flags appear or disappear based on a boolean value.

```bash
studio ls "[--all]" "[--long]"
```

Creates two boolean parameters: `all` and `long`.

```bash
# Usage
all: true, long: false
# Executes: ls --all

# Usage
all: true, long: true
# Executes: ls --all --long

# Usage
all: false, long: false
# Executes: ls
```

### Flag Naming

The flag name is derived from the template:

- `[--verbose]` → parameter: `verbose`
- `[--all]` → parameter: `all`
- `[-a]` → parameter: `a`

Parameter names automatically:
- Remove dashes
- Convert to lowercase
- Allow underscores

So `[--dry-run]` becomes parameter `dry_run` or `dryrun`.

## Descriptions

Add descriptions after `#` to document parameters.

```bash
studio echo "{message # The text to echo}"
```

Creates a parameter with description "The text to echo".

### Description Syntax

```bash
# Basic description
"{message # A message to display}"

# Optional with description
"[name # Optional user name]"

# Array with description
"{files... # List of files to process}"

# Boolean with description
"[--verbose # Enable verbose output]"
```

### Cursor Compatibility

Cursor has a bug with spaces in descriptions. Use underscores:

```bash
# Claude Desktop (works)
"{message # A message to display}"

# Cursor (required)
"{message#message_to_display}"
```

## Multiple Templates in One Argument

You can use multiple templates in a single argument:

```bash
studio echo "{{greeting}} {{name}}!"
```

```bash
# Usage
greeting: "Hello", name: "World"
# Executes: echo "Hello World!"
```

### Complex Example

```bash
studio curl "https://api.example.com/{{version}}/{{endpoint}}?key={{api_key}}"
```

```bash
# Usage
version: "v1", endpoint: "users", api_key: "secret123"
# Executes: curl "https://api.example.com/v1/users?key=secret123"
```

## Template in URLs

Templates work well in URLs:

```bash
studio curl "https://wttr.in/{city}?format={format}"
```

```bash
# Usage
city: "London", format: "3"
# Executes: curl "https://wttr.in/London?format=3"
```

## Mixing Static and Dynamic Content

Combine literal text with templates:

```bash
studio git commit -m "{message # Commit message}"
```

```bash
# Usage
message: "Fix bug in parser"
# Executes: git commit -m "Fix bug in parser"
```

## Field Naming Rules

Parameter names must follow these rules:

- **Allowed characters**: Letters, numbers, underscores, dashes
- **Case-insensitive**: `{userName}` and `{username}` are the same
- **Dash/underscore equivalence**: `{user-name}` and `{user_name}` are the same
- **No spaces**: `{user name}` is invalid

### Valid Names

```bash
{message}      ✓
{user_name}    ✓
{user-name}    ✓
{fileName}     ✓
{file123}      ✓
{_private}     ✓
```

### Invalid Names

```bash
{user name}    ✗ (space)
{user@name}    ✗ (special char)
{123file}      ✗ (starts with number)
{}             ✗ (empty)
```

## Edge Cases

### Empty Arguments

```bash
# Valid: empty text arguments
studio echo "" "{message}" ""

# Invalid: empty field names
studio echo "{}"  # Error
```

### Nested Templates

Not supported:

```bash
# Invalid
studio echo "{{name{suffix}}}"  # Error
```

### Escaped Braces

To include literal braces, don't use templates:

```bash
# Literal braces
studio echo "{ not a template }"

# This works because spaces around braces mean it's not a template
```

## Complete Examples

### Example 1: File Operations

```bash
studio cat "{file_path # Path to file to read}"
```

### Example 2: Multiple Required Args

```bash
studio cp "{source # Source file}" "{dest # Destination file}"
```

### Example 3: Optional Flags

```bash
studio ls "[--all]" "[--long]" "[--human-readable]" "{path}"
```

### Example 4: Array with Options

```bash
studio grep "{pattern}" "[--ignore-case]" "{files... # Files to search}"
```

### Example 5: Complex Command

```bash
studio ffmpeg -i "{input # Input video}" \
  "[--verbose]" \
  -c:v "{codec # Video codec}" \
  "{output # Output file}"
```

### Example 6: API Request

```bash
studio curl -X "{method # HTTP method}" \
  -H "Content-Type: application/json" \
  "[--data {payload}]" \
  "https://api.example.com/{endpoint}"
```

## Schema Generation

Templates automatically generate JSON schemas:

```bash
# Template
studio echo "{message}" "[name]"

# Generated Schema
{
  "type": "object",
  "properties": {
    "message": {
      "type": "string",
      "description": null
    },
    "name": {
      "type": "string",
      "description": null
    }
  },
  "required": ["message"]
}
```

### Array Schema

```bash
# Template
studio tar -czf archive.tar.gz "{files...}"

# Generated Schema
{
  "type": "object",
  "properties": {
    "files": {
      "type": "array",
      "items": {
        "type": "string"
      }
    }
  },
  "required": ["files"]
}
```

### Boolean Schema

```bash
# Template
studio ls "[--all]"

# Generated Schema
{
  "type": "object",
  "properties": {
    "all": {
      "type": "boolean"
    }
  },
  "required": []
}
```

## Best Practices

### 1. Use Descriptive Names

```bash
# Good
"{user_id # The user's ID}"

# Bad
"{id}"
```

### 2. Add Descriptions

```bash
# Good
"{endpoint # API endpoint path}"

# Bad
"{endpoint}"
```

### 3. Use Optional for Flags

```bash
# Good
"[--verbose]"

# Bad (don't make flags required)
"{--verbose}"
```

### 4. Array for Variable Arguments

```bash
# Good
"[files... # Files to process]"

# Bad (can't handle multiple files)
"[file1] [file2] [file3]"
```

### 5. Group Related Args

```bash
# Good
studio command "{required_arg}" "[optional_arg]" "[--flag]"

# Bad (hard to understand)
studio command "[--flag]" "{required_arg}" "[optional_arg]"
```

## Syntax Quick Reference

```bash
# Required string
{name}
{{name}}

# Optional string
[name]

# Required array (1+ items)
{name...}

# Optional array (0+ items)
[name...]

# Boolean flag
[--flag]

# With description
{name # description}

# Multiple in one arg
"{{greeting}} {{name}}"
```

## Next Steps

- [Examples](/examples/basic) - See templates in action
- [CLI Options](/reference/cli-options) - Command-line flags
- [Architecture](/reference/architecture) - How parsing works
