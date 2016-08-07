#!/bin/bash

CMD = famicom-disasm famicom-dump-image famicom-cpu famicom-ppu famicom

.PHONY: build fmt clean

build:
	@$(foreach cmd, $(CMD), echo go build -o build/$(cmd) cmd/$(cmd)/main.go ; go build -o build/$(cmd) cmd/$(cmd)/main.go ;)

fmt:
	go fmt github.com/mzp/famicom/...

test:
	go test github.com/mzp/famicom/...

clean:
	rm -rf build
