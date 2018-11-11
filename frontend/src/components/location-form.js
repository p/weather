import preventDefaultWrapper from '@rq/prevent-default-wrapper'
import React from 'react'
import PropTypes from 'prop-types'

export default class LocationForm extends React.Component {
  constructor(props) {
    super(props)

    this.state = {
      location: null,
    }
  }

  render() {
    return (
      <form>
        <input
          type="text"
          defaultValue={this.state.location || ''}
          onChange={this.location_did_change.bind(this)}
        />
        <input
          type="submit"
          onClick={preventDefaultWrapper(this.location_did_submit.bind(this))}
        />
      </form>
    )
  }

  location_did_change(e) {
    this.setState({ location: e.target.value })
  }

  location_did_submit(e) {
    if (this.props.location_did_submit) {
      this.props.location_did_submit(this.state.location)
    }
  }
}

LocationForm.propTypes = {
  location_did_submit: PropTypes.func,
}
