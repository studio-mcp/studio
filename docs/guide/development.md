# Development

This guide covers how to develop, test, and contribute to Studio MCP itself.

## Setting Up Development Environment

### Prerequisites

- Node.js 18 or later
- Bun (recommended) or npm
- Git

### Clone the Repository

```bash
git clone https://github.com/studio-mcp/studio.git
cd studio
```

### Install Dependencies

```bash
# Using bun (recommended)
bun install

# Or using npm
npm install
```

### Project Structure

```
studio/
├── src/
│   ├── index.ts              # Main entry point
│   ├── studio/
│   │   ├── types.ts          # Type definitions
│   │   ├── parse.ts          # Template parsing
│   │   ├── schema.ts         # JSON schema generation
│   │   └── render.ts         # Command rendering
│   └── __tests__/            # Test files
├── package.json
├── tsconfig.json
└── README.md
```

## Development Commands

### Build

Compile TypeScript to JavaScript:

```bash
bun run build
# or
npm run build
```

Output goes to `dist/` directory.

### Run in Development Mode

Run without building:

```bash
bun run dev echo "{message}"
# or
npm run dev echo "{message}"
```

### Run Tests

Execute the test suite:

```bash
# Run all tests
bun run test

# Watch mode
bun run test:watch

# Specific test file
bun test src/__tests__/parse.test.ts
```

### Linting and Formatting

```bash
# Lint code
bun run lint

# Format code
bun run format
```

## Testing Your Changes

### Unit Tests

Studio uses Bun's built-in test runner. Create tests in `src/__tests__/`:

```typescript
import { describe, test, expect } from 'bun:test'
import { parseTemplate } from '../studio/parse'

describe('parseTemplate', () => {
  test('parses required arguments', () => {
    const template = parseTemplate(['echo', '{message}'])
    expect(template.fields).toHaveLength(1)
    expect(template.fields[0].name).toBe('message')
    expect(template.fields[0].required).toBe(true)
  })
})
```

### Integration Testing

Test with MCP Inspector:

```bash
# Install inspector
npm install -g @modelcontextprotocol/inspector

# Run with dev build
bun run inspector:echo

# Or custom command
mcp-inspector bun run dev echo "{message}"
```

### Manual Testing

Test the built version:

```bash
# Build first
bun run build

# Test the built binary
node dist/index.js echo "{message}"
```

## Code Architecture

### Template Parsing Flow

```
User Input → Tokenization → Field Extraction → Schema Generation → Command Execution
```

#### 1. Tokenization (src/studio/parse.ts)

Converts command arguments into tokens:

```typescript
// Input: ["echo", "{message}", "world"]
// Output: [
//   TextToken { value: "echo" },
//   FieldToken { name: "message", required: true },
//   TextToken { value: "world" }
// ]
```

#### 2. Schema Generation (src/studio/schema.ts)

Creates JSON schema from tokens:

```typescript
// Input: FieldToken { name: "message", required: true }
// Output: {
//   type: "object",
//   properties: {
//     message: { type: "string" }
//   },
//   required: ["message"]
// }
```

#### 3. Command Rendering (src/studio/render.ts)

Builds final command with parameters:

```typescript
// Input: template + { message: "Hello" }
// Output: ["echo", "Hello", "world"]
```

### Key Types

From `src/studio/types.ts`:

```typescript
interface Template {
  args: Token[]
  fields: FieldToken[]
}

interface FieldToken {
  kind: 'field'
  name: string
  description?: string
  required: boolean
  array: boolean
  boolean: boolean
}

interface TextToken {
  kind: 'text'
  value: string
}
```

## Adding New Features

### Example: Adding a New Template Syntax

1. **Update types** in `src/studio/types.ts`
2. **Add parsing logic** in `src/studio/parse.ts`
3. **Update schema generation** in `src/studio/schema.ts`
4. **Add rendering logic** in `src/studio/render.ts`
5. **Write tests** in `src/__tests__/`
6. **Update documentation**

### Example: Adding a New CLI Flag

1. **Update parseArgs()** in `src/index.ts`
2. **Add to help text**
3. **Implement functionality**
4. **Add tests**
5. **Update README**

