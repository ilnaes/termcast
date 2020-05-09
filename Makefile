GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
BINARY_NAME=termcast
CLI_PATH=cmd/cli/main.go

all: build
build: 
		$(GOBUILD) -o $(BINARY_NAME) $(CLI_PATH)
test: 
		$(GOTEST) -v ./...
clean: 
		$(GOCLEAN)
		rm -f $(BINARY_NAME)
run:
		make build
		./$(BINARY_NAME)
