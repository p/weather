package main

// A weather backend, its primary purpose is to cache weather data
// from weather services.

import (
  "bytes"
  "errors"
  "fmt"
  "github.com/gin-gonic/gin"
  "os"
  "strconv"
  //"io"
  //"io/ioutil"
  "encoding/gob"
  log "github.com/sirupsen/logrus"
  //"regexp"
  "encoding/json"
  "time"

  owm "github.com/briandowns/openweathermap"
  bolt "github.com/coreos/bbolt"
  //"html/template"
)
import "github.com/jasonwinn/geocoder"

//import "net/http"

var owm_api_key string
var db *bolt.DB

type resolved_location struct {
  Lat       float64
  Lng       float64
  CreatedAt int64
}

type current_conditions struct {
  Temp    float64 `json:"temp"`
  TempMin float64 `json:"temp_min"`
  TempMax float64 `json:"temp_max"`

  CreatedAt int64 `json:"created_at"`
}

func persist(bucket_name string, key string, data interface{}) error {
  store := new(bytes.Buffer)
  enc := gob.NewEncoder(store)
  err := enc.Encode(data)
  if err != nil {
    return errors.New("Could not encode: " + err.Error())
  }

  err = db.Update(func(tx *bolt.Tx) error {
    b := tx.Bucket([]byte(bucket_name))
    err := b.Put([]byte(key), store.Bytes())
    return err
  })
  if err != nil {
    return errors.New("Could not store: " + err.Error())
  }
  return nil
}

func resolve_location(location string) (*resolved_location, error) {
  var coords []byte
  err := db.View(func(tx *bolt.Tx) error {
    b := tx.Bucket([]byte("geocodes"))
    coords = b.Get([]byte(location))
    return nil
  })
  if err != nil {
    return nil, errors.New("Could not look up location: " + err.Error())
  }

  var resloc resolved_location
  if coords == nil {
    lat, lng, err := geocoder.Geocode(location)
    lat = lat
    lng = lng
    if err != nil {
      return nil, errors.New("Could not geocode " + location + ": " + err.Error())
    }

    resloc = resolved_location{lat, lng, time.Now().UnixNano()}

    err = persist("geocodes", location, &resloc)
    if err != nil {
      return nil, errors.New("Could not persist: " + err.Error())
    }

    log.Debug(fmt.Sprintf("Geocoded %s to %f,%f", location, resloc.Lat, resloc.Lng))
  } else {
    store := bytes.NewBuffer(coords)
    dec := gob.NewDecoder(store)
    err := dec.Decode(&resloc)
    if err != nil {
      return nil, errors.New("Could not decode: " + err.Error())
    }

    log.Debug(fmt.Sprintf("Retrieved %s from cache as %f,%f", location, resloc.Lat, resloc.Lng))
  }
  return &resloc, nil
}

func get_current_weather(location string,
  resloc resolved_location) (*current_conditions, error) {
  var cc current_conditions
  var encoded []byte

  err := db.View(func(tx *bolt.Tx) error {
    b := tx.Bucket([]byte("current_conditions"))
    encoded = b.Get([]byte(location))
    return nil
  })
  if err != nil {
    return nil, errors.New("Could not look up location: " + err.Error())
  }

  if encoded == nil {
    w, err := owm.NewCurrent("F", "FI", owm_api_key)
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

    cc = current_conditions{
      w.Main.Temp,
      w.Main.TempMin,
      w.Main.TempMax,
      time.Now().UnixNano(),
    }
    err = persist("current_conditions", location, &cc)
    if err != nil {
      return nil, errors.New("Could not persist: " + err.Error())
    }

    log.Debug(fmt.Sprintf("Fetched weather for %s", location))
  } else {
    store := bytes.NewBuffer(encoded)
    dec := gob.NewDecoder(store)
    err := dec.Decode(&cc)
    if err != nil {
      return nil, errors.New("Could not decode: " + err.Error())
    }

    log.Debug(fmt.Sprintf("Retrieved cached weather for %s", location))
  }
  return &cc, nil
}

func list_locations(c *gin.Context) {
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

func get_conditions(c *gin.Context) {
  location := c.Param("location")
  resloc, err := resolve_location(location)
  if err != nil {
    c.String(500, err.Error())
    return
  }
  cc, err := get_current_weather(location, *resloc)
  if err != nil {
    c.String(500, "Could not get weather: "+err.Error())
    return
  }

  render_json(c, cc)
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
      "geocodes", "current_conditions", "forecasts"}

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

  // Disable Console Color
  // gin.DisableConsoleColor()

  debug := os.Getenv("DEBUG")
  if debug == "" {
    gin.SetMode(gin.ReleaseMode)
    log.SetLevel(log.WarnLevel)
  } else {
    log.SetLevel(log.DebugLevel)
  }

  owm_api_key = os.Getenv("OWM_API_KEY")
  if owm_api_key == "" {
    panic("Must have OWM_API_KEY set")
  }

  geocoder_key := os.Getenv("MAPQUEST_API_KEY")
  if geocoder_key == "" {
    panic("Must have MAPQUEST_API_KEY sset")
  }
  geocoder.SetAPIKey(geocoder_key)

  // Creates a gin router with default middleware:
  // logger and recovery (crash-free) middleware
  router := gin.Default()

  //router.LoadHTMLGlob("views/*.html")

  //router.Use(gin.Recovery())

  router.GET("/locations", list_locations)
  router.GET("/locations/:location/current", get_conditions)

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
}
