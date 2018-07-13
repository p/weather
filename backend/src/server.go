package main

// A weather backend, its primary purpose is to cache weather data
// from weather services.

// gob & deserialization:
// http://www.funcmain.com/gob_encoding_an_interface

import (
  "bytes"
  "errors"
  "fmt"
  "github.com/gin-gonic/gin"
  "os"
  "strconv"
  "io"
  "io/ioutil"
  "encoding/gob"
  "net/http"
  log "github.com/sirupsen/logrus"
  "regexp"
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

var online bool
var owm_api_key string
var db *bolt.DB

type resolved_location struct {
  Lat       float64
  Lng       float64
  CreatedAt int64
}

func (resloc resolved_location) GetCreatedAt() int64 {
  return resloc.CreatedAt
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
  Time    float64 `json:"time"`
  TempMin float64 `json:"temp_min"`
  TempMax float64 `json:"temp_max"`
  ConditionName string `json:"condition_name"`
  ConditionDescription string `json:"condition_description"`
}

type presented_forecast struct {
  DailyForecasts []presented_daily_forecast `json:"daily_forecasts"`
  CreatedAt      float64 `json:"created_at"`
}

func present_forecast(f *forecast) presented_forecast {
  presented_daily := make([]presented_daily_forecast, 0, len(f.DailyForecasts))
  for _, v := range f.DailyForecasts {
    presented_daily = append(presented_daily, presented_daily_forecast{
      float64(v.Time) / 1e9,
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

func persist(bucket_name string, key string, data persistable) error {
  store := new(bytes.Buffer)
  enc := gob.NewEncoder(store)
  err := enc.Encode(&data)
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

func lookup(bucket_name string, key string) (interface{}, error) {
  var encoded []byte
  err := db.View(func(tx *bolt.Tx) error {
    b := tx.Bucket([]byte(bucket_name))
    encoded = b.Get([]byte(key))
    return nil
  })
  if err != nil {
    return nil, errors.New("Error retrieving data from db: " + err.Error())
  }
  
  if encoded == nil {
    return nil, nil
  }

  store := bytes.NewBuffer(encoded)
  dec := gob.NewDecoder(store)
  var data persistable
  err = dec.Decode(&data)
  if err != nil {
    return nil, errors.New("Could not decode: " + err.Error())
  }
  
  return data, nil
}

func resolve_location(location string) (*resolved_location, error) {
  data, err := lookup("geocodes", location)
  log.Debug("location: "+location)
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

    resloc = resolved_location{lat, lng, time.Now().UnixNano()}

    err = persist("geocodes", location, &resloc)
    if err != nil {
      return nil, errors.New("Could not persist: " + err.Error())
    }

    log.Debug(fmt.Sprintf("Geocoded %s to %f,%f", location, resloc.Lat, resloc.Lng))
  } else {
    resloc = *(data.(* resolved_location))

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

  data, err := lookup(bucket_name, location)
  if err != nil {
    return nil, errors.New("Error retrieving from cache: " + err.Error())
  }
  if data != nil {
    typed := data.(persistable)

    log.Debug(fmt.Sprintf("Retrieved cached data for %s", location))

    if !online || time.Now().UnixNano()-typed.GetCreatedAt() <= current_age {
      p = typed
    }
  }

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

func get_wu_forecast(location string,
  resloc resolved_location) (*forecast, error) {

  p, err := get_weather_with_cache(
    location,
    resloc,
    "wu_forecasts",
    wu_forecast_retriever)

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

func get_wu_forecast_route(c *gin.Context) {
  location := c.Param("location")
  resloc, err := resolve_location(location)
  if err != nil {
    c.String(500, err.Error())
    return
  }
  f, err := get_wu_forecast(location, *resloc)
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

type wu_credentials struct {
  ApiKey string
  UpdatedAt int64
}

func (x wu_credentials) GetCreatedAt() int64 {
  return x.UpdatedAt
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
      "geocodes", "current_conditions", "forecasts", "wu_forecasts", "config"}

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

  router.GET("/locations", list_locations)
  router.GET("/locations/:location/current", get_conditions_route)
  router.GET("/locations/:location/forecast", get_forecast_route)
  router.GET("/locations/:location/forecast/wu", get_wu_forecast_route)

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
  
  a:=ioutil.ReadAll
  a=a
  
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
          time.Now().UnixNano(),
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
      time.Now().UnixNano(),
    }
    err = persist("config", "wu_api_key", &wu_creds)
    if err != nil {
      return "", errors.New("Could not persist wu api key: " + err.Error())
    }
  }
  return api_key, err
}

type WuForecastResponseMetadata struct {
  Language string `json:"language"`
  TransactionId string `json:"transaction_id"`
  Version string `json:"version"`
  Latitude float64 `json:"latitude"`
  Longitude float64 `json:"longitude"`
  Units string `json:"units"`
  ExpireTimeGmt int64 `json:"expire_time_gmt"`
  StatusCode int `json:"status_code"`
}

type WuForecastResponseForecast struct {
  Class string `json:"class"`
  ExpireTimeGmt int64 `json:"expire_time_gmt"`
  FcstValid int64 `json:"fcst_valid"`
  FcstValidLocal string `json:"fcst_valid_local"`
  Num int `json:"num"`
  MaxTemp json.RawMessage `json:"max_temp"`
  MinTemp json.RawMessage `json:"min_temp"`
}

type WuForecastResponse struct {
  Metadata WuForecastResponseMetadata `json:"metadata"`
  Forecasts []WuForecastResponseForecast `json:"forecasts"`
  
}

func wu_forecast_retriever(resloc resolved_location) (persistable, error) {
  api_key, err := get_wu_api_key_cached()
  if err != nil {
      return nil, errors.New("Error retrieving wu api key:" + err.Error())
  }
  log.Debug(api_key)
  url := fmt.Sprintf("https://api.weather.com/v1/geocode/%f/%f/forecast/daily/10day.json?apiKey=%s&units=e", resloc.Lat, resloc.Lng, api_key)
  res, err := http.Get(url)
  if err != nil {
      return nil, errors.New("Error retrieving wu forecast:" + err.Error())
  }
  defer res.Body.Close()

  var payload WuForecastResponse
  dec := json.NewDecoder(res.Body)
  err = dec.Decode(&payload)
  if err != nil {
      return nil, errors.New("Could not decode wu forecast:" + err.Error())
  }
  
if (log.GetLevel() == log.DebugLevel) {
  fmt.Printf("%# v", pretty.Formatter(payload))
  }
  
  dailies := make([]daily_forecast, 0)
  
  f := forecast{
    dailies,
    time.Now().UnixNano(),
  }
  
  return &f, nil
}
