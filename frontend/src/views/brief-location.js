import SingleDayTemp from '../blocks/single-day-temp'
import {
  DailyForecastPropTypes,
  LocationPropTypes,
  DayPartPropTypes,
} from '../data/prop-types'
import PrecipType from '../blocks/precip-type'
import { network_flag, unim } from '../util'
import { data_age } from '../util'
import PropTypes from 'prop-types'
import moment from 'moment'
import { Link } from '@rq/react-easy-router'
import Immutable from 'seamless-immutable'
import preventDefaultWrapper from '@rq/prevent-default-wrapper'
import _ from 'underscore'
import React from 'react'
import Store from '../store'
import Current from '../components/current'

export default class BriefLocationView extends React.Component {
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
                    {moment(forecast.start_timestamp * 1000).format('dddd')}
                  </div>
                  <div>
                    {moment(forecast.start_timestamp * 1000).format('MMM D')}
                  </div>
                </div>

                {forecast.day
                  ? this.render_day_part_forecast('day', forecast.day, forecast)
                  : this.render_day_part_forecast(
                      'night',
                      forecast.night,
                      forecast,
                    )}
              </div>
            ))}
          </div>
        )}
      </div>
    )
  }

  render_day_part_forecast(day_part_name, forecast, full_forecast) {
    return (
      <div className="forecast-row" key={forecast.time}>
        <div className={'forecast-' + day_part_name}>
          <div className="forecast-temp">
            <SingleDayTemp forecast={full_forecast} />
            {'\xb0'}
          </div>
          <div className="forecast-precip">
            {forecast.precip_probability > 10 ? (
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
        </div>
      </div>
    )
  }

  format_short_forecast(name, dpf) {
    return (
      <p>
        {name}:{' '}
        <b>
          {dpf.temp}&deg;, {dpf.precip_type}: {dpf.precip_probability}%
        </b>{' '}
        {dpf.narrative}
      </p>
    )
  }
}

BriefLocationView.propTypes = {
  location_query: PropTypes.string.isRequired,

  location: LocationPropTypes.isRequired,

  daily_forecasts: PropTypes.arrayOf(DailyForecastPropTypes).isRequired,

  current: Current.propTypes.current.isRequired,
}
