# Debugging

When your Studio MCP tools aren't working as expected, debugging helps identify the problem. This guide covers debugging techniques and common issues.

## Debug Mode

Studio provides built-in debug logging to help troubleshoot issues.

### Enable Debug Output

Add the `--debug` flag to enable debug logging:

```bash
studio --debug echo "{message}"
```

This prints debug information to stderr, including:
- Template parsing details
- Generated JSON schema
- Parameter validation
- Command execution

### Log to File

For persistent logs, use the `--log` flag:

```bash
studio --log /tmp/studio.log echo "{message}"
```

Or combine with debug:

```bash
studio --debug --log /tmp/studio.log echo "{message}"
```

## Debugging in MCP Clients

### Claude Desktop

Add debug flags to your configuration:

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

Then check `/tmp/studio-echo.log` for debug output.

**macOS Log Location:**
```bash
tail -f ~/Library/Application\ Support/Claude/logs/mcp*.log
```

### Cursor

Add debug logging:

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
        "/tmp/cursor-studio.log",
        "echo",
        "{message}"
      ]
    }
  }
}
```

### VSCode

Use workspace folder for logs:

```json
{
  "mcp.servers": {
    "echo": {
      "type": "stdio",
      "command": "npx",
      "args": [
        "-y",
        "@studio-mcp/studio",
        "--debug",
        "--log",
        "${workspaceFolder}/.studio.log",
        "echo",
        "{message}"
      ]
    }
  }
}
```

## Testing Commands Directly

Before adding to MCP clients, test commands in the terminal:

```bash
# Test basic execution
npx -y @studio-mcp/studio echo "{message}"

# This starts an MCP server on stdio
# You won't see output directly, but you can test the server works
```

For interactive testing, use the MCP Inspector:

```bash
npm install -g @modelcontextprotocol/inspector
mcp-inspector npx -y @studio-mcp/studio echo "{message}"
```

This opens a web interface to test your tool.

## Common Issues

### Issue: Tool not appearing in client

**Symptoms:**
- Tool doesn't show up in Claude Desktop/Cursor
- No errors visible

**Debugging steps:**

1. Check JSON syntax:
   ```bash
   # Use a JSON validator
   cat claude_desktop_config.json | jq .
   ```

2. Verify configuration file location:
   ```bash
   # macOS Claude Desktop
   cat ~/Library/Application\ Support/Claude/claude_desktop_config.json

   # Cursor
   cat ~/.cursor/mcp.json
   ```

3. Check client logs:
   ```bash
   # Claude Desktop (macOS)
   ls -lt ~/Library/Application\ Support/Claude/logs/
   tail -f ~/Library/Application\ Support/Claude/logs/mcp*.log
   ```

4. Restart the client completely

### Issue: Command fails to execute

**Symptoms:**
- Error messages in client
- Tool appears but doesn't work

**Debugging steps:**

1. Test the command directly:
   ```bash
   # Test without Studio
   echo "Hello"

   # Test with Studio
   npx -y @studio-mcp/studio --debug echo "{message}"
   ```

2. Check command availability:
   ```bash
   which echo
   which curl
   which git
   ```

3. Verify parameters:
   ```bash
   # Enable debug logging
   studio --debug --log /tmp/debug.log echo "{message}"

   # Check the log
   cat /tmp/debug.log
   ```

4. Test with MCP Inspector:
   ```bash
   mcp-inspector npx -y @studio-mcp/studio echo "{message}"
   ```

### Issue: Template parsing errors

**Symptoms:**
- "Invalid template" errors
- Parameters not recognized

**Debugging steps:**

1. Check template syntax:
   ```bash
   # Required: {name} or {{name}}
   studio echo "{message}"

   # Optional: [name]
   studio echo "[message]"

   # Array: [name...] or {name...}
   studio echo "[args...]"
   ```

2. Verify quote escaping:
   ```bash
   # In JSON config, escape quotes
   "args": ["echo", "{message}"]  # Correct
   ```

3. Test template parsing:
   ```bash
   studio --debug echo "{message}"
   # Check output for template parsing details
   ```

### Issue: Cursor spaces bug

**Symptoms:**
- Errors with descriptions containing spaces
- Works in Claude but not Cursor

**Solution:**

Use underscores instead of spaces in Cursor:

```json
// Claude Desktop (works)
"{message # A message to echo}"

