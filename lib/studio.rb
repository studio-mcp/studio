# frozen_string_literal: true

require "open3"
require "mcp"
require "mcp/transports/stdio"

require_relative "studio/version"
require_relative "studio/blueprint"
require_relative "studio/tool"

module Studio # rubocop:todo Style/Documentation
  class Error < StandardError; end

  def self.serve(argv)
    if argv.empty?
      warn "Usage: studio <command> [blueprint_args...]"
      exit 1
    end

    blueprint = Blueprint.new(argv)
    tool_class = Tool.for_blueprint(blueprint)

    server = MCP::Server.new(
      name: "studio_server",
      tools: [tool_class]
    )

    transport = MCP::Transports::StdioTransport.new(server)
    transport.open
  rescue StandardError => e
    warn "Studio error: #{e.message}"
    exit 1
  end
end
