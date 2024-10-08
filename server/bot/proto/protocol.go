package proto

import (
	"bufio"
	"net"
	"time"
)

type Connection struct {
	conn     net.Conn
	ByteCode ByteCode
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
	if err := c.SetReadDeadline(time.Second * 60); err != nil {
		return nil
	}

	pktType, err := readInt32(c)
	if err != nil {
		return nil
	}

	pktLen, err := readInt32(c)
	if err != nil {
		return nil
	}

	pkt := NewEmptyPacket(PacketType(pktType))

	for i := 0; i < int(pktLen); i++ {
		argType, err := readInt32(c)
		if err != nil {
			return nil
		}

		key, err := readString(c)
		if err != nil {
			return nil
		}

		value, err := readData(c, ArgType(argType))
		if err != nil {
			return nil
		}

		pkt.Args[key] = value
	}

	return pkt
}
func (c *Connection) WritePacket(pkt *Packet) error {
	if err := writeInt32(c, int32(pkt.Type)); err != nil {
		return err
	}

	if err := writeInt32(c, int32(len(pkt.Args))); err != nil {
		return err
	}

	for k, v := range pkt.Args {
		if err := writeInt32(c, int32(GetDataArgType(v))); err != nil {
			return err
		}

		if err := writeString(c, k); err != nil {
			return err
		}

		if err := writeData(c, v); err != nil {
			return err
		}
	}
	return nil
}

func (c *Connection) Close() error {
	return c.conn.Close()
}
