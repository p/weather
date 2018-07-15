package main

import(
  "encoding/json"
  "net/http"
  "fmt"
  "errors"
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
  MaxTemp *int `json:"max_temp"`
  MinTemp *int `json:"min_temp"`
}

type WuForecast10Response struct {
  Metadata WuForecastResponseMetadata `json:"metadata"`
  Forecasts []WuForecastResponseForecast `json:"forecasts"`
  
}

type WuClient struct {
  api_key string
  http_client http.Client
}

func NewWuClient(api_key string) (*WuClient, error) {
  client := WuClient{
    api_key,
    http.Client{},
  }
  return &client, nil
}

func (c *WuClient) doGetForecast10(url string) (*WuForecast10Response, error) {
  req, err := http.NewRequest("GET", url, nil)
  if err != nil {
      return nil, errors.New("Could not send wu request:" + err.Error())
    return nil, err
  }
  
  res, err := c.http_client.Do(req)
  if err != nil {
      return nil, errors.New("Could not read wu response:" + err.Error())
  }
  
  defer res.Body.Close()
  
  var payload WuForecast10Response
  dec := json.NewDecoder(res.Body)
  err = dec.Decode(&payload)
  if err != nil {
      return nil, errors.New("Could not decode wu forecast:" + err.Error())
  }

  return &payload, nil
}

func (c *WuClient) GetForecast10ByLocation(lat float64, lng float64) (*WuForecast10Response, error) {
  url := fmt.Sprintf("https://api.weather.com/v1/geocode/%f/%f/forecast/daily/10day.json?apiKey=%s&units=e", lat, lng, c.api_key)
  return c.doGetForecast10(url)
}
