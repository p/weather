require 'dotenv'
require 'logger'

Dotenv.load

error_logger = Logger.new(STDERR)

$: << 'lib'

require 'app'

App.error_logger(error_logger)
run App
