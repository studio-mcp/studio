# frozen_string_literal: true

require "spec_helper"

RSpec.describe Studio do
  it "has a version number" do
    expect(Studio::VERSION).not_to be_nil
  end

  describe ".serve" do
    it "raises error with empty argv" do
      expect { described_class.serve([]) }.to exit_with_code(1).and output(/Usage: studio/).to_stderr
    end
  end

  describe "inspector integration tests" do
    before do
      unless system("node --version | grep -q 'v24'")
        skip("These tests require Node with a version that supports the inspector")
      end
    rescue StandardError
      skip("These tests require Node with a version that supports the inspector")
    end

    it "can echo a message" do
      expect(`bin/echo 'Hello, world!'`).to eq("Hello, world!")
    end

    it "can handle a failing message" do
      expect(system("bin/echo")).to be_falsey
    end
  end
end
