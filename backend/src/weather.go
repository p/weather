package main

import(
"fmt"
log "github.com/sirupsen/logrus"
  owm "github.com/briandowns/openweathermap"
  "github.com/kr/pretty"
"regexp"
"io"
"io/ioutil"
"net/http"
"errors"
)

type current_conditions struct {
  Temp    float64 `json:"temp"`
  TempMin float64 `json:"temp_min"`
  TempMax float64 `json:"temp_max"`

  UpdatedAt float64 `json:"updated_at"`
}

func (cc current_conditions) GetUpdatedAt() float64 {
  return cc.UpdatedAt
}

type day_part_forecast struct {
  Time                 float64 `json:"time"`
  Temp                 float64 `json:"temp"`
  PrecipProbability    int     `json:"precip_probability"`
  PrecipType           string  `json:"precip_type"`
  ConditionName        string  `json:"condition_name"`
  ConditionDescription string  `json:"condition_description"`
}

type daily_forecast struct {
  Time  float64            `json:"time"`
  Day   *day_part_forecast `json:"day"`
  Night *day_part_forecast `json:"night"`
}

type forecast struct {
  DailyForecasts []daily_forecast `json:"daily_forecasts"`
  UpdatedAt      float64          `json:"updated_at"`
}

func (f forecast) GetUpdatedAt() float64 {
  return f.UpdatedAt
}

func get_weather_with_cache(
  location string,
  resloc resolved_location,
  bucket_name string,
  retriever func(resloc resolved_location) (persistable, error),
  network NetworkUse,
) (persistable, error) {
  var p persistable
  var err error

  if network != NetworkForce {
    data, err := lookup(bucket_name, location)
    if err != nil {
      return nil, errors.New("Error retrieving from cache: " + err.Error())
    }
    if data != nil {
      typed := data.(persistable)

      log.Debug(fmt.Sprintf("Retrieved cached data for %s", location))

      if !online || network == NetworkSkip || now()-typed.GetUpdatedAt() <= current_age {
        p = typed
      }
    }
  }

  if network != NetworkSkip {
    if p == nil {
      if !online {
        return nil, errors.New("Cannot get weather - running in offline mode")
      }

      p, err = retriever(resloc)
      if err != nil {
        return nil, errors.New("Could not retrieve: " + err.Error())
      }

      err = persist(bucket_name, location, p)
      if err != nil {
        return nil, errors.New("Could not persist: " + err.Error())
      }

      log.Debug(fmt.Sprintf("Fetched data for %s", location))
    }
  }

  if p == nil {
    return nil, errors.New("Could not retrieve weather")
  }

  return p, nil
}

func current_retriever(resloc resolved_location) (persistable, error) {
  w, err := owm.NewCurrent("F", "en", owm_api_key)
  if err != nil {
    return nil, errors.New("Could not make current: " + err.Error())
  }

  err = w.CurrentByCoordinates(
    &owm.Coordinates{
      Longitude: resloc.Lng,
      Latitude:  resloc.Lat,
    },
  )
  if err != nil {
    return nil, errors.New("Could not get current weather: " + err.Error())
  }

  if log.GetLevel() == log.DebugLevel {
    fmt.Printf("%# v", pretty.Formatter(w))
  }

  p := current_conditions{
    w.Main.Temp,
    w.Main.TempMin,
    w.Main.TempMax,
    now(),
  }
  return &p, nil
}

func get_current_weather(location string,
  resloc resolved_location) (*current_conditions, error) {

  p, err := get_weather_with_cache(
    location,
    resloc,
    "current_conditions",
    current_retriever, 0)

  if err != nil {
    return nil, err
  }

  return p.(*current_conditions), nil
}

func forecast_retriever(resloc resolved_location) (persistable, error) {
  w, err := owm.NewForecast("5", "F", "en", owm_api_key)
  if err != nil {
    return nil, errors.New("Could not make forecast: " + err.Error())
  }

  err = w.DailyByCoordinates(
    &owm.Coordinates{
      Longitude: resloc.Lng,
      Latitude:  resloc.Lat,
    },
    5,
  )
  if err != nil {
    return nil, errors.New("Could not get forecast: " + err.Error())
  }

  if log.GetLevel() == log.DebugLevel {
    fmt.Printf("%# v", pretty.Formatter(w))
  }

  l := w.ForecastWeatherJson.(*owm.Forecast5WeatherData).List
  dailies := make([]daily_forecast, 0)
  for _, v := range l {
    dailies = append(dailies, daily_forecast{
      float64(v.Dt),
      &day_part_forecast{
        float64(v.Dt),
        v.Main.TempMax,
        0,
        "",
        v.Weather[0].Main,
        v.Weather[0].Description,
      },
      &day_part_forecast{
        float64(v.Dt) + 12*3600,
        v.Main.TempMin,
        0,
        "",
        "",
        "",
      },
    })
  }

  p := forecast{
    dailies,
    now(),
  }
  return &p, nil
}

