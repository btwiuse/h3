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

## Connect to public demo server at https://h3.k0s.io:32443

```
$ env HOST=h3.k0s.io PORT=32443 h3 client
2022/09/17 18:45:44 dialing https://h3.k0s.io:32443/echo (UDP)
2022/09/17 18:45:45 new conn [::]:36805
BTW I USE ARCH
BTW I USE ARCH
```

## TODO

- [ ] Figure how how to expose the HTTP/3 server through traefik ingress (needs help)
