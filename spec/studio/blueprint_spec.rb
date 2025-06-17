# frozen_string_literal: true

require "spec_helper"
require "studio/blueprint"
require "studio/tool"

RSpec.describe Studio::Blueprint do
  describe "#argv" do
    it "parses simple command with explicit args" do
      blueprint = described_class.argv(%w[git status [args...]])
      expect(blueprint.base_command).to eq("git")
      expect(blueprint.input_schema).to eq(
        {
          type: "object",
          properties: {
            args: {
              type: "array",
              items: { type: "string" },
              description: "Additional command line arguments"
            }
          },
          required: ["args"]
        }
      )
    end

    it "parses simple command without args" do
      blueprint = described_class.argv(%w[git status])
      expect(blueprint.base_command).to eq("git")
      expect(blueprint.input_schema).to eq(
        {
          type: "object",
          properties: {}
        }
      )
    end

    it "parses blueprinted command" do
      blueprint = described_class.argv(["rails", "generate",
                                        "{{generator#A valid rails generator like scaffold, model}}"])
      expect(blueprint.base_command).to eq("rails")
      expect(blueprint.input_schema).to eq(
        {
          type: "object",
          properties: {
            generator: { type: "string", description: "A valid rails generator like scaffold, model" }
          },
          required: ["generator"]
        }
      )
    end

    it "parses blueprinted command with spaces between tokens" do
      blueprint = described_class.argv(["rails", "generate",
                                        "{{ generator # A valid rails generator like scaffold, model }}"])
      expect(blueprint.base_command).to eq("rails")
      expect(blueprint.input_schema).to eq(
        {
          type: "object",
          properties: {
            generator: { type: "string",
                         description: "A valid rails generator like scaffold, model" }
          },
          required: ["generator"]
        }
      )
    end

    it "parses blueprinted command without description" do
      blueprint = described_class.argv(["echo", "{{text}}"])
      expect(blueprint.base_command).to eq("echo")
      expect(blueprint.tool_name).to eq("echo")
      expect(blueprint.tool_description).to eq("Run the shell command `echo {{text}}`")
      expect(blueprint.input_schema).to eq({
                                             type: "object",
                                             properties: {
                                               text: { type: "string" }
                                             },
                                             required: ["text"]
                                           })
    end

    it "parses mixed blueprints with and without descriptions" do
      blueprint = described_class.argv(["command", "{{arg1#Custom description}}", "{{arg2}}"])
      expect(blueprint.base_command).to eq("command")
      expect(blueprint.tool_name).to eq("command")
      expect(blueprint.tool_description).to eq("Run the shell command `command {{arg1}} {{arg2}}`")
      expect(blueprint.input_schema).to eq({
                                             type: "object",
                                             properties: {
                                               arg1: { type: "string", description: "Custom description" },
                                               arg2: { type: "string" }
                                             },
                                             required: %w[arg1 arg2]
                                           })
    end

    it "prioritizes explicit description over default" do
      blueprint = described_class.argv(["echo", "{{text#Explicit description}}", "{{text}}"])
      expect(blueprint.base_command).to eq("echo")
      expect(blueprint.tool_name).to eq("echo")
      expect(blueprint.tool_description).to eq("Run the shell command `echo {{text}} {{text}}`")
      expect(blueprint.input_schema).to eq({
                                             type: "object",
                                             properties: {
                                               text: { type: "string", description: "Explicit description" }
                                             },
                                             required: ["text"]
                                           })
    end

    it "parses array arguments with explicit syntax" do
      blueprint = described_class.argv(["echo", "[files...#Files to echo]"])
      expect(blueprint.base_command).to eq("echo")
      expect(blueprint.tool_name).to eq("echo")
      expect(blueprint.tool_description).to eq("Run the shell command `echo [files...]`")
      expect(blueprint.input_schema).to eq({
                                             type: "object",
                                             properties: {
                                               files: {
                                                 type: "array",
                                                 items: { type: "string" },
                                                 description: "Files to echo"
                                               }
                                             },
                                             required: ["files"]
                                           })
    end

    it "parses array arguments without description" do
      blueprint = described_class.argv(["ls", "[paths...]"])
      expect(blueprint.base_command).to eq("ls")
      expect(blueprint.tool_name).to eq("ls")
      expect(blueprint.tool_description).to eq("Run the shell command `ls [paths...]`")
      expect(blueprint.input_schema).to eq({
                                             type: "object",
                                             properties: {
                                               paths: {
                                                 type: "array",
                                                 items: { type: "string" },
                                                 description: "Additional command line arguments"
                                               }
                                             },
                                             required: ["paths"]
                                           })
    end

    it "parses [name] without ellipsis as optional string field" do
      blueprint = described_class.argv(["echo", "[not-array]"])
      expect(blueprint.base_command).to eq("echo")
      expect(blueprint.tool_name).to eq("echo")
      expect(blueprint.tool_description).to eq("Run the shell command `echo [not-array]`")
      expect(blueprint.input_schema).to eq({
                                             type: "object",
                                             properties: {
                                               "not-array": { type: "string" }
                                             }
                                           })
    end

    it "parses mixed string and array arguments" do
      blueprint = described_class.argv(["command", "{{flag#Command flag}}", "[files...#Files to process]"])
      expect(blueprint.base_command).to eq("command")
      expect(blueprint.tool_name).to eq("command")
      expect(blueprint.tool_description).to eq("Run the shell command `command {{flag}} [files...]`")
      expect(blueprint.input_schema).to eq({
                                             type: "object",
                                             properties: {
                                               flag: { type: "string", description: "Command flag" },
                                               files: {
                                                 type: "array",
                                                 items: { type: "string" },
                                                 description: "Files to process"
                                               }
                                             },
                                             required: %w[flag files]
                                           })
    end
  end

  describe "#tool_name" do
    it "converts command to valid tool name" do
      blueprint = described_class.argv(["git-flow"])
      expect(blueprint.tool_name).to eq("git_flow")
    end
  end

  describe "#tool_description" do
    it "generates description for simple command without args" do
      blueprint = described_class.argv(["git"])
      expect(blueprint.tool_description).to eq("Run the shell command `git`")
    end

    it "generates description for simple command with explicit args" do
      blueprint = described_class.argv(["git", "[args...]"])
      expect(blueprint.tool_description).to eq("Run the shell command `git [args...]`")
    end

    it "generates description for blueprinted command" do
      blueprint = described_class.argv(["rails", "generate", "{{generator#A rails generator}}"])
      expect(blueprint.tool_description).to eq("Run the shell command `rails generate {{generator}}`")
    end
  end

  describe "#execute" do
    it "executes simple command" do
      blueprint = described_class.argv(%w[echo hello])
      full_command = blueprint.build_command_args
      output, success = Studio::Tool.execute(*full_command)

      expect(success).to be true
      expect(output).to eq("hello\n")
    end

    it "handles command errors" do
      blueprint = described_class.argv(["false"])
      full_command = blueprint.build_command_args
      output, success = Studio::Tool.execute(*full_command)

      expect(output).to eq("")
      expect(success).to be false
    end

    it "handles blueprint arguments with spaces" do
      blueprint = described_class.argv(["echo", "{{text#text to echo}}"])
      full_command = blueprint.build_command_args({ "text" => "Hello World" })
      output, success = Studio::Tool.execute(*full_command)

      expect(success).to be true
      expect(output).to eq("Hello World\n")
    end

    it "handles blueprint arguments without descriptions" do
      blueprint = described_class.argv(["echo", "{{text}}"])
      full_command = blueprint.build_command_args({ "text" => "Hello World" })
      output, success = Studio::Tool.execute(*full_command)

      expect(success).to be true
      expect(output).to eq("Hello World\n")
    end

    it "handles mixed blueprint with and without descriptions" do
      blueprint = described_class.argv(["echo", "{{greeting#The greeting}}", "{{name}}"])
      full_command = blueprint.build_command_args({ "greeting" => "Hello", "name" => "World" })
      output, success = Studio::Tool.execute(*full_command)

      expect(success).to be true
      expect(output).to eq("Hello World\n")
    end

    it "handles blueprint arguments with spaces in mixed content" do
      blueprint = described_class.argv(["echo", "simon says {{text#text for simon to say}}"])
      full_command = blueprint.build_command_args({ "text" => "Hello World" })
      output, success = Studio::Tool.execute(*full_command)

      expect(success).to be true
      expect(output).to eq("simon says Hello World\n")
    end

    it "handles blueprint arguments with special shell characters" do
      blueprint = described_class.argv(["echo", "{{text#text to echo}}"])
      full_command = blueprint.build_command_args({ "text" => "Hello & World; echo pwned" })
      output, success = Studio::Tool.execute(*full_command)

      expect(success).to be true
      expect(output).to eq("Hello & World; echo pwned\n")
    end

    it "handles multiple blueprint substitutions in one argument" do
      blueprint = described_class.argv(["echo", "{{greeting}} {{name}}!"])
      full_command = blueprint.build_command_args({ "greeting" => "Hello", "name" => "World" })
      output, success = Studio::Tool.execute(*full_command)

      expect(success).to be true
      expect(output).to eq("Hello World!\n")
    end

    it "handles blueprint in the middle of argument" do
      blueprint = described_class.argv(["echo", "--message={{text#message content}}"])
      full_command = blueprint.build_command_args({ "text" => "Hello World" })
      output, success = Studio::Tool.execute(*full_command)

      expect(success).to be true
      expect(output).to eq("--message=Hello World\n")
    end

    it "handles blueprint with prefix and suffix" do
      blueprint = described_class.argv(["echo", "prefix-{{text#middle part}}-suffix"])
      full_command = blueprint.build_command_args({ "text" => "Hello World" })
      output, success = Studio::Tool.execute(*full_command)

      expect(success).to be true
      expect(output).to eq("prefix-Hello World-suffix\n")
    end

    it "handles mixed blueprint and non-blueprint arguments" do
      blueprint = described_class.argv(["echo", "static", "{{dynamic#dynamic content}}", "more-static"])
      full_command = blueprint.build_command_args({ "dynamic" => "Hello World" })
      output, success = Studio::Tool.execute(*full_command)

      expect(success).to be true
      expect(output).to eq("static Hello World more-static\n")
    end

    it "preserves shell safety with complex blueprint values" do
      blueprint = described_class.argv(["echo", "Result: {{text#text content}}"])
      full_command = blueprint.build_command_args({ "text" => "$(echo 'dangerous'); echo 'safe'" })
      output, success = Studio::Tool.execute(*full_command)

      expect(success).to be true
      expect(output).to eq("Result: $(echo 'dangerous'); echo 'safe'\n")
    end

    it "handles array arguments" do
      blueprint = described_class.argv(["echo", "[files...#Files to echo]"])
      full_command = blueprint.build_command_args({ "files" => ["file1.txt", "file2.txt"] })
      output, success = Studio::Tool.execute(*full_command)

      expect(success).to be true
      expect(output).to eq("file1.txt file2.txt\n")
    end

    it "handles empty array arguments" do
      blueprint = described_class.argv(["echo", "hello", "[files...]"])
      full_command = blueprint.build_command_args({ "files" => [] })
      output, success = Studio::Tool.execute(*full_command)

      expect(success).to be true
      expect(output).to eq("hello\n")
    end

    it "handles mixed string and array arguments" do
      blueprint = described_class.argv(["echo", "{{prefix#Prefix text}}", "[files...#Files to list]"])
      full_command = blueprint.build_command_args({ "prefix" => "Files:", "files" => ["a.txt", "b.txt"] })
      output, success = Studio::Tool.execute(*full_command)

      expect(success).to be true
      expect(output).to eq("Files: a.txt b.txt\n")
    end

    it "handles optional arguments when provided" do
      blueprint = described_class.argv(["echo", "hello", "[name]"])
      full_command = blueprint.build_command_args({ "name" => "world" })
      output, success = Studio::Tool.execute(*full_command)

      expect(success).to be true
      expect(output).to eq("hello world\n")
    end

    it "skips optional arguments when not provided" do
      blueprint = described_class.argv(["echo", "hello", "[name]"])
      full_command = blueprint.build_command_args({})
      output, success = Studio::Tool.execute(*full_command)

      expect(success).to be true
      expect(output).to eq("hello\n")
    end

    it "handles mixed required and optional arguments" do
      blueprint = described_class.argv(["echo", "{{required#Required text}}", "[optional#Optional text]"])

      # With both provided
      full_command = blueprint.build_command_args({ "required" => "hello", "optional" => "world" })
      output, success = Studio::Tool.execute(*full_command)
      expect(success).to be true
      expect(output).to eq("hello world\n")

      # With only required provided
      full_command = blueprint.build_command_args({ "required" => "hello" })
      output, success = Studio::Tool.execute(*full_command)
      expect(success).to be true
      expect(output).to eq("hello\n")
    end
  end
end
