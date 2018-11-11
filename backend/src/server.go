package main

// A weather backend, its primary purpose is to cache weather data
// from weather services.

// gob & deserialization:
// http://www.funcmain.com/gob_encoding_an_interface

import (
  //"bytes"
  "encoding/gob"
  "encoding/json"
  "fmt"
  "os"
  "regexp"
  "strconv"

  "github.com/gin-gonic/gin"
  log "github.com/sirupsen/logrus"
  "gopkg.in/weather.v0"

  //"time"

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

var debug bool
var online bool
var owm_api_key string
var db *bolt.DB

// err can be nil here
func return_500(c *gin.Context, msg string, err error) {
  var full_msg string
  if err != nil {
    full_msg = msg + ": " + err.Error()
  } else {
    full_msg = msg
  }
  log.Warn(full_msg)
  if debug {
    c.String(500, full_msg)
  } else {
    c.String(500, msg)
  }
}
func render_json(c *gin.Context, data interface{}) {
  payload, err := json.Marshal(data)
  if err != nil {
    return_500(c, "Could not jsonify response", err)
    return
  }

  c.Writer.Header().Set("content-type", "application/json")
  set_cors_headers(c)
  c.String(200, string(payload))
}

func set_cors_headers(c *gin.Context) {
  c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
  c.Writer.Header().Set("Access-Control-Allow-Methods", "GET")
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

  err = create_buckets()
  if err != nil {
    panic(err)
  }
  err = check_schema()
  if err != nil {
    panic(err)
  }

  gob.Register(&resolved_location{})
  gob.Register(&current_conditions{})
  gob.Register(&forecast{})
  gob.Register(&wu_credentials{})
  gob.Register(&weather.Forecast10Response{})

  // Disable Console Color
  // gin.DisableConsoleColor()

  _debug := os.Getenv("DEBUG")
  if _debug == "" {
    gin.SetMode(gin.ReleaseMode)
    log.SetLevel(log.WarnLevel)
    debug = false
  } else {
    log.SetLevel(log.DebugLevel)
    debug = true
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

  define_routes(router)

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
  log.Info(fmt.Sprintf("Listening on port %d", iport))
  router.Run(fmt.Sprintf(":%d", iport))
  // router.Run(":3000") for a hard coded port

  a := pretty.Formatter
  a = a
}
