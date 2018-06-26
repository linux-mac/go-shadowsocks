PREFIX := ss
CLIENT := $(GOPATH)/bin/$(PREFIX)-client
SERVER := $(GOPATH)/bin/$(PREFIX)-server

.PHONY: clean compress

all: $(CLIENT) $(SERVER)

$(CLIENT): common/*.go $(PREFIX)-client/*.go
	cd $(PREFIX)-client;go install

$(SERVER): common/*.go $(PREFIX)-server/*.go
	cd $(PREFIX)-server;go install

client: $(CLIENT)

server: $(SERVER)

compress: $(CLIENT) $(SERVER)
	cd $(GOPATH)/bin; \
	tar cvzf $(PREFIX)-client.tar.gz $(PREFIX)-client; \
	tar cvzf $(PREFIX)-server.tar.gz $(PREFIX)-server

clean:
	rm -f $(CLIENT) $(SERVER)
