# Installation

Studio MCP can be installed globally, used via npx, or run from source for development.

## Prerequisites

- **Node.js**: Version 18 or later
- **npm**: Comes with Node.js (or use yarn, pnpm, bun)

Check your versions:

```bash
node --version  # Should be v18.0.0 or higher
npm --version
```

## Installation Methods

### Method 1: npx (No Installation)

The easiest way to use Studio is with npx, which requires no installation:

```bash
npx -y @studio-mcp/studio echo "{message}"
```

**Pros:**
- No global packages to manage
- Always uses the latest version
- No installation step needed

**Cons:**
- Slightly slower first run (caches afterward)
- Requires `-y` flag in MCP configs

**Best for:** Quick testing, MCP client configurations

### Method 2: Global Installation

Install Studio globally for faster execution:

```bash
npm install -g @studio-mcp/studio
```

Verify installation:

```bash
studio --version
```

Use directly:

```bash
studio echo "{message}"
```

**Pros:**
- Faster execution (no npx overhead)
- Cleaner command (just `studio`)
- Better for local development

**Cons:**
- Requires manual updates
- Global package management
- May need sudo on some systems

**Best for:** Frequent local use, development

### Method 3: Local Installation

Install in a project:

```bash
npm install @studio-mcp/studio
```

Use via package.json scripts:

```json
{
  "scripts": {
    "mcp:echo": "studio echo \"{message}\""
  }
}
```

**Best for:** Project-specific tools

### Method 4: From Source

Clone and build for development:

```bash
git clone https://github.com/studio-mcp/studio.git
cd studio
npm install
npm run build
```

Run with bun (development):

```bash
bun run dev echo "{message}"
```

Or link globally:

```bash
npm link
studio echo "{message}"
```

**Best for:** Contributing to Studio, debugging

## Platform-Specific Notes

### macOS

Standard installation works out of the box:

```bash
npm install -g @studio-mcp/studio
```

Claude Desktop config location:
```
~/Library/Application Support/Claude/claude_desktop_config.json
```

### Windows

Use an elevated command prompt or PowerShell:

```bash
npm install -g @studio-mcp/studio
```

Claude Desktop config location:
```
%APPDATA%\Claude\claude_desktop_config.json
```

### Linux

You may need to use sudo for global installation:

```bash
sudo npm install -g @studio-mcp/studio
```

Or use a node version manager (nvm) to avoid sudo:

```bash
# Install nvm first
curl -o- https://raw.githubusercontent.com/nvm-sh/nvm/v0.39.0/install.sh | bash

# Install Node via nvm
nvm install 20
nvm use 20

# Now install Studio without sudo
npm install -g @studio-mcp/studio
```

## Package Managers

Studio works with all major Node package managers:

### npm (default)

```bash
npm install -g @studio-mcp/studio
```

### yarn

```bash
yarn global add @studio-mcp/studio
```

### pnpm

```bash
pnpm add -g @studio-mcp/studio
```

### bun

```bash
bun install -g @studio-mcp/studio
```

## Updating

### Update npx (automatic)

npx automatically uses the latest version. To force update the cache:

```bash
npx clear-npx-cache
npx -y @studio-mcp/studio --version
```

### Update global installation

```bash
npm update -g @studio-mcp/studio
```

Check the installed version:

```bash
studio --version
```

### Update local installation

```bash
npm update @studio-mcp/studio
```

## Uninstalling

### Remove global installation

```bash
npm uninstall -g @studio-mcp/studio
```

### Remove npx cache

```bash
npm cache clean --force
```

### Remove local installation

```bash
npm uninstall @studio-mcp/studio
```

## Troubleshooting

### Command not found

If `studio` is not found after global installation:

1. Check if npm global bin is in your PATH:
   ```bash
   npm config get prefix
   ```

2. Add to PATH (add to ~/.bashrc or ~/.zshrc):
   ```bash
   export PATH="$(npm config get prefix)/bin:$PATH"
   ```

3. Reload your shell:
   ```bash
   source ~/.bashrc  # or ~/.zshrc
   ```

### Permission errors on macOS/Linux

Use nvm instead of sudo:

```bash
curl -o- https://raw.githubusercontent.com/nvm-sh/nvm/v0.39.0/install.sh | bash
nvm install 20
npm install -g @studio-mcp/studio
```

### npx slow first run

The first `npx` run downloads and caches the package. Subsequent runs are faster. For even faster execution, use global installation instead.

## Verifying Installation

Test your installation:

```bash
# Check version
studio --version

# Run help
studio --help

# Test a simple command
studio echo "Hello, Studio!"
```

You should see the MCP protocol messages on stdout (this is normal for MCP servers).

## Next Steps

- [Getting Started](/guide/getting-started) - Create your first tool
- [Configuration](/guide/config-claude) - Set up with MCP clients
- [Template Syntax](/reference/template-syntax) - Learn the template language
