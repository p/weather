import PropTypes from 'prop-types'
import moment from 'moment'
import { Link } from '@rq/react-easy-router'
import Immutable from 'seamless-immutable'
import preventDefaultWrapper from '@rq/prevent-default-wrapper'
import _ from 'underscore'
import React from 'react'

export default function ForecastDayOfWeek(props) {
  const { forecast } = props
  return moment(forecast.start_timestamp * 1000).format('dddd')
}

ForecastDayOfWeek.propTypes = {
  forecast: PropTypes.shape({
    start_timestamp: PropTypes.number.isRequired,
  }),
}