## Debugging Development

### Enable Debug Logging

```typescript
// In src/index.ts
const DEBUG = process.env.DEBUG === 'true'

if (DEBUG) {
  console.error('[DEBUG]', debugInfo)
}
```

Run with debug:

```bash
DEBUG=true bun run dev echo "{message}"
```

### Use MCP Inspector

The inspector is invaluable for development:

```bash
# Run current dev build
mcp-inspector bun run dev echo "{message}"

# Test specific scenarios
mcp-inspector bun run dev --debug --log /tmp/dev.log curl "{url}"
```

### Test with Real MCP Clients

Create a dev configuration for Claude Desktop:

```json
{
  "mcpServers": {
    "studio-dev": {
      "command": "bun",
      "args": [
        "run",
        "/path/to/studio/src/index.ts",
        "--debug",
        "--log",
        "/tmp/studio-dev.log",
        "echo",
        "{message}"
      ]
    }
  }
}
```

## Running Examples

The repository includes example configurations:

```bash
# Test with the example config
cat .mcp.json

# Run with inspector
bun run inspector:echo
bun run inspector:git
```

## Contributing

### Before Submitting

1. **Run tests**: `bun run test`
2. **Check types**: `bun run build`
3. **Format code**: `bun run format`
4. **Update docs** if needed
5. **Add tests** for new features

### Pull Request Process

1. Fork the repository
2. Create a feature branch: `git checkout -b feature/amazing-feature`
3. Make your changes
4. Add tests
5. Commit: `git commit -m 'Add amazing feature'`
6. Push: `git push origin feature/amazing-feature`
7. Open a Pull Request

### Commit Message Convention

Use conventional commits:

```
feat: add support for regex patterns in templates
fix: handle empty array parameters correctly
docs: update template syntax reference
test: add tests for boolean flags
chore: update dependencies
```

## Publishing (Maintainers)

### Version Bump

```bash
# Patch (0.0.x)
npm version patch

# Minor (0.x.0)
npm version minor

# Major (x.0.0)
npm version major
```

### Build and Publish

```bash
# Build
bun run build

# Publish to npm
npm publish

# Or with tag
npm publish --tag beta
```

### Release Checklist

- [ ] Update version in package.json
- [ ] Update CHANGELOG.md
- [ ] Run full test suite
- [ ] Build successfully
- [ ] Test built package locally
- [ ] Create git tag
- [ ] Push to GitHub
- [ ] Publish to npm
- [ ] Create GitHub release

## Development Tips

### Quick Iteration

Use watch mode for tests:

```bash
bun run test:watch
```

### Test Multiple Scenarios

Create a test script:

```bash
#!/bin/bash
# test-scenarios.sh

echo "Test 1: Required arg"
bun run dev echo "{message}"

echo "Test 2: Optional arg"
bun run dev echo "[message]"

echo "Test 3: Array arg"
bun run dev git "[args...]"
```

### Use Type Checking

TypeScript helps catch errors early:

```bash
# Check types without building
tsc --noEmit
```

### Debug with VSCode

Create `.vscode/launch.json`:

```json
{
  "version": "0.2.0",
  "configurations": [
    {
      "type": "node",
      "request": "launch",
      "name": "Debug Studio",
      "program": "${workspaceFolder}/src/index.ts",
      "args": ["echo", "{message}"],
      "runtimeExecutable": "bun",
      "runtimeArgs": ["run"],
      "console": "integratedTerminal"
    }
  ]
}
```

## Resources

- [TypeScript Handbook](https://www.typescriptlang.org/docs/handbook/intro.html)
- [Bun Documentation](https://bun.sh/docs)
- [MCP Specification](https://modelcontextprotocol.io/)
- [Zod Documentation](https://zod.dev/)

## Getting Help

- GitHub Discussions for questions
- GitHub Issues for bugs
- Pull Requests for contributions

## Next Steps

- [Architecture Reference](/reference/architecture) - Deep dive into internals
- [Template Syntax](/reference/template-syntax) - Complete syntax reference
- [Examples](/examples/basic) - Real-world examples
