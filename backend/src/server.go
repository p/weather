package main

// A weather backend, its primary purpose is to cache weather data
// from weather services.

import (
	//"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"os"
	"strconv"
	//"sync"
	//"io"
	//"io/ioutil"
		"encoding/gob"
	"log"
	"regexp"
	//"strings"
	"time"

	owm "github.com/briandowns/openweathermap"
	bolt "github.com/coreos/bbolt"
	"html/template"
)
import "github.com/jasonwinn/geocoder"
//import "net/http"

var owm_api_key string

type entry struct {
	text        string
	received_at int64
}

var entries []entry

func twiml(forward_number string, from_number string, text string) string {
	twiml_template := `
<?xml version='1.0' encoding='UTF-8'?>
<Response>
    <Message to='%s'>[OTPBASE:%s] %s</Message>
</Response>
`
	return fmt.Sprintf(twiml_template, forward_number, from_number, text)
}

var http_user, http_password string
var forward_number string
var ticker *time.Ticker
var code_regexp *regexp.Regexp
var apps_template *template.Template
var db *bolt.DB

func conditions(c *gin.Context) {
	location := c.Param("location")
	var coords []byte
	db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("geocodes"))
		coords = b.Get([]byte(location))
		return nil
	})
	
	if coords == nil {
		lat, lng, err := geocoder.Geocode(location)
		lat=lat
		lng=lng
		if err != nil {
			c.String(500, "Could not geocode " + location + ": "+err.Error())
			return
		}
		
		db.Update(func(tx *bolt.Tx) error {
			b := tx.Bucket([]byte("geocodes"))
			err := b.Put([]byte(location), coords)
			return err
		})
	} else {
	}
	
    w, err := owm.NewCurrent("F", "FI", owm_api_key)
    if err != nil {
        log.Fatalln(err)
    }

    w.CurrentByCoordinates(
            &owm.Coordinates{
                Longitude: -112.07,
                Latitude: 33.45,
            },
    )

	c.String(200, "r")
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
		b, err := tx.CreateBucketIfNotExists([]byte("geocodes"))
		if err != nil {
			log.Fatal("Cannot create geocodes bucket")
		}
		b, err = tx.CreateBucketIfNotExists([]byte("conditions"))
		if err != nil {
			log.Fatal("Cannot create conditions bucket")
		}
		b = b
		return nil
	})
	
	gob.Register(&owm.Coordinates{})

	// Disable Console Color
	// gin.DisableConsoleColor()

	debug := os.Getenv("DEBUG")
	if debug == "" {
		gin.SetMode(gin.ReleaseMode)
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

	router.GET("/:location", conditions)

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
