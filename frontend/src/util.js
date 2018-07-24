import moment from 'moment'

export function  data_age(struct) {
    if (struct) {
      let d = new Date().getTime() / 1000 - struct.updated_at
      return moment.duration(d, 'seconds').humanize() + ' ago'
    } else {
      return null
    }
  }
