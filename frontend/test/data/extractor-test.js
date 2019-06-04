import _ from 'underscore'
import * as u from '../../src/util'
import { assert } from 'chai'

import { extract_today_hourly } from '../../src/data/extractor'

const nyc = require('../fixtures/nyc.json')

describe('extract_today_hourly', function() {
  const forecasts = nyc.hourly_forecasts

  const expected_timestamps = [
    '2019-05-22T08:00:00-04:00',
    '2019-05-22T09:00:00-04:00',
    '2019-05-22T10:00:00-04:00',
    '2019-05-22T11:00:00-04:00',
    '2019-05-22T12:00:00-04:00',
    '2019-05-22T13:00:00-04:00',
    '2019-05-22T14:00:00-04:00',
    '2019-05-22T15:00:00-04:00',
    '2019-05-22T16:00:00-04:00',
    '2019-05-22T17:00:00-04:00',
    '2019-05-22T18:00:00-04:00',
    '2019-05-22T19:00:00-04:00',
    '2019-05-22T20:00:00-04:00',
    '2019-05-22T21:00:00-04:00',
    '2019-05-22T22:00:00-04:00',
    '2019-05-22T23:00:00-04:00',
    '2019-05-23T00:00:00-04:00',
    '2019-05-23T01:00:00-04:00',
  ]

  it('works for just past midnight', () => {
    let actual = extract_today_hourly(forecasts)
    let actual_timestamps = _.map(actual, forecast => forecast.start_at)
    assert.deepEqual(actual_timestamps, expected_timestamps)
  })
})
