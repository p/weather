class HourlyForecastPresenter
  def initialize(forecast)
    @forecast = forecast
  end

  attr_reader :forecast

  def to_hash
    {
      expire_timestamp: forecast.expire_timestamp,
      start_timestamp: forecast.start_timestamp,
      start_at: Time.at(forecast.start_timestamp).iso8601,
      temp: forecast.temp,
      precip_probability: forecast.precip_probability,
      precip_type: forecast.precip_type,
      phrase: forecast.phrase,
    }
  end
end
