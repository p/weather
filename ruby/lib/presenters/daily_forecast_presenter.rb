class DailyForecastPresenter
  def initialize(forecast)
    @forecast = forecast
  end

  attr_reader :forecast

  def to_hash
    {
      expires_at: forecast.expire_time_gmt,
      day: forecast.day_forecast && DayPartForecastPresenter.new(forecast.day_forecast).to_hash,
      night: DayPartForecastPresenter.new(forecast.night_forecast).to_hash,
    }
  end
end
