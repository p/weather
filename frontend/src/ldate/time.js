const iso_8601_regexp = /^(\d\d\d\d)-(\d\d)-(\d\d)T(\d\d):(\d\d):(\d\d)/

export default class LocalTime {
  constructor(time_str) {
    if (typeof time_str != 'string'){
      throw new Error(`Argument is not a string: ${time_str}`)
    }
    let m = time_str.match(iso_8601_regexp)
    if (m) {
      this.year = m[1]
      this.month = m[2]
      this.day = m[3]
      this.hour = m[4]
      this.minute = m[5]
      this.second = m[6]
    } else {
      throw new Error('Invalid time format: ' + time_str)
    }
  }
}
