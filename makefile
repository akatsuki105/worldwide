TARGET = worldwide

ifeq ($(OS),Windows_NT)
    TARGET = worldwide.exe
endif

.PHONY: all
all:
	go build -o $(TARGET) -ldflags "-X main.version=$(shell git describe --tags)" ./cmd/

.PHONY: clean
clean:
	rm -f $(TARGET)