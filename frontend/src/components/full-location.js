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
import { merge_hourly_into_daily_forecasts } from '../data/merger'
import React from 'react'
import Store from '../store'
import Current from '../components/current'

export default
@connect(props => ({
  daily_forecasts: [['weather', props.location_query, 'daily_forecasts'], unim],
  hourly_forecasts: [
    ['weather', props.location_query, 'hourly_forecasts'],
    unim,
  ],
  current: [['weather', props.location_query, 'current'], unim],
  location: [['weather', props.location_query, 'location'], unim],
}))
class FullLocation extends React.Component {
  render() {
    return (
      <FullLocationView
        location_query={this.props.location_query}
        daily_forecasts={merge_hourly_into_daily_forecasts(
          this.props.daily_forecasts,
          this.props.hourly_forecasts,
        )}
        current={this.props.current}
        location={this.props.location}
      />
    )
  }
}

FullLocation.propTypes = {
  location_query: PropTypes.string.isRequired,
}
