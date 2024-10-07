package user

type versionCommand struct {
	session *TelnetSession
}

func (c *versionCommand) Execute(args []string) error {
	return c.session.Printf("Version: 0.1-SNAPSHOT\n")
}

func (s *TelnetSession) RegisterCommands() {
	s.manager.AddCommand("version", "show the version of the application", "show the version of the application", &versionCommand{s})
}
