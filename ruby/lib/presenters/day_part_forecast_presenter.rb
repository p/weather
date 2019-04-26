class DayPartForecastPresenter
  def initialize(forecast)
    @forecast = forecast
  end

  attr_reader :forecast

  def to_hash
    {
      temp: forecast.temp,
    }
  end
end
