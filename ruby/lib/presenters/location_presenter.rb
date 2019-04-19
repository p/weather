class LocationPresenter
  def initialize(location)
    @answer = location
  end

  attr_reader :location

  def to_hash
    {
      lat: location.lat,
      lng: location.lng,
    }
  end
end
