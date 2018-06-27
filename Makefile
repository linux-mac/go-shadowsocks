prefix:=ss
client_darwin:="$(GOPATH)/bin/$(prefix)-client-darwin-amd64"
client_linux:="$(GOPATH)/bin/$(prefix)-client-linux-amd64"
server_darwin:="$(GOPATH)/bin/$(prefix)-server-darwin-amd64"
server_linux:="$(GOPATH)/bin/$(prefix)-server-linux-amd64"

.PHONY: clean compress

all: client server

client: common/*.go $(prefix)-client/*.go
	cd $(prefix)-client; \
	go build -o $(client_darwin) client.go; \
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o $(client_linux) client.go

server: common/*.go $(prefix)-client/*.go
	cd $(prefix)-server; \
	go build -o $(server_darwin) server.go; \
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o $(server_linux) server.go

compress: all
	cd $(GOPATH)/bin; \
	tar czf $(prefix)-client-darwin-amd64.tar.gz $(prefix)-client-darwin-amd64; \
	tar czf $(prefix)-client-linux-amd64.tar.gz $(prefix)-client-linux-amd64; \
	tar czf $(prefix)-server-darwin-amd64.tar.gz $(prefix)-server-darwin-amd64; \
	tar czf $(prefix)-server-linux-amd64.tar.gz $(prefix)-server-linux-amd64;

clean:
	rm -f $(client_linux) $(client_darwin) $(server_linux) $(server_darwin)
