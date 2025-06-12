# Studio

Make any CLI into a single tool MCP server.

> Bright Studio Apt – No Walls, No Rules – 800/mo OBO.
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
$ claude mcp add git-log studio git log --one-line -n 20 --ref "{{branch}}"
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
      "args": ["rails", "generate", "{{generator # A valid rails generator like scaffold, model}}"]
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

**ping**

```json
{"jsonrpc":"2.0","id":"1","method":"ping"}
```

or not... just play ping pong if you wish.

 **initialize**

```json
{"jsonrpc": "2.0","id": "2","method": "initialize","params": {"protocolVersion": "2024-11-05","capabilities": {},"clientInfo": {"name": "test-client","version": "1.0.0"}}}
```

that's better...

 **list tools**

```json
{"jsonrpc":"2.0","id":"3","method":"tools/list"}
```

Oh look, a hammer...

**use tool**

```json
{"jsonrpc": "2.0","id": "4","method": "tools/call","params": {"name": "echo","arguments": {"args": ["hello", "world"]}}}
```

### Blueprints

Use blueprints (templates) to keep your studio tidy.

```bash
studio rails generate "{{generator # A valid rails generator like scaffold, model}}"
```

This creates an Studio server with one argument: `generator`.
Everything after the `#` will be used as the description of the argument.

Blueprints use the format: `{{name # description}}`.

- `name`: The argument name that will be exposed in the MCP tool schema
- `description`: A description of what the argument should contain. May contain spaces.

This is a simple studio, not one of those fancy 1 bedroom flats. Blueprint types, flags, arrays? The landlord will upgrade the place for free eventually, right? Right?

## Development

**Keys to the City**. Run `bin/setup` to move in. Run `bin/rake` for a quick inspection.

Bump `version.rb`, tag it `vX.Y.Z`, push to GitHub. Trusted publishing will be happy to receive your rent check.

After checking out the repo, run `bin/setup` to install dependencies. Then, run `rake spec` to run the tests. You can also run `bin/console` for an interactive prompt that will allow you to experiment.

To install this gem onto your local machine, run `bundle exec rake install`. To release a new version, update the version number in `version.rb`, tag the release vX.Y.Z, then push the tag to github. Trusted publishing will do the rest.

## Home Is Where You Make It

This is your studio too. Bugs, features, ideas? Swing by the repo:
🏠 https://github.com/martinemde/studio

## Lease Terms: MIT

Move in under the standard terms, no fine print. Full text here: [MIT License](https://opensource.org/licenses/MIT).
