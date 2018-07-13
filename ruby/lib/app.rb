require 'byebug'
require 'geocoder'
require 'json'
require 'sinatra'
require 'daybreak'

$db = Daybreak::DB.new(ENV['DB_PATH'] || 'ruby.db')

class App < Sinatra::Base
  get '/locations' do
    locations = $db['locations'] || []
    content_type :json
    JSON.generate(locations)
  end

  get '/locations/:location/current' do |location|
    coords = resolve(location)
    content_type :json
    JSON.generate(coords)
  end

  get '/locations/:location/forecast' do |location|
  end

  private def resolve(location)
    coords = $db["geocode:#{location}"]
    if coords.nil?
      result = Geocoder.search(location).first
      coords = {
        lat: result.geometry['location']['lat'],
        lng: result.geometry['location']['lng'],
      }
      $db["geocode:#{location}"] = coords
    end
    coords
  end
end
