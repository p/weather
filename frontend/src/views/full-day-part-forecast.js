import ForecastDayOfWeek from '../format/forecast-day-of-week'
import Temp from '../format/temp'
import ForecastDate from '../format/forecast-date'
import ForecastTime from '../format/forecast-time'
import {
  DayPartWithHourlyPropTypes
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

export default class FullDayPartForecastView extends React.Component {
  render() {
      let {day_part_name,forecast} = this.props
    return (
      <div className={'forecast-' + day_part_name}>
        <div className="forecast-temp"><Temp temp={forecast.temp}/></div>
        <div className="forecast-precip">
          {forecast.precip_probability ? (
            <div>
              <div>{forecast.precip_probability}%</div>
              <div>
                <PrecipType
                  precip_type={forecast.precip_type}
                  start_timestamp={forecast.start_timestamp}
                />
              </div>
            </div>
          ) : (
            ''
          )}
        </div>
        <div className="forecast-blurb">{forecast.narrative}</div>
        <div>
        {_.map(forecast.hourly,hfc=><div key={hfc.start_timestamp}>
        <ForecastTime forecast={hfc}/>
        <Temp temp={hfc.temp}/>
        </div>)}
        </div>
      </div>
    )
  }
}

FullDayPartForecastView.propTypes = {
    day_part_name: PropTypes.string.isRequired,
  forecast: DayPartWithHourlyPropTypes.isRequired,
}
