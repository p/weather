import PropTypes from 'prop-types'
import moment from 'moment'
import { Link } from '@rq/react-easy-router'
import Immutable from 'seamless-immutable'
import preventDefaultWrapper from '@rq/prevent-default-wrapper'
import _ from 'underscore'
import React from 'react'
import Store from '../store'
import { data_age } from '../util'

export default function PrecipType(props) {
  let date = new Date(props.start_timestamp * 1000)
  let month = date.getMonth() + 1
  if (month >= 4 && month <= 11 && props.precip_type == 'rain') {
    return null
  } else {
    return props.precip_type
  }
}

PrecipType.propTypes = {
  precip_type: PropTypes.string.isRequired,
  start_timestamp: PropTypes.number.isRequired,
}
