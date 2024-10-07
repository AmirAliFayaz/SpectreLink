package user

import (
	"SpectreLink/log"
	"github.com/google/shlex"
	"github.com/google/uuid"
	"github.com/jessevdk/go-flags"
	"github.com/tester2024/telnet"
	"github.com/tester2024/telnet/options"
	"golang.org/x/crypto/ssh/terminal"
	"time"
)

type TelnetSession struct {
	terminal *terminal.Terminal
	conn     *telnet.Connection
	manager  *flags.Parser
	uid      uuid.UUID
}

func (s *TelnetSession) ReadCommand() ([]string, error) {
	data, err := s.Prompt("> ")
	if err != nil {
		return nil, err
	}

	args, err := shlex.Split(data)
	if err != nil {
		return nil, err
	}

	return s.manager.ParseArgs(args)
}

func (s *TelnetSession) Error(err string) error {
	return s.Printf("[red b]%v\n[/red b]\n", err)
}

func (s *TelnetSession) Destroy(fn func(uid uuid.UUID)) {
	s.conn.Close()

	s.manager = nil
	s.terminal = nil
	s.conn = nil

	fn(s.uid)
}

func (s *TelnetSession) ReadKey() error {
	_, err := s.terminal.ReadLine()
	if err != nil {
		return err
	}
	return nil
}

func NewUserSession(conn *telnet.Connection, uid uuid.UUID) *TelnetSession {
	return &TelnetSession{
		conn:     conn,
		uid:      uid,
		manager:  flags.NewNamedParser("", (flags.Default^flags.PrintErrors)|flags.AllowBoolValues),
		terminal: terminal.NewTerminal(conn, ""),
	}
}

func (s *TelnetSession) Handle() {
	for {
		if s.conn == nil || s.terminal == nil {
			return
		}

		naws := s.conn.OptionHandlers[telnet.TeloptNAWS].(*options.NAWSHandler)
		if err := s.terminal.SetSize(int(naws.Width), int(naws.Height)); err != nil {
			log.Exception(err)
			return
		}

		if err := s.Title("SpectreLink | Admin"); err != nil {
			log.Exception(err)
			return
		}
		time.Sleep(time.Second)
	}
}
