require 'sinatra'

class WildcardIndex < Sinatra::Base
  get '/*' do
    File.read('html/index.html')
  end
end
