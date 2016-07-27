#!/bin/bash

CMD = famicom-disasm famicom-dump-image famicom-cpu

.PHONY: build fmt clean

build:
	$(foreach cmd, $(CMD), go build -o build/$(cmd) cmd/$(cmd)/main.go ;)

fmt:
	go fmt github.com/mzp/famicom/...

clean:
	rm -rf build
