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

        {this.props.daily_forecasts&& (
          <div>
            {_.map(this.props.daily_forecasts, forecast => (
              <div key={forecast.time} className="forecast-row" key={forecast.time}>
                <div className="forecast-date">
                  <div>{moment(forecast.time * 1000).format('dddd')}</div>
                  <div>{moment(forecast.time * 1000).format('MMM D')}</div>
                </div>

                {forecast.day &&
                  this.render_day_part_forecast('day', forecast.day)}
                {this.render_day_part_forecast('night', forecast.night)}
              </div>
            ))}
            <p>Updated: {data_age(this.props.forecast)}</p>
          </div>
        )}
      </div>
    )
  }

  render_day_part_forecast(day_part_name, forecast) {
    return (
      <div className={'forecast-' + day_part_name}>
        <div className="forecast-temp">{forecast.temp.toString() + '\xb0'}</div>
        <div className="forecast-precip">
          {forecast.precip_probability ? (
            <div>
              <div>{forecast.precip_probability}%</div>
              <div>{forecast.precip_type}</div>
            </div>
          ) : (
            ''
          )}
        </div>
        <div className="forecast-blurb">{forecast.narrative}</div>
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

const DayPartPropTypes =PropTypes.shape({
        temp: PropTypes.number.isRequired,
        precip_probability: PropTypes.number.isRequired,
        precip_type: PropTypes.string.isRequired,
      })

FullLocationView.propTypes = {
  location_query: PropTypes.string.isRequired,

  location: PropTypes.shape({
    city: PropTypes.string.isRequired,
    state_abbr: PropTypes.string.isRequired,
  }),

  daily_forecasts: PropTypes.arrayOf(
    PropTypes.shape({
      // UTC timestamp
      time: PropTypes.number.isRequired,
      // UTC timestamp
      expires_at: PropTypes.number.isRequired,
      day: DayPartPropTypes,
      night: DayPartPropTypes,
    }),
  ),

  current: Current.propTypes.current,
}
