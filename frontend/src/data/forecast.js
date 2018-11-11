import ReactTimeout from 'react-timeout'
import { network_flag, unim } from '../util'
import { connect } from 'nuclear-js-react-addons-chefsplate'
import PropTypes from 'prop-types'
import moment from 'moment'
import { Link } from '@rq/react-easy-router'
import Immutable from 'seamless-immutable'
import preventDefaultWrapper from '@rq/prevent-default-wrapper'
import _ from 'underscore'
import React from 'react'
import Store from '../store'
import actions from '../actions'

@connect(props => ({
  network: [['network'], unim],
  forecast: [['forecast', props.location_query], unim],
}))
@ReactTimeout
export default class Forecast extends React.Component {
  componentDidMount() {
    if (!this.props.forecast) {
      actions.fetch_forecast(
        this.props.location_query,
        network_flag(this.props.network.up),
      )
    }
    this.props.setInterval(function() {
      actions.fetch_forecast(this.props.location_query)
    }, 10 * 60 * 1000)
  }

  render() {
    if (this.props.forecast) {
      return <div>{this.props.children}</div>
    } else {
      return <div>Loading...</div>
    }
  }
}

Forecast.propTypes = {
  location_query: PropTypes.string.isRequired,
}
