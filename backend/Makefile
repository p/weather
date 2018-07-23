GOPATH?=$$HOME/go

bare_sources=server.go weather.go location.go persist.go schema.go types.go owm.go wu.go util.go
sources=src/server.go src/weather.go src/location.go src/persist.go src/schema.go src/types.go src/owm.go src/wu.go src/util.go

all: b

d:
	go get github.com/codegangsta/gin
	go-get-deps

b:
	mkdir -p tmp
	go build -o tmp/weather $(sources)

r:
	env `cat .env` DEBUG=1 ${GOPATH}/bin/gin --build src --port 8093 --bin tmp/gin-bin run $(bare_sources)

rr:
	env `cat .env` DEBUG=1 go run $(sources)

fmt:
	for f in src/*.go; do go fmt $$f && sed -i -e 's/	/  /g' $$f; done
