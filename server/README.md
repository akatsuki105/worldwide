# Server

`worldwide` contains an HTTP server, and the user can give various instructions to it through HTTP requests.

## Start

Normal execution will not start the HTTP server. You can start the HTTP server by running it with the port option specified.

```sh
go run ./cmd -p 8888 ./PM_PRISM.gbc # start HTTP server on localhost:8888
```

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

**mute**

Mute emulator(toggle)

```sh
curl localhost:8888/mute
```

## Debug commands

**debug/status**

Get general status

```sh
curl localhost:8888/debug/status
```

```json
{
    "A":"0x00", "F":"0xa0",
    "B":"0x00", "C":"0x00",
    "D":"0x00", "E":"0x04",
    "H":"0xff", "L":"0xfe",
    "PC":"0x4a10","SP":"0xc0f8",
    "IE":"0x000f","IF":"0x00e1","IME":"0x01",
    "LCDC":"0xcf","STAT":"0x40","LY":"0x90","LYC":"0xc7",
    "DoubleSpeed":"true"
}
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

**debug/cartridge**

Get cartridge info

```sh
curl localhost:8888/debug/cartridge
```

```json
{
    "title":"PM_PRISM",
    "cartridge_type":"MBC3+TIMER+RAM+BATTERY",
    "rom_size":"2MB",
    "ram_size":"32KB"
}
```

**debug/io**

Get IO registers(`0xff00-0xffff`) at 1-second intervals using Websocket.

IO registers is sent in arraybuffer. Please refer to [io.html](./io.html) for how to display it.

```sh
wscat -c ws://localhost:8888/debug/io
```

**debug/tileview/bank0**

Get tile data at 1-second intervals using Websocket.

Tile data is sent in binary format. Please refer to [tileview.html](./tileview.html) for how to display it.

```sh
wscat -c ws://localhost:8888/debug/tileview/bank0 # or ws://localhost:8888/debug/tileview/bank1
```

**debug/sprview**

Get sprite data at 1-second intervals using Websocket.

Sprite data is sent in binary format. Please refer to [sprview.html](./sprview.html) for how to display it.

```sh
wscat -c ws://localhost:8888/debug/sprview
```

