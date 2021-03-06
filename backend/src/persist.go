package main

import (
  "bytes"
  "encoding/gob"
  "errors"

  bolt "github.com/coreos/bbolt"
)

type persistable interface {
  GetUpdatedAt() float64
  GetExpiresAt() float64
}

func persist(bucket_name string, key string, data interface{}) error {
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
    if b == nil {
      return errors.New("Bucket " + bucket_name + " does not exist")
    }
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
  var data interface{}
  err = dec.Decode(&data)
  if err != nil {
    return nil, errors.New("Could not decode: " + err.Error())
  }

  return data, nil
}

func lookup_expiring(bucket_name string, key string) (interface{ persistable }, error) {
  data, err := lookup(bucket_name, key)
  if err != nil {
    return nil, err
  }
  if data == nil {
    return nil, nil
  }
  p := data.(persistable)

  if p.GetExpiresAt() < now() {
    return nil, nil
  }
  return p, nil
}

const start_of_2020 = float64(1577854800)
