# frozen_string_literal: true

require "mcp"
require "open3"

module Studio
  class Tool # rubocop:todo Style/Documentation
    # Execute the command with provided arguments
    def self.execute(*command)
      output, status = Open3.capture2e(*command)
      [output, status.success?]
    rescue StandardError => e
      ["Studio error: #{e.message}", false]
    end

    def self.for_blueprint(blueprint)
      Class.new(MCP::Tool) do
        tool_name blueprint.tool_name
        description blueprint.tool_description
        input_schema blueprint.input_schema

        define_singleton_method :call do |server_context:, **arguments| # rubocop:todo Lint/UnusedBlockArgument
          full_command = blueprint.build_command_args(arguments)
          output, success = Tool.execute(*full_command)
          MCP::Tool::Response.new([{ type: "text", text: output }], !success)
        end
      end
    end
  end
end
