import { connect } from 'nuclear-js-react-addons-chefsplate'
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

export default
@connect(props => ({
  forecast: [['forecast', props.location_query, 'forecast'], unim],
  current: [['forecast', props.location_query, 'current'], unim],
  location: [['forecast', props.location_query, 'location'], unim],
}))
class BriefLocation extends React.Component {
  render() {
    //console.log(this.props.forecast)
    return (
      <div>
        <h2>
          {this.props.location
            ? this.props.location.city + ', ' + this.props.location.state_abbr
            : this.props.location_query}
        </h2>

        {this.props.current && <Current current={this.props.current} />}

        {this.props.forecast && (
          <div>
            {_.map(this.props.forecast.daily_forecasts, forecast =>
              forecast.day
                ? this.render_row('day', forecast)
                : this.render_row('night', forecast),
            )}
          </div>
        )}
      </div>
    )
  }

  render_row(day_part_name, forecast) {
    forecast = ForecastPresenter(forecast)
    return (
      <div className="forecast-row" key={forecast.time}>
        <div className="forecast-date">
          <div>{moment(forecast.time * 1000).format('dddd')}</div>
          <div>{moment(forecast.time * 1000).format('MMM D')}</div>
        </div>

        <div className={'forecast-' + day_part_name}>
          <div className="forecast-temp">
            {forecast.temp.toString() + '\xb0'}
          </div>
          <div className="forecast-precip">
            {forecast.precip_probability > 10 ? (
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

function ForecastPresenter(forecast) {
  forecast = Object.assign({}, forecast)
  if (!forecast.temp) {
    forecast.temp = forecast.day ? forecast.day.temp : forecast.night.temp
  }
  return forecast
}
