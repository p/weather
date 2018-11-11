import reactor from './reactor'
import {
  Provider,
} from 'nuclear-js-react-addons-chefsplate'
import { Router } from '@rq/react-easy-router'
import routes from './routes'
import React from 'react'
import ReactDOM from 'react-dom'
import { hot } from 'react-hot-loader'
import history from './history'

import './store'

class Root extends React.Component {
  render() {
    return <Provider reactor={reactor}>
    <Router history={history} routes={routes} />
    </Provider>
  }
}

export default hot(module)(Root)
