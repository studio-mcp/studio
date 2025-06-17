# frozen_string_literal: true

require "shellwords"

module Studio
  class Blueprint # rubocop:todo Style/Documentation
    attr_reader :argv, :base_command, :compiled_args

    def self.argv(argv)
      new(argv)
    end

    def initialize(argv)
      @argv = argv.dup
      @base_command = @argv.shift
      @properties = {}
      @compiled_args = compile_template_args(@argv)
    end

    # Build the command with provided arguments
    def build_command_args(args = [])
      build_command(args)
    end

    # Get the tool name for MCP
    def tool_name
      @base_command.gsub(/[^a-zA-Z0-9_]/, "_")
    end

    # Get the tool description for MCP
    def tool_description
      "Run the shell command `#{format_command}`"
    end

    # Build display format directly from compiled args
    # rubocop:todo Metrics/PerceivedComplexity
    # rubocop:todo Metrics/AbcSize
    def format_command # rubocop:todo Metrics/AbcSize, Metrics/PerceivedComplexity
      display_parts = [@base_command]

      @compiled_args.each do |compiled_arg|
        processed_arg = +""
        compiled_arg.each do |part|
          case part
          when String
            processed_arg << part
          when Hash
            if part[:type] == "array"
              processed_arg = "[#{part[:name]}...]"
              break
            elsif part[:required]
              processed_arg << "{{#{part[:name]}}}"
            else
              processed_arg << "[#{part[:name]}]"
            end
          end
        end
        display_parts << processed_arg
      end

      display_parts = display_parts.map { |part| part.include?(" ") ? %("#{part}") : part }
      display_parts.join(" ")
    end
    # rubocop:enable Metrics/AbcSize
    # rubocop:enable Metrics/PerceivedComplexity

    # Get the input schema for MCP tool
    def input_schema
      # Build properties without the internal required field
      clean_properties = @properties.transform_values do |prop|
        prop.except(:required)
      end

      schema = {
        type: "object",
        properties: clean_properties
      }

      required_fields = @properties.select { |_k, v| v[:required] }.keys.map(&:to_s)
      schema[:required] = required_fields if required_fields.any?

      schema
    end

    private

    # Lexer: tokenizes a single shell word into tokens;
    # [:text, content]
    # [:field, opts]
    # rubocop:todo Metrics/PerceivedComplexity
    # rubocop:todo Metrics/AbcSize
    def lex(word, &) # rubocop:todo Metrics/AbcSize, Metrics/PerceivedComplexity
      return enum_for(:lex, word) unless block_given?

      # Split on template boundaries while capturing the delimiters
      parts = word.split(/(\{\{[^}]*\}\}|\[[^\]]*\])/)

      parts.each do |part|
        next if part.empty?

        if part.match?(/^\{\{.*\}\}$/)
          # Extract field content without the braces
          field_content = part[2...-2]
          name, description = field_content.split("#", 2)
          name = name.strip
          type = "string"
          required = true
          yield [:field, { name: name, description: description, content: part, required: required, type: type }]
        elsif part.match?(/^\[.*\]$/) && part == word.strip # only accept arrays that are the entire shell word
          field_content = part[1...-1]
          name, description = field_content.split("#", 2)
          name = name.strip

          if name.end_with?("...")
            name = name[0...-3]
            type = "array"
            required = true
          else
            type = "string"
            required = false
          end

          yield [:field, { name: name, description: description, content: part, required: required, type: type }]
        else
          yield [:text, { content: part }]
        end
      end
    end
    # rubocop:enable Metrics/AbcSize
    # rubocop:enable Metrics/PerceivedComplexity

    def compile_template_args(args)
      args.map { |arg| compile_template_arg(arg) }
    end

    def compile_template_arg(arg)
      parts = []
      lex(arg) do |token_type, opts|
        case token_type
        when :text
          parts << opts[:content]
        when :field, :array
          add_property(opts)
          parts << opts
        end
      end
      parts
    end

    # update properties, allowing that a variable may be used twice, not always with the description
    def add_property(opts) # rubocop:todo Metrics/AbcSize
      name = opts[:name]
      description = opts[:description]
      type = opts[:type]
      required = opts[:required]

      name = name.to_sym

      @properties[name] ||= {}
      @properties[name][:type] = type
      @properties[name][:required] = required # Keep this for internal tracking
      if type == "array"
        @properties[name][:items] = { type: "string" }
        @properties[name][:description] = description&.strip || "Additional command line arguments"
      elsif description
        @properties[name][:description] = description.strip
      end
    end

    # rubocop:todo Metrics/PerceivedComplexity
    # rubocop:todo Metrics/AbcSize
    def build_command(args) # rubocop:todo Metrics/AbcSize, Metrics/PerceivedComplexity
      command_parts = [@base_command]

      # Process compiled arguments
      @compiled_args.each do |compiled_arg| # rubocop:todo Metrics/BlockLength
        processed_arg = +""
        skip_arg = false

        compiled_arg.each do |part|
          case part
          when String
            processed_arg << part
          when Hash
            if part[:type] == "array"
              # For array arguments, we expect them to be passed as the whole argument
              value = args[part[:name]] || args[part[:name].to_sym] || []
              value = Array(value) # Ensure it's an array
              command_parts.concat(value) if value.any?
              processed_arg = nil # Skip adding this as a single argument
              break
            else
              value = args[part[:name]] || args[part[:name].to_sym]
              if value.nil? || value == ""
                if part[:required] # rubocop:todo Metrics/BlockNesting
                  processed_arg << ""
                else
                  # Skip the entire argument if it's optional and not provided
                  skip_arg = true
                  break
                end
              else
                processed_arg << value.to_s
              end
            end
          end
        end

        command_parts << processed_arg if processed_arg && !skip_arg
      end

      command_parts.compact
    end
    # rubocop:enable Metrics/AbcSize
    # rubocop:enable Metrics/PerceivedComplexity
  end
end
