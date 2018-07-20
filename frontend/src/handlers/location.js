import PropTypes from 'prop-types'
import moment from 'moment'
import {Link} from 'react-easy-router'
import Immutable from 'seamless-immutable'
import preventDefaultWrapper from '@rq/prevent-default-wrapper'
import _ from 'underscore'
import connectToStores from 'alt-utils/lib/connectToStores';
import React from 'react';
import Store from '../store';

export default class Location extends React.Component {
  constructor(props) {
    super(props)
    this.state = {}
  }
  
  static getStores() {
    return [Store]
  }
    
  static getPropsFromStores() {
    return Store.getState()
  }
  
  componentDidMount() {
    this.load_data('current', 10*60*1000)
    this.load_data('forecast', 2*3600*1000)
  }
  
  load_data(key) {
    let url_key = key
    if (key == 'forecast') {
      url_key = 'forecast/wu'
    }
    fetch('http://localhost:8093/locations/' + this.props.params.location + '/' + url_key + '?network=2')
    .then(resp => resp.json())
    .then(payload => {
      let state_delta = {}
      state_delta[key] = payload
      console.log(state_delta)
      this.setState(state_delta)
    })
  }
  
  load_data_periodically(key, interval) {
    this.load_data(key)
    setTimeout(this.load_data_periodically.bind(this, key, interval), interval)
  }

  render() {
    console.log(this.state.forecast)
    return <div>
      <h2>{this.props.params.location}</h2>
      
      {this.state.current &&
        <div>
      <p>Now: {this.state.current.temp}&deg;</p>
      <p>High: {this.state.current.temp_max}&deg;</p>
      <p>Low: {this.state.current.temp_min}&deg;</p>
      <p>Updated: {this.data_age('current')}</p>
      </div>}
      
      {this.state.forecast &&
        <div>
        {_.map(this.state.forecast.daily_forecasts, forecast => 
          <div className='forecast-row' key={forecast.time}>
      
        <div className='forecast-date'>
          <div>{moment(forecast.time*1000).format('dddd')}</div>
          <div>{moment(forecast.time*1000).format('MMM D')}</div>
        </div>
        
    {forecast.day &&
      <div className='forecast-day'>
        <div className='forecast-temp'>{forecast.day.temp.toString() + '\xb0'}</div>
        <div className='forecast-precip'>
          <div>{forecast.day.precip_probability}%</div>
          <div>{forecast.day.precip_type}</div>
        </div>
        <div className='forecast-blurb'>
          {forecast.day.condition_description}
        </div>
      </div>
    }
      <div className='forecast-night'>
        <div className='forecast-temp'>{forecast.night.temp.toString() + '\xb0'}</div>
        <div className='forecast-precip'>
          <div>{forecast.night.precip_probability}%</div>
          <div>{forecast.night.precip_type}</div>
        </div>
        <div className='forecast-blurb'>
          {forecast.night.condition_description}
        </div>
      </div>
      </div>)}
      <p>Updated: {this.data_age('forecast')}</p>
      </div>}
    </div>
  }
  
  format_short_forecast(name, dpf) {
      return <p>{name}:
        {' '}
        <b>{dpf.temp}&deg;,
        {' '}
        {dpf.precip_type}:
        {' '}
        {dpf.precip_probability}%</b>
        {' '}
        {dpf.condition_description}</p>
  }
  
  data_age(key) {
    if (this.state[key]) {
      let d = new Date().getTime()/1000 - this.state[key].updated_at
      return moment.duration(d, 'seconds').humanize() + ' ago'
    } else {
      return null
    }
  }
  
  format_date(timestamp) {
    let date = new Date(timestamp*1000)
    return moment(date).format('dddd, MMM D')
  }
}

Location.propTypes = {
  params: PropTypes.shape({
    location: PropTypes.string,
  }),
}
