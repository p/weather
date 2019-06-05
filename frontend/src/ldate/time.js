import moment from 'moment'

const iso_8601_regexp = /^(\d\d\d\d)-(\d\d)-(\d\d)T(\d\d):(\d\d):(\d\d)/

export default class LocalTime {
  constructor(time_str) {
    if (typeof time_str != 'string'){
      throw new Error(`Argument is not a string: ${time_str}`)
    }
    let m = time_str.match(iso_8601_regexp)
    if (m) {
      this.year = parseInt(m[1])
      this.month = parseInt(m[2])
      this.day = parseInt(m[3])
      this.hour = parseInt(m[4])
      this.minute = parseInt(m[5])
      this.second = parseInt(m[6])
    } else {
      throw new Error('Invalid time format: ' + time_str)
    }
  }
  
  format(format_str){
    return moment(new Date(this.year, this.month-1, this.day, this.hour,this.minute,this.second)).format(format_str)
  }
}
