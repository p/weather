# Weather Backend

A server to cache weather data.

## Usage

The server recognizes the following environment variables at runtime:

- DEBUG: enable Gin debug mode
- DB_PATH: Path to database file for storing cached weather data
- PORT: port number to bind to (default is 8093)
- OWM_API_KEY: OpenWeatherMap API key
- MAPQUEST_API_KEY: MapQuest API key (for geocoding)

## License

2 clause BSD
