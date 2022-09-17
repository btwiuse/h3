FROM btwiuse/k0s as k0s

FROM btwiuse/arch:golang as h3

WORKDIR /h3

ADD . /h3

RUN go mod tidy

RUN GOBIN=/usr/local/bin go install .

FROM btwiuse/arch

COPY --from=k0s /usr/bin/k0s /usr/local/bin/

COPY --from=h3 /usr/local/bin/h3 /usr/local/bin/

RUN pacman -Sy --noconfirm --needed --overwrite='*' mkcert

RUN mkcert -install

WORKDIR /h3

RUN mkcert localhost

CMD k0s agent
