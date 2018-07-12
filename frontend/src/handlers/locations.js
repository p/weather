import PropTypes from 'prop-types'
import {Link} from 'react-easy-router'
import Immutable from 'seamless-immutable'
import preventDefaultWrapper from '@rq/prevent-default-wrapper'
import _ from 'underscore'
import connectToStores from 'alt-utils/lib/connectToStores';
import React from 'react';
import LocationForm from '../location-form'
import Store from '../store';

@connectToStores
export default class Locations extends React.Component {
  constructor(props) {
    super(props)
    
    this.state = {
      locations: props.locations,
    }
  }
  
  componentDidMount() {
    fetch('http://localhost:8093/locations')
    .then(resp => resp.json())
    .then(payload => {
      this.setState({locations: payload})
    })
  }
  
  static getStores() {
    return [Store]
  }
    
  static getPropsFromStores() {
    return Store.getState()
  }

  render() {
    return <div>
      Weather for:
      <LocationForm
        location_did_submit={this.location_did_submit.bind(this)}
      />
      <p>{this.state.location}</p>
      <h2>Locations</h2>
      {_.map(this.state.locations, (location) => (
          <p key={location}>
          <Link to='Location' params={{location:location}}>{location}</Link>
          &nbsp;
            <a href='#' onClick={preventDefaultWrapper(this.remove_location.bind(this, location))}>
              Remove</a>
          </p>
        ))}
    </div>
  }

  location_did_submit(location) {
    this.setState({location: location})
    this.add_location(location)
  }
  
  add_location(location) {
    let locations = this.load_locations()

    if (!_.contains(locations, location)) {
      locations.push(location)
      this.save_locations(locations)
    }
    this.save_locations_to_state(locations)
  }
  
  remove_location(location) {
    let locations = this.load_locations()
    
    if (_.contains(locations, location)) {
      locations = _.without(locations, location)
      this.save_locations(locations)
    }
    
    this.save_locations_to_state(locations)
  }
  
  load_locations() {
    let locations = localStorage.getItem('locations') || '{}'
    try {
      locations = JSON.parse(locations)
    } catch(e) {
      console.log('Cannot parse locations: ' + e)
      locations = {}
    }
    return locations
  }
  
  save_locations(locations) {
    localStorage.setItem('locations', JSON.stringify(locations))
  }
  
  save_locations_to_state(locations) {
    // TODO make immutable
    locations = Immutable(Object.assign({}, locations, {location: 1}))
    this.setState({locations})
  }
}

Locations.propTypes = {
  locations: PropTypes.arrayOf(PropTypes.object),
}
