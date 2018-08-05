package main

import (
  "time"
)

func int_ptr_to_float_ptr(v *int) *float64 {
  if v == nil {
    return nil
  } else {
    q := float64(*v)
    return &q
  }
}

func now() float64 {
  return float64(time.Now().UnixNano()) / 1e9
}
