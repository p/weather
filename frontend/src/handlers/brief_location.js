import BriefLocation from '../components/brief-location'
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

@mapProps({
  params: {
    location: unescape,
  },
})
export default class BriefLocationHandler extends React.Component {
  render() {
    return (
      <Forecast location_query={this.props.params.location}>
        <BriefLocation location_query={this.props.params.location} />
        <Link
          to="FullLocation"
          params={{ location: this.props.params.location }}
        >
          Full View
        </Link>
      </Forecast>
    )
  }
}
