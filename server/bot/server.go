package bot

import (
	"SpectreLink/bot/proto"
	"SpectreLink/log"
	"net"
	"sync"
)

type Server struct {
	bots *sync.Map
}

func NewServer() *Server {
	return &Server{
		bots: new(sync.Map),
	}
}

func (s *Server) ListenAndServe(wg *sync.WaitGroup) {
	defer wg.Done()
	
	listen, err := net.Listen("tcp", ":2024")
	if err != nil {
		panic(err)
	}
	
	for {
		conn, err := listen.Accept()
		if err != nil {
			panic(err)
		}
		
		go s.handleConnection(conn)
	}
}

func (s *Server) handleConnection(c net.Conn) {
	conn := proto.NewProtoConnection(c)
	defer conn.Close()
	
	bot := NewBot(conn)
	
	if err := bot.ReadHandshake(); err != nil {
		log.Exception(err)
		return
	}
	
	addr := c.RemoteAddr().String()
	host, _, _ := net.SplitHostPort(addr)
	
	s.bots.Store(host, &bot)
	
	log.Infof("Bot connected: %s", host)
	
	for packet := bot.ReadPacket(); packet != nil; packet = bot.ReadPacket() {
		log.Infof("Packet: %v", packet)
	}
	
}
