# h3

WebTransport over HTTP/3 echo server & client

install

```
$ go install github.com/btwiuse/h3@latest
```

start server

```
$ env HOST=localhost PORT=8443 CERT=./localhost.pem KEY=./localhost-key.pem h3 server
2022/09/17 16:16:10 listening on https://localhost:8443 (UDP)
2022/09/17 16:16:13 new conn 127.0.0.1:45394
```

start client

```
$ env HOST=localhost PORT=8443 h3 client
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
