PREFIX := ss
CLIENT := $(GOPATH)/bin/$(PREFIX)-client
SERVER := $(GOPATH)/bin/$(PREFIX)-server

.PHONY: clean

all: $(CLIENT) $(SERVER)

$(CLIENT): common/*.go $(PREFIX)-client/*.go
	cd $(PREFIX)-client;go install

$(SERVER): common/*.go $(PREFIX)-server/*.go
	cd $(PREFIX)-server;go install

client: $(CLIENT)

server: $(SERVER)

clean:
	rm -f $(CLIENT) $(SERVER)
