import { Router } from 'react-easy-router'
import routes from './routes'
import React from 'react'
import ReactDOM from 'react-dom'
import { hot } from 'react-hot-loader'
import history from './history'

class Root extends React.Component {
  render() {
    return <Router history={history} routes={routes} />
  }
}

export default hot(module)(Root)
