# frozen_string_literal: true

require "studio"

RSpec.configure do |config|
  # Enable flags like --only-failures and --next-failure
  config.example_status_persistence_file_path = ".rspec_status"

  # Disable RSpec exposing methods globally on `Module` and `main`
  config.disable_monkey_patching!

  config.expect_with :rspec do |c|
    c.syntax = :expect
  end
end

# Custom matcher for exit codes
RSpec::Matchers.define :exit_with_code do |expected_code|
  supports_block_expectations

  match do |block|
    begin
      block.call
      @actual_code = 0
    rescue SystemExit => e
      @actual_code = e.status
    end

    @actual_code == expected_code
  end

  failure_message do |_block|
    "expected block to exit with code #{expected_code}, but exited with #{@actual_code}"
  end
end
