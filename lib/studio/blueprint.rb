# frozen_string_literal: true

require "shellwords"

module Studio
  class Blueprint # rubocop:todo Style/Documentation
    attr_reader :argv, :base_command, :blueprint_args

    def self.argv(argv)
      new(argv)
    end

    def initialize(argv)
      @argv = argv.dup
      @base_command = @argv.shift
      @blueprint_args = parse_blueprint_args(@argv)
    end

    # Execute the command with provided arguments
    def execute(args = [])
      full_command = build_command(args)

      # Spawn the process and capture output
      stdout, stderr, status = Open3.capture3(*full_command)

      {
        stdout: stdout,
        stderr: stderr,
        exit_code: status.exitstatus,
        success: status.success?
      }
    rescue StandardError => e
      {
        stdout: "",
        stderr: "Studio error: #{e.message}",
        exit_code: 1,
        success: false
      }
    end

    # Get the tool name for MCP
    def tool_name
      @base_command.gsub(/[^a-zA-Z0-9_]/, "_")
    end

    # Get the tool description for MCP
    def tool_description
      if @blueprint_args.any?
        blueprint_desc = @blueprint_args.map { |arg| "{{#{arg[:name]}}}" }.join(" ")
        "Run the shell command `#{@base_command} #{blueprint_desc}`"
      else
        "Run the shell command `#{@base_command} [args]`"
      end
    end

    # Get the input schema for MCP tool
    def input_schema
      schema = {
        type: "object",
        properties: {}
      }

      if @blueprint_args.empty?
        schema[:properties][:args] = {
          type: "array",
          items: { type: "string" },
          description: "Additional command line arguments"
        }
      else
        # Add blueprint arguments to schema
        @blueprint_args.each do |blueprint_arg|
          schema[:properties][blueprint_arg[:name].to_sym] = {
            type: "string",
            description: blueprint_arg[:description]
          }
        end
      end

      schema
    end

    private

    def parse_blueprint_args(args)
      blueprint_args = []

      args.each do |arg|
        # First check for templates with descriptions: {{name#description}}
        arg.scan(/\{\{(\w+)\s*#\s*(.+?)\}\}/) do |name, description|
          blueprint_args << { name: name, description: description.strip }
        end

        # Then check for templates without descriptions: {{name}}
        # Only add if we haven't already found this name with a description
        arg.scan(/\{\{(\w+)\}\}/) do |name|
          name = name.first
          unless blueprint_args.any? { |ba| ba[:name] == name }
            blueprint_args << { name: name, description: "the {{#{name}}} arg" }
          end
        end
      end

      blueprint_args
    end

    # rubocop:todo Metrics/PerceivedComplexity
    # rubocop:todo Metrics/AbcSize
    # rubocop:todo Metrics/CyclomaticComplexity
    def build_command(args) # rubocop:todo Metrics/CyclomaticComplexity, Metrics/AbcSize, Metrics/PerceivedComplexity
      command_parts = [@base_command]

      # Add blueprint arguments if they exist
      if @blueprint_args.any?
        @argv.each do |arg|
          # Check if this argument contains any blueprint patterns
          processed_arg = arg.dup

          # Replace templates with descriptions: {{name#description}}
          processed_arg.gsub!(/\{\{(\w+)\s*#\s*.+?\}\}/) do |_match|
            name = ::Regexp.last_match(1)
            value = args[name.to_s] || args[name.to_sym]
            value.to_s # Convert to string in case of nil
          end

          # Replace templates without descriptions: {{name}}
          processed_arg.gsub!(/\{\{(\w+)\}\}/) do |_match|
            name = ::Regexp.last_match(1)
            value = args[name.to_s] || args[name.to_sym]
            value.to_s # Convert to string in case of nil
          end

          command_parts << processed_arg
        end
      else
        # Add original blueprint arguments
        command_parts.concat(@argv)
      end

      # Add additional args
      if args.is_a?(Hash) && (args["args"] || args[:args])
        additional_args = args["args"] || args[:args]
        command_parts.concat(additional_args) if additional_args
      elsif args.is_a?(Array)
        command_parts.concat(args)
      end

      command_parts.compact
    end
    # rubocop:enable Metrics/CyclomaticComplexity
    # rubocop:enable Metrics/AbcSize
    # rubocop:enable Metrics/PerceivedComplexity
  end
end
