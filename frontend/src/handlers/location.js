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
    fetch('http://localhost:8093/locations/' + this.props.params.location + '/' + key)
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
    return <div>
      <h2>{this.props.params.location}</h2>
      
      {this.state.current &&
        <div>
      <p>Now: {this.state.current.temp}&deg;</p>
      <p>Min: {this.state.current.temp_min}&deg;</p>
      <p>Max: {this.state.current.temp_max}&deg;</p>
      <p>Updated: {this.data_age('current')}</p>
      </div>}
      
      {this.state.forecast &&
        <div>
        <ul>
        {_.map(this.state.forecast.daily_forecasts, forecast => <li key={forecast.time}>
      <p>Now: {forecast.temp}&deg;</p>
      <p>Min: {forecast.temp_min}&deg;</p>
      <p>Max: {forecast.temp_max}&deg;</p>
      </li>)}
      </ul>
      <p>Updated: {this.data_age('forecast')}</p>
      </div>}
    </div>
  }
  
  data_age(key) {
    if (this.state[key]) {
      let d = new Date().getTime()/1000 - this.state[key].created_at
      return moment.duration(d, 'seconds').humanize() + ' ago'
    } else {
      return null
    }
  }
}

Location.propTypes = {
  params: PropTypes.shape({
    location: PropTypes.string,
  }),
}
