class HourlyForecastPresenter
  def initialize(forecast)
    @answer = forecast
  end

  attr_reader :forecast

  def to_hash
    {
    }
  end
end
