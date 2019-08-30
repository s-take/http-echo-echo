FROM golang:1.11-alpine AS build

RUN apk --no-cache add gcc g++ make ca-certificates git
WORKDIR /go/src/github.com/s-take/http-echo-echo

COPY main.go main.go

RUN go install ./...

# 上記でできたバイナリだけをalpineに入れる
FROM alpine:3.10
WORKDIR /usr/bin
COPY --from=build /go/bin .
CMD ["http-echo-echo"]
