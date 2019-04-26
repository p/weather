class HourlyForecastPresenter
  def initialize(forecast)
    @forecast = forecast
  end

  attr_reader :forecast

  def to_hash
    {
      expires_at: forecast.expire_time_gmt,
    }
  end
end
