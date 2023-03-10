### Bulder
FROM golang:alpine as builder
RUN apk update
RUN apk add git ca-certificates upx

WORKDIR /usr/src/app

COPY go.mod .
COPY go.sum .
COPY .env bin/main

RUN go mod tidy
# install dependencies

COPY . .

RUN GO111MODULE=on CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -ldflags="-s -w" -o bin/main ./main.go; \
    upx --best --lzma bin/main
# compile & pack

### Executable Image
FROM scratch

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /etc/passwd /etc/passwd
COPY --from=builder /usr/src/app/bin/main ./main
COPY  .env .

EXPOSE 3000

ENTRYPOINT ["./main"]