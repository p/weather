import reactor from './reactor'
import { Store, toImmutable } from 'nuclear-js'

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

reactor.registerStores({
  network: NetworkStore,
})

export default NetworkStore
