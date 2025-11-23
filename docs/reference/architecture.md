# Architecture

This document explains how Studio MCP works internally, from template parsing to command execution.

## Overview

Studio MCP is a lightweight bridge between the Model Context Protocol and command-line tools. It consists of several key components working together:

```
┌─────────────────────────────────────────────────────────────┐
│                        MCP Client                           │
│                   (Claude, Cursor, etc.)                    │
└────────────────────────┬────────────────────────────────────┘
                         │ JSON-RPC over stdio
                         │
┌────────────────────────▼────────────────────────────────────┐
│                      Studio MCP                             │
│  ┌──────────────────────────────────────────────────────┐  │
│  │ 1. Parse Arguments                                   │  │
│  │    - Extract flags (--debug, --log)                  │  │
│  │    - Extract command template                        │  │
│  └──────────────────────┬───────────────────────────────┘  │
│                         │                                   │
│  ┌──────────────────────▼───────────────────────────────┐  │
│  │ 2. Template Parsing                                  │  │
│  │    - Tokenize shell words                            │  │
│  │    - Identify fields: {name}, [name], etc.           │  │
│  │    - Extract descriptions                            │  │
│  └──────────────────────┬───────────────────────────────┘  │
│                         │                                   │
│  ┌──────────────────────▼───────────────────────────────┐  │
│  │ 3. Schema Generation                                 │  │
│  │    - Convert fields to JSON schema                   │  │
│  │    - Define types (string, array, boolean)           │  │
│  │    - Mark required vs optional                       │  │
│  └──────────────────────┬───────────────────────────────┘  │
│                         │                                   │
│  ┌──────────────────────▼───────────────────────────────┐  │
│  │ 4. MCP Server Registration                           │  │
│  │    - Create tool with schema                         │  │
│  │    - Set up JSON-RPC handlers                        │  │
│  │    - Start stdio transport                           │  │
│  └──────────────────────┬───────────────────────────────┘  │
│                         │                                   │
│  ┌──────────────────────▼───────────────────────────────┐  │
│  │ 5. Request Handling                                  │  │
│  │    - Receive tool call from client                   │  │
│  │    - Validate parameters with Zod                    │  │
│  │    - Render command with values                      │  │
│  └──────────────────────┬───────────────────────────────┘  │
│                         │                                   │
│  ┌──────────────────────▼───────────────────────────────┐  │
│  │ 6. Command Execution                                 │  │
│  │    - Spawn shell command                             │  │
│  │    - Capture stdout/stderr                           │  │
│  │    - Return output to client                         │  │
│  └──────────────────────────────────────────────────────┘  │
└────────────────────────┬────────────────────────────────────┘
                         │ spawn()
                         │
┌────────────────────────▼────────────────────────────────────┐
│                     Shell Command                           │
│                 (echo, curl, git, etc.)                     │
└─────────────────────────────────────────────────────────────┘
```

## Component Details

### 1. Argument Parsing

**Location:** `src/index.ts:parseArgs()`

Parses command-line arguments into Studio options and command template.

**Input:**
```javascript
process.argv = [
  'node', 'studio.js',
  '--debug',
  '--log', '/tmp/studio.log',
  'echo', '{message}'
]
```

**Output:**
```javascript
{
  debug: true,
  logFile: '/tmp/studio.log',
  commandArgs: ['echo', '{message}']
}
```

**Algorithm:**
1. Skip first two args (node, script)
2. Parse flags until non-flag found or `--` encountered
3. Remaining args become command template

### 2. Template Parsing

**Location:** `src/studio/parse.ts`

Converts command template into structured tokens.

#### Tokenization

**Function:** `fromArgs(args: string[]): Template`

**Input:**
```javascript
['echo', '{message}', 'world']
```

**Process:**
1. Iterate through each argument
2. Check for template patterns:
   - `{name}` → FieldToken (required)
   - `[name]` → FieldToken (optional)
   - `[--flag]` → FieldToken (boolean)
   - Otherwise → TextToken
3. Parse descriptions after `#`
4. Handle arrays with `...` suffix

**Output:**
```javascript
{
  args: [
    { kind: 'text', value: 'echo' },
    { kind: 'field', name: 'message', required: true },
    { kind: 'text', value: 'world' }
  ],
  fields: [
    { kind: 'field', name: 'message', required: true }
  ]
}
```

#### Field Patterns

Parsed using these regex patterns:

```typescript
// Required: {name} or {{name}}
/^\{\{?(\w+)(?:\s*#\s*(.+?))?\}?\}$/

// Optional: [name]
/^\[(\w+)(?:\s*#\s*(.+?))?\]$/

// Array: {name...} or [name...]
/^\{(\w+)\.\.\.(?:\s*#\s*(.+?))?\}$/
/^\[(\w+)\.\.\.(?:\s*#\s*(.+?))?\]$/

// Boolean: [--flag]
/^\[(--?\w+)(?:\s*#\s*(.+?))?\]$/
```

### 3. Schema Generation

