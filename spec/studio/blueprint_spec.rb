# frozen_string_literal: true

require "spec_helper"

RSpec.describe Studio::Blueprint do
  describe "#initialize" do
    it "parses simple command" do
      blueprint = described_class.new(%w[echo hello])
      expect(blueprint.base_command).to eq("echo")
      expect(blueprint.blueprint_args).to be_empty
    end

    it "parses blueprinted command" do
      blueprint = described_class.new(["rails", "generate",
                                       "{{generator#A valid rails generator like scaffold, model}}"])
      expect(blueprint.base_command).to eq("rails")
      expect(blueprint.blueprint_args).to contain_exactly(
        { name: "generator", description: "A valid rails generator like scaffold, model" }
      )
    end

    it "parses blueprinted command with spaces around hash" do
      blueprint = described_class.new(["rails", "generate",
                                       "{{generator # A valid rails generator like scaffold, model}}"])
      expect(blueprint.base_command).to eq("rails")
      expect(blueprint.blueprint_args).to contain_exactly(
        { name: "generator", description: "A valid rails generator like scaffold, model" }
      )
    end

    it "parses blueprinted command without description" do
      blueprint = described_class.new(["echo", "{{text}}"])
      expect(blueprint.base_command).to eq("echo")
      expect(blueprint.blueprint_args).to contain_exactly(
        { name: "text", description: "the {{text}} arg" }
      )
    end

    it "parses mixed blueprints with and without descriptions" do
      blueprint = described_class.new(["command", "{{arg1#Custom description}}", "{{arg2}}"])
      expect(blueprint.base_command).to eq("command")
      expect(blueprint.blueprint_args).to contain_exactly(
        { name: "arg1", description: "Custom description" },
        { name: "arg2", description: "the {{arg2}} arg" }
      )
    end

    it "prioritizes explicit description over default" do
      blueprint = described_class.new(["echo", "{{text#Explicit description}}", "{{text}}"])
      expect(blueprint.base_command).to eq("echo")
      expect(blueprint.blueprint_args).to contain_exactly(
        { name: "text", description: "Explicit description" }
      )
    end
  end

  describe "#tool_name" do
    it "converts command to valid tool name" do
      blueprint = described_class.new(["git-flow"])
      expect(blueprint.tool_name).to eq("git_flow")
    end
  end

  describe "#tool_description" do
    it "generates description for simple command" do
      blueprint = described_class.new(["git"])
      expect(blueprint.tool_description).to eq("Run the shell command `git [args]`")
    end

    it "generates description for blueprinted command" do
      blueprint = described_class.new(["rails", "generate", "{{generator#A rails generator}}"])
      expect(blueprint.tool_description).to eq("Run the shell command `rails {{generator}}`")
    end
  end

  describe "#execute" do
    it "executes simple command" do
      blueprint = described_class.new(%w[echo hello])
      result = blueprint.execute

      expect(result[:success]).to be true
      expect(result[:stdout]).to eq("hello\n")
      expect(result[:exit_code]).to eq(0)
    end

    it "handles command errors" do
      blueprint = described_class.new(["false"])
      result = blueprint.execute

      expect(result[:success]).to be false
      expect(result[:exit_code]).to eq(1)
    end

    it "handles blueprint arguments with spaces" do
      blueprint = described_class.new(["echo", "{{text#text to echo}}"])
      result = blueprint.execute({ "text" => "Hello World" })

      expect(result[:success]).to be true
      expect(result[:stdout]).to eq("Hello World\n")
      expect(result[:exit_code]).to eq(0)
    end

    it "handles blueprint arguments without descriptions" do
      blueprint = described_class.new(["echo", "{{text}}"])
      result = blueprint.execute({ "text" => "Hello World" })

      expect(result[:success]).to be true
      expect(result[:stdout]).to eq("Hello World\n")
      expect(result[:exit_code]).to eq(0)
    end

    it "handles mixed blueprint with and without descriptions" do
      blueprint = described_class.new(["echo", "{{greeting#The greeting}}", "{{name}}"])
      result = blueprint.execute({ "greeting" => "Hello", "name" => "World" })

      expect(result[:success]).to be true
      expect(result[:stdout]).to eq("Hello World\n")
      expect(result[:exit_code]).to eq(0)
    end

    it "handles blueprint arguments with spaces in mixed content" do
      blueprint = described_class.new(["echo", "simon says {{text#text to echo}}"])
      result = blueprint.execute({ "text" => "Hello World" })

      expect(result[:success]).to be true
      expect(result[:stdout]).to eq("simon says Hello World\n")
      expect(result[:exit_code]).to eq(0)
    end

    it "handles blueprint arguments with special shell characters" do
      blueprint = described_class.new(["echo", "{{text#text to echo}}"])
      result = blueprint.execute({ "text" => "Hello & World; echo pwned" })

      expect(result[:success]).to be true
      expect(result[:stdout]).to eq("Hello & World; echo pwned\n")
      expect(result[:exit_code]).to eq(0)
    end

    it "handles multiple blueprint substitutions in one argument" do
      blueprint = described_class.new(["echo", "{{greeting#greeting}} {{name#name}}!"])
      result = blueprint.execute({ "greeting" => "Hello", "name" => "World" })

      expect(result[:success]).to be true
      expect(result[:stdout]).to eq("Hello World!\n")
      expect(result[:exit_code]).to eq(0)
    end

    it "handles blueprint in the middle of argument" do
      blueprint = described_class.new(["echo", "--message={{text#message content}}"])
      result = blueprint.execute({ "text" => "Hello World" })

      expect(result[:success]).to be true
      expect(result[:stdout]).to eq("--message=Hello World\n")
      expect(result[:exit_code]).to eq(0)
    end

    it "handles blueprint with prefix and suffix" do
      blueprint = described_class.new(["echo", "prefix-{{text#middle part}}-suffix"])
      result = blueprint.execute({ "text" => "Hello World" })

      expect(result[:success]).to be true
      expect(result[:stdout]).to eq("prefix-Hello World-suffix\n")
      expect(result[:exit_code]).to eq(0)
    end

    it "handles mixed blueprint and non-blueprint arguments" do
      blueprint = described_class.new(["echo", "static", "{{dynamic#dynamic content}}", "more-static"])
      result = blueprint.execute({ "dynamic" => "Hello World" })

      expect(result[:success]).to be true
      expect(result[:stdout]).to eq("static Hello World more-static\n")
      expect(result[:exit_code]).to eq(0)
    end

    it "preserves shell safety with complex blueprint values" do
      blueprint = described_class.new(["echo", "Result: {{text#text content}}"])
      result = blueprint.execute({ "text" => "$(rm -rf /); echo 'dangerous'" })

      expect(result[:success]).to be true
      expect(result[:stdout]).to eq("Result: $(rm -rf /); echo 'dangerous'\n")
      expect(result[:exit_code]).to eq(0)
    end
  end
end
