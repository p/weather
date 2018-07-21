package main

// A weather backend, its primary purpose is to cache weather data
// from weather services.

// gob & deserialization:
// http://www.funcmain.com/gob_encoding_an_interface

import (
  //"bytes"
  "encoding/gob"
  "encoding/json"
  "errors"
  "fmt"
  "github.com/gin-gonic/gin"
  log "github.com/sirupsen/logrus"
  "os"
  "regexp"
  "strconv"
  "time"

  bolt "github.com/coreos/bbolt"
  //"html/template"
  "github.com/jasonwinn/geocoder"
  "github.com/kr/pretty"
)

//import "net/http"

const current_age = 10 * 60
const forecast_age = 60 * 60

type NetworkUse int

const (
  NetworkDefault NetworkUse = iota
  NetworkForce
  NetworkSkip
)

var online bool
var owm_api_key string
var db *bolt.DB

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
    c.String(500, "Problem: "+err.Error())
    return
  }

  render_json(c, locations)
}

func get_conditions_route(c *gin.Context) {
  location := c.Param("location")
  resloc, err := resolve_location(location)
  if err != nil {
    c.String(500, err.Error())
    return
  }
  network, err := get_network_flag(c)
  if err != nil {
    c.String(500, err.Error())
    return
  }
  cc, err := get_current_weather(location, *resloc, network)
  if err != nil {
    c.String(500, "Could not get weather: "+err.Error())
    return
  }

  render_json(c, cc)
}

func get_forecast_route(c *gin.Context) {
  location := c.Param("location")
  resloc, err := resolve_location(location)
  if err != nil {
    c.String(500, err.Error())
    return
  }
  f, err := get_forecast(location, *resloc, 0)
  if err != nil {
    c.String(500, "Could not get weather: "+err.Error())
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
    c.String(500, err.Error())
    return
  }
  resloc, err := resolve_location(location)
  if err != nil {
    c.String(500, err.Error())
    return
  }
  f, err := get_wu_forecast(location, *resloc, network)
  if err != nil {
    c.String(500, "Could not get weather: "+err.Error())
    return
  }

  render_json(c, f)
}

func get_wu_forecast_raw_route(c *gin.Context) {
  location := c.Param("location")

  log.Debug(location)
  data, err := lookup("wu_forecasts_raw", location)
  if err != nil {
    c.String(500, err.Error())
    return
  }
  if data != nil {
    render_json(c, data)
    return
  }

  resloc, err := resolve_location(location)
  if err != nil {
    c.String(500, err.Error())
    return
  }
  f, err := get_wu_forecast(location, *resloc, NetworkDefault)
  f = f
  if err != nil {
    c.String(500, "Could not get weather: "+err.Error())
    return
  }

  data, err = lookup("wu_forecasts_raw", location)
  if err != nil {
    c.String(500, err.Error())
    return
  }
  if data != nil {
    render_json(c, data)
    return
  }

  c.String(500, "No wu cached data after retrieving a forecast")
}

func location_route(c *gin.Context) {
  location := c.Param("location")
  network, err := get_network_flag(c)
  if err != nil {
    c.String(500, err.Error())
    return
  }
  resloc, err := resolve_location(location)
  if err != nil {
    c.String(500, err.Error())
    return
  }
  f, err := get_location_everything(location, *resloc, network)
  if err != nil {
    c.String(500, "Could not get weather: "+err.Error())
    return
  }

  render_json(c, f)
  f = f
}

func render_json(c *gin.Context, data interface{}) {
  payload, err := json.Marshal(data)
  if err != nil {
    c.String(500, "Could not jsonify: "+err.Error())
    return
  }

  c.Writer.Header().Set("content-type", "application/json")
  set_cors_headers(c)
  c.String(200, string(payload))
}

func set_cors_headers(c *gin.Context) {
  c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
  c.Writer.Header().Set("Access-Control-Allow-Method", "*")
}

func main() {
  var err error

  wu_api_key_regexp = regexp.MustCompile(WU_API_KEY_REGEXP)

  db_path := os.Getenv("DB_PATH")
  if db_path == "" {
    db_path = "weather.db"
  }
  db, err = bolt.Open(db_path, 0600, nil)
  if err != nil {
    log.Fatal("Error opening database")
  }
  defer db.Close()

  db.Update(func(tx *bolt.Tx) error {
    buckets := []string{
      "geocodes", "current_conditions", "forecasts", "wu_forecasts",
      "wu_forecasts_raw", "config"}

    for index, bucket := range buckets {
      b, err := tx.CreateBucketIfNotExists([]byte(bucket))
      if err != nil {
        log.Fatal("Cannot create " + bucket + " bucket")
      }

      b = b
      index = index
    }
    return nil
  })

  gob.Register(&resolved_location{})
  gob.Register(&current_conditions{})
  gob.Register(&forecast{})
  gob.Register(&wu_credentials{})
  gob.Register(&WuForecast10Response{})

  // Disable Console Color
  // gin.DisableConsoleColor()

  debug := os.Getenv("DEBUG")
  if debug == "" {
    gin.SetMode(gin.ReleaseMode)
    log.SetLevel(log.WarnLevel)
  } else {
    log.SetLevel(log.DebugLevel)
  }

  offline := os.Getenv("OFFLINE")
  if offline == "" {
    online = true

    owm_api_key = os.Getenv("OWM_API_KEY")
    if owm_api_key == "" {
      panic("Must have OWM_API_KEY set")
    }

    geocoder_key := os.Getenv("MAPQUEST_API_KEY")
    if geocoder_key == "" {
      panic("Must have MAPQUEST_API_KEY sset")
    }
    geocoder.SetAPIKey(geocoder_key)
  } else {
    online = false
  }

  // Creates a gin router with default middleware:
  // logger and recovery (crash-free) middleware
  router := gin.Default()

  //router.LoadHTMLGlob("views/*.html")

  //router.Use(gin.Recovery())

  router.GET("/locations", list_locations_route)
  router.GET("/locations/:location", location_route)
  router.GET("/locations/:location/current", get_conditions_route)
  router.GET("/locations/:location/forecast", get_forecast_route)
  router.GET("/locations/:location/forecast/wu", get_wu_forecast_route)
  router.GET("/locations/:location/forecast/wu/raw", get_wu_forecast_raw_route)

  // By default it serves on :8080 unless a
  // PORT environment variable was defined.
  port := os.Getenv("PORT")
  var iport int
  if port == "" {
    iport = 8093
  } else {
    iport, err = strconv.Atoi(port)
    if err != nil {
      log.Fatal(err)
    }
  }
  router.Run(fmt.Sprintf(":%d", iport))
  // router.Run(":3000") for a hard coded port

  a := pretty.Formatter
  a = a
}

func int_ptr_to_float_ptr(v *int) *float64 {
  if v == nil {
    return nil
  } else {
    q := float64(*v)
    return &q
  }
}

func shortcast_maybe(v *WuForecastResponseDaypart) string {
  if v == nil {
    return ""
  } else {
    return v.Shortcast
  }
}

func narrative_maybe(v *WuForecastResponseDaypart) string {
  if v == nil {
    return ""
  } else {
    return v.Narrative
  }
}

func now() float64 {
  return float64(time.Now().UnixNano()) / 1e9
}

type location_everything struct {
  Location resolved_location  `json:"location"`
  Current  current_conditions `json:"current"`
  Forecast forecast           `json:"forecast"`
}

func get_location_everything(location string, resloc resolved_location,
  network NetworkUse) (*location_everything, error) {
  cc, err := get_current_weather(location, resloc, network)
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
