package main

import (
	"flag"
	"net"
	"net/http"
	"time"
)

var (
	originProviderHost  = "https://registry.terraform.io"
	privateProviderHost = ""

	timeout    = 5 * time.Second
	httpClient = &http.Client{
		Transport: &http.Transport{
			Dial: func(netw, addr string) (net.Conn, error) {
				conn, err := net.DialTimeout(netw, addr, timeout)
				if err != nil {
					return nil, err
				}

				return NewTimeoutConn(conn, timeout), nil
			},
			ResponseHeaderTimeout: timeout,
		},
	}
)

func init() {
	flag.StringVar(&privateProviderHost, "private_host", "https://localhost", "The host for private providers")
}

func main() {
	flag.Parse()

	startServer()
}

type TimeoutConn struct {
	conn    net.Conn
	timeout time.Duration
}

func NewTimeoutConn(conn net.Conn, timeout time.Duration) *TimeoutConn {
	return &TimeoutConn{
		conn:    conn,
		timeout: timeout,
	}
}

func (c *TimeoutConn) Read(b []byte) (n int, err error) {
	c.SetReadDeadline(time.Now().Add(c.timeout))
	return c.conn.Read(b)
}

func (c *TimeoutConn) Write(b []byte) (n int, err error) {
	c.SetWriteDeadline(time.Now().Add(c.timeout))
	return c.conn.Write(b)
}

func (c *TimeoutConn) Close() error {
	return c.conn.Close()
}

func (c *TimeoutConn) LocalAddr() net.Addr {
	return c.conn.LocalAddr()
}

func (c *TimeoutConn) RemoteAddr() net.Addr {
	return c.conn.RemoteAddr()
}

func (c *TimeoutConn) SetDeadline(t time.Time) error {
	return c.conn.SetDeadline(t)
}

func (c *TimeoutConn) SetReadDeadline(t time.Time) error {
	return c.conn.SetReadDeadline(t)
}

func (c *TimeoutConn) SetWriteDeadline(t time.Time) error {
	return c.conn.SetWriteDeadline(t)
}
