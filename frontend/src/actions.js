import reactor from './reactor'

export default {
  fetch_network(){
    fetch(API_URL + '/network')
      .then(resp => resp.json())
      .then(payload => {
        reactor.dispatch('receive_network', payload)
      })
  },
  
  fetch_forecast(location_query,network_flag){
    fetch(
      API_URL +
        '/locations/' +
        location_query +
        '?network=' +
        network_flag,
    )
      .then(resp => resp.json())
      .then(payload => {
        reactor.dispatch('receive_forecast', {location_query, forecast:payload})
      })
  }
}
