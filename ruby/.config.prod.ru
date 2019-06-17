require 'dotenv'

Dotenv.load

$: << 'lib'

require 'app'
require 'wildcard_index'

map '/api' do
  run App
end

use Rack::Static, urls: %w(/static /icons), root: 'html'

run WildcardIndex
