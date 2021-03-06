import actions from '../actions'
import { unim } from '../util'
import { connect } from 'nuclear-js-react-addons-chefsplate'
import history from '../history'
import PropTypes from 'prop-types'
import { Link } from '@rq/react-easy-router'
import Immutable from 'seamless-immutable'
import preventDefaultWrapper from '@rq/prevent-default-wrapper'
import _ from 'underscore'
import React from 'react'
import LocationForm from '../components/location-form'
import Store from '../store'

export default
@connect(props => ({
  locations: [['locations'], unim],
}))
class Locations extends React.Component {
  constructor(props) {
    super(props)
    this.state = {
      location: '',
    }
  }
  render() {
    return (
      <div>
        Weather for:
        <LocationForm
          location_did_submit={this.location_did_submit.bind(this)}
        />
        <p>{this.state.location}</p>
        <h2>Locations</h2>
        {_.map(this.props.locations, location => (
          <p key={location}>
            <Link to="BriefLocation" params={{ location: location }}>
              {location}
            </Link>
            &nbsp;
            <a
              href="#"
              onClick={preventDefaultWrapper(
                this.remove_location.bind(this, location),
              )}
            >
              Remove
            </a>
          </p>
        ))}
      </div>
    )
  }

  location_did_submit(location) {
    history.push('/' + location)
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
    } catch (e) {
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
    locations = Immutable(Object.assign({}, locations, { location: 1 }))
    this.setState({ locations })
  }
}

Locations.propTypes = {
  locations: PropTypes.arrayOf(PropTypes.object),
}
