import { assert } from 'chai'
import rdate from '../../src/rdate'

describe('Time', () => {
  describe('local', () => {
    it('is local', () => {
      let actual = rdate.Time.local(2010, 8, 4, 3, 7, 9)
      let expected = new Date(2010, 8, 4, 3, 7, 9)
      assert.deepEqual(actual.to_js_date(), expected)
      console.log(actual)
    })
  })

  describe('utc', () => {
    it('is utc', () => {
      let actual = rdate.Time.utc(2010, 8, 4, 3, 7, 9)
      let expected = 1280891229
      assert.equal(actual.timestamp(), expected)
    })
  })
})
