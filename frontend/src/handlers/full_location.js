import FullLocation from '../components/full-location'
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
class FullLocationHandler extends React.Component {
  render() {
    return (
      <Forecast location_query={this.props.params.location}>
        <FullLocation location_query={this.props.params.location} />
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
