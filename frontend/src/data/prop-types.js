import PropTypes from 'prop-types'

// primitives

const LocalTimePropTypes = PropTypes.shape({
  year: PropTypes.number.isRequired,
  month: PropTypes.number.isRequired,
  day: PropTypes.number.isRequired,
  hour: PropTypes.number.isRequired,
  minute: PropTypes.number.isRequired,
  second: PropTypes.number.isRequired,
})

// incoming data

export const HourlyForecastPropTypes = PropTypes.shape({
  temp: PropTypes.number.isRequired,
  precip_probability: PropTypes.number.isRequired,
  precip_type: PropTypes.string.isRequired,
  phrase: PropTypes.string.isRequired,
})

export const DayPartPropTypes = PropTypes.shape({
  temp: PropTypes.number.isRequired,
  precip_probability: PropTypes.number.isRequired,
  precip_type: PropTypes.string.isRequired,
  phrase: PropTypes.string.isRequired,
  narrative: PropTypes.string.isRequired,
})

export const DayPartWithHourlyPropTypes = PropTypes.shape({
  temp: PropTypes.number.isRequired,
  precip_probability: PropTypes.number.isRequired,
  precip_type: PropTypes.string.isRequired,
  narrative: PropTypes.string.isRequired,
  hourly: PropTypes.arrayOf(HourlyForecastPropTypes).isRequired,
})

export const DailyForecastPropTypes = PropTypes.shape({
  // UTC timestamp
  start_timestamp: PropTypes.number.isRequired,
  // UTC timestamp
  expire_timestamp: PropTypes.number.isRequired,
  day: DayPartPropTypes,
  night: DayPartPropTypes,
})

export const LocationPropTypes = PropTypes.shape({
  city: PropTypes.string.isRequired,
  state_abbr: PropTypes.string.isRequired,
})

export const DailyWithHourlyForecastPropTypes = PropTypes.shape({
  // UTC timestamp
  start_timestamp: PropTypes.number.isRequired,
  // UTC timestamp
  expire_timestamp: PropTypes.number.isRequired,
  day: DayPartWithHourlyPropTypes,
  night: DayPartWithHourlyPropTypes,
})

// transformed

export const TransformedHourlyForecastPropTypes = PropTypes.shape({
  temp: PropTypes.number.isRequired,
  precip_probability: PropTypes.number.isRequired,
  precip_type: PropTypes.string.isRequired,
  phrase: PropTypes.string.isRequired,
  start_ltime: LocalTimePropTypes.isRequired,
  phrase: PropTypes.string.isRequired,
})
