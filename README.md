# h3

WebTransport over HTTP/3 echo server & client

## Install

```
$ go install github.com/btwiuse/h3@latest
```

## Start local server (needs [mkcert](https://mkcert.dev))

```
$ mkcert -install && mkcert localhost
The local CA is already installed in the system trust store! 👍
The local CA is already installed in the Firefox and/or Chrome/Chromium trust store! 👍

Created a new certificate valid for the following names 📜
 - "localhost"

The certificate is at "./localhost.pem" and the key at "./localhost-key.pem" ✅

It will expire on 17 December 2024 🗓

$ env HOST=localhost PORT=8443 CERT=./localhost.pem KEY=./localhost-key.pem h3 server
2022/09/17 16:16:10 listening on https://localhost:8443 (UDP)
2022/09/17 16:16:13 new conn 127.0.0.1:45394
```

## Connect to local server

```
$ env HOST=localhost PORT=8443 h3 client
2022/09/17 16:16:13 dialing https://localhost:8443/echo (UDP)
2022/09/17 16:16:13 new conn [::]:45394
BTW I USE ARCH
BTW I USE ARCH
```

## Connect to public demo server at [h3.k0s.io](https://h3.k0s.io)

```
$ env HOST=h3.k0s.io PORT=443 h3 client
2022/09/17 18:45:44 dialing https://h3.k0s.io:443/echo (UDP)
2022/09/17 18:45:45 new conn [::]:36805
BTW I USE ARCH
BTW I USE ARCH
```

## Test HTTP/3 against demo server with `curl`

```
$ curl3 -s https://h3.k0s.io -I --http3
HTTP/3 200
content-type: text/plain; charset=utf-8
x-content-type-options: nosniff

$ curl3 -s https://h3.k0s.io -I
HTTP/2 200
alt-svc: h3=":443"
content-type: text/plain; charset=utf-8
date: Sat, 17 Sep 2022 15:44:28 GMT
x-content-type-options: nosniff
content-length: 12
```

## TODO

- [ ] Figure how how to expose the HTTP/3 server through traefik ingress, rather
      that `hostPort` (help wanted)
