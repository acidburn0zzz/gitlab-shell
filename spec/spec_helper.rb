ROOT_PATH = File.expand_path(File.join(File.dirname(__FILE__), ".."))

if ENV['COVERALLS']
  require 'coveralls'
  Coveralls.wear!
else
  require 'simplecov'
  SimpleCov.start
end

require 'vcr'
require 'webmock'
require_relative '../lib/gitlab_excon'

VCR.configure do |c|
  c.cassette_library_dir = 'spec/vcr_cassettes'
  c.hook_into :excon
  c.configure_rspec_metadata!
end
