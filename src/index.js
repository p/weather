import {Router} from 'react-easy-router'
import {createHashHistory, useBasename} from 'history'
import routes from './routes'
import React from 'react';
import ReactDOM from 'react-dom';

const history = createHashHistory({basename: '/'})

ReactDOM.render(
  <Router
    history={history}
    routes={routes}/>
, document.getElementById('root'));
