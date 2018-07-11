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
  "github.com/jasonwinn/geocoder"
  "github.com/kr/pretty"
)

//import "net/http"

const current_age = 10 * 60 * 1e9
const forecast_age = 60 * 60 * 1e9

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

func (cc current_conditions) GetCreatedAt() int64 {
  return cc.CreatedAt
}

type presented_current_conditions struct {
  Temp    float64 `json:"temp"`
  TempMin float64 `json:"temp_min"`
  TempMax float64 `json:"temp_max"`

  CreatedAt float64 `json:"created_at"`
}

func present_current_conditions(cc *current_conditions) presented_current_conditions {
  return presented_current_conditions{
    cc.Temp,
    cc.TempMin,
    cc.TempMax,
    float64(cc.CreatedAt) / 1e9,
  }
}

type daily_forecast struct {
  Time    int64
  Temp    float64
  TempMin float64
  TempMax float64
  ConditionName string
  ConditionDescription string
}

type forecast struct {
  DailyForecasts []daily_forecast
  CreatedAt      int64 `json:"created_at"`
}

func (f forecast) GetCreatedAt() int64 {
  return f.CreatedAt
}

type presented_daily_forecast struct {
  Time    float64
  Temp    float64
  TempMin float64
  TempMax float64
  ConditionName string
  ConditionDescription string
}

type presented_forecast struct {
  DailyForecasts []presented_daily_forecast
  CreatedAt      float64 `json:"created_at"`
}

func present_forecast(f *forecast) presented_forecast {
  presented_daily := make([]presented_daily_forecast, 0, len(f.DailyForecasts))
  for _, v := range f.DailyForecasts {
    presented_daily = append(presented_daily, presented_daily_forecast{
      float64(v.Time) / 1e9,
      v.Temp,
      v.TempMin,
      v.TempMax,
      v.ConditionName,
      v.ConditionDescription,
    })
  }
  return presented_forecast{
    presented_daily,
    float64(f.CreatedAt) / 1e9,
  }
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

type persistable interface {
  GetCreatedAt() int64
}

func get_weather_with_cache(
  location string,
  resloc resolved_location,
  bucket_name string,
  retriever func(resloc resolved_location) (persistable, error),
) (persistable, error) {
  var p persistable
  var encoded []byte

  err := db.View(func(tx *bolt.Tx) error {
    b := tx.Bucket([]byte(bucket_name))
    encoded = b.Get([]byte(location))
    return nil
  })
  if err != nil {
    return nil, errors.New("Error retrieving from cache: " + err.Error())
  }

  if encoded != nil {
    store := bytes.NewBuffer(encoded)
    dec := gob.NewDecoder(store)
    var temp persistable
    err := dec.Decode(&temp)
    if err != nil {
      return nil, errors.New("Could not decode: " + err.Error())
    }
    typed := temp.(persistable)

    log.Debug(fmt.Sprintf("Retrieved cached data for %s", location))

    if time.Now().UnixNano()-typed.GetCreatedAt() <= current_age {
      p = typed
    }
  }

  if p == nil {
    p, err = retriever(resloc)
    if err != nil {
      return nil, errors.New("Could not retrieve: " + err.Error())
    }

    err = persist(bucket_name, location, &p)
    if err != nil {
      return nil, errors.New("Could not persist: " + err.Error())
    }

    log.Debug(fmt.Sprintf("Fetched data for %s", location))
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

if (log.GetLevel() == log.DebugLevel) {
  fmt.Printf("%# v", pretty.Formatter(w))
  }

  p := current_conditions{
    w.Main.Temp,
    w.Main.TempMin,
    w.Main.TempMax,
    time.Now().UnixNano(),
  }
  return &p, nil
}

func get_current_weather(location string,
  resloc resolved_location) (*current_conditions, error) {

  p, err := get_weather_with_cache(
    location,
    resloc,
    "current_conditions",
    current_retriever)

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

if (log.GetLevel() == log.DebugLevel) {
  fmt.Printf("%# v", pretty.Formatter(w))
  }

  l := w.ForecastWeatherJson.(*owm.Forecast5WeatherData).List
  dailies := make([]daily_forecast, 0, len(l))
  for _, v := range l {
    dailies = append(dailies, daily_forecast{
      int64(v.Dt) * 1e9,
      v.Main.Temp,
      v.Main.TempMin,
      v.Main.TempMax,
      v.Weather[0].Main,
      v.Weather[0].Description,
    })
  }

  p := forecast{
    dailies,
    time.Now().UnixNano(),
  }
  return &p, nil
}

func get_forecast(location string,
  resloc resolved_location) (*forecast, error) {

  p, err := get_weather_with_cache(
    location,
    resloc,
    "forecasts",
    forecast_retriever)

  if err != nil {
    return nil, err
  }

  return p.(*forecast), nil
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

func get_conditions_route(c *gin.Context) {
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

  pcc := present_current_conditions(cc)

  render_json(c, pcc)
}

func get_forecast_route(c *gin.Context) {
  location := c.Param("location")
  resloc, err := resolve_location(location)
  if err != nil {
    c.String(500, err.Error())
    return
  }
  f, err := get_forecast(location, *resloc)
  if err != nil {
    c.String(500, "Could not get weather: "+err.Error())
    return
  }

  pf := present_forecast(f)

  render_json(c, pf)
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
  gob.Register(&current_conditions{})
  gob.Register(&forecast{})

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
  router.GET("/locations/:location/current", get_conditions_route)
  router.GET("/locations/:location/forecast", get_forecast_route)

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
