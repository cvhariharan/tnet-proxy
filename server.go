package proxy

import (
	"fmt"
	"io"
	"log"
	"net"
	"time"

	"golang.org/x/sync/errgroup"
	"tailscale.com/tsnet"
)

const IdleTimeout = 3600

type ProxyServer interface {
	ListenAndServe() error
}

type ForwardProxy struct {
	target     string
	listenPort string
	hostname   string
	listener   net.Listener
	proto      string
}

// NewForwardProxy creates a proxy that copies packets
// between src and dest
func NewForwardProxy(listenPort int, target string, hostname, proto string) ProxyServer {
	return &ForwardProxy{
		target:     target,
		listenPort: fmt.Sprintf(":%d", listenPort),
		hostname:   hostname,
		proto:      proto,
	}
}

// ListenAndServe creates a tsnet listener and starts the proxy
func (f *ForwardProxy) ListenAndServe() error {
	// Create ts listener
	s := new(tsnet.Server)
	s.Hostname = f.hostname
	defer s.Close()

	ln, err := s.Listen(f.proto, f.listenPort)
	if err != nil {
		return fmt.Errorf("could not listen on port %s: %v", f.listenPort, err)
	}
	f.listener = ln

	// Set an idle timeout (default 1 hour)
	f.listen(IdleTimeout)
	return nil
}

// listen checks for new connections on the listener and spins up additional
// goroutines to copy data between src and dest
func (f *ForwardProxy) listen(idleTimeout int) {
	for {
		srcConn, err := f.listener.Accept()
		if err != nil {
			// Just log the error and continue accepting new connections
			log.Println(err)
			continue
		}
		src := &Conn{srcConn, idleTimeout}

		go func(src net.Conn) {
			err := f.handle(src)
			if err != nil {
				log.Println(err)
			}
		}(src)
	}
}

func (f *ForwardProxy) handle(src net.Conn) error {
	defer src.Close()

	destConn, err := net.DialTimeout(f.proto, f.target, 5*time.Second)
	if err != nil {
		return fmt.Errorf("could not open connection to %s: %v", f.target, err)
	}
	dest := &Conn{destConn, IdleTimeout}
	defer dest.Close()

	var eg errgroup.Group
	eg.Go(func() error {
		return forward(src, dest)
	})

	eg.Go(func() error {
		return forward(dest, src)
	})

	err = eg.Wait()
	if err != nil {
		return fmt.Errorf("error while forwarding packets between src and dest: %v", err)
	}
	return nil
}

func forward(src, dest net.Conn) error {
	_, err := io.Copy(dest, src)
	if err != nil {
		return fmt.Errorf("could not forward data from %s to %s: %v", src.LocalAddr().String(), dest.LocalAddr().String(), err)
	}
	if cw, ok := dest.(*net.TCPConn); ok {
		cw.CloseWrite()
	}
	return nil
}
