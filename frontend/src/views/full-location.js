import ForecastDayOfWeek from '../format/forecast-day-of-week'
import ForecastDate from '../format/forecast-date'
import {
  DailyWithHourlyForecastPropTypes,
  LocationPropTypes,
  DailyForecastPropTypes,
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
import FullDayPartForecastView from './full-day-part-forecast'

export default class FullLocationView extends React.Component {
  render() {
    return (
      <div>
        <h2>
          {this.props.location
            ? this.props.location.city + ', ' + this.props.location.state_abbr
            : this.props.location_query}
        </h2>

        {this.props.current && <Current current={this.props.current} />}

        {this.props.daily_forecasts && (
          <div>
            {_.map(this.props.daily_forecasts, forecast => (
              <div key={forecast.start_timestamp} className="forecast-row">
                <div className="forecast-date">
                  <div>
                    <ForecastDayOfWeek forecast={forecast} />
                  </div>
                  <div>
                    <ForecastDate forecast={forecast} />
                  </div>
                </div>

                {forecast.day &&
                  <FullDayPartForecastView forecast={forecast.day}
                  day_part_name='day'/>
                  }
                  <FullDayPartForecastView forecast={forecast.night}
                  day_part_name='night'/>
              </div>
            ))}
          </div>
        )}
      </div>
    )
  }
}

FullLocationView.propTypes = {
  location_query: PropTypes.string.isRequired,

  location: LocationPropTypes.isRequired,

  daily_forecasts: PropTypes.arrayOf(DailyWithHourlyForecastPropTypes)
    .isRequired,

  current: Current.propTypes.current.isRequired,
}
