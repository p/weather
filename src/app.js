import preventDefaultWrapper from '@rq/prevent-default-wrapper'
import _ from 'underscore'
import connectToStores from 'alt-utils/lib/connectToStores';
import React from 'react';
import Location from './location'
import Store from './store';

@connectToStores
export default class App extends React.Component {
  constructor(props) {
    super(props)
    
    this.state = {
      locations: props.locations,
    }
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
      <Location
        location_did_submit={this.location_did_submit.bind(this)}
      />
      <p>{this.state.location}</p>
      <h2>Locations</h2>
      {_.map(this.state.locations, (location) => (
          <p key={location}>{location}
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
    let locations = localStorage.getItem('locations') || '[]'
    try {
      locations = JSON.parse(locations)
    } catch(e) {
      console.log('Cannot parse locations: ' + e)
      locations = []
    }
    if (!_.contains(locations, location)) {
      locations.push(location)
      localStorage.setItem('locations', JSON.stringify(locations))
    }
    // TODO make immutable
    locations = Object.assign({}, locations)
    locations[location] = 1
    this.setState({locations})
  }
  
  remove_location(location) {
    let locations = localStorage.getItem('locations') || '[]'
    try {
      locations = JSON.parse(locations)
    } catch(e) {
      console.log('Cannot parse locations: ' + e)
      locations = []
    }
    
    if (_.contains(locations, location)) {
      locations = _.without(locations, location)
      localStorage.setItem('locations', JSON.stringify(locations))
    }
    
    // TODO make immutable
    locations = Object.assign({}, locations)
    locations[location] = 1
    this.setState({locations})
  }
}
