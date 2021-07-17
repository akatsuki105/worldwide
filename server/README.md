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

**debug/tileview**

Get tile data at 1-second intervals using Websocket.

Tile data is sent in binary format. Please refer to [tileview.html](./tileview.html) for how to display it.

```sh
wscat -c ws://localhost:8888/debug/tileview
```

**debug/sprview**

Get sprite data at 1-second intervals using Websocket.

Sprite data is sent in binary format. Please refer to [sprview.html](./sprview.html) for how to display it.

```sh
wscat -c ws://localhost:8888/debug/sprview
```

