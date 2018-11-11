import PropTypes from 'prop-types'
import moment from 'moment'
import { Link } from '@rq/react-easy-router'
import Immutable from 'seamless-immutable'
import preventDefaultWrapper from '@rq/prevent-default-wrapper'
import _ from 'underscore'
import React from 'react'
import Store from '../store'
import { data_age } from '../util'

export default function Current(props) {
  return (
    <div>
      <p>Now: {props.current.temp}&deg;</p>
      <p>High: {props.current.temp_max}&deg;</p>
      <p>Low: {props.current.temp_min}&deg;</p>
      <p>Updated: {data_age(props.current)}</p>
    </div>
  )
}

Current.propTypes = {
  current: PropTypes.shape({
    temp: PropTypes.number.isRequired,
    temp_min: PropTypes.number.isRequired,
    temp_max: PropTypes.number.isRequired,
    updated_at: PropTypes.number.isRequired,
  }),
}
