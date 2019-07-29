require 'dotenv'
require 'logger'

Dotenv.load

error_logger = Logger.new(STDERR)

$: << 'lib'

require 'app'
require 'wildcard_index'

App.error_logger(error_logger)

map '/api' do
  run App
end

use Rack::Static, urls: %w(/static /icons), root: 'html'

run WildcardIndex
