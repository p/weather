import PropTypes from 'prop-types'
import moment from 'moment'
import { Link } from 'react-easy-router'
import Immutable from 'seamless-immutable'
import preventDefaultWrapper from '@rq/prevent-default-wrapper'
import _ from 'underscore'
import React from 'react'
import Store from '../store'

let NETWORK = NODE_ENV == 'production' ? 0 : 2

export default class BaseLocation extends React.Component {
  constructor(props) {
    super(props)
    this.state = {}
  }

  static getStores() {
    return [Store]
  }

  static getPropsFromStores() {
    return Store.getState()
  }

  componentDidMount() {
    this.load_data('current', 10 * 60 * 1000)
    //this.load_data('forecast', 2*3600*1000)
  }

  load_data(key) {
    let url_key = key
    if (key == 'forecast') {
      url_key = 'forecast/wu'
    }
    fetch(
      API_URL +
        '/locations/' +
        this.props.params.location +
        '?network=' +
        NETWORK,
    )
      .then(resp => resp.json())
      .then(payload => {
        //let state_delta = {}
        //state_delta[key] = payload
        //console.log(state_delta)
        this.setState(payload)
      })
  }

  load_data_periodically(key, interval) {
    this.load_data(key)
    setTimeout(this.load_data_periodically.bind(this, key, interval), interval)
  }

  format_date(timestamp) {
    let date = new Date(timestamp * 1000)
    return moment(date).format('dddd, MMM D')
  }
}

BaseLocation.propTypes = {
  params: PropTypes.shape({
    location: PropTypes.string,
  }),
}
