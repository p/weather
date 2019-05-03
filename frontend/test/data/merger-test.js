import {assert} from 'chai'

import {merge_hourly_into_daily_forecasts} from '../../src/data/merger'

const mhd = merge_hourly_into_daily_forecasts


describe('merge_hourly_into_daily_forecasts', function() {
  describe('1 daily & 1 hourly', ()=>{
const daily_fcs = [{
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
}]

const hourly_fcs = [{
  expire_timestamp: 1,
  start_timestamp: 1,
  temp: 21,
}]

const expected=[{
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
    hourly:[],
  },
}]
  
    it('works', function() {
      let actual=mhd(daily_fcs,hourly_fcs)
      assert.deepEqual(actual,expected)
    });
  })
});
