# h3

WebTransport over HTTP/3 echo server & client

install

```
$ go install github.com/btwiuse/h3@latest
```

start local server

```
$ env HOST=localhost PORT=8443 CERT=./localhost.pem KEY=./localhost-key.pem h3 server
2022/09/17 16:16:10 listening on https://localhost:8443 (UDP)
2022/09/17 16:16:13 new conn 127.0.0.1:45394
```

connect to local server

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

connect to public demo server at https://h3.k0s.io:32443

```
$ env HOST=h3.k0s.io PORT=32443 h3 client
2022/09/17 18:45:44 dialing https://h3.k0s.io:32443/echo (UDP)
2022/09/17 18:45:45 new conn [::]:36805
btw
btw
i
i
use
use
arch
arch
```
