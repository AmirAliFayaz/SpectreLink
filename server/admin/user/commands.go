package user

type versionCommand struct {
	session *TelnetSession
}

func (c *versionCommand) Execute(args []string) error {
	return c.session.Printf("[green]Version: 0.1-SNAPSHOT[/green]\n")
}

type clearCommand struct {
	session *TelnetSession
}

func (c *clearCommand) Execute(args []string) error {
	return c.session.SendBanner()
}

func (s *TelnetSession) RegisterCommands() {
	s.manager.AddCommand("version", "show the version of the application", "show the version of the application", &versionCommand{s})
	s.manager.AddCommand("clear", "clear the screen", "clear the screen", &clearCommand{s})
}
