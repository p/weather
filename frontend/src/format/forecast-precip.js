import PrecipType from './precip-type'
import PropTypes from 'prop-types'
import moment from 'moment'
import { Link } from '@rq/react-easy-router'
import Immutable from 'seamless-immutable'
import preventDefaultWrapper from '@rq/prevent-default-wrapper'
import _ from 'underscore'
import React from 'react'

export default function ForecastPrecip(props) {
  const { forecast } = props
  return <span>
  {forecast.precip_probability}%
  <PrecipType precip_type={forecast.precip_type} start_timestamp={forecast.start_timestamp}/>
  </span>
}

ForecastPrecip.propTypes = {
  forecast: PropTypes.shape({
  start_timestamp: PropTypes.number.isRequired,
  precip_probability: PropTypes.number.isRequired,
  precip_type: PropTypes.string.isRequired,
  }),
}
