all: b

b:
	mkdir -p tmp
	go build -o tmp/weather src/server.go

r:
	env `cat .env` DEBUG=1 tmp/weather
