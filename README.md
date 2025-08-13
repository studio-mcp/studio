# Studio MCP

Make any CLI into an AI tool with `studio` MCP server.

> Bright Studio Apt – Just enough, nothing more – $0/mo OBO.
>
> Text 404 for more details!

## What's Included?

`studio` is the simplest possible way to add CLI tools that your AI Agent can use right now. Built on the [Model Context Protocol](https://modelcontextprotocol.io/), `studio` spawns a single tool mcp server with your command.

The tool turns everything after `studio` command into an MCP tool that can be called by Cursor, Claude, etc.

`studio` is great for patching CLIs into Claude, debugging MCPs or providing custom scripts to your LLM without having to deal with MCP.

It uses a very simple Mustache-like template syntax to provide inputs and descriptions telling the LLM how to use your command.

```sh
$ npx -y @studio-mcp/studio command "{ required_argument # Description of argument }" "[optional_args... # any array of arguments]"
```

`studio` turns this into an input schema for the MCP tool so that tool calls know what to send:

```json
{
  "type": "object",
  "properties": {
    "required_argument": { "type": "string", "description": "Description of argument" },
    "optional_args": { "type": "array", "items": { "type": "string" }, "description": "any array of arguments" }
  },
  "required": ["required_argument"]
}
```

You can run almost any command, but you might need to put the full path for scripts or commands installed via other package managers. Just run `which cmd` to find it's full path.

Since you're just renting the place, please be a good tenant and don't `rm -rf` anything.

## Move-In

These install instructions are like my lease agreement: full of gotchas.
Have your lawyer read it over. (You do have a lawyer right?)

You can install to your system with `npm`, use `npx` directly, or install with `go install github.com/studio-mcp/studio@latest`

```sh
npm install -g @studio-mcp/studio
```

Or download directly from [GitHub Releases](https://github.com/studio-mcp/studio/releases/latest) and add to your PATH yourself.

## Unpack (it's an apartment metaphor)

Most MCPs don't run in your shell environment 😭 You'll probably need to add the full path somewhere.

We'll use the MacOS `say` command as an example command. If you're not on a Mac, use `echo` (it's worse than useless, but it's easy to understand).

### Claude Desktop

Go to the Claude Desktop settings and click Developer > Edit Config.
It should open your Claude Desktop MCP configuration. (e.g. `~/Library/Application Support/Claude/claude_desktop_config.json` on Mac):

```json
{
  "mcpServers": {
    "say": {
      "command": "npx",
      "args": ["-y", "@studio-mcp/studio", "say", "-v", "siri", "{speech # A concise message to say outloud}"]
    },
  }
}
```

### Cursor

Add to your `~/.cursor/mcp.json` (in your home or project directory) or go to Tools section of the Cursor UI.

*Note:* A bug in Cursor breaks args with spaces.

```json
{
  "mcpServers": {
    "say": {
      "command": "npx",
      "args": ["-y", "@studio-mcp/studio", "say", "-v", "siri", "{speech#say_outloud}"]
    },
  }
}
```

### VSCode

It's a lot of the same here. Don't get confused betwee studio and stdio (that's how I got the name).

```json
{
  "mcp": {
    "servers": {
      "echo": {
        "type": "stdio",
        "command": "npx",

        "args": ["-y", "@studio-mcp/studio", "echo", "{text#What do you want to say?}"]
      }
    }
  }
}
```

## Blueprint Syntax

Studio uses blueprints (templates) to keep your studio tidy.

```bash
studio say -v "{voice# Choose your Mac say voice}" "[args...#Any additional args]"
```

This creates a Studio server with two arguments: `voice` and `args`.
Everything after the `#` will be used as the description for the LLM to understand.

Blueprint templates are a modified mustache format with descriptions: `{name # description}` but they also add shell like optional `[--flag]` booleans, `[optional # an optional string]` and `[args... # array with 0 or more strings]` for additional args:

- `{name}`: Required string argument
- `[name]`: Optional string argument
- `[name...]`: Optional array argument (spreads as multiple command line args)
- `[--flag]`: Optional boolean named `flag` that prints `--flag` only when true.
- `{name...}`: Required array (1 or more arguments required).

Inside a tag, there is a name and description:

- `name`: The argument name that will be shown in the MCP tool schema. Only letter numbers and underscores (dashes and underscores are interchangeable, case-insensitive).
- `description`: A description of what the argument should contain. Reads everything after the `#` to the end of the template tag.

Note: Double curly braces `{{name}}` are still supported for backward compatibility, but single braces `{name}` are preferred and used throughout the docs.

#### What about {cool_template_feature: string /[A-Z]+/ # Fancy tags}?

This is a simple studio, not one of those fancy 1 bedroom flats.

Maybe the landlord will get around to it at some point (but your rent will go up).

## Utilities Included

To build and test locally:

```bash
make
make test
studio echo "{text # what you want said to you}"
```

### Did something break?

The landlord _definitely_ takes care of the place...

- more than none tests
- files! lots of 'em!
- maybe even some test coverage
- you still need Proof of Renters Insurance

Uncovered portions are tenant's responsibility. (no one appreciates how hard it is for us landlords)

## Home Is Where You Make It

This is your studio too. Bugs, bedbugs, features, ideas? Swing by the repo during open-house:

🏠 https://github.com/studio-mcp/studio

## Lease Terms: MIT

Move in under standard terms, no fine print. Full text here: [MIT License](https://opensource.org/licenses/MIT).