**Location:** `src/studio/schema.ts:generateInputSchema()`

Converts template fields into JSON Schema for MCP.

**Input:**
```javascript
fields: [
  { name: 'message', required: true, array: false, boolean: false },
  { name: 'count', required: false, array: false, boolean: false }
]
```

**Output:**
```javascript
{
  type: 'object',
  properties: {
    message: { type: 'string', description: null },
    count: { type: 'string', description: null }
  },
  required: ['message']
}
```

**Field Type Mapping:**

| Template | JSON Schema Type |
|----------|------------------|
| `{name}` | `{ type: 'string' }` |
| `[name]` | `{ type: 'string' }` |
| `{name...}` | `{ type: 'array', items: { type: 'string' } }` |
| `[--flag]` | `{ type: 'boolean' }` |

### 4. Zod Schema Conversion

**Location:** `src/index.ts:schemaToZod()`

Converts JSON Schema to Zod schema for runtime validation.

**Input:**
```javascript
{
  type: 'object',
  properties: {
    message: { type: 'string' }
  },
  required: ['message']
}
```

**Output:**
```javascript
z.object({
  message: z.string()
})
```

**Type Conversion:**

```typescript
// String
{ type: 'string' } → z.string()

// Array
{ type: 'array', items: { type: 'string' } } → z.array(z.string())

// Boolean
{ type: 'boolean' } → z.boolean()

// Optional (not in required)
Not in required → .optional()
```

### 5. MCP Server Setup

**Location:** `src/index.ts:main()`

Creates and configures the MCP server.

**Steps:**

1. **Create server instance:**
   ```typescript
   const server = new McpServer({
     name: 'studio-mcp',
     version: '0.0.0'
   })
   ```

2. **Derive tool name from command:**
   ```typescript
   const toolName = commandName.replace(/-/g, '_')
   // 'git-status' → 'git_status'
   ```

3. **Register tool:**
   ```typescript
   server.tool(
     toolName,
     inputSchema,
     async (params) => {
       // Handle tool call
     }
   )
   ```

4. **Start transport:**
   ```typescript
   const transport = new StdioServerTransport()
   await server.connect(transport)
   ```

### 6. Command Rendering

**Location:** `src/studio/render.ts:buildCommandArgs()`

Renders the final command with parameter values.

**Input:**
```javascript
template: {
  args: [
    { kind: 'text', value: 'echo' },
    { kind: 'field', name: 'message' },
    { kind: 'text', value: 'world' }
  ]
},
params: { message: 'Hello' }
```

**Process:**

1. Iterate through template tokens
2. For TextToken: add value as-is
3. For FieldToken:
   - String: add parameter value
   - Array: spread values
   - Boolean: add flag if true

**Output:**
```javascript
['echo', 'Hello', 'world']
```

#### Array Rendering

```javascript
// Template
template: [
  { kind: 'text', value: 'git' },
  { kind: 'field', name: 'args', array: true }
]

// Parameters
params: { args: ['status', '--short'] }

// Output
['git', 'status', '--short']
```

#### Boolean Rendering

```javascript
// Template
template: [
  { kind: 'text', value: 'ls' },
  { kind: 'field', name: 'all', boolean: true, value: '--all' }
]

// Parameters (all: true)
params: { all: true }
// Output: ['ls', '--all']

// Parameters (all: false)
params: { all: false }
// Output: ['ls']
```

### 7. Command Execution

**Location:** `src/index.ts` (tool handler)

Executes the rendered command and returns output.

**Process:**

```typescript
const { spawn } = require('child_process')

const [cmd, ...args] = renderedCommand
const child = spawn(cmd, args, {
  stdio: ['ignore', 'pipe', 'pipe']
})

let stdout = ''
let stderr = ''

child.stdout.on('data', (data) => {
  stdout += data.toString()
})

child.stderr.on('data', (data) => {
  stderr += data.toString()
})

child.on('close', (code) => {
  if (code === 0) {
    return { content: [{ type: 'text', text: stdout }] }
  } else {
    return { isError: true, content: stderr }
  }
})
```

## Data Flow Example

Complete flow for `studio echo "{message}"`:

### Step 1: Parse Arguments

```
Input: ['studio', 'echo', '{message}']
Output: commandArgs = ['echo', '{message}']
```

### Step 2: Parse Template

```
Input: ['echo', '{message}']
Output: {
  args: [
    { kind: 'text', value: 'echo' },
    { kind: 'field', name: 'message', required: true }
  ],
  fields: [
    { name: 'message', required: true }
  ]
}
```

### Step 3: Generate Schema

```
Input: fields = [{ name: 'message', required: true }]
Output: {
  type: 'object',
  properties: {
    message: { type: 'string' }
  },
  required: ['message']
}
```

### Step 4: Register Tool

```
Tool name: 'echo'
Tool schema: { ... }
Tool handler: async (params) => { ... }
```

### Step 5: Receive Call

