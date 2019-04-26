import moment from 'moment'

export function data_age(struct) {
  if (struct) {
    let d = new Date().getTime() / 1000 - struct.updated_at
    return moment.duration(d, 'seconds').humanize() + ' ago'
  } else {
    return null
  }
}

export function unim(any) {
  if (any && any.toJS) {
    return any.toJS()
  } else {
    return any
  }
}

export function network_flag(up) {
  return up ? 0 : 2
}

export function merge(dest, src){
  return Object.assign({},dest,src)
}

export function make_hash(){
  let hash={}
  for(let i =0;i<arguments.length;i+=2){
    hash[arguments[i]] = arguments[i+1]
  }
  return hash
  
}
