package user

import (
	"SpectreLink/log"
	"errors"
	"github.com/google/shlex"
	"github.com/jessevdk/go-flags"
	"github.com/tester2024/telnet"
	"github.com/tester2024/telnet/options"
	"golang.org/x/term"
	"time"
)

type TelnetSession struct {
	terminal      *term.Terminal
	conn          *telnet.Connection
	manager       *flags.Parser
	uid           string
	Width, Height uint16
}

func (s *TelnetSession) Write(p []byte) (n int, err error) {
	return s.terminal.Write(p)
}

func (s *TelnetSession) ReadCommand() ([]string, error) {
	if s.manager == nil || s.terminal == nil || s.conn == nil {
		return nil, errors.New("invalid session")
	}

	if err := s.Printf(" [bg#2f2f2f][#26c6da] ✧ %s 💀 SpecterLink [/#26c6da][/bg#2f2f2f][#2f2f2f]►[/#2f2f2f] ", "root"); err != nil {
		return nil, err
	}

	data, err := s.Promptf("")
	if err != nil {
		return nil, err
	}

	args, err := shlex.Split(data)
	if err != nil {
		return nil, err
	}

	return s.manager.ParseArgs(args)
}
func (s *TelnetSession) Destroy(fn func(uid string)) {
	if s.conn != nil {
		s.conn.Close()
	}

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

func NewUserSession(conn *telnet.Connection, uid string) *TelnetSession {
	terminal := term.NewTerminal(conn, "")
	return &TelnetSession{
		conn:     conn,
		uid:      uid,
		manager:  flags.NewNamedParser("", (flags.Default^flags.PrintErrors)|flags.AllowBoolValues),
		terminal: terminal,
	}
}

func (s *TelnetSession) SetSize(width, height uint16) error {
	if width == 0 || height == 0 {
		return nil
	}
	if width == s.Width && height == s.Height {
		return nil
	}

	log.Infof("Setting size to %d x %d", width, height)

	s.Width, s.Height = width, height

	return s.terminal.SetSize(int(width), int(height))
}

func (s *TelnetSession) Handle() {
	for {
		if s.conn == nil || s.terminal == nil {
			return
		}

		s.UpdateSize()

		if err := s.Titlef("SpectreLink | Admin"); err != nil {
			log.Exception(err)
			return
		}

		time.Sleep(time.Second)
	}
}

func (s *TelnetSession) UpdateSize() {
	naws := s.conn.OptionHandlers[telnet.TeloptNAWS].(*options.NAWSHandler)
	if err := s.SetSize(naws.Width, naws.Height); err != nil {
		log.Exception(err)
		return
	}
}
