import PropTypes from 'prop-types'
import {Link} from 'react-easy-router'
import Immutable from 'seamless-immutable'
import preventDefaultWrapper from '@rq/prevent-default-wrapper'
import _ from 'underscore'
import connectToStores from 'alt-utils/lib/connectToStores';
import React from 'react';
import Store from '../store';

export default class Location extends React.Component {
  constructor(props) {
    super(props)
  }
  
  static getStores() {
    return [Store]
  }
    
  static getPropsFromStores() {
    return Store.getState()
  }

  render() {
    return <div>
      Weather for
    {this.props.params.location}
    </div>
  }
}

Location.propTypes = {
  params: PropTypes.shape({
    location: PropTypes.string,
  }),
}