// Cursor (use this)
"{message#message_to_echo}"
```

### Issue: npx command not found

**Symptoms:**
- "command not found: npx"
- Tool fails to start

**Debugging steps:**

1. Verify Node.js installation:
   ```bash
   node --version
   npm --version
   npx --version
   ```

2. Check PATH:
   ```bash
   echo $PATH
   which npx
   ```

3. Use full path:
   ```bash
   # Find npx location
   which npx
   # Example: /usr/local/bin/npx

   # Use in config
   "command": "/usr/local/bin/npx"
   ```

4. Install globally instead:
   ```bash
   npm install -g @studio-mcp/studio
   # Then use
   "command": "studio"
   ```

### Issue: Permission denied

**Symptoms:**
- "EACCES" or "Permission denied" errors
- Command can't execute

**Debugging steps:**

1. Check file permissions:
   ```bash
   ls -la /path/to/script
   chmod +x /path/to/script
   ```

2. Verify user permissions:
   ```bash
   # Test as current user
   echo "test" > /tmp/test.txt
   ```

3. Check command permissions:
   ```bash
   which <command>
   ls -la $(which <command>)
   ```

## Verbose Logging

For maximum debugging information:

```bash
studio --debug --log /tmp/studio-verbose.log echo "{message}" 2>&1 | tee /tmp/studio-stderr.log
```

This captures:
- Debug output to `/tmp/studio-verbose.log`
- Stderr to `/tmp/studio-stderr.log`
- Stdout (MCP protocol) to terminal

## MCP Inspector

The MCP Inspector is the best way to debug MCP tools:

### Install

```bash
npm install -g @modelcontextprotocol/inspector
```

### Use

```bash
mcp-inspector npx -y @studio-mcp/studio echo "{message}"
```

This opens a web interface where you can:
- See the generated tool schema
- Test tool execution
- View input/output
- Check errors and logs

### Example Session

```bash
# Start inspector
mcp-inspector npx -y @studio-mcp/studio \
  --debug \
  --log /tmp/studio.log \
  curl "https://api.github.com/{endpoint}"

# Open browser to http://localhost:5173
# Test the tool with different endpoints
# Check debug logs in /tmp/studio.log
```

## Debugging Checklist

When a tool isn't working:

- [ ] Verify JSON syntax in config file
- [ ] Check command works in terminal
- [ ] Test with `--debug` flag
- [ ] Review log file with `--log`
- [ ] Use MCP Inspector
- [ ] Check client-specific logs
- [ ] Verify command is in PATH
- [ ] Test with simple example first
- [ ] Check for Cursor spaces bug
- [ ] Restart MCP client completely

## Debug Log Example

A typical debug log shows:

```
[DEBUG] Parsing template: echo "{message # text to echo}"
[DEBUG] Tokens: [
  { type: "text", value: "echo" },
  { type: "field", name: "message", description: "text to echo", required: true }
]
[DEBUG] Generated schema: {
  "type": "object",
  "properties": {
    "message": {
      "type": "string",
      "description": "text to echo"
    }
  },
  "required": ["message"]
}
[DEBUG] Executing: echo "Hello, World!"
[DEBUG] Exit code: 0
```

## Getting Help

If you're still stuck:

1. Check the [GitHub Issues](https://github.com/studio-mcp/studio/issues)
2. Create a minimal reproduction case
3. Include debug logs
4. Specify your platform and client
5. Share your configuration (sanitized)

## Next Steps

- [Template Syntax](/reference/template-syntax) - Verify template correctness
- [Examples](/examples/basic) - See working examples
- [Architecture](/reference/architecture) - Understand how Studio works
