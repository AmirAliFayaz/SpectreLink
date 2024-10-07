package bot

import (
	"SpectreLink/bot/proto"
	"SpectreLink/log"
	"errors"
	"io"
	"time"
)

type Bot struct {
	*proto.Connection
	ByteCode proto.ByteCode
}

var ErrInvalidHandshake = errors.New("invalid handshake")

func NewBot(connection *proto.Connection) *Bot {
	return &Bot{Connection: connection}
}

func (b *Bot) ReadHandshake() error {
	buffer := make([]byte, 2)
	
	if err := b.SetReadDeadline(time.Second * 10); err != nil {
		return err
	}
	
	if _, err := io.ReadFull(b.Reader, buffer); err != nil {
		return err
	}
	
	b.ByteCode = proto.ByteCode(buffer[1])
	
	if buffer[0] != 0x0 || (b.ByteCode != proto.BigEndian && b.ByteCode != proto.LittleEndian) {
		return ErrInvalidHandshake
	}
	
	log.Infof("Handshake: %v-%v", buffer, b.ByteCode)
	
	return nil
}
