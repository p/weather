import _ from 'underscore'
import date from 'date-fns'
import {LocalTime} from '../ldate'

// Extracts one day worth of forecasts out of the hourly forecasts.
//
// If the first forecast is for a time before 7 pm, "today" is taken to
// be the day of the first forecast, and forecasts up to 2 am of the following
// day are returned. Otherwise "today" is taken to be the next day, and
// forecasts from 8 am up to but not including 2 am of the day after are returned.
export function extract_today_hourly(hourly_forecasts){
    let out=[]
  
  if (hourly_forecasts.length==0){
    return out
  }
  
  let first=hourly_forecasts[0]
  let time = new LocalTime(first.start_at)
    let i =0
  
  
  if (time.hour>=19){
    // skip forecasts until midnight
    while(i<hourly_forecasts.length){
      let forecast=hourly_forecasts[i]
      
      let time = new LocalTime(forecast.start_at)
      
      if (time.hour<8){
        break
      }
      
      ++i
    }
  }
  
    // skip very early forecasts - prior to 8 am
    while(i<hourly_forecasts.length){
      let forecast=hourly_forecasts[i]
      
      let time = new LocalTime(forecast.start_at)
      
      if (time.hour>=8){
        break
      }
      
      ++i
    }

    // use forecasts until midnight
    while(i<hourly_forecasts.length){
      let forecast=hourly_forecasts[i]
      
      let time = new LocalTime(forecast.start_at)
      
      if (time.hour<2){
        break
      }
      
      out.push(forecast)
      
      ++i
    }

    // use forecasts until 2 am
    while(i<hourly_forecasts.length){
      let forecast=hourly_forecasts[i]
      
      let time = new LocalTime(forecast.start_at)
      
      if (time.hour>=2){
        break
      }
      
      out.push(forecast)
      
      ++i
    }
    
    return out
}

// 12 am: one day forecast, 8 am-2 am
// 9 am: one day forecast, 10 am-2am
// 7 pm: two day forecast, 7 pm-2 am & 8 am-2 am
// 11 pm: two day forecast, 11 pm-2 am & 8 am-2 am
