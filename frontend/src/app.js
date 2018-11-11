import { Link } from '@rq/react-easy-router'
import Immutable from 'seamless-immutable'
import preventDefaultWrapper from '@rq/prevent-default-wrapper'
import _ from 'underscore'
import React from 'react'
import Store from './store'

export default class AppBase extends React.Component {
  constructor(props) {
    super(props)
  }

  render() {
    return <div>{this.props.children}</div>
  }
}
