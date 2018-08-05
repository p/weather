package main

type current_conditions struct {
  Temp       float64 `json:"temp"`
  TempMin    float64 `json:"temp_min"`
  TempMax    float64 `json:"temp_max"`
  WwirPhrase string  `json:"wwir_phrase"`

  UpdatedAt float64 `json:"updated_at"`
  ExpiresAt float64 `json:"expires_at"`
}

func (cc current_conditions) GetUpdatedAt() float64 {
  return cc.UpdatedAt
}

func (cc current_conditions) GetExpiresAt() float64 {
  return cc.ExpiresAt
}

type day_part_forecast struct {
  Time              float64 `json:"time"`
  Temp              float64 `json:"temp"`
  PrecipProbability int     `json:"precip_probability"`
  PrecipType        string  `json:"precip_type"`
  //ShortNarrative string  `json:"short_narrative"`
  Narrative string `json:"narrative"`
}

type daily_forecast struct {
  Time              float64            `json:"time"`
  Day               *day_part_forecast `json:"day"`
  Night             *day_part_forecast `json:"night"`
  PrecipProbability int                `json:"precip_probability"`
  PrecipType        string             `json:"precip_type"`
  Narrative         string             `json:"narrative"`
}

type forecast struct {
  DailyForecasts []daily_forecast `json:"daily_forecasts"`
  UpdatedAt      float64          `json:"updated_at"`
  ExpiresAt      float64          `json:"expires_at"`
}

func (f forecast) GetUpdatedAt() float64 {
  return f.UpdatedAt
}

func (f forecast) GetExpiresAt() float64 {
  return f.ExpiresAt
}
