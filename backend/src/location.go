package main

import (
  "errors"
  "fmt"
  "github.com/jasonwinn/geocoder"
  log "github.com/sirupsen/logrus"
)

type resolved_location struct {
  Query string `json:"query"`
  Lat       float64 `json:"lat"`
  Lng       float64 `json:"lng"`
  City      string  `json:"city"`
  State     string  `json:"state"`
  UpdatedAt float64 `json:"updated_at"`
}

func (resloc resolved_location) GetUpdatedAt() float64 {
  return resloc.UpdatedAt
}

func resolve_location(location string) (*resolved_location, error) {
  data, err := lookup("geocodes", location)
  log.Debug("location: " + location)
  log.Debug(data)
  var resloc resolved_location

  if err != nil {
    return nil, errors.New("Could not look up location: " + err.Error())
  }

  if data == nil {
    if !online {
      return nil, errors.New("Cannot geocode - running in offline mode")
    }

    lat, lng, err := geocoder.Geocode(location)
    lat = lat
    lng = lng
    if err != nil {
      return nil, errors.New("Could not geocode " + location + ": " + err.Error())
    }

    resloc = resolved_location{
    location,
      lat,
      lng,
      "",
      "",
      now(),
    }

    err = persist("geocodes", location, &resloc)
    if err != nil {
      return nil, errors.New("Could not persist: " + err.Error())
    }

    log.Debug(fmt.Sprintf("Geocoded %s to %f,%f", location, resloc.Lat, resloc.Lng))
  } else {
    resloc = *(data.(*resolved_location))

    log.Debug(fmt.Sprintf("Retrieved %s from cache as %f,%f", location, resloc.Lat, resloc.Lng))
  }
  return &resloc, nil
}
