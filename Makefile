BIN_FILE=vfs
all: clean build run
build:
	@go build -o "${BIN_FILE}" ./cmd/main.go
clean:
	@go clean
	@rm -f ./"${BIN_FILE}"
run:
	./"${BIN_FILE}"
help:
	@echo "make: rebuild the application, clean up any previous builds, and run the virtual file system."
	@echo "make build: run the virtual file system."
	@echo "make clean: delete the target file."
	@echo "make run: run the virtual file system."

.PHONY: all build clean run help