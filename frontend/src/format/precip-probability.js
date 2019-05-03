import PropTypes from 'prop-types'
import moment from 'moment'
import { Link } from '@rq/react-easy-router'
import Immutable from 'seamless-immutable'
import preventDefaultWrapper from '@rq/prevent-default-wrapper'
import _ from 'underscore'
import React from 'react'
import Store from '../store'
import { data_age } from '../util'

export default function PrecipProbability(props) {
  const { precip_probability } = props
  return `${precip_probability}%`
}

PrecipProbability.propTypes = {
  precip_probability: PropTypes.number.isRequired,
}
