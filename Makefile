.PHONY: all build run clean update list test mocks

all: # Runs build and run
	go run main.go
	go build main.go

run: # Runs the application locally
	go run main.go

build: # Build the go application
	go build main.go

clean: # Add missing and remove unused go modules
	go mod tidy

update: # Install and update go modules
	go get -u ./...

list: # List modules that are being used
	go install github.com/icholy/gomajor@latest | gomajor list\

test: # Runs all the tests in the application and returns if they passed or failed, along with a coverage percentage
	go install github.com/mfridman/tparse@latest | go mod tidy
	go test -json -cover ./... | tparse -all -pass

mocks: # Install mock module and updates all mocks files
	go install github.com/vektra/mockery/v2@latest
	mockery --with-expecter --all --output mocks	

coverage:
	go test -coverprofile=coverage.out -covermode=count ./...
	go tool cover -func coverage.out   