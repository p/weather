GOPATH?=$$HOME/go

all: b

d:
	go get github.com/codegangsta/gin
	go-get-deps

b:
	mkdir -p tmp
	go build -o tmp/weather src/server.go

r:
	env `cat .env` DEBUG=1 ${GOPATH}/bin/gin --build src --port 8093 --bin tmp/gin-bin run server.go wu.go weather.go location.go

fmt:
	for f in src/*.go; do go fmt $$f && sed -i -e 's/	/  /g' $$f; done
