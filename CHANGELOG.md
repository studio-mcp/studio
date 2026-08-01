# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

- Update to the [2026-07-28 MCP specification](https://modelcontextprotocol.io/specification/2026-07-28)
  by moving from `@modelcontextprotocol/sdk` v1 to the v2 SDK
  (`@modelcontextprotocol/server`). 2025-era clients keep working: the stdio
  entry point negotiates the protocol era from the opening message.
- Advertise tool input schemas as JSON Schema 2020-12 (requires zod 4).
- Annotate the tool with `readOnlyHint`/`destructiveHint`/`openWorldHint` so
  clients know studio runs an arbitrary command.
- Send a tool list cache hint (`ttlMs`/`cacheScope`) — studio's tool list is
  fixed by its command line and never changes while the process is alive.
- Require Node.js 20 or newer.

## [0.3.0]

- Switch to typescript for faster install, rather than downloading a binary after install

## [0.2.0]

- Support single curly braces `{name}` for required fields to simplify templating

## [0.0.2] - 2025-06-28

- Support for `--help` and `--debug` flags to studio-mcp itself.
- Improved handling of templated command arguments.

## [0.0.1] - 2025-06-28

### Added

- Initial TypeScript port of `studio` from Ruby
- MCP server implementation using @modelcontextprotocol/sdk
- Templating
- Support for required `{{name}}` and optional `[name]` string arguments
- Support for array arguments `[name...]`
- CLI interface compatible with Claude Desktop, Cursor, and other MCP clients
