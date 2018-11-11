import reactor from './reactor'
import { Store, toImmutable } from 'nuclear-js'

// network

let NetworkStore = Store({
  getInitialState() {
    return toImmutable({})
  },

  initialize() {
    this.on('receive_network', receive_network)
  }
})

function receive_network(state, network){
  return state.merge(network)
}

// forecast

let ForecastStore = Store({
  getInitialState() {
    return toImmutable({})
  },

  initialize() {
    this.on('receive_forecast', receive_forecast)
  }
})

function receive_forecast(state, {location_query, forecast}){
  let delta = {}
  delta[location_query]=forecast
  return state.merge(delta)
}

reactor.registerStores({
  network: NetworkStore,
  forecast: ForecastStore,
})
