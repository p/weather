package main

import (
  "encoding/json"
  "errors"
  "fmt"
  log "github.com/sirupsen/logrus"
  "net/http"
)

type WuForecastResponseMetadata struct {
  Language      string  `json:"language"`
  TransactionId string  `json:"transaction_id"`
  Version       string  `json:"version"`
  Latitude      float64 `json:"latitude"`
  Longitude     float64 `json:"longitude"`
  Units         string  `json:"units"`
  ExpireTimeGmt int64   `json:"expire_time_gmt"`
  StatusCode    int     `json:"status_code"`
}

type WuForecastResponseDaypart struct {
  FcstValid         int64  `json:"fcst_valid"`
  FcstValidLocal    string `json:"fcst_valid_local"`
  DayInd            string `json:"day_ind"`
  ThunderEnum       int    `json:"thunder_enum"`
  ThunderEnumPhrase string `json:"thunder_enum_phrase"`
  DaypartName       string `json:"daypart_name"`
  LongDaypartName   string `json:"long_daypart_name"`
  AltDaypartName    string `json:"alt_daypart_name"`
  Num               int    `json:"num"`
  Temp              int    `json:"temp"`
  Hi                int    `json:"hi"`
  Wc                int    `json:"wc"`
  Pop               int    `json:"pop"`
  IconExtd          int    `json:"icon_extd"`
  IconCode          int    `json:"icon_code"`
  Wxman             string `json:"wxman"`
  Phrase12Char      string `json:"phrase_12char"`
  Phrase22Char      string `json:"phrase_22char"`
  Phrase32Char      string `json:"phrase_32char"`
  SubphrasePt1      string `json:"subphrase_pt1"`
  SubphrasePt2      string `json:"subphrase_pt2"`
  SubphrasePt3      string `json:"subphrase_pt3"`
  PrecipType string `json:"precip_type"`
  Rh int `json:"rh"`
  Wspd int `json:"wspd"`
  Wdir int `json:"wdir"`
  WdirCardinal string `json:"wdir_cardinal"`
  Clds int `json:"clds"`
  PopPhrase string `json:"pop_phrase"`
  TempPhrase string `json:"temp_phrase"`
  AccumulationPhrase string `json:"accumulation_phrase"`
  WindPhrase string `json:"wind_phrase"`
  Shortcast         string `json:"shortcast"`
  Narrative         string `json:"narrative"`
  Qpf float64 `json:"qpf"`
  // may be int
  SnowQpf float64 `json:"snow_qpf"`
  SnowRange string `json:"snow_range"`
  SnowPhrase string `json:"snow_phrase"`
  SnowCode string `json:"snow_code"`
  VocalKey string `json:"vocal_key"`
  // this was always null, don't know type
  QualifierCode *string `json:"wind_phrase"`
  Qualifier *string `json:"qualifier"`
  UvIndexRaw float64 `json:"uv_index_raw"`
  UvIndex int `json:"uv_index"`
  UvWarning int `json:"uv_warning"`
  UvDesc string `json:"uv_desc"`
  GolfIndex int `json:"golf_index"`
  GolfCategory string `json:"golf_category"`
}

type WuForecastResponseForecast struct {
  Class          string                     `json:"class"`
  ExpireTimeGmt  int64                      `json:"expire_time_gmt"`
  FcstValid      int64                      `json:"fcst_valid"`
  FcstValidLocal string                     `json:"fcst_valid_local"`
  Num            int                        `json:"num"`
  MaxTemp        *int                       `json:"max_temp"`
  MinTemp        *int                       `json:"min_temp"`
  // was always null
  Torcon *string `json:"torcon"`
  // was always null
  Stormcon *string `json:"stormcon"`
  // was always null
  Blurb *string `json:"blurb"`
  // was always null
  BlurbAuthor *string `json:"blurb_author"`
  LunarPhaseDay int `json:"lunar_phase_day"`
  Dow string `json:"dow"`
  LunarPhase string `json:"lunar_phase"`
  LunarPhaseCode string `json:"lunar_phase_code"`
  Sunrise string `json:"sunrise"`
  Sunset string `json:"sunset"`
  Moonrise string `json:"moonrise"`
  Moonset string `json:"moonset"`
  // assume *string
  QualifierCode *string `json:"qualifier_code"`
  // assume *string
  Qualifier string `json:"qualifier"`
  Narrative string `json:"narrative"`
  Qpf float64 `json:"qpf"`
  // may be int
  SnowQpf float64 `json:"snow_qpf"`
  SnowRange string `json:"snow_range"`
  SnowPhrase string `json:"snow_phrase"`
  SnowCode string `json:"snow_code"`
  Night          *WuForecastResponseDaypart `json:"night"`
  Day            *WuForecastResponseDaypart `json:"day"`
}

type WuForecast10Response struct {
  Metadata  WuForecastResponseMetadata   `json:"metadata"`
  Forecasts []WuForecastResponseForecast `json:"forecasts"`
}

type WuClient struct {
  api_key     string
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
  log.Debug(url)
  return c.doGetForecast10(url)
}
