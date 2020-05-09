GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GORUN=$(GOCMD) run
GOTEST=$(GOCMD) test
BINARY_NAME=termcast
CLI_PATH=cmd/cli/main.go
SERVER_PATH=cmd/server/main.go

all: cli
start:
		$(GORUN) $(SERVER_PATH)
cli: 
		$(GOBUILD) -o $(BINARY_NAME) $(CLI_PATH)
test: 
		$(GOTEST) -v ./...
clean: 
		$(GOCLEAN)
		rm -f $(BINARY_NAME)
run:
		make build
		./$(BINARY_NAME)
