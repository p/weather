class DaybreakCache
  def initialize(db)
    @db = db
  end

  def get(key)
    @db[key]
  end

  def set(key, value)
    @db[key] = value
    @db.flush
  end
end
