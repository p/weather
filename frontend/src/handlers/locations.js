import actions from '../actions'
import { unim } from '../util'
import { connect } from 'nuclear-js-react-addons-chefsplate'
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
import Locations from '../components/locations'
import { mapProps } from '@rq/react-map-props'

@mapProps({
  params: {
    location: unescape,
  },
})
@connect(props => ({
  locations: [['locations'], unim],
}))
export default class LocationsHandler extends React.Component {
  componentDidMount() {
    actions.fetch_locations()
  }
  render() {
    return <Locations />
  }
}
