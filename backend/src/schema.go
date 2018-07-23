package main

import (
  bolt "github.com/coreos/bbolt"
  log "github.com/sirupsen/logrus"
)

const schema_version = 1
const dep_schema_version = 1

func recreate_buckets(tx *bolt.Tx, buckets []string) error {
  for index, bucket := range buckets {
    err := tx.DeleteBucket([]byte(bucket))
    if err != nil {
      return err
    }
    b, err := tx.CreateBucket([]byte(bucket))
    if err != nil {
      return err
    }

    b = b
    index = index
  }
  return nil

}

func create_buckets() error {
  err := db.Update(func(tx *bolt.Tx) error {
    buckets := []string{
      "geocodes", "current_conditions", "forecasts", "wu_forecasts",
      "wu_currents_raw",
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
  return err
}

func check_schema() error {
  err := db.Update(func(tx *bolt.Tx) error {
    p, err := lookup("config", "schema_version")
    if err != nil {
      return err
    }
    if p != nil {
      v := p.(int)
      if v < schema_version {
        err := recreate_buckets(tx, []string{"current_conditions", "forecasts"})
        if err != nil {
          return err
        }
      }
    }
    return nil
  })
  persist("config", "schema_version", schema_version)
  return err

}
