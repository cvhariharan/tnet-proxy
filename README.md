# tnet-proxy
A proxy server that exposes targets as machines on Tailscale. Useful for attaching custom ACL policies.

## Installation
### Go
```bash
go install github.com/cvhariharan/tnet-proxy/cmd/tnet@latest
```
### Docker
```
docker run -e TS_AUTHKEY=ts-authkey -p 8000:8000 ghcr.io/cvhariharan/tnet-proxy -port 8000 -target docker-container-name:9000 -hostname example
```

## Usage
`TS_AUTHKEY` environment variable should be set with an auth key obtained from Tailscale admin console.
```bash
> go build -o tnet-proxy cmd/tnet/main.go
> ./tnet-proxy --help
Usage of ./tnet-proxy:
  -hostname string
        hostname used on the tailnet (default "tnet-proxy")
  -port int
        port that the proxy listens on (default 8000)
  -proto string
        protocol tcp/udp (default "tcp")
  -target string
        target to proxy (eg: localhost:8080, :8080)
```
`-hostname` sets the Tailscale hostname. Using MagicDNS, the service will be accessible using this hostname.

The proxy listens for incoming connections on `-port` and routes it to the `-target`.
