require 'hashie/mash'
require 'open-uri'
require 'forwardable'
require 'byebug'
require 'geocoder'
require 'json'
require 'sinatra'
require 'daybreak'

$db = Daybreak::DB.new(ENV['DB_PATH'] || 'ruby.db')

Geocoder.configure(
  lookup: :mapquest,
  api_key: ENV['MAPQUEST_API_KEY'],
)

class ResolvedLocation
  extend Forwardable

  def initialize(info)
    @info = Hashie::Mash.new(info)
  end

  attr_reader :info

  def_delegators :info, :lat, :lng, :city, :state, :updated_at

  def wu_current_url
    "https://www.wunderground.com/weather/us/ny/new-york"
  end
end

class App < Sinatra::Base
  private def wu_api_request(path, resloc)
    attempt = 1
    payload = nil
    while true
      begin
        url = "https://api.weather.com/v1/geocode/#{resloc.lat}/#{resloc.lng}/#{path}.json?apiKey=#{api_key}&units=e"
        puts url
        payload = JSON.parse(open(url).read)
        break
      rescue OpenURI::HTTPError => e
        if attempt > 1
          raise
        elsif e.message =~ /\A401\b/
          # This is probably not correct as my 401 was from a wrong URL
          reset_api_key
          attempt += 1
        else
          raise
        end
      end
    end
    payload
  end

  private def reset_api_key
    @api_key = $db['wu_api_key'] = nil
  end

  private def api_key
    @api_key ||= wu_api_key("https://www.wunderground.com/weather/us/ny/new-york")
  end

  get '/locations' do
    locations = $db['locations'] || []
    content_type :json
    JSON.generate(locations)
  end

  get '/locations/:location/current' do |location|
    resloc = resolve(location)
    api_key = wu_api_key(resloc.wu_current_url)
    payload = wu_api_request('observations/current', resloc)
    response = {
      location: resloc.info,
      current: payload,
    }
    content_type :json
    JSON.generate(response)
  end

  get '/locations/:location/forecast' do |location|
    resloc = resolve(location)
    api_key = wu_api_key(resloc.wu_current_url)
    url = "https://api.weather.com/v1/geocode/#{resloc.lat}/#{resloc.lng}/forecast/daily/10day.json?apiKey=#{api_key}&units=e"
    payload = JSON.parse(open(url).read)
    forecasts = payload['forecasts'].map do |forecast|
      {
        time: forecast['fcst_valid'],
        day: map_forecast(forecast['day']),
        night: map_forecast(forecast['night']),
      }
    end
    content_type :json
    JSON.generate(forecasts)
  end

  private def resolve(location)
    coords = $db["geocode:#{location}"]
    if coords.nil?
      result = Geocoder.search(location).first
      coords = {
        #lat: result.geometry['location']['lat'],
        #lng: result.geometry['location']['lng'],
        lat: result.coordinates.first,
        lng: result.coordinates.last,
        city: result.city,
        state: result.state,
        updated_at: Time.now.to_f,
      }
      $db["geocode:#{location}"] = coords
      $db.flush
    end
    ResolvedLocation.new(coords)
  end

  private def wu_api_key(url)
    api_key = $db['wu_api_key']
    if api_key.nil?
      contents = open(url).read
      if contents =~ /apiKey=(\w+)/
        api_key = $1
        $db['wu_api_key'] = api_key
        $db.flush
      else
        raise "Could not find the api key in #{url}"
      end
    end
    api_key
  end

  private def map_forecast(forecast)
    {
      temp: forecast['temp'],
      condition_name: forecast['shortcast'],
      condition_description: forecast['narrative'],
      precip_prob: forecast['pop'],
    }
  end
end
