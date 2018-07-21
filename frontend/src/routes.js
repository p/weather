import Handlers from './handlers'
import { AppBase } from './app'

export default {
  Locations: { path: '/', component: Handlers.Locations, wrapper: AppBase },
  Location: {
    path: '/:location',
    component: Handlers.Location,
    wrapper: AppBase,
  },
}
