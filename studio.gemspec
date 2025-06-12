# frozen_string_literal: true

require_relative "lib/studio/version"

Gem::Specification.new do |spec|
  spec.name = "studio"
  spec.version = Studio::VERSION
  spec.authors = ["Martin Emde"]
  spec.email = ["me@martinemde.com"]

  spec.summary = "Make any CLI command an MCP server"
  spec.description = <<-DESCRIPTION.lines.map(&:strip).join(" ").strip
    Studio transforms any CLI command into a StdIO Model Context Protocol (MCP) server.
    It exposes templated command-line tools so they can be executed more easily and safely than free inputs.
  DESCRIPTION
  spec.homepage = "https://github.com/martinemde/studio"
  spec.license = "MIT"
  spec.required_ruby_version = ">= 3.2.0"

  spec.metadata["homepage_uri"] = spec.homepage
  spec.metadata["source_code_uri"] = "https://github.com/martinemde/studio"
  spec.metadata["changelog_uri"] = "https://github.com/martinemde/studio/blob/main/CHANGELOG.md"
  spec.metadata["rubygems_mfa_required"] = "true"

  # Specify which files should be added to the gem when it is released.
  # The `git ls-files -z` loads the files in the RubyGem that have been added into git.
  gemspec = File.basename(__FILE__)
  spec.files = IO.popen(%w[git ls-files -z], chdir: __dir__, err: IO::NULL) do |ls|
    ls.readlines("\x0", chomp: true).reject do |f|
      (f == gemspec) ||
        f.start_with?(*%w[bin/ test/ spec/ features/ .git .github appveyor Gemfile])
    end
  end
  spec.bindir = "exe"
  spec.executables = spec.files.grep(%r{\Aexe/}) { |f| File.basename(f) }
  spec.require_paths = ["lib"]

  # Dependencies
  spec.add_dependency "mcp", "~> 0.1"
end
