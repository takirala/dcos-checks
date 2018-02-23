BUILD_DIR=build
EXECUTABLE_NAME=dcos-checks

.PHONY: all
all: clean build test

.PHONY: build
build:
	go build -o $(BUILD_DIR)/$(EXECUTABLE_NAME)

.PHONY: test
test:
	bash -c "./script/test.sh"

.PHONY: clean
clean:
	rm -rf $(BUILD_DIR)
