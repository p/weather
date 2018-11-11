package main

import (
  "errors"
  "fmt"
  //"github.com/kr/pretty"
  log "github.com/sirupsen/logrus"
  //"io"
  //"io/ioutil"
  //"net/http"
  //"regexp"
  //"strings"
)

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
    data, err := lookup_expiring(bucket_name, location)
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
  if network==NetworkSkip{
    return nil, errors.New("Could not retrieve weather (in offline mode)")
    }else{
    return nil, errors.New("Could not retrieve weather (in online mode)")
    }
  }

  return p, nil
}

func get_current_weather(location string,
  resloc resolved_location,
  service string,
  network NetworkUse) (*current_conditions, error) {

  var cr func(resloc resolved_location) (persistable, error)
  if service == "wu" {
    cr = current_retriever_wu
  } else {
    cr = current_retriever_owm
  }

  p, err := get_weather_with_cache(
    location,
    resloc,
    "current_conditions",
    cr, network)

  if err != nil {
    return nil, err
  }

  return p.(*current_conditions), nil
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

type location_everything struct {
  Location resolved_location  `json:"location"`
  Current  current_conditions `json:"current"`
  Forecast forecast           `json:"forecast"`
}

func get_location_everything(location string, resloc resolved_location,
  network NetworkUse) (*location_everything, error) {
  cc, err := get_current_weather(location, resloc, "wu", network)
  if err != nil {
    return nil, err
  }
  f, err := get_wu_forecast(location, resloc, network)
  if err != nil {
    return nil, err
  }

  return &location_everything{
    resloc,
    *cc,
    *f,
  }, nil

}