func get_forecast(location string,
  resloc resolved_location, network NetworkUse) (*forecast, error) {

  p, err := get_weather_with_cache(
    location,
    resloc,
    "forecasts",
    forecast_retriever, network)

  if err != nil {
    return nil, err
  }

  return p.(*forecast), nil
}

func get_wu_forecast(location string,
  resloc resolved_location, network NetworkUse) (*forecast, error) {

  p, err := get_weather_with_cache(
    location,
    resloc,
    "wu_forecasts",
    wu_forecast_retriever, network)

  if err != nil {
    return nil, err
  }

  return p.(*forecast), nil
}

type wu_credentials struct {
  ApiKey    string
  UpdatedAt float64
}

func (x wu_credentials) GetUpdatedAt() float64 {
  return x.UpdatedAt
}

const WU_API_KEY_REGEXP = "apiKey=([a-zA-Z0-9]+)[^A-Za-z0-9]"

var wu_api_key_regexp *regexp.Regexp

const WU_API_KEY_URL = "https://www.wunderground.com/weather/us/ny/new-york"

func get_wu_api_key() (string, error) {
  req, err := http.NewRequest("GET", WU_API_KEY_URL, nil)
  if err != nil {
    return "", err
  }
  req.Header.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64; rv:52.9) Gecko/20100101 Goanna/3.4 Firefox/52.9 PaleMoon/27.9.2")
  client := &http.Client{}
  res, err := client.Do(req)
  if err != nil {
    return "", errors.New("Could not retrieve wu api key: " + err.Error())
  }

  defer res.Body.Close()

  //content, err := ioutil.ReadAll(res.Body)
  //ioutil.WriteFile("x.html",content ,0644)
  //log.Debug(string(content))
  //log.Debug(len(content))

  a := ioutil.ReadAll
  a = a

  window := ""
  buf := make([]byte, 32768)
  for {
    n, err := res.Body.Read(buf)
    if n > 0 {
      if len(window) > 20 {
        window = window[len(window)-20:len(window)] + string(buf[:n])
      } else {
        window = window + string(buf[:n])
      }
      //log.Debug(string(buf[:n]))

      if len(window) > 20 || err != nil {
        matches := wu_api_key_regexp.FindStringSubmatch(window)
        if len(matches) > 0 {
          api_key := matches[1]
          wu_creds := wu_credentials{
            api_key,
            now(),
          }
          persist("config", "wu_credentials", &wu_creds)
          return api_key, nil
        }
        window = window[:20]
      }
    }

    if err == io.EOF {
      break
    } else if err != nil {
      return "", errors.New("Error reading while retrieving wu api key:" + err.Error())
    }
  }

  return "", errors.New("Could not find wu api key while looking for it")
}

func get_wu_api_key_cached() (string, error) {
  raw_wu_creds, err := lookup("config", "wu_api_key")
  var api_key string
  if err != nil {
    return "", errors.New("Could not retrieve wu api key: " + err.Error())
  }

  if raw_wu_creds != nil {
    wu_creds := raw_wu_creds.(*wu_credentials)
    api_key = wu_creds.ApiKey
  }

  if api_key == "" {
    api_key, err = get_wu_api_key()
    if err != nil {
      return "", err
    }
    wu_creds := wu_credentials{
      api_key,
      now(),
    }
    err = persist("config", "wu_api_key", &wu_creds)
    if err != nil {
      return "", errors.New("Could not persist wu api key: " + err.Error())
    }
  }
  return api_key, err
}

func wu_forecast_retriever(resloc resolved_location) (persistable, error) {
  api_key, err := get_wu_api_key_cached()
  if err != nil {
    return nil, errors.New("Error retrieving wu api key:" + err.Error())
  }
  log.Debug(api_key)

  c, err := NewWuClient(api_key)
  if err != nil {
    return nil, err
  }
  payload, err := c.GetForecast10ByLocation(
    resloc.Lat, resloc.Lng)
  if err != nil {
    return nil, err
  }

  dailies := make([]daily_forecast, 0)
  for _, v := range payload.Forecasts {
    dailies = append(dailies, daily_forecast{
      float64(v.FcstValid),
      convert_wu_forecast(v.Day),
      convert_wu_forecast(v.Night),
    })
  }

  f := forecast{
    dailies,
    now(),
  }

  return &f, nil
}

func convert_wu_forecast(v *WuForecastResponseDaypart) *day_part_forecast {
  if v == nil {
    return nil
  }
  return &day_part_forecast{
    float64(v.FcstValid),
    float64(v.Temp),
    v.Pop,
    v.PrecipType,
    v.Shortcast,
    v.Narrative,
  }
}
