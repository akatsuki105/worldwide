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

.PHONY: test

TEST0 = ./test/gb-test-roms/cpu_instrs/
TEST1 = ./test/gb-test-roms/instr_timing/

.SILENT:
test:
	go run ./cmd/ --test="$(TEST0)actual.jpg" $(TEST0)cpu_instrs.gb
	-diff "$(TEST0)actual.jpg" "$(TEST0)expected.jpg" && echo "gb-test-roms/cpu_instrs OK"

	go run ./cmd/ --test="$(TEST1)actual.jpg" $(TEST1)instr_timing.gb
	-diff "$(TEST1)actual.jpg" "$(TEST1)expected.jpg" && echo "gb-test-roms/instr_timing OK"

	rm -f $(TEST0)actual.jpg $(TEST1)actual.jpg