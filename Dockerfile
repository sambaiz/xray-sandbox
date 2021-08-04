FROM golang:1.16 AS builder
WORKDIR /go/src/github.com/sambaiz/xray-sandbox
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o app main.go

FROM alpine:latest  
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /go/src/github.com/sambaiz/xray-sandbox ./
CMD ["./app"]  