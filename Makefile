.PHONY: help build test version sync-version release dev clean verify-release pre-release-checks install-dev uninstall-dev claude claude-remove test-inspector

# Default target
help: ## Show this help message
	@echo "Available targets:"
	@awk 'BEGIN {FS = ":.*##"; printf "\033[36m%-15s\033[0m %s\n", "Target", "Description"} /^[a-zA-Z_-]+:.*?##/ { printf "\033[36m%-15s\033[0m %s\n", $$1, $$2 }' $(MAKEFILE_LIST)

# Get version from git tag or default
VERSION ?= $(shell git describe --tags --abbrev=0 2>/dev/null | sed 's/^v//' || echo "0.0.1")
COMMIT ?= $(shell git rev-parse --short HEAD)
DATE ?= $(shell date -u +"%Y-%m-%dT%H:%M:%SZ")

build: ## Build the Go binary
	go build -ldflags "-X main.Version=$(VERSION) -X main.commit=$(COMMIT) -X main.date=$(DATE)" -o bin/studio-mcp

install-dev: build ## Build and install binary to /usr/local/bin for dev testing
	@echo "Installing studio-mcp to /usr/local/bin..."
	sudo cp bin/studio-mcp /usr/local/bin/studio-mcp-dev
	@echo "✅ Installed as 'studio-mcp-dev' (to avoid conflicts with released version)"
	@echo "📋 Test with: studio-mcp-dev --version"
	@echo "📋 Uninstall with: make uninstall-dev"

uninstall-dev: ## Remove dev binary from /usr/local/bin
	@echo "Removing studio-mcp-dev from /usr/local/bin..."
	sudo rm -f /usr/local/bin/studio-mcp-dev
	@echo "✅ Uninstalled studio-mcp-dev"

claude: install-dev ## Install echo server into Claude Desktop MCP config
	@echo "Installing echo server into Claude Desktop..."
	@CLAUDE_CONFIG="$$HOME/Library/Application Support/Claude/claude_desktop_config.json"; \
	if [ ! -f "$$CLAUDE_CONFIG" ]; then \
		echo "Creating new Claude Desktop config..."; \
		mkdir -p "$$(dirname "$$CLAUDE_CONFIG")"; \
		echo '{"mcpServers":{}}' > "$$CLAUDE_CONFIG"; \
	fi; \
	echo "Updating MCP configuration..."; \
	node -e "const fs = require('fs'); \
	const configPath = '$$CLAUDE_CONFIG'; \
	const config = JSON.parse(fs.readFileSync(configPath, 'utf8')); \
	config.mcpServers = config.mcpServers || {}; \
	config.mcpServers.echo = { \
		command: '/usr/local/bin/studio-mcp-dev', \
		args: ['echo', '{{text # a message that will be echoed back to you}}'] \
	}; \
	fs.writeFileSync(configPath, JSON.stringify(config, null, 2));" && \
	echo "✅ Echo server installed in Claude Desktop" && \
	echo "📋 Restart Claude Desktop to load the new server" && \
	echo "📋 Remove with: make claude-remove"

claude-remove: ## Remove echo server from Claude Desktop MCP config
	@echo "Removing echo server from Claude Desktop..."
	@CLAUDE_CONFIG="$$HOME/Library/Application Support/Claude/claude_desktop_config.json"; \
	if [ ! -f "$$CLAUDE_CONFIG" ]; then \
		echo "⚠️  Claude Desktop config not found"; \
		exit 0; \
	fi; \
	node -e "const fs = require('fs'); \
	const configPath = '$$CLAUDE_CONFIG'; \
	const config = JSON.parse(fs.readFileSync(configPath, 'utf8')); \
	if (config.mcpServers && config.mcpServers.echo) { \
		delete config.mcpServers.echo; \
		fs.writeFileSync(configPath, JSON.stringify(config, null, 2)); \
		console.log('✅ Echo server removed from Claude Desktop'); \
	} else { \
		console.log('⚠️  Echo server not found in configuration'); \
	}" && \
	echo "📋 Restart Claude Desktop to apply changes"

test: ## Run tests
	go test ./...

test-inspector: ## Run MCP inspector tests
	cd cmd && go test -v ./inspector_test.go

version: ## Show current version information
	@echo "Version: $(VERSION)"
	@echo "Commit:  $(COMMIT)"
	@echo "Date:    $(DATE)"
	@echo "Package.json version: $(shell grep '"version"' package.json | cut -d'"' -f4)"
	@echo ""
	@echo "Binary version:"
	@if [ -f bin/studio-mcp ]; then ./bin/studio-mcp --version; else echo "Binary not built. Run 'make build' first."; fi

