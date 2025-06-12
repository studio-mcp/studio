# frozen_string_literal: true

require "spec_helper"

RSpec.describe Studio do
  it "has a version number" do
    expect(Studio::VERSION).not_to be_nil
  end

  describe ".serve" do
    it "raises error with empty argv" do
      expect { described_class.serve([]) }.to exit_with_code(1)
    end
  end
end
