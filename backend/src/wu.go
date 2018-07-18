package main

import (
  "encoding/json"
  "errors"
  "fmt"
  log "github.com/sirupsen/logrus"
  "net/http"
)

// For wunderground api, night follows day.
// This means when retrieving a forecast in the middle of a day,
// there is data for night but not day part
// of the day on which the forecast is retrieved.

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
  // UTC timestamp for the forecast, e.g. 1531782000
  FcstValid int64 `json:"fcst_valid"`
  // ISO8601 time for the forecast, e.g. "2018-07-16T19:00:00-0400"
  FcstValidLocal string `json:"fcst_valid_local"`
  // "D" for day, "N" for night
  DayInd string `json:"day_ind"`
  // "Tonight", "Tomorrow", "Wednesday"
  DaypartName string `json:"daypart_name"`
  // "Monday night", "Tuesday night", "Tuesday"
  LongDaypartName string `json:"long_daypart_name"`
  // This is sometimes the same as DaypartName and sometimes same as
  // LongDaypartName
  AltDaypartName string `json:"alt_daypart_name"`
  // Number of this forecast in the returned data, starting with 1.
  // Forecasts for day parts (this struct) and days overall have separate numbering.
  // Forecast for today has num=1. If today only has a night day part,
  // that night's forecast would have num=1 as well.
  // Forecast for tomorrow will have num=2, tomorrow's day num=2,
  // tomorrow's night num=3. The day after tomorrow will have num=3 for the
  // entire day, num=4 for the day part, num=5 for the night part.
  Num                int     `json:"num"`
  Temp               int     `json:"temp"`
  Hi                 int     `json:"hi"`
  Wc                 int     `json:"wc"`
  Pop                int     `json:"pop"`
  PopPhrase          string  `json:"pop_phrase"`
  IconExtd           int     `json:"icon_extd"`
  IconCode           int     `json:"icon_code"`
  Wxman              string  `json:"wxman"`
  Phrase12Char       string  `json:"phrase_12char"`
  Phrase22Char       string  `json:"phrase_22char"`
  Phrase32Char       string  `json:"phrase_32char"`
  SubphrasePt1       string  `json:"subphrase_pt1"`
  SubphrasePt2       string  `json:"subphrase_pt2"`
  SubphrasePt3       string  `json:"subphrase_pt3"`
  PrecipType         string  `json:"precip_type"`
  Rh                 int     `json:"rh"`
  Wspd               int     `json:"wspd"`
  Wdir               int     `json:"wdir"`
  WdirCardinal       string  `json:"wdir_cardinal"`
  Clds               int     `json:"clds"`
  TempPhrase         string  `json:"temp_phrase"`
  AccumulationPhrase string  `json:"accumulation_phrase"`
  WindPhrase         string  `json:"wind_phrase"`
  Shortcast          string  `json:"shortcast"`
  Narrative          string  `json:"narrative"`
  ThunderEnum        int     `json:"thunder_enum"`
  ThunderEnumPhrase  string  `json:"thunder_enum_phrase"`
  Qpf                float64 `json:"qpf"`
  // may be int
  SnowQpf    float64 `json:"snow_qpf"`
  SnowRange  string  `json:"snow_range"`
  SnowPhrase string  `json:"snow_phrase"`
  SnowCode   string  `json:"snow_code"`
  VocalKey   string  `json:"vocal_key"`
  // this was always null, don't know type
  QualifierCode *string `json:"wind_phrase"`
  Qualifier     *string `json:"qualifier"`
  UvIndexRaw    float64 `json:"uv_index_raw"`
  UvIndex       int     `json:"uv_index"`
  UvWarning     int     `json:"uv_warning"`
  UvDesc        string  `json:"uv_desc"`
  GolfIndex     int     `json:"golf_index"`
  GolfCategory  string  `json:"golf_category"`
}

type WuForecastResponseForecast struct {
  // Type of forecast, "fod_long_range_daily" for this data
  Class string `json:"class"`
  // UTC timestamp: 1531769805
  ExpireTimeGmt int64 `json:"expire_time_gmt"`
  // UTC timestamp: 1531911600
  FcstValid int64 `json:"fcst_valid"`
  // ISO8601 time: "2018-07-18T07:00:00-0400"
  FcstValidLocal string `json:"fcst_valid_local"`
  // Number of this forecast in the returned data, starting with 1.
  // Forecasts for days (this struct) and day parts have separate numbering.
  // Forecast for today has num=1. If today only has a night day part,
  // that night's forecast would have num=1 as well.
  // Forecast for tomorrow will have num=2, tomorrow's day num=2,
  // tomorrow's night num=3. The day after tomorrow will have num=3 for the
  // entire day, num=4 for the day part, num=5 for the night part.
  Num int `json:"num"`
  // Same as Day.Temp, can be null if there is no day data
  // which will happen when a forecast is retreived late enough in the day
  MaxTemp *int `json:"max_temp"`
  // Same as Night.Temp
  MinTemp int `json:"min_temp"`
  // was always null
  Torcon *string `json:"torcon"`
  // was always null
  Stormcon *string `json:"stormcon"`
  // was always null
  Blurb *string `json:"blurb"`
  // was always null
  BlurbAuthor   *string `json:"blurb_author"`
  LunarPhaseDay int     `json:"lunar_phase_day"`
  // Day of week, e.g. "Monday", "Tuesday"
  Dow            string `json:"dow"`
  LunarPhase     string `json:"lunar_phase"`
  LunarPhaseCode string `json:"lunar_phase_code"`
  Sunrise        string `json:"sunrise"`
  Sunset         string `json:"sunset"`
  Moonrise       string `json:"moonrise"`
  Moonset        string `json:"moonset"`
  // assume *string
  QualifierCode *string `json:"qualifier_code"`
  // assume *string
  Qualifier string  `json:"qualifier"`
  Narrative string  `json:"narrative"`
  Qpf       float64 `json:"qpf"`
  // may be int
  SnowQpf    float64                    `json:"snow_qpf"`
  SnowRange  string                     `json:"snow_range"`
  SnowPhrase string                     `json:"snow_phrase"`
  SnowCode   string                     `json:"snow_code"`
  Night      *WuForecastResponseDaypart `json:"night"`
  Day        *WuForecastResponseDaypart `json:"day"`
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
