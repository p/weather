import moment from 'moment'

export function  data_age(value) {
    if (value) {
      let d = new Date().getTime() / 1000 - value.updated_at
      return moment.duration(d, 'seconds').humanize() + ' ago'
    } else {
      return null
    }
  }
