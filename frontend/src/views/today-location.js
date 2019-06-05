import ForecastPrecip from '../format/forecast-precip'
import ForecastDayOfWeek from '../format/forecast-day-of-week'
import ForecastDate from '../format/forecast-date'
import SingleDayTemp from '../format/single-day-temp'
import {extract_today_hourly} from '../data/extractor'
import {
  TransformedHourlyForecastPropTypes,
  DailyForecastPropTypes,
  LocationPropTypes,
  DayPartPropTypes,
} from '../data/prop-types'
import PrecipType from '../format/precip-type'
import { network_flag, unim } from '../util'
import { data_age } from '../util'
import PropTypes from 'prop-types'
import { Link } from '@rq/react-easy-router'
import Immutable from 'seamless-immutable'
import preventDefaultWrapper from '@rq/prevent-default-wrapper'
import _ from 'underscore'
import React from 'react'
import Store from '../store'
import Current from '../components/current'

export default class TodayLocationView extends React.Component {
  render() {
    const extracted_hfcs = extract_today_hourly(this.props.hourly_forecasts)
    const first_hfc=extracted_hfcs[0]
    
    return (
      <div>
        <h2>
          {this.props.location
            ? this.props.location.city + ', ' + this.props.location.state_abbr
            : this.props.location_query}
        </h2>
        
        <div>
          {first_hfc.start_ltime.format('ddd MMM D')}
        </div>

        {extracted_hfcs && (
          <table>
              <tbody>
            {_.map(extracted_hfcs, forecast => (
                <tr key={forecast.start_timestamp} className="forecast-row">
                  <td>
                    {forecast.start_ltime.format('h a')}
                  </td>
                    <td>{forecast.temp}&deg;</td>
                    <td><ForecastPrecip forecast={forecast}/></td>
                  <td>
                    {forecast.phrase}
                  </td>
                </tr>
            ))}
          </tbody></table>
        )}
      </div>
    )
  }
}

TodayLocationView.propTypes = {
  location_query: PropTypes.string.isRequired,

  location: LocationPropTypes.isRequired,

  daily_forecasts: PropTypes.arrayOf(DailyForecastPropTypes).isRequired,
  hourly_forecasts: PropTypes.arrayOf(TransformedHourlyForecastPropTypes).isRequired,

  current: Current.propTypes.current.isRequired,
}
