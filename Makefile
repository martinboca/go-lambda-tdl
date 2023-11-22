GOCMD=go
GOBUILD=$(GOCMD) build
BUILDDIR=build

CLIENT_BINARY_NAME=client
SERVER_BINARY_NAME=server

all: build_client build_server

build_client:
	$(GOBUILD) -o $(BUILDDIR)/$(CLIENT_BINARY_NAME) ./client

build_server:
	$(GOBUILD) -o $(BUILDDIR)/$(SERVER_BINARY_NAME) ./server

clean:
	rm -f $(BUILDDIR)/$(CLIENT_BINARY_NAME)
	rm -f $(BUILDDIR)/$(SERVER_BINARY_NAME)

$(shell mkdir -p $(BUILDDIR))
