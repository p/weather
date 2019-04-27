import reactor from './reactor'
import { merge } from './util'

export default {
  fetch_network() {
    if (NODE_ENV == 'production') {
      setTimeout(function() {
        reactor.dispatch('receive_network', { up: true })
      }, 0)
    } else {
      fetch(API_URL + '/network')
        .then(resp => resp.json())
        .then(payload => {
          reactor.dispatch('receive_network', payload)
        })
    }
  },

  fetch_weather(location_query, network_flag) {
    fetch(API_URL + '/locations/' + location_query + '?network=' + network_flag)
      .then(resp => resp.json())
      .then(payload => {
        reactor.dispatch('receive_weather', { location_query, payload })
      })
  },

  fetch_locations() {
    fetch(API_URL + '/locations')
      .then(resp => resp.json())
      .then(payload => {
        reactor.dispatch('receive_locations', payload)
      })
  },
}
