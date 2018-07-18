package main

type presented_current_conditions struct {
  Temp    float64 `json:"temp"`
  TempMin float64 `json:"temp_min"`
  TempMax float64 `json:"temp_max"`

  CreatedAt float64 `json:"created_at"`
}

func present_current_conditions(cc *current_conditions) presented_current_conditions {
  return presented_current_conditions{
    cc.Temp,
    cc.TempMin,
    cc.TempMax,
    float64(cc.CreatedAt) / 1e9,
  }
}

type presented_day_part_forecast struct {
  Time                 int64   `json:"time"`
  Temp                 float64 `json:"temp"`
  PrecipProbability int `json:"precip_probability"`
  PrecipType string `json:"precip_type"`
  ConditionName        string  `json:"condition_name"`
  ConditionDescription string  `json:"condition_description"`
}

type presented_daily_forecast struct {
  Time  float64                      `json:"time"`
  Night *presented_day_part_forecast `json:"night"`
  Day   *presented_day_part_forecast `json:"day"`
}

type presented_forecast struct {
  DailyForecasts []presented_daily_forecast `json:"daily_forecasts"`
  CreatedAt      float64                    `json:"created_at"`
}

func present_forecast(f *forecast) presented_forecast {
  presented_daily := make([]presented_daily_forecast, 0, len(f.DailyForecasts))
  for _, v := range f.DailyForecasts {
    presented_daily = append(presented_daily, presented_daily_forecast{
      float64(v.Time) / 1e9,
      present_day_part_forecast(v.Day),
      present_day_part_forecast(v.Night),
    })
  }
  return presented_forecast{
    presented_daily,
    float64(f.CreatedAt) / 1e9,
  }
}

func present_day_part_forecast(v *day_part_forecast) *presented_day_part_forecast {
  if v == nil {
    return nil
  }
  return &presented_day_part_forecast{
    v.Time / 1e9,
    v.Temp,
    v.PrecipProbability,
    v.PrecipType,
    v.ConditionName,
    v.ConditionDescription,
  }
}
