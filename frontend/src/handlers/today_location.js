import TodayLocation from '../components/today-location'
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
import Forecast from '../data/forecast'
import { mapProps } from '@rq/react-map-props'

export default
@mapProps({
  params: {
    location: unescape,
  },
})
class TodayLocationHandler extends React.Component {
  render() {
    return (
      <Forecast location_query={this.props.params.location}>
        <TodayLocation location_query={this.props.params.location} />
        <Link
          to="FullLocation"
          params={{ location: this.props.params.location }}
        >
          Full View
        </Link>
          {' '}
        <Link
          to="BriefLocation"
          params={{ location: this.props.params.location }}
        >
          Brief View
        </Link>
      </Forecast>
    )
  }
}
