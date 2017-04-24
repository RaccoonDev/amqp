// +build gofuzz

package amqp

import (
	"io"
	"net"
	"time"
)

func Fuzz(data []byte) int {
	conn, err := New(&MyConn{data: data, done: make(chan bool)}, ConnSASLPlain("listen", "3aCXZYFcuZA89xe6lZkfYJvOPnTGipA3ap7NvPruBhI="))
	if err != nil {
		return 0
	}
	defer conn.Close()

	s, err := conn.NewSession()
	if err != nil {
		return 0
	}

	r, err := s.NewReceiver(LinkSource("source"), LinkCredit(2))
	if err != nil {
		return 0
	}

	_, err = r.Receive()
	if err != nil {
		return 0
	}

	return 1
}

type MyConn struct {
	data []byte
	done chan bool
}

func (c *MyConn) Read(b []byte) (n int, err error) {
	if len(c.data) == 0 {
		return 0, io.EOF
	}
	n = copy(b, c.data)
	c.data = c.data[n:]
	return
}

func (c *MyConn) Write(b []byte) (n int, err error) {
	return len(b), nil
}

func (c *MyConn) Close() error {
	close(c.done)
	return nil
}

func (c *MyConn) LocalAddr() net.Addr {
	return &net.TCPAddr{net.IP{127, 0, 0, 1}, 49706, ""}
}

func (c *MyConn) RemoteAddr() net.Addr {
	return &net.TCPAddr{net.IP{127, 0, 0, 1}, 49706, ""}
}

func (c *MyConn) SetDeadline(t time.Time) error {
	return nil
}

func (c *MyConn) SetReadDeadline(t time.Time) error {
	return nil
}

func (c *MyConn) SetWriteDeadline(t time.Time) error {
	return nil
}