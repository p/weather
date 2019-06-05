import _ from 'underscore'
import {LocalTime} from './ldate'
import * as u from './util'

export function transform_forecasts(payload ) {
  payload = u.merge(payload, {
    hourly_forecasts: _.map(payload.hourly_forecasts,transform_hourly_forecast)
  })
    return payload
}

export function transform_hourly_forecast(hfc){
    return u.merge(hfc,
    {start_ltime:new LocalTime(hfc.start_at)})
}
