import { createHashHistory, useBasename } from 'history'

const history = createHashHistory({ basename: '/' })

export default history
