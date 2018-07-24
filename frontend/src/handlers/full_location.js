import {data_age } from '../util'
import PropTypes from 'prop-types'
import moment from 'moment'
import { Link } from 'react-easy-router'
import Immutable from 'seamless-immutable'
import preventDefaultWrapper from '@rq/prevent-default-wrapper'
import _ from 'underscore'
import React from 'react'
import Store from '../store'
import BaseLocation from './base_location'
import Current from '../components/current'

export default class FullLocation extends BaseLocation {
  render() {
    console.log(this.state.forecast)
    return (
      <div>
        <h2>
          {this.state.location
            ? this.state.location.city + ', ' + this.state.location.state
            : this.props.params.location}
        </h2>

        {this.state.current && 
          <Current current={this.state.current}/>
        }

        {this.state.forecast && (
          <div>
            {_.map(this.state.forecast.daily_forecasts, forecast => (
              <div className="forecast-row" key={forecast.time}>
                <div className="forecast-date">
                  <div>{moment(forecast.time * 1000).format('dddd')}</div>
                  <div>{moment(forecast.time * 1000).format('MMM D')}</div>
                </div>

                {forecast.day &&
                  this.render_day_part_forecast('day', forecast.day)}
                {this.render_day_part_forecast('night', forecast.night)}
              </div>
            ))}
            <p>Updated: {data_age(this.state.forecast)}</p>
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
