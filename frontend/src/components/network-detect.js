import actions from '../actions'
import { unim } from '../util'
import PropTypes from 'prop-types'
import moment from 'moment'
import { Link } from '@rq/react-easy-router'
import Immutable from 'seamless-immutable'
import preventDefaultWrapper from '@rq/prevent-default-wrapper'
import _ from 'underscore'
import React from 'react'
import Store from '../store'
import { data_age } from '../util'
import { connect } from 'nuclear-js-react-addons-chefsplate'

export default
@connect(props => ({
  network: [['network'], unim],
}))
class NetworkDetect extends React.Component {
  componentDidMount() {
    if (!('up' in this.props.network)) {
      actions.fetch_network()
    }
  }
  render() {
    if ('up' in this.props.network) {
      return <div>{this.props.children}</div>
    } else {
      return <div>Loading...</div>
    }
  }
}

NetworkDetect.propTypes = {
  network: PropTypes.shape({
    up: PropTypes.bool,
  }),
}
