import reactor from './reactor'

export default {
  fetch_network(){
    fetch(API_URL + '/network')
      .then(resp => resp.json())
      .then(payload => {
        reactor.dispatch('receive_network', payload)
      })
  }
}
