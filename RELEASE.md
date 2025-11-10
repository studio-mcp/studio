# Studio MCP Release Process

This document describes the release process for studio, a TypeScript/npm package.

## Overview

Studio uses a tag-based release workflow:

- Push a git tag → GitHub Actions automatically builds, tests, and publishes to npm
- No manual version bumping required
- Single source of truth: git tags

## Quick Release

Create the tag in the GitHub UI, or:

```bash
# Create and push tag (replace X.Y.Z with desired version)
git tag -a v1.2.3 -m "Release v1.2.3"
git push origin v1.2.3
```

GitHub Actions will:

1. Extract version from tag
2. Update package.json version
3. Install dependencies
4. Run all tests
5. Build TypeScript → JavaScript
6. Publish to npm with provenance
7. Create GitHub release with notes

## Detailed Release Process

### Prerequisites

- Clean working directory
- All tests passing locally
- On main branch (or release branch)

### Step 1: Test Locally

```bash
# Run tests
bun run test

# Build
bun run build

# Verify binary works
node dist/index.js --version
```

### Step 2: Create and Push Tag

```bash
# Create annotated tag
git tag -a v1.2.3 -m "Release v1.2.3"

# Push tag to trigger release
git push origin v1.2.3

# Or do both in one command
git tag -a v1.2.3 -m "Release v1.2.3" && git push origin v1.2.3
```

### Step 3: Monitor Release

GitHub Actions will automatically:

- Run CI tests
- Build the package
- Publish to npm
- Create GitHub release

Monitor progress at:

- **GitHub Actions**: <https://github.com/studio-mcp/studio/actions>
- **GitHub Releases**: <https://github.com/studio-mcp/studio/releases>
- **npm Package**: <https://www.npmjs.com/package/@studio-mcp/studio>

### Step 4: Verify Release

```bash
# Test npm installation
npm install -g @studio-mcp/studio@latest
studio --version

# Or with npx
npx @studio-mcp/studio@latest --version
```

## Trusted Publishing

This project uses [trusted publishing](https://docs.npmjs.com/trusted-publishers).

### One-time setup on npmjs.com

1. Go to your package settings: <https://www.npmjs.com/package/@studio-mcp/studio/access>
2. Find the "Trusted Publisher" section
3. Click "GitHub Actions"
4. Configure:
   - **Organization or user**: `studio-mcp`
   - **Repository**: `studio`
   - **Workflow filename**: `release.yml`
   - **Environment name**: `release`
5. Click "Add Trusted Publisher"

## Rollback

If a release has issues:

```bash
# Deprecate the problematic version on npm
npm deprecate @studio-mcp/studio@X.Y.Z "Reason for deprecation"

# Publish a patch release
git tag vX.Y.Z+1
git push origin vX.Y.Z+1
```

### Dry run (test without publishing)

Test locally before releasing:

```bash
# Install dependencies
bun install

# Run tests
bun run test

# Build
bun run build

# Preview what would be published
npm pack --dry-run

# Or create a tarball to inspect
npm pack
tar -tzf studio-mcp-studio-*.tgz
```

## Development

### Local builds

```bash
# Build
bun run build

# Watch mode
bun run test:watch

# Development mode (runs source directly)
bun run dev echo "{{text}}"
```

### Testing with MCP Inspector

```bash
# List tools
bun run inspector:echo

# Test git command
bun run inspector:git
```
