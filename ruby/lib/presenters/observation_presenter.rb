class ObservationPresenter
  def initialize(obs)
    @obs = obs
  end

  attr_reader :obs

  def to_hash
    {
      expire_time_gmt: obs.expire_time_gmt,
      temp: obs.temp,
      feels_like: obs.feels_like,
      phrase: obs.phrase_32char,
    }
  end
end
