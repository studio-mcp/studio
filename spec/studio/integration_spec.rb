# frozen_string_literal: true

require "spec_helper"
require "open3"
require "json"
require "shellwords"

RSpec.describe "Studio MCP Server Integration" do # rubocop:todo RSpec/DescribeClass
  # Helper method to send MCP request and get response
  def send_mcp_request(command_args, request)
    request_json = JSON.generate(request)
    # Properly escape command arguments for shell
    escaped_args = command_args.map { |arg| Shellwords.escape(arg) }.join(" ")
    stdout, stderr, status = Open3.capture3(
      "bundle exec ruby exe/studio #{escaped_args}",
      stdin_data: request_json
    )

    expect(status).to be_success, "Server failed: #{stderr}"
    expect(stdout).not_to be_empty, "No response received"

    JSON.parse(stdout)
  end

  # Helper method to test the serve method by spawning a subprocess
  def test_serve_with_args(args, requests)
    # Start the server process
    escaped_args = args.map { |arg| Shellwords.escape(arg) }.join(" ")
    stdin, stdout, stderr, wait_thr = Open3.popen3("bundle exec ruby exe/studio #{escaped_args}")

    responses = []
    begin
      # Send each request and collect responses
      requests.each do |request|
        request_json = JSON.generate(request)
        stdin.puts(request_json)
        stdin.flush

        # Read response
        response_line = stdout.gets&.strip
        responses << JSON.parse(response_line) if response_line && !response_line.empty?
      end
    ensure
      # Clean up
      stdin.close
      stdout.close
      stderr.close
      begin
        Process.kill("TERM", wait_thr.pid)
      rescue StandardError
        nil
      end
      wait_thr.join
    end

    responses
  end

  describe "Basic MCP Protocol" do
    it "responds to ping requests" do
      request = {
        "jsonrpc" => "2.0",
        "id" => "1",
        "method" => "ping"
      }

      response = send_mcp_request(%w[echo hello], request)

      expect(response["jsonrpc"]).to eq("2.0")
      expect(response["id"]).to eq("1")
      expect(response["result"]).to eq({})
    end

    it "responds to initialize requests" do
      request = {
        "jsonrpc" => "2.0",
        "id" => "1",
        "method" => "initialize",
        "params" => {
          "protocolVersion" => "2024-11-05",
          "capabilities" => {},
          "clientInfo" => {
            "name" => "test-client",
            "version" => "1.0.0"
          }
        }
      }

      response = send_mcp_request(%w[echo hello], request)

      expect(response["jsonrpc"]).to eq("2.0")
      expect(response["id"]).to eq("1")
      expect(response["result"]).to have_key("protocolVersion")
      expect(response["result"]).to have_key("capabilities")
      expect(response["result"]).to have_key("serverInfo")
      expect(response["result"]["serverInfo"]["name"]).to eq("studio_server")
    end
  end

  describe "Tools functionality" do
    context "with simple echo command" do
      it "lists available tools" do
        request = {
          "jsonrpc" => "2.0",
          "id" => "2",
          "method" => "tools/list"
        }

        response = send_mcp_request(%w[echo hello], request)

        expect(response["jsonrpc"]).to eq("2.0")
        expect(response["id"]).to eq("2")
        expect(response["result"]["tools"]).to be_an(Array)
        expect(response["result"]["tools"].length).to eq(1)

        tool = response["result"]["tools"].first
        expect(tool["name"]).to eq("echo")
        expect(tool["description"]).to eq("Run the shell command `echo [args]`")
        expect(tool["inputSchema"]["type"]).to eq("object")
        expect(tool["inputSchema"]["properties"]).to have_key("args")
      end

      it "executes simple echo command" do
        request = {
          "jsonrpc" => "2.0",
          "id" => "3",
          "method" => "tools/call",
          "params" => {
            "name" => "echo",
            "arguments" => {
              "args" => %w[hello world]
            }
          }
        }

        response = send_mcp_request(%w[echo hello], request)

        expect(response["jsonrpc"]).to eq("2.0")
        expect(response["id"]).to eq("3")
        expect(response["result"]["content"]).to be_an(Array)
        expect(response["result"]["content"].first["type"]).to eq("text")
        expect(response["result"]["content"].first["text"]).to eq("hello hello world\n")
        expect(response["result"]["isError"]).to be false
      end
    end

    context "with blueprinted echo command" do
      it "lists blueprinted tool with proper schema" do
        request = {
          "jsonrpc" => "2.0",
          "id" => "4",
          "method" => "tools/list"
        }

        response = send_mcp_request(["echo", "{{text#the text to echo}}"], request)

        expect(response["jsonrpc"]).to eq("2.0")
        expect(response["id"]).to eq("4")

        tool = response["result"]["tools"].first
        expect(tool["name"]).to eq("echo")
        expect(tool["description"]).to eq("Run the shell command `echo {{text}}`")

        schema = tool["inputSchema"]
        expect(schema["properties"]).to have_key("text")
        expect(schema["properties"]["text"]["type"]).to eq("string")
        expect(schema["properties"]["text"]["description"]).to eq("the text to echo")
        expect(schema["properties"]).not_to have_key("args")
      end

      it "executes blueprinted command" do
        request = {
          "jsonrpc" => "2.0",
          "id" => "5",
          "method" => "tools/call",
          "params" => {
            "name" => "echo",
            "arguments" => {
              "text" => "Hello Blueprint!"
            }
          }
        }

        response = send_mcp_request(["echo", "{{text#the text to echo}}"], request)

        expect(response["jsonrpc"]).to eq("2.0")
        expect(response["id"]).to eq("5")
        expect(response["result"]["content"].first["type"]).to eq("text")
        expect(response["result"]["content"].first["text"]).to eq("Hello Blueprint!\n")
        expect(response["result"]["isError"]).to be false
      end

      it "executes blueprinted command with additional args" do
        request = {
          "jsonrpc" => "2.0",
          "id" => "6",
          "method" => "tools/call",
          "params" => {
            "name" => "echo",
            "arguments" => {
              "text" => "Hello",
              "args" => %w[World from args]
            }
          }
        }

        response = send_mcp_request(["echo", "{{text#the text to echo}}"], request)

        expect(response["jsonrpc"]).to eq("2.0")
        expect(response["id"]).to eq("6")
        expect(response["result"]["content"].first["type"]).to eq("text")
        expect(response["result"]["content"].first["text"]).to eq("Hello World from args\n")
        expect(response["result"]["isError"]).to be false
      end

      it "handles blueprint arguments with spaces" do
        request = {
          "jsonrpc" => "2.0",
          "id" => "7",
          "method" => "tools/call",
          "params" => {
            "name" => "echo",
            "arguments" => {
              "text" => "Hello World with spaces"
            }
          }
        }

        response = send_mcp_request(["echo", "{{text#the text to echo}}"], request)

        expect(response["jsonrpc"]).to eq("2.0")
        expect(response["id"]).to eq("7")
        expect(response["result"]["content"].first["type"]).to eq("text")
        expect(response["result"]["content"].first["text"]).to eq("Hello World with spaces\n")
        expect(response["result"]["isError"]).to be false
      end

      it "handles mixed blueprint with spaces" do
        request = {
          "jsonrpc" => "2.0",
          "id" => "8",
          "method" => "tools/call",
          "params" => {
            "name" => "echo",
            "arguments" => {
              "text" => "Hello World"
            }
          }
        }

        response = send_mcp_request(["echo", "simon says {{text#the text to echo}}"], request)

        expect(response["jsonrpc"]).to eq("2.0")
        expect(response["id"]).to eq("8")
        expect(response["result"]["content"].first["type"]).to eq("text")
        expect(response["result"]["content"].first["text"]).to eq("simon says Hello World\n")
        expect(response["result"]["isError"]).to be false
      end

      it "handles blueprint definitions with spaces around hash" do
        request = {
          "jsonrpc" => "2.0",
          "id" => "9",
          "method" => "tools/call",
          "params" => {
            "name" => "echo",
            "arguments" => {
              "text" => "Hello World"
            }
          }
        }

        response = send_mcp_request(["echo", "{{text # the text to echo}}"], request)

        expect(response["jsonrpc"]).to eq("2.0")
        expect(response["id"]).to eq("9")
        expect(response["result"]["content"].first["type"]).to eq("text")
        expect(response["result"]["content"].first["text"]).to eq("Hello World\n")
        expect(response["result"]["isError"]).to be false
      end

      it "handles blueprints without descriptions" do
        request = {
          "jsonrpc" => "2.0",
          "id" => "10",
          "method" => "tools/list"
        }

        response = send_mcp_request(["echo", "{{text}}"], request)

        tool = response["result"]["tools"].first
        expect(tool["name"]).to eq("echo")
        expect(tool["description"]).to eq("Run the shell command `echo {{text}}`")

        schema = tool["inputSchema"]
        expect(schema["properties"]).to have_key("text")
        expect(schema["properties"]["text"]["type"]).to eq("string")
        expect(schema["properties"]["text"]["description"]).to eq("the {{text}} arg")
      end

      it "executes blueprints without descriptions" do
        request = {
          "jsonrpc" => "2.0",
          "id" => "11",
          "method" => "tools/call",
          "params" => {
            "name" => "echo",
            "arguments" => {
              "text" => "Hello Blueprint!"
            }
          }
        }

        response = send_mcp_request(["echo", "{{text}}"], request)

        expect(response["jsonrpc"]).to eq("2.0")
        expect(response["id"]).to eq("11")
        expect(response["result"]["content"].first["type"]).to eq("text")
        expect(response["result"]["content"].first["text"]).to eq("Hello Blueprint!\n")
        expect(response["result"]["isError"]).to be false
      end
    end
  end

  describe "Error handling" do
    it "handles command errors gracefully" do
      request = {
        "jsonrpc" => "2.0",
        "id" => "12",
        "method" => "tools/call",
        "params" => {
          "name" => "false",
          "arguments" => {}
        }
      }

      response = send_mcp_request(["false"], request)

      expect(response["jsonrpc"]).to eq("2.0")
      expect(response["id"]).to eq("12")
      expect(response["result"]["isError"]).to be true
      expect(response["result"]["content"]).to be_an(Array)
    end

    it "handles nonexistent tools" do
      request = {
        "jsonrpc" => "2.0",
        "id" => "13",
        "method" => "tools/call",
        "params" => {
          "name" => "nonexistent",
          "arguments" => {}
        }
      }

      response = send_mcp_request(%w[echo hello], request)

      expect(response["jsonrpc"]).to eq("2.0")
      expect(response["id"]).to eq("13")
      expect(response).to have_key("error")
    end
  end

  describe "Complex commands" do
    it "works with git commands" do
      request = {
        "jsonrpc" => "2.0",
        "id" => "14",
        "method" => "tools/list"
      }

      response = send_mcp_request(%w[git status], request)

      tool = response["result"]["tools"].first
      expect(tool["name"]).to eq("git")
      expect(tool["description"]).to eq("Run the shell command `git [args]`")
    end

    it "works with multiple blueprint arguments" do
      request = {
        "jsonrpc" => "2.0",
        "id" => "15",
        "method" => "tools/list"
      }

      response = send_mcp_request([
                                    "rails", "generate",
                                    "{{generator#Rails generator name}}",
                                    "{{name#Resource name}}"
                                  ], request)

      tool = response["result"]["tools"].first
      expect(tool["name"]).to eq("rails")
      expect(tool["description"]).to eq("Run the shell command `rails {{generator}} {{name}}`")

      schema = tool["inputSchema"]
      expect(schema["properties"]).to have_key("generator")
      expect(schema["properties"]).to have_key("name")
      expect(schema["properties"]).not_to have_key("args")
    end
  end

  describe "Serve functionality (persistent server)" do
    describe "with simple echo command" do
      let(:args) { %w[echo test] }

      it "handles ping and tools/list" do
        requests = [
          { "jsonrpc" => "2.0", "id" => "1", "method" => "ping" },
          { "jsonrpc" => "2.0", "id" => "2", "method" => "tools/list" }
        ]

        responses = test_serve_with_args(args, requests)

        # Ping response
        expect(responses[0]["jsonrpc"]).to eq("2.0")
        expect(responses[0]["id"]).to eq("1")
        expect(responses[0]["result"]).to eq({})

        # Tools list response
        expect(responses[1]["jsonrpc"]).to eq("2.0")
        expect(responses[1]["id"]).to eq("2")
        expect(responses[1]["result"]["tools"]).to be_an(Array)
        expect(responses[1]["result"]["tools"].length).to eq(1)

        tool = responses[1]["result"]["tools"].first
        expect(tool["name"]).to eq("echo")
        expect(tool["description"]).to eq("Run the shell command `echo [args]`")
      end

      it "executes tools" do
        requests = [
          {
            "jsonrpc" => "2.0",
            "id" => "3",
            "method" => "tools/call",
            "params" => {
              "name" => "echo",
              "arguments" => { "args" => %w[hello world] }
            }
          }
        ]

        responses = test_serve_with_args(args, requests)

        expect(responses[0]["jsonrpc"]).to eq("2.0")
        expect(responses[0]["id"]).to eq("3")
        expect(responses[0]["result"]["content"].first["text"]).to eq("test hello world\n")
        expect(responses[0]["result"]["isError"]).to be false
      end
    end

    describe "with blueprinted command" do
      let(:args) { ["echo", "{{message#The message to echo}}"] }

      it "lists blueprinted tools correctly" do
        requests = [
          { "jsonrpc" => "2.0", "id" => "4", "method" => "tools/list" }
        ]

        responses = test_serve_with_args(args, requests)

        tool = responses[0]["result"]["tools"].first
        expect(tool["name"]).to eq("echo")
        expect(tool["description"]).to eq("Run the shell command `echo {{message}}`")

        schema = tool["inputSchema"]
        expect(schema["properties"]).to have_key("message")
        expect(schema["properties"]["message"]["description"]).to eq("The message to echo")
      end

      it "executes blueprinted tools" do
        requests = [
          {
            "jsonrpc" => "2.0",
            "id" => "5",
            "method" => "tools/call",
            "params" => {
              "name" => "echo",
              "arguments" => {
                "message" => "Hello Blueprint!",
                "args" => %w[with extra]
              }
            }
          }
        ]

        responses = test_serve_with_args(args, requests)

        expect(responses[0]["jsonrpc"]).to eq("2.0")
        expect(responses[0]["id"]).to eq("5")
        expect(responses[0]["result"]["content"].first["text"]).to eq("Hello Blueprint! with extra\n")
        expect(responses[0]["result"]["isError"]).to be false
      end
    end

    describe "error conditions" do
      it "handles invalid commands gracefully" do
        requests = [
          {
            "jsonrpc" => "2.0",
            "id" => "6",
            "method" => "tools/call",
            "params" => {
              "name" => "false",
              "arguments" => {}
            }
          }
        ]

        responses = test_serve_with_args(["false"], requests)

        expect(responses[0]["jsonrpc"]).to eq("2.0")
        expect(responses[0]["id"]).to eq("6")
        expect(responses[0]["result"]["isError"]).to be true
      end
    end
  end
end
