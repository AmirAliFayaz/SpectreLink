package telnet

import (
	"SpectreLink/admin/user"
	"SpectreLink/log"
	"errors"
	"github.com/google/uuid"
	"github.com/jessevdk/go-flags"
	"github.com/tester2024/telnet"
	"github.com/tester2024/telnet/options"
	"io"
	"sync"
)

type Server struct {
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
				if err := session.Messagef("[#b6e3ff]" + flagsErr.Message + "[/#b6e3ff]\n"); err != nil {
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
	uid, err := uuid.NewRandom()
	if err != nil {
		panic(err)
	}
	
	sess := user.NewUserSession(conn, uid)
	
	s.sessions.Store(uid, sess)
	
	return sess
}

func (s *Server) deleteSession(uid uuid.UUID) {
	s.sessions.Delete(uid)
}

func (s *Server) ListenAddr() any {
	return s.server.Address
}

func NewTelnetServer() *Server {
	t := &Server{
		sessions: new(sync.Map),
	}
	
	t.server = telnet.NewServer(
		":1337",
		t,
		options.EchoOption,
		options.SuppressGoAheadOption,
		options.NAWSOption,
	)
	return t
}
