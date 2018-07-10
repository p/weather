import PropTypes from 'prop-types'
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
    fetch('http://localhost:8093/locations/' + this.props.params.location + '/current')
    .then(resp => resp.json())
    .then(payload => {
      this.setState({weather: payload})
    })
  }

  render() {
    return <div>
      <h2>{this.props.params.location}</h2>
      {this.state.weather &&
        <div>
      <p>Now: {this.state.weather.temp}&deg;</p>
      <p>Min: {this.state.weather.temp_min}&deg;</p>
      <p>Max: {this.state.weather.temp_max}&deg;</p>
      </div>}
    </div>
  }
}

Location.propTypes = {
  params: PropTypes.shape({
    location: PropTypes.string,
  }),
}
