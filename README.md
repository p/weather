# Weather App

This is a weather service with a web frontend and a React Native frontend
to come later. The goals of the application are:

1. Be speedy. The frontend shouldn't need to load megabytes of
images and javascript to show a weather forecast.
2. Have each view/page be as informative as possible.
In particular, the main view for a location should offer most of the
information that the user is ilkely to need. This means
current conditions and a concise forecast.
3. Offer good location selection and management. I travel a fair bit
and need weather in multiple cities.
3. Be ad free.

## Components

- `backend` - Go service to retrieve forecasts from upstream weather
services and cache them.
- `frontend` - Web frontend, built with React.
- `ruby` - Ruby backend, for prototyping.

## License

2-clause BSD
