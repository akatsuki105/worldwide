# Server

`worldwide` contains an HTTP server, and the user can give various instructions to it through HTTP requests.

## Start

Normal execution will not start the HTTP server. You can start the HTTP server by running it with the port option specified.

```sh
go run ./cmd -p 8888 ./PM_PRISM.gbc # start HTTP server on localhost:8888
```

Note: On this API server, GET method changes emulator's state!

## Commands

**pause**

Pause emulator

```sh
curl localhost:8888/pause
```

**continue**

Continue emulator from pause state

```sh
curl localhost:8888/continue
```

**reset**

Reset emulator. Debug information(e.g. breakpoints, history-mode) is not reset.

```sh
curl localhost:8888/reset
```

**quit**

Quit emulator. Exitcode is 0 and savefile is written.

```sh
curl localhost:8888/quit
```

**mute**

Mute emulator(toggle)

```sh
curl localhost:8888/mute
```

## Debug commands

**debug/register(GET)**

Get value from CPU-registers and other important registers

```sh
curl localhost:8888/debug/register
```

```jsonc
// application/json
{
    "A":"0x00", "F":"0xa0",
    "B":"0x00", "C":"0x00",
    "D":"0x00", "E":"0x04",
    "H":"0xff", "L":"0xfe",
    "PC":"0x4a10","SP":"0xc0f8",
    "IE":"0x0f", "IF":"0xe1", "IME":"0x01",
    "Halt":"true", "DoubleSpeed":"true"
}
```

**debug/register(POST)**

Set a value into register

```sh
# target: register's name (a, b, c, d, e, h, l, f, sp, pc, ime)
# value: hex value (e.g. 0x0486)
curl -X POST -d '{"target":"ime", "value":"0x1"}' -H "Content-Type: application/json" localhost:8888/debug/register
```

**debug/break(POST)**

Set a breakpoint

```sh
curl -X POST -d '{"addr":"0x0486"}' -H "Content-Type: application/json" localhost:8888/debug/break
```

**debug/break(GET)**

List breakpoints

```sh
curl localhost:8888/debug/break
```

```sh
[0x0486, 0x0490] # text/plain
```

**debug/break(DELETE)**

Delete a breakpoint

```sh
curl -X DELETE "localhost:8888/debug/break?addr=0x0486"
```

**debug/read1**

Read a byte from memory

```sh
curl "localhost:8888/debug/read1?addr=0x0150"
```

```sh
0x12 # text/plain
```
**debug/read2**

Read two bytes from memory

```sh
curl "localhost:8888/debug/read2?addr=0x0150"
```

```sh
0x1411 # text/plain
```

**debug/history(POST)**

Start recording the history of the executed instructions.

Specify how many past histories to record in the `history` parameter.(Max 100) 

The larger the number, the greater the load on the emulator CPU.

```sh
curl -X POST -d '{"history":"0x20"}' -H "Content-Type: application/json" localhost:8888/debug/history
```

**debug/history(GET)**

Displays the history of records started by POST requests.

```sh
curl localhost:8888/debug/history
```

```sh
# text/plain
0x0048: JP a16
0xcda1: PUSH AF
0xcda2: PUSH HL
...
0x048a: JR NZ r8
0x0485: HALT
0x0486: LD A (a16)
```

**debug/cartridge**

Get cartridge info

```sh
curl localhost:8888/debug/cartridge
```

```jsonc
// application/json
{
    "title":"PM_PRISM",
    "cartridge_type":"MBC3+TIMER+RAM+BATTERY",
    "rom_size":"2MB",
    "ram_size":"32KB"
}
```

**debug/disasm**

Disassemble instructions

```sh
curl localhost:8888/debug/disasm
```

```jsonc
// application/json
{
    "pc":"0x0486",
    "mnemonic":"LD A (a16)"
}
```

**debug/trace**

Trace a number of instructions

```sh
curl "localhost:8888/debug/trace?step=20" # trace 20 instructions
```

```sh
# text/plain
0x16cf: INC L
0x16d0: LD (HL) D
0x16d1: INC L
0x16d2: POP DE
0x16d3: LD (HL) E
...
0x16dd: LD (HL) E
0x16de: INC L
0x16df: LD (HL) D
0x16e0: INC L
0x16e1: POP DE
0x16e2: LD (HL) E
```

**debug/io(Websocket)**

Get IO registers(`0xff00-0xffff`) at 100-milisecond intervals using Websocket.

IO registers is sent in arraybuffer. Please refer to [io.html](./io.html) for how to display it.

```sh
wscat -c ws://localhost:8888/debug/io
```

**debug/tileview/bank0(Websocket)**

Get tile data at 100-milisecond intervals using Websocket.

Tile data is sent in binary format. Please refer to [tileview.html](./tileview.html) for how to display it.

```sh
wscat -c ws://localhost:8888/debug/tileview/bank0 # or ws://localhost:8888/debug/tileview/bank1
```

**debug/sprview(Websocket)**

Get sprite data at 100-milisecond intervals using Websocket.

Sprite data is sent in binary format. Please refer to [sprview.html](./sprview.html) for how to display it.

```sh
wscat -c ws://localhost:8888/debug/sprview
```
