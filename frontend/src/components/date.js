import PropTypes from 'prop-types'
import moment from 'moment'
import { Link } from '@rq/react-easy-router'
import Immutable from 'seamless-immutable'
import preventDefaultWrapper from '@rq/prevent-default-wrapper'
import _ from 'underscore'
import React from 'react'
import Store from '../store'
import { data_age } from '../util'

export default function Date(props) {
  let date = new Date(props.timestamp * 1000)
  return moment(date).format('dddd, MMM D')
}

Date.propTypes = {
  timestamp: PropTypes.number.isRequired,
}
