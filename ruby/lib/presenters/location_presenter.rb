class LocationPresenter
  def initialize(location)
    @location = location
  end

  attr_reader :location

  def to_hash
    {
      lat: location.lat,
      lng: location.lng,
      city: location.city,
      state_abbr: location.state_abbr,
    }
  end
end