sync-version: ## Sync package.json version with latest git tag (for testing)
	./scripts/sync-version.sh

bump-version: ## Bump version (usage: make bump-version VERSION=1.2.3)
ifndef VERSION
	$(error VERSION is required. Usage: make bump-version VERSION=1.2.3)
endif
	git tag -a v$(VERSION) -m "Release v$(VERSION)"

verify-release: ## Verify the release setup locally
	@echo "🔍 Verifying release setup..."
	@echo ""
	@echo "1. Checking Go build..."
	$(MAKE) build
	@echo "✅ Go build successful"
	@echo ""
	@echo "2. Running tests..."
	$(MAKE) test
	@echo "✅ Tests passed"
	@echo ""
	@echo "3. Testing binary..."
	./bin/studio-mcp --version
	@echo "✅ Binary works"
	@echo ""
	@echo "4. Checking GoReleaser config..."
	goreleaser check || (echo "❌ GoReleaser config invalid. Install with: brew install goreleaser" && exit 1)
	@echo "✅ GoReleaser config valid"
	@echo ""
	@echo "5. Testing NPM postinstall script..."
	@if [ -f scripts/postinstall.js ]; then \
		echo "✅ Postinstall script exists"; \
	else \
		echo "❌ Postinstall script missing"; exit 1; \
	fi
	@echo ""
	@echo "6. Simulating GoReleaser build (dry-run)..."
	goreleaser build --snapshot --clean || (echo "❌ GoReleaser build failed" && exit 1)
	@echo "✅ GoReleaser build simulation successful"
	@echo ""
	@echo "🎉 All verification checks passed!"

pre-release-checks: ## Run pre-release checks
ifndef VERSION
	$(error VERSION is required. Usage: make pre-release-checks VERSION=1.2.3)
endif
	@echo "🔍 Pre-release checks for v$(VERSION)..."
	@echo ""
	@echo "1. Checking if version already exists..."
	@if git tag -l | grep -q "^v$(VERSION)$$"; then \
		echo "❌ Version v$(VERSION) already exists!"; \
		echo "Existing tags:"; git tag -l | tail -5; \
		exit 1; \
	else \
		echo "✅ Version v$(VERSION) is new"; \
	fi
	@echo ""
	@echo "2. Checking working directory is clean..."
	@if [ -n "$$(git status --porcelain)" ]; then \
		echo "❌ Working directory is not clean:"; \
		git status --short; \
		echo "Please commit or stash changes first"; \
		exit 1; \
	else \
		echo "✅ Working directory is clean"; \
	fi
	@echo ""
	@echo "3. Running verification..."
	$(MAKE) verify-release
	@echo ""
	@echo "🎉 All pre-release checks passed for v$(VERSION)!"

release: ## Create a release (usage: make release VERSION=1.2.3)
ifndef VERSION
	$(error VERSION is required. Usage: make release VERSION=1.2.3)
endif
	@echo "🚀 Starting release process for v$(VERSION)..."
	$(MAKE) pre-release-checks VERSION=$(VERSION)
	@echo ""
	@echo "Creating and pushing tag v$(VERSION)..."
	$(MAKE) bump-version VERSION=$(VERSION)
	git push origin v$(VERSION)
	@echo ""
	@echo "🎉 Release v$(VERSION) initiated!"
	@echo "📋 The GitHub Actions workflow will now:"
	@echo "  • Sync package.json version to $(VERSION)"
	@echo "  • Build cross-platform binaries with GoReleaser"
	@echo "  • Create GitHub release with binaries + checksums"
	@echo "  • Publish NPM package that fetches those binaries"
	@echo ""
	@echo "📋 Monitor progress:"
	@echo "  1. GitHub Actions: https://github.com/studio-mcp/studio-mcp/actions"
	@echo "  2. GitHub Release: https://github.com/studio-mcp/studio-mcp/releases"
	@echo "  3. NPM Package: https://www.npmjs.com/package/studio-mcp"
	@echo "  4. Test install: npm install studio-mcp@$(VERSION)"

dev: ## Run in development mode
	go run -ldflags "-X main.Version=dev -X main.commit=$(COMMIT) -X main.date=$(DATE)" . $(ARGS)

clean: ## Clean build artifacts
	rm -rf bin/studio-mcp dist/ studio-mcp
