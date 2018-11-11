import Handlers from './handlers'
import { AppBase } from './app'

export default {
  Locations: { path: '/', component: Handlers.Locations, wrapper: AppBase },
  BriefLocation: {
    path: '/:location',
    component: Handlers.BriefLocation,
    wrapper: AppBase,
    options: {segmentValueCharset: 'a-zA-Z0-9, %-'},
  },
  FullLocation: {
    path: '/:location/full',
    component: Handlers.FullLocation,
    wrapper: AppBase,
  },
}
