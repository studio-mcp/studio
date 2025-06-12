# frozen_string_literal: true

require "mcp"

module Studio
  class Tool # rubocop:todo Style/Documentation
    # rubocop:todo Metrics/AbcSize
    def self.for_blueprint(blueprint) # rubocop:todo Metrics/AbcSize
      Class.new(MCP::Tool) do
        tool_name blueprint.tool_name
        description blueprint.tool_description
        input_schema blueprint.input_schema

        define_singleton_method :call do |server_context:, **arguments| # rubocop:todo Lint/UnusedBlockArgument
          result = blueprint.execute(arguments)

          # Log stderr if present
          warn "[Studio] #{result[:stderr]}" unless result[:stderr].empty?

          if result[:success]
            MCP::Tool::Response.new([{
                                      type: "text",
                                      text: result[:stdout]
                                    }])
          else
            MCP::Tool::Response.new([{
                                      type: "text",
                                      text: result[:stdout]
                                    }], true)
          end
        end
      end
    end
    # rubocop:enable Metrics/AbcSize
  end
end
