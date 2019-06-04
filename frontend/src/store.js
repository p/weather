import _ from 'underscore'
import {LocalTime} from './ldate'
import * as u from './util'
import reactor from './reactor'
import { Store, toImmutable } from 'nuclear-js'

// network

let NetworkStore = Store({
  getInitialState() {
    return toImmutable({})
  },

  initialize() {
    this.on('receive_network', receive_network)
  },
})

function receive_network(state, network) {
  return state.merge(network)
}

// locations

let LocationsStore = Store({
  getInitialState() {
    return toImmutable([])
  },

  initialize() {
    this.on('receive_locations', receive_locations)
  },
})

function receive_locations(state, locations) {
  return locations
}

// forecast

let WeatherStore = Store({
  getInitialState() {
    return toImmutable({})
  },

  initialize() {
    this.on('receive_weather', receive_weather)
  },
})

function receive_weather(state, { location_query, payload }) {
  payload = u.merge(payload, {
    hourly_forecasts: _.map(payload.hourly_forecasts,hfc=>u.merge(hfc,
    {start_ltime:new LocalTime(hfc.start_at)}))})
  let new_payload = u.merge(state[location_query] || {}, payload)
  let delta = u.make_hash(location_query, new_payload)
  return state.merge(delta)
}

reactor.registerStores({
  network: NetworkStore,
  weather: WeatherStore,
  locations: LocationsStore,
})
