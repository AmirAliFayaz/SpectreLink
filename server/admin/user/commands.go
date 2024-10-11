package user

import (
	"SpectreLink/log"
)

type versionCommand struct {
	session *TelnetSession
}

type clearCommand struct {
	session *TelnetSession
}

type helpCommand struct {
	session *TelnetSession
}

type exitCommand struct {
	session *TelnetSession
}

type attackMethods struct {
	UDPFlood string `short:"u" long:"udp" description:"udp flood"`
}

func (c *exitCommand) Execute(_ []string) error {
	c.session.Destroy(func(uid string) {
		log.Infof("User %s exited", uid)
	})
	return nil
}

func (c *helpCommand) Execute(_ []string) error {
	c.session.manager.WriteHelp(c.session)
	return nil
}

func (c *versionCommand) Execute(_ []string) error {
	return c.session.Messagef("[green]Version: 1.0-SNAPSHOT[/green]\n")
}

func (c *clearCommand) Execute(_ []string) error {
	return c.session.SendBanner()
}
