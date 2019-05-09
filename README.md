# ws

ws is a simple command line websocket client designed for exploring and debugging websocket servers. ws includes readline-style keyboard shortcuts, persistent history, and colorization.

![Example usage recording](https://hashrocket-production.s3.amazonaws.com/uploads/blog/misc/ws/ws.gif)

## Installation

```
go get -u github.com/hashrocket/ws
```

## Usage

Simply run ws with the destination URL. For security some sites check the origin header. ws will automatically send the destination URL as the origin. If this doesn't work you can specify it directly with the `--origin` parameter.

```
$ ws ws://localhost:3000/ws
> {"type": "echo", "payload": "Hello, world"}
< {"type":"echo","payload":"Hello, world"}
> {"type": "broadcast", "payload": "Hello, world"}
< {"type":"broadcast","payload":"Hello, world"}
< {"type":"broadcastResult","payload":"Hello, world","listenerCount":1}
> ^D
EOF
```

use command **@audio:** to upload local audio file to the server as a websocket binary message

```
$ ws ws://localhost:8080/endpoint
> {"audio_format":{"type":"MP3"}, "begin_silence_in_milli":2000, "end_silence_in_milli":2000}
> @audio:/Users/lls/Downloads/bug-audio.mp3
> EOS
< {"type":"success","request_id":"bj9pmf5tgsakv6qkuee0","data":{"endpoint_detected":false,"endpoint_position_in_milli":0}}
```