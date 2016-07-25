#!/bin/bash

CMD = famicom-disasm famicom-dump-image

.PHONY: build clean

build:
	$(foreach cmd, $(CMD), go build -o build/$(cmd) cmd/$(cmd)/main.go ;)

clean:
	rm -rf build
