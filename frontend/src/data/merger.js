export function merge_hourly_into_daily_forecasts(daily_forecasts,hourly_forecasts){
    return daily_forecasts
    
}

// for each hourlyforecast,
// determine which daily forecast it goes with,
// assuming 7 am/7 pm boundaries or looking at first daily forecast boundary.
// write hourly forecast into daily forecast list for corresponding daily forecast
// order is maintained throughout
// first day may be missing the daily forecast since forecast would start with night one?
