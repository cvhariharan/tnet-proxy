FROM golang:1.22-alpine3.19 as Builder
WORKDIR /app
COPY . .
RUN go build -o tnet-proxy cmd/tnet/main.go

FROM alpine:3.19.1
COPY --from=Builder /app/tnet-proxy /usr/bin/tnet-proxy
WORKDIR /app
ENTRYPOINT ["/usr/bin/tnet-proxy"]