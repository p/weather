package main

import (
  "errors"
  //"fmt"
  //"github.com/kr/pretty"
  log "github.com/sirupsen/logrus"
  "gopkg.in/weather.v0"
  "io"
  "io/ioutil"
  "net/http"
  "regexp"
  "strings"
)

func current_retriever_wu(resloc resolved_location) (persistable, error) {
  api_key, err := get_wu_api_key_cached()
  if err != nil {
    return nil, errors.New("Error retrieving wu api key:" + err.Error())
  }
  log.Debug(api_key)

  c := weather.NewClient(api_key)
  current, err := c.GetCurrentByLocation(
    resloc.Lat, resloc.Lng, "e")
  if err != nil {
    return nil, err
  }

  persist("wu_currents_raw", resloc.Query, current)

  wwir, err := c.GetWwirByLocation(
    resloc.Lat, resloc.Lng, "e")
  if err != nil {
    return nil, err
  }

  persist("wu_wwirs_raw", resloc.Query, wwir)

  p := current_conditions{
    float64(current.Observation.Imperial.Temp),
    float64(current.Observation.Imperial.TempMin24hour),
    float64(current.Observation.Imperial.TempMax24hour),
    wwir.Forecast.Phrase,
    now(),
  }
  return &p, nil
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

  c := weather.NewClient(api_key)
  payload, err := c.GetForecast10ByLocation(
    resloc.Lat, resloc.Lng, "")
  if err != nil {
    return nil, err
  }

  persist("wu_forecasts_raw", resloc.Query, payload)

  dailies := make([]daily_forecast, 0)
  for _, v := range payload.Forecasts {
    var dpv weather.DaypartForecast
    if v.Day != nil {
      dpv = *v.Day
    } else {
      dpv = v.Night
    }
    dailies = append(dailies, daily_forecast{
      float64(v.FcstValid),
      convert_wu_forecast(v.Day),
      convert_wu_forecast(&v.Night),
      dpv.Pop,
      dpv.PrecipType,
      extract_top_level_narrative(v.Narrative),
    })
  }

  f := forecast{
    dailies,
    now(),
  }

  return &f, nil
}

func extract_narrative(v weather.DaypartForecast) string {
  n := v.Narrative
  if v.WindPhrase != "" {
    n = strings.Replace(n, " "+v.WindPhrase, "", 1)
  }
  if v.TempPhrase != "" {
    n = strings.Replace(n, " "+v.TempPhrase, "", 1)
  }
  if v.PopPhrase != "" {
    n = strings.Replace(n, " "+v.PopPhrase, "", 1)
  }
  return n
}

func extract_top_level_narrative(n string) string {
  return n
}

func convert_wu_forecast(v *weather.DaypartForecast) *day_part_forecast {
  if v == nil {
    return nil
  }
  return &day_part_forecast{
    float64(v.FcstValid),
    float64(v.Temp),
    v.Pop,
    v.PrecipType,
    //v.Shortcast,
    extract_narrative(*v),
  }
}
