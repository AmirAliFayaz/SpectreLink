package proto

import (
	"bufio"
	"net"
	"time"
)

type Connection struct {
	conn net.Conn
	*bufio.Writer
	*bufio.Reader
}

func NewProtoConnection(conn net.Conn) *Connection {
	return &Connection{
		conn:   conn,
		Writer: bufio.NewWriter(conn),
		Reader: bufio.NewReader(conn),
	}
}

func (c *Connection) SetReadDeadline(duration time.Duration) error {
	return c.conn.SetReadDeadline(time.Now().Add(duration))
}

func (c *Connection) SetWriteDeadline(duration time.Duration) error {
	return c.conn.SetWriteDeadline(time.Now().Add(duration))
}

func (c *Connection) ReadPacket() *Packet {
	return nil
}

func (c *Connection) Close() error {
	return c.conn.Close()
}
