BINARY_NAME=markdown.app
APP_NAME=markdown
VERSION=1.0.0

## build: build binary and package app
build:
	rm -rf ${BINARY_NAME}
	rm -f markdown
	fyne package -appVersion ${VERSION} -name ${APP_NAME} -release

## run: builds and runs app
run:
	go run .

## clean: runs go clean and deletes binaries
clean:
	@echo "Cleaning..."
	@go clean
	@rm -rf ${BINARY_NAME}
	@echo "Done!"

## tests: runs all tests
test:
	go test -v ./...