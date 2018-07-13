require 'dotenv'

Dotenv.load

$: << 'lib'

require 'app'

run App
