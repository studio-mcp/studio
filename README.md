# Studio

Make any CLI into a single tool MCP server.

> Bright Studio Apt – No Walls, No Rules – $0/mo OBO.
>
> Wired and wild. No questions, no backups. rm -rf compatible. Cash only. Basement-ish. BYO everything. Just enough, nothing more.
>
> Text 404 to /dev/null for more details!

## Install

```sh
$ gem install studio
```

### Claude Code

```sh
#                <serv-name> studio <cmd> <arguments or blueprints>
$ claude mcp add echo-server studio echo "{{text#What do you want to say?}}"
```

```sh
$ claude mcp add git-log studio git log --one-line -n 20 "{{branch}}"
```

### Cursor

Make sure `studio` is installed, then

[![Install MCP Server](https://cursor.com/deeplink/mcp-install-dark.svg)](https://cursor.com/install-mcp?name=echo&config=eyJjb21tYW5kIjoic3R1ZGlvIHt7dGV4dCNXaGF0IGRvIHlvdSB3YW50IHRvIHNheT99fSJ9)

Add to your `.cursor/mcp.json` manually:

```json
{
  "mcpServers": {
    "echo": {
      "command": "studio",
      "args": ["echo", "{{text#What do you want to say?}}"]
    },
    "rails": {
      "command": "studio",
      "args": ["rails", "generate", "{{generator # A valid rails generator like scaffold, model}}", "[args#Any extra args needed]"]
    }
  }
}
```

### VSCode

_Pay attention, one says `stdio`, the other says `studio`! (uh oh, what have I done?)_

```json
    "mcp": {
        "servers": {
            "echo": {
                "type": "stdio",
                "command": "studio",
                "args": ["echo", "{{text#What do you want to say?}}"]
            }
        }
    }
```

### Do you speak MCP?

```bash
$ studio echo
```

Now you're speaking to the mcp server, so say something smart.

```json
// ping
{"jsonrpc":"2.0","id":"1","method":"ping"}
```

or just play ping pong...

**initialize**

```json
{"jsonrpc": "2.0","id": "2","method": "initialize","params": {"protocolVersion": "2024-11-05","capabilities": {},"clientInfo": {"name": "test-client","version": "1.0.0"}}}
```

**tool list**

```json
{"jsonrpc":"2.0","id":"3","method":"tools/list"}
```

**tool call**

```json
{"jsonrpc": "2.0","id": "4","method": "tools/call","params": {"name": "echo","arguments": {"args": ["hello", "world"]}}}
```

### Blueprint

Studio uses blueprints (templates) to keep your studio tidy.

```bash
studio rails generate "{{generator # A valid rails generator like scaffold, model}}" "[args# Any additional args needed for the generator]
```

This creates an Studio server with two arguments: `generator` and `args`.
Everything after the `#` will be used as the description for the LLM to understand.

Blueprints use the format: `{{name # description}}` and `[name # description]` for string and array arguments.

- `name`: The argument name that will be exposed in the MCP tool schema. May not contain spaces.
- `description`: A description of what the argument should contain. May contain spaces.

This is a simple studio, not one of those fancy 1 bedroom flats.
Blueprint types, flags, validation??? The landlord will probably upgrade the place for free eventually... right?

## Development

**Keys to the City**. Run `bin/setup` to move in. Run `bin/rake` for a quick inspection.

Bump `version.rb`, tag it `vX.Y.Z`, push to GitHub. Trusted publishing will be happy to receive your rent check.

To install this gem onto your local machine, run `bundle exec rake install`.

## Home Is Where You Make It

This is your studio too. Bugs, features, ideas? Swing by the repo:

🏠 https://github.com/martinemde/studio

## Lease Terms: MIT

Move in under standard terms, no fine print. Full text here: [MIT License](https://opensource.org/licenses/MIT).
