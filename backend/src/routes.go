package main

import (
  //"bytes"

  "errors"

  "github.com/gin-gonic/gin"
  log "github.com/sirupsen/logrus"

  //"time"

  bolt "github.com/coreos/bbolt"
  detector "gopkg.in/network-detect.v0"
  //"html/template"
)

func list_locations_route(c *gin.Context) {
  var locations []string
  err := db.View(func(tx *bolt.Tx) error {
    b := tx.Bucket([]byte("geocodes"))
    b.ForEach(func(k, v []byte) error {
      locations = append(locations, string(k))
      return nil
    })
    return nil
  })
  if err != nil {
    return_500(c, "Problem listing locations", err)
    return
  }

  render_json(c, locations)
}

func get_conditions_route(c *gin.Context) {
  location := c.Param("location")
  resloc, err := resolve_location(location)
  if err != nil {
    return_500(c, "Could not resolve location: "+location, err)
    return
  }
  network, err := get_network_flag(c)
  if err != nil {
    return_500(c, "Problem getting the network flag", err)
    return
  }
  cc, err := get_current_weather(location, *resloc, "wu", network)
  if err != nil {
    return_500(c, "Could not get weather for location: "+location, err)
    return
  }

  render_json(c, cc)
}

func get_forecast_route(c *gin.Context) {
  location := c.Param("location")
  resloc, err := resolve_location(location)
  if err != nil {
    return_500(c, "Could not resolve location: "+location, err)
    return
  }
  f, err := get_forecast(location, *resloc, 0)
  if err != nil {
    return_500(c, "Could not get weather for location: "+location, err)
    return
  }

  render_json(c, f)
}

func get_network_flag(c *gin.Context) (NetworkUse, error) {
  raw_network := c.Query("network")
  var network NetworkUse
  switch raw_network {
  case "":
    network = NetworkDefault
  case "0":
    network = NetworkDefault
  case "1":
    network = NetworkForce
  case "2":
    network = NetworkSkip
  default:
    return NetworkDefault, errors.New("Invalid network value: " + raw_network)
  }
  return network, nil
}

func get_wu_forecast_route(c *gin.Context) {
  location := c.Param("location")
  network, err := get_network_flag(c)
  if err != nil {
    return_500(c, "Problem getting the network flag", err)
    return
  }
  resloc, err := resolve_location(location)
  if err != nil {
    return_500(c, "Could not resolve location: "+location, err)
    return
  }
  f, err := get_wu_forecast(location, *resloc, network)
  if err != nil {
    return_500(c, "Could not get weather for location: "+location, err)
    return
  }

  render_json(c, f)
}

func get_wu_forecast_raw_route(c *gin.Context) {
  location := c.Param("location")

  log.Debug(location)
  data, err := lookup("wu_forecasts_raw", location)
  if err != nil {
    return_500(c, "Could not get raw forecast for location: "+location, err)
    return
  }
  if data != nil {
    render_json(c, data)
    return
  }

  resloc, err := resolve_location(location)
  if err != nil {
    return_500(c, "Could not resolve location: "+location, err)
    return
  }
  f, err := get_wu_forecast(location, *resloc, NetworkDefault)
  f = f
  if err != nil {
    return_500(c, "Could not get weather for location: "+location, err)
    return
  }

  data, err = lookup("wu_forecasts_raw", location)
  if err != nil {
    return_500(c, "Could not get raw forecast for location: "+location, err)
    return
  }
  if data != nil {
    render_json(c, data)
    return
  }

  return_500(c, "No wu cached data after retrieving a forecast", nil)
}

func location_route(c *gin.Context) {
  location := c.Param("location")
  network, err := get_network_flag(c)
  if err != nil {
    return_500(c, "Problem getting the network flag", err)
    return
  }
  resloc, err := resolve_location(location)
  if err != nil {
    return_500(c, "Could not resolve location: "+location, err)
    return
  }
  f, err := get_location_everything(location, *resloc, network)
  if err != nil {
    return_500(c, "Could not get weather for location: "+location, err)
    return
  }

  render_json(c, f)
  f = f
}

var nd *detector.NetworkDetector

func network_route(c *gin.Context) {
  if nd == nil {
    and := detector.NewNetworkDetector()
    nd = &and
  }
  up, err := nd.Up()
  if err != nil {
    return_500(c, "Could not figure out network status", err)
    return
  }
  ns := network_status{
    up}
  render_json(c, ns)

}

func define_routes(router *gin.Engine) {
  router.GET("/locations", list_locations_route)
  router.GET("/locations/:location", location_route)
  router.GET("/locations/:location/current", get_conditions_route)
  router.GET("/locations/:location/forecast", get_forecast_route)
  router.GET("/locations/:location/forecast/wu", get_wu_forecast_route)
  router.GET("/locations/:location/forecast/wu/raw", get_wu_forecast_raw_route)
  router.GET("/network", network_route)
}
