# Studio MCP Release Process

This document contains instructions for managing releases of the studio project, which uses a dual-publish system: Go binaries distributed via GitHub Releases and an npm wrapper package.

## Architecture Overview

The project maintains:
- **Go codebase**: Main application in Go
- **npm wrapper**: Downloads and executes appropriate binary for the platform
- **Dual distribution**: Available via `go install`, `npm install`, and direct binary download

## Release Components

### 1. Version Management
- `scripts/sync-version.sh`: Syncs package.json version with git tags
- `package.json`: npm package metadata and version
- Git tags: Source of truth for versions (format: `v1.2.3`)

### 2. Build & Release
- `.goreleaser.yml`: Cross-platform binary builds and GitHub releases
- `.github/workflows/release.yml`: CI/CD pipeline for releases
- `install.js`: Downloads platform-specific binary during npm install
- `run.js`: Executes the downloaded binary

## Release Process

### Step 1: Prepare Release
```bash
# Ensure you're on main branch with clean working directory
git checkout main
git pull origin main

# Test that everything builds locally
go build -o studio ./main.go
./studio --version

# Run tests if available
go test ./...
```

### Step 2: Version Bump
```bash
# Update version and create tag (replace X.Y.Z with desired version)
./scripts/sync-version.sh X.Y.Z

# This script:
# - Updates package.json version to X.Y.Z
# - Creates git tag vX.Y.Z
# - Shows next steps
```

### Step 3: Commit and Push
```bash
# Commit version changes
git add package.json
git commit -m "Bump version to X.Y.Z"

# Push changes and tag
git push origin main
git push origin vX.Y.Z
```

### Step 4: Automated Release
The GitHub Actions workflow (`.github/workflows/release.yml`) automatically:

1. **Triggers on tag push** (`v*` pattern)
2. **Runs GoReleaser** which:
   - Builds binaries for Linux, macOS, Windows (amd64/arm64)
   - Creates GitHub release with archives
   - Includes package.json and package-lock.json in release
3. **Extracts individual binaries** for npm compatibility:
   - `studio-linux` (from Linux archive)
   - `studio-macos` (from Darwin archive)
   - `studio-win.exe` (from Windows archive)
4. **Publishes to npm** using NODE_AUTH_TOKEN

### Step 5: Verify Release
```bash
# Check GitHub release was created
open https://github.com/studio-mcp/studio/releases/latest

# Test npm installation
npm install -g @studio-mcp/studio@latest
studio --version

# Test Go installation
go install github.com/studio-mcp/studio@latest
studio --version
```

## Required Secrets

GitHub repository must have these secrets configured:
- `NPM_TOKEN`: npm authentication token for publishing

## Binary Naming Convention

- **GoReleaser creates**: `studio_Linux_x86_64.tar.gz`, `studio_Darwin_x86_64.tar.gz`, `studio_Windows_x86_64.zip`
- **GitHub Actions extracts to**: `studio-linux`, `studio-macos`, `studio-win.exe`
- **npm wrapper expects**: Platform-specific binaries with consistent naming

## Troubleshooting

### Release Failed
- Check GitHub Actions logs for specific errors
- Ensure all required secrets are set
- Verify goreleaser configuration: `goreleaser check`

### npm Package Issues
- Binary naming must match between install.js and GitHub releases
- Ensure package.json files array includes necessary files
- Check Node.js version compatibility in engines field

### Version Sync Issues
- Ensure sync-version.sh is executable: `chmod +x scripts/sync-version.sh`
- Git tags must follow `v*` pattern for workflow trigger
- package.json version should not include 'v' prefix

## Manual Release

If automated release fails:

```bash
# Build and release manually with goreleaser
goreleaser release --clean

# Publish npm package manually
npm publish
```

## Testing Before Release

```bash
# Test goreleaser configuration
goreleaser check

# Test version sync script
./scripts/sync-version.sh

# Build locally to verify
go build -o studio ./main.go
./studio --version
```

## File Responsibilities

- `main.go`: Entry point, version info
- `cmd/root.go`: CLI implementation
- `internal/`: Core application logic
- `scripts/sync-version.sh`: Version management
- `.goreleaser.yml`: Build and release configuration
- `.github/workflows/release.yml`: CI/CD pipeline
- `package.json`: npm package definition
- `install.js`: Binary download logic
- `run.js`: Binary execution wrapper
