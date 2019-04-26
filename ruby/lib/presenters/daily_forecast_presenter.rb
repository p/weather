class DailyForecastPresenter
  def initialize(forecast)
    @forecast = forecast
  end

  attr_reader :forecast

  def to_hash
    {
      expire_timestamp: forecast.expire_timestamp,
      start_timestamp: forecast.start_timestamp,
      day: forecast.day_forecast && DayPartForecastPresenter.new(forecast.day_forecast).to_hash,
      night: DayPartForecastPresenter.new(forecast.night_forecast).to_hash,
    }
  end
end
