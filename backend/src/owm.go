package main

import (
  "errors"
  "fmt"

  owm "github.com/briandowns/openweathermap"
  "github.com/kr/pretty"
  log "github.com/sirupsen/logrus"
  //"io"
  //"io/ioutil"
  //"net/http"
  //"regexp"
  //"strings"
)

func current_retriever_owm(resloc resolved_location) (persistable, error) {
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
    "",
    now(),
    now() + 300,
  }
  return &p, nil
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
        //v.Weather[0].Main,
        v.Weather[0].Description,
      },
      &day_part_forecast{
        float64(v.Dt) + 12*3600,
        v.Main.TempMin,
        0,
        "",
        //"",
        "",
      },
      0,
      "",
      "",
    })
  }

  p := forecast{
    dailies,
    now(),
    now() + 900,
  }
  return &p, nil
}
