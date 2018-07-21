import alt from './alt'

class Actions {
  updateTodo(id, text) {
    return { id, text }
  }
}

export default alt.createActions(Actions)
