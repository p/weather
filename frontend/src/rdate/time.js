export default class Time {
  constructor(timestamp, tz) {
    this._timestamp = timestamp
    this.tz = tz
  }
  
  utc() {
    return new Time(this._timestamp, 'UTC')
  }
  
  local() {
    return new Time(this._timestamp, 'xx')
  }
  
  timestamp(){
    return this._timestamp
  }
  
  to_js_date() {
    return new Date(this._timestamp*1000)
  }
  
  static utc(...args){
    if(args.length>1){
      --args[1]
    }
    const date=new Date(...args)
    let ts=date.getTime()/1000
    let tzoffset = date.getTimezoneOffset()*60
    return new Time(ts-tzoffset)
  }
  
  static local(...args) {
    const date = new Date(...args)
    return new Time(date.getTime() / 1000)
  }
}

Time.toJSDate = Time.to_js_date
