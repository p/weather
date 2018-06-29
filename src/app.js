import React from 'react';
import Location from './location'

export default class App extends React.Component {
  constructor(props) {
    super(props)

    this.state = {
      location: null,
    }
  }

  render() {
    return <div>
      Weather for:
      <Location
        location_did_change={this.location_did_change.bind(this)}
      />
      <p>{this.state.location}</p>
    </div>
  }

  location_did_change(location) {
    this.setState({location: location})
  }
}
