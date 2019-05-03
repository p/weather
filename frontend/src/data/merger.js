import _ from 'underscore'
import * as u from '../util'

export function merge_hourly_into_daily_forecasts(
  daily_forecasts,
  hourly_forecasts,
) {
  if (daily_forecasts.length == 0) {
    return []
  }

  function day_part(forecasts, index) {
    let daily_index = parseInt(index / 2)
    let day = forecasts[daily_index]
    if (index % 2) {
      return day.night
    } else {
      return day.day
    }
  }

  let out_dfcs = _.map(daily_forecasts, dfc =>
    u.merge(dfc, {
      day: dfc.day && u.merge(dfc.day, { hourly: [] }),
      night: u.merge(dfc.night, { hourly: [] }),
    }),
  )

  let daily_index = 0
  if (!daily_forecasts[0].day) {
    ++daily_index
  }
  let current_dfc = day_part(out_dfcs, daily_index)

  _.each(hourly_forecasts, hfc => {
    console.log(daily_forecasts.length*2 , daily_index,
    hfc.start_timestamp,
    day_part(daily_forecasts, daily_index + 1).start_timestamp)
    
    if (
      daily_forecasts.length*2 > daily_index &&
      hfc.start_timestamp >=
        day_part(daily_forecasts, daily_index + 1).start_timestamp
    ) {
      ++daily_index
      current_dfc = day_part(out_dfcs, daily_index)
    }

    current_dfc.hourly.push(hfc)
  })

  return out_dfcs
}

// for each hourlyforecast,
// determine which daily forecast it goes with,
// assuming 7 am/7 pm boundaries or looking at first daily forecast boundary.
// write hourly forecast into daily forecast list for corresponding daily forecast
// order is maintained throughout
// first day may be missing the daily forecast since forecast would start with night one?
