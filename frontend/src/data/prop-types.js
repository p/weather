import PropTypes from 'prop-types'

export const DayPartPropTypes = PropTypes.shape({
  temp: PropTypes.number.isRequired,
  precip_probability: PropTypes.number.isRequired,
  precip_type: PropTypes.string.isRequired,
  narrative: PropTypes.string.isRequired,
})

export const DailyForecastPropTypes = PropTypes.shape({
  // UTC timestamp
  start_timestamp: PropTypes.number.isRequired,
  // UTC timestamp
  expire_timestamp: PropTypes.number.isRequired,
  day: DayPartPropTypes,
  night: DayPartPropTypes,
})

export const LocationPropTypes=PropTypes.shape({
    city: PropTypes.string.isRequired,
    state_abbr: PropTypes.string.isRequired,
  })
  