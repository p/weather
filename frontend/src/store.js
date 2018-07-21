import alt from './alt'
import Actions from './actions'

class Store {
  constructor() {
    this.bindListeners({
      updateTodo: Actions.updateTodo,
    })

    let locations = localStorage.getItem('locations') || '[]'
    try {
      locations = JSON.parse(locations)
    } catch (e) {
      locations = []
    }

    this.state = {
      locations: locations,
    }
  }

  updateTodo(todo) {
    this.setState({ todos: this.state.todos.concat(todo) })
  }
}

export default alt.createStore(Store, 'Store')
