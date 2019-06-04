import { assert } from 'chai'
import {LocalTime} from '../../src/ldate'

describe('LocalTime', ()=>{
      let time =new LocalTime("2019-06-03T14:20:33-04:00")
  describe('year',()=>{
    it('is the year',()=>{
      assert.equal(time.year,2019)
    })
  })
  describe('month',()=>{
    it('is the month',()=>{
      assert.equal(time.month,6)
    })
  })
  describe('day',()=>{
    it('is the day',()=>{
      assert.equal(time.day,3)
    })
  })
  describe('hour',()=>{
    it('is the hour',()=>{
      assert.equal(time.hour,14)
    })
  })
  describe('minute',()=>{
    it('is the minute',()=>{
      assert.equal(time.minute,20)
    })
  })
  describe('second',()=>{
    it('is the second',()=>{
      assert.equal(time.second,33)
    })
  })
})
