FROM golang:1.21 AS build

WORKDIR /app
COPY . .

RUN go env -w GO111MODULE=on \
    && go env -w GOPROXY=https://goproxy.cn,https://goproxy.io,direct \
    && go mod tidy \
    && GOOS=linux CGO_ENABLED=0 go build -ldflags="-s -w" -o main

FROM alpine:latest

WORKDIR /app
RUN apk add openssl bash

COPY --from=build /app/main .
COPY templates templates
COPY static static
COPY ca.cnf gen.cert.sh gen.root.sh .
RUN chmod +x gen.cert.sh gen.root.sh

EXPOSE 8000

CMD ["./main"]