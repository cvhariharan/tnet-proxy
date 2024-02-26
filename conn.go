package proxy

import (
	"net"
	"time"
)

// Conn wraps net.Conn and provides idle timeouts (seconds) for reads and writes
type Conn struct {
	net.Conn
	IdleTimeout int
}

// Write is a wrapper around net.Conn write and sets write deadline
func (c *Conn) Write(b []byte) (n int, err error) {
	// Reset write deadline
	if c.IdleTimeout > 0 {
		c.SetWriteDeadline(time.Now().Add(time.Duration(c.IdleTimeout) * time.Second))
	}
	return c.Conn.Write(b)
}

// Read is a wrapper around net.Conn read and sets read deadline
func (c *Conn) Read(b []byte) (n int, err error) {
	// Reset Read deadline
	if c.IdleTimeout > 0 {
		c.SetReadDeadline(time.Now().Add(time.Duration(c.IdleTimeout) * time.Second))
	}
	return c.Conn.Read(b)
}
