package main

import(
  "encoding/json"
)

type WuForecastResponseMetadata struct {
  Language string `json:"language"`
  TransactionId string `json:"transaction_id"`
  Version string `json:"version"`
  Latitude float64 `json:"latitude"`
  Longitude float64 `json:"longitude"`
  Units string `json:"units"`
  ExpireTimeGmt int64 `json:"expire_time_gmt"`
  StatusCode int `json:"status_code"`
}

type WuForecastResponseForecast struct {
  Class string `json:"class"`
  ExpireTimeGmt int64 `json:"expire_time_gmt"`
  FcstValid int64 `json:"fcst_valid"`
  FcstValidLocal string `json:"fcst_valid_local"`
  Num int `json:"num"`
  MaxTemp json.RawMessage `json:"max_temp"`
  MinTemp json.RawMessage `json:"min_temp"`
}

type WuForecastResponse struct {
  Metadata WuForecastResponseMetadata `json:"metadata"`
  Forecasts []WuForecastResponseForecast `json:"forecasts"`
  
}
