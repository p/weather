import { assert } from 'chai'

import { merge_hourly_into_daily_forecasts } from '../../src/data/merger'

const mhd = merge_hourly_into_daily_forecasts

describe('merge_hourly_into_daily_forecasts', function() {
  describe('1 daily & 1 hourly into day', () => {
    const daily_fcs = [
      {
        start_timestamp: 1,
        expire_timestamp: 1,
        day: {
          start_timestamp: 1,
          temp: 20,
        },
        night: {
          start_timestamp: 2,
          temp: 20,
        },
      },
    ]

    const hourly_fcs = [
      {
        expire_timestamp: 1,
        start_timestamp: 1,
        temp: 21,
      },
    ]

    const expected = [
      {
        start_timestamp: 1,
        expire_timestamp: 1,
        day: {
          start_timestamp: 1,
          temp: 20,
          hourly: hourly_fcs,
        },
        night: {
          start_timestamp: 2,
          temp: 20,
          hourly: [],
        },
      },
    ]

    it('works', function() {
      let actual = mhd(daily_fcs, hourly_fcs)
      assert.deepEqual(actual, expected)
    })
  })
  
  describe('1 daily & 1 hourly into night', () => {
    const daily_fcs = [
      {
        start_timestamp: 1,
        expire_timestamp: 1,
        day: {
          start_timestamp: 1,
          temp: 20,
        },
        night: {
          start_timestamp: 2,
          temp: 20,
        },
      },
    ]

    const hourly_fcs = [
      {
        expire_timestamp: 1,
        start_timestamp: 2,
        temp: 21,
      },
    ]

    const expected = [
      {
        start_timestamp: 1,
        expire_timestamp: 1,
        day: {
          start_timestamp: 1,
          temp: 20,
          hourly: [],
        },
        night: {
          start_timestamp: 2,
          temp: 20,
          hourly: hourly_fcs,
        },
      },
    ]

    it('works', function() {
      let actual = mhd(daily_fcs, hourly_fcs)
      assert.deepEqual(actual, expected)
    })
  })
  
  describe('1 daily & 1 hourly for day not aligned with daily start time', () => {
    const daily_fcs = [
      {
        start_timestamp: 1,
        expire_timestamp: 1,
        day: {
          start_timestamp: 1,
          temp: 20,
        },
        night: {
          start_timestamp: 3,
          temp: 20,
        },
      },
    ]

    const hourly_fcs = [
      {
        expire_timestamp: 1,
        start_timestamp: 2,
        temp: 21,
      },
    ]

    const expected = [
      {
        start_timestamp: 1,
        expire_timestamp: 1,
        day: {
          start_timestamp: 1,
          temp: 20,
          hourly: hourly_fcs,
        },
        night: {
          start_timestamp: 3,
          temp: 20,
          hourly: [],
        },
      },
    ]

    it('works', function() {
      let actual = mhd(daily_fcs, hourly_fcs)
      assert.deepEqual(actual, expected)
    })
  })
  
  describe('2 daily & 1 hourly for second day', () => {
    const daily_fcs = [
      {
        start_timestamp: 2,
        expire_timestamp: 1,
        night: {
          start_timestamp: 2,
          temp: 20,
        },
      },
      {
        start_timestamp: 3,
        expire_timestamp: 7,
        day: {
          start_timestamp: 3,
          temp: 20,
        },
        night: {
          start_timestamp: 4,
          temp: 20,
        },
      },
    ]

    const hourly_fcs = [
      {
        expire_timestamp: 1,
        start_timestamp: 3,
        temp: 21,
      },
    ]

    const expected = [
      {
        start_timestamp: 2,
        expire_timestamp: 1,
        day: undefined,
        night: {
          start_timestamp: 2,
          temp: 20,
          hourly:[],
        },
      },
      {
        start_timestamp: 3,
        expire_timestamp: 7,
        day: {
          start_timestamp: 3,
          temp: 20,
          hourly: hourly_fcs,
        },
        night: {
          start_timestamp: 4,
          temp: 20,
          hourly:[],
        },
      },
    ]

    it('works', function() {
      let actual = mhd(daily_fcs, hourly_fcs)
      assert.deepEqual(actual, expected)
    })
  })
  
  describe('2 daily & 1 hourly for night not aligned with first start time', () => {
    const daily_fcs = [
      {
        start_timestamp: 2,
        expire_timestamp: 1,
        night: {
          start_timestamp: 2,
          temp: 20,
        },
      },
      {
        start_timestamp: 4,
        expire_timestamp: 7,
        day: {
          start_timestamp: 4,
          temp: 20,
        },
        night: {
          start_timestamp: 5,
          temp: 20,
        },
      },
    ]

    const hourly_fcs = [
      {
        expire_timestamp: 1,
        start_timestamp: 3,
        temp: 21,
      },
    ]

    const expected = [
      {
        start_timestamp: 2,
        expire_timestamp: 1,
        day:undefined,
        night: {
          start_timestamp: 2,
          temp: 20,
          hourly: hourly_fcs,
        },
      },
      {
        start_timestamp: 4,
        expire_timestamp: 7,
        day: {
          start_timestamp: 4,
          temp: 20,
          hourly: [],
        },
        night: {
          start_timestamp: 5,
          temp: 20,
          hourly: [],
        },
      },
    ]

    it('works', function() {
      let actual = mhd(daily_fcs, hourly_fcs)
      assert.deepEqual(actual, expected)
    })
  })
})
