import FullLocationView from '../views/full-location'
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
import { mapProps } from '@rq/react-map-props'

export default
@connect(props => ({
  forecast: [['forecast', props.location_query, 'forecast'], unim],
  current: [['forecast', props.location_query, 'current'], unim],
  location: [['forecast', props.location_query, 'location'], unim],
}))
class FullLocation extends React.Component {
  render() {
    return (
      <FullLocationView
    location_query={this.props.location_query}
    forecast={this.props.forecast}
    current={this.props.current}
    location={this.props.location}
    />
    )
  }
}

FullLocation.propTypes = {
  location_query: PropTypes.string.isRequired,
}