```
MCP client sends:
{
  "jsonrpc": "2.0",
  "method": "tools/call",
  "params": {
    "name": "echo",
    "arguments": { "message": "Hello" }
  }
}
```

### Step 6: Render Command

```
Input: template + { message: 'Hello' }
Output: ['echo', 'Hello']
```

### Step 7: Execute

```
spawn('echo', ['Hello'])
stdout: "Hello\n"
Exit code: 0
```

### Step 8: Return Result

```
MCP server responds:
{
  "jsonrpc": "2.0",
  "result": {
    "content": [
      { "type": "text", "text": "Hello\n" }
    ]
  }
}
```

## Key Design Decisions

### 1. Stdio Transport

Studio uses stdio for MCP communication because:
- Simplest integration with MCP clients
- No network configuration needed
- Process isolation and security
- Standard for CLI tools

### 2. Shell Word Parsing

Arguments are parsed as shell words:
- Supports quoted strings: `"hello world"`
- Handles escape sequences
- Preserves spaces in quotes
- Compatible with shell syntax

### 3. Template Syntax

Design choices:
- **Familiar:** Mustache-like syntax developers know
- **Minimal:** Only braces and brackets, no complex DSL
- **Shell-compatible:** Works in JSON strings
- **Descriptive:** Inline descriptions with `#`

### 4. Schema Generation

Automatic schema generation provides:
- No manual JSON schema writing
- Guaranteed schema/template sync
- Type safety with Zod validation
- Clear parameter documentation

### 5. Single Tool Per Server

Each Studio instance serves one tool because:
- Simple mental model
- Easy configuration
- Independent tool lifecycles
- No cross-tool complexity

## Performance Characteristics

### Startup Time

```
Component               Time
─────────────────────  ─────
Node.js startup         ~50ms
Parse arguments         <1ms
Parse template          <1ms
Generate schema         <1ms
Start MCP server        ~10ms
─────────────────────  ─────
Total                   ~60ms
```

### Runtime Overhead

Per tool call:
```
Component               Time
─────────────────────  ─────
Validate params         <1ms
Render command          <1ms
Spawn process           ~5ms
Execute command         varies
Parse output            <1ms
─────────────────────  ─────
Studio overhead         ~7ms
```

The command execution time dominates; Studio adds minimal overhead.

## Dependencies

### Core Dependencies

- **@modelcontextprotocol/sdk** - MCP protocol implementation
- **zod** - Schema validation

### Why So Few?

Minimal dependencies provide:
- Fast installation
- Small bundle size
- Fewer security vulnerabilities
- Easy auditing
- Stable foundation

## Error Handling

### Template Parsing Errors

```typescript
// Invalid template
"{}"  → Error: Empty field name
"{a{b}}"  → Error: Nested templates not supported
```

### Validation Errors

```typescript
// Missing required parameter
Required: message
Provided: {}
Error: "message is required"
```

### Execution Errors

```typescript
// Command not found
Command: nonexistent
Error: spawn ENOENT

// Command failed
Exit code: 1
Stderr: "error message"
```

## Security Considerations

### Command Injection Prevention

Studio does NOT prevent command injection. It's designed to:
- Run the specified command with parameters
- Not interpret or execute shell syntax in parameters

MCP clients should sanitize user input before calling tools.

### Process Isolation

Each command runs as a subprocess:
- Inherits parent environment
- No shell interpretation (uses spawn, not exec)
- Isolated from MCP server process

### Permission Model

Studio runs with the permissions of the user who started it:
- No privilege escalation
- Standard file system permissions apply
- Commands run as the invoking user

## Testing Architecture

### Unit Tests

Test individual components:
- Template parsing
- Schema generation
- Command rendering

**Location:** `src/__tests__/*.test.ts`

### Integration Tests

Test complete flows:
- MCP protocol handling
- Command execution
- Error scenarios

### Test Utilities

```typescript
// Parse test
import { fromArgs } from '../studio/parse'
const template = fromArgs(['echo', '{message}'])

// Schema test
import { generateInputSchema } from '../studio/schema'
const schema = generateInputSchema(template.fields)

// Render test
import { buildCommandArgs } from '../studio/render'
const command = buildCommandArgs(template, { message: 'test' })
```

## Extensibility Points

While Studio is intentionally minimal, these extension points exist:

### Custom Transports

Replace StdioServerTransport:
```typescript
import { SSEServerTransport } from '@modelcontextprotocol/sdk'
const transport = new SSEServerTransport(...)
```

### Custom Validation

Extend schema generation:
```typescript
// Add custom constraints
const schema = generateInputSchema(fields)
schema.properties.email.format = 'email'
```

### Middleware

Wrap tool handler:
```typescript
server.tool(toolName, schema, async (params) => {
  // Pre-processing
  const result = await originalHandler(params)
  // Post-processing
  return result
})
```

## Next Steps

- [Template Syntax](/reference/template-syntax) - Syntax details
- [Development](/guide/development) - Contributing guide
- [Examples](/examples/advanced) - Advanced patterns
