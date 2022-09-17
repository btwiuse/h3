# h3

WebTransport over HTTP/3 echo server & client

server

```
$ env PORT=8443 go run ./server/
2022/09/17 16:16:10 listening on https://localhost:8443 (UDP)
2022/09/17 16:16:13 new conn 127.0.0.1:45394
```

client

```
$ env PORT=8443 go run ./client/
2022/09/17 16:16:13 dialing https://localhost:8443/echo (UDP)
2022/09/17 16:16:13 new conn [::]:45394
btw
btw
i
i
use
use
arch
arch
```
