import { DailyForecastPropTypes, DayPartPropTypes } from '../data/prop-types'
import PropTypes from 'prop-types'
import moment from 'moment'
import { Link } from '@rq/react-easy-router'
import Immutable from 'seamless-immutable'
import preventDefaultWrapper from '@rq/prevent-default-wrapper'
import _ from 'underscore'
import React from 'react'
import Store from '../store'
import { data_age } from '../util'

export default function SingleDayTemp(props) {
  if (props.forecast.day) {
    return props.forecast.day.temp
  } else {
    return props.forecast.night.temp
  }
}

SingleDayTemp.propTypes = {
  forecast: DailyForecastPropTypes.isRequired,
}
