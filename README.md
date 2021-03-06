# famicom toolkit

![](https://raw.githubusercontent.com/mzp/famicom/master/docs/images/logo.png)

Famicom(NES) toolkit written by golang.

## Install

```
export GOPATH=/path/to/gopath
go get github.com/mzp/famicom/cmd/...
```

## Current tools
### famicom-dump-image

![](https://raw.githubusercontent.com/mzp/famicom/master/docs/images/famicom-dump-image.png)

Dump NES pallet data to PNG image.

### famicom-disasm

![](https://raw.githubusercontent.com/mzp/famicom/master/docs/images/famicom-disasm.png)

Disasm 6502 machine code.

### famicom-cpu

![](https://raw.githubusercontent.com/mzp/famicom/master/docs/images/famicom-cpu.png)

6502 emulator. It can execute all instruction, and show register, memory.

### famicom-ppu

![](https://raw.githubusercontent.com/mzp/famicom/master/docs/images/famicom-ppu.png)

PPU emulator. It load `*.nes` file and show argumented text.

## License

MIT License
