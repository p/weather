require 'child_process_helper'
require 'daybreak_cache'
require 'weathercom'
require 'hashie/mash'
require 'open-uri'
require 'forwardable'
require 'byebug'
require 'geocoder'
require 'json'
require 'sinatra'
require 'daybreak'

Dir[File.join(File.dirname(__FILE__), 'presenters', '*.rb')].each do |path|
  require path[File.dirname(__FILE__).length+1...path.length].sub(/\.rb$/, '')
end

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

Location = Struct.new(
  :wc_client,
  :lat, :lng, :city, :state_abbr,
) do
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

  private def wc_client
    @wc_client ||= Weathercom::Client.new(cache: DaybreakCache.new($db))
  end

  private def geocode(query)
    wc_client.cached_geocode(query, ttl: 86400*100)
  end

  get '/locations' do
    locations = $db['locations'] || []
    content_type :json
    render_json(locations)
  end

  get '/locations/:location' do |location|
    loc = geocode(location)
    obs = loc.current_observation
    response = {
      location: LocationPresenter.new(loc).to_hash,
      current: ObservationPresenter.new(obs).to_hash,
    }
    content_type :json
    render_json(response)
  end

  get '/locations/:location/current' do |location|
    loc = geocode(location)
    obs = loc.current_observation
    response = {
      location: LocationPresenter.new(loc).to_hash,
      current: ObservationPresenter.new(obs).to_hash,
    }
    content_type :json
    render_json(response)
  end

  get '/locations/:location/hourly' do |location|
    loc = geocode(location)
    forecast = loc.hourly_forecast
    response = {
      location: LocationPresenter.new(loc).to_hash,
      current: HourlyForecastPresenter.new(forecast).to_hash,
    }
    content_type :json
    render_json(response)
  end

  get '/locations/:location/forecast' do |location|
    loc = geocode(location)
    forecast = loc.daily_forecast_10
    forecasts = payload['forecasts'].map do |forecast|
      {
        time: forecast['fcst_valid'],
        day: map_forecast(forecast['day']),
        night: map_forecast(forecast['night']),
      }
    end
    content_type :json
    render_json(forecasts)
  end

  get '/network' do
    up = false
    output = ChildProcessHelper.check_output(%w(ip a))
    output.split("\n").each do |line|
      if line =~ /^\d+: (eth|wlan)\d+:/
        if line =~ /\bNO-CARRIER\b/
          next
        elsif line =~ /\bUP\b/
          up = true
          break
        end
      end
    end

    render_json(up: up)
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

  def render_json(payload)
    set_cors_headers
    content_type :json
    JSON.generate(payload)
  end

  private def set_cors_headers
    response.headers["Access-Control-Allow-Origin"] = "*"
    response.headers["Access-Control-Allow-Methods"] = "GET,POST,PUT,PATCH,DELETE,OPTIONS"
    response.headers["Access-Control-Allow-Headers"] = "content-type"
  end

  options '*' do
    set_cors_headers
    ''
  end
end
