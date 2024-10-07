package user

import (
	"fmt"
	"github.com/iskaa02/qalam/bbcode"
)

func (s *TelnetSession) Printf(format string, a ...any) error {
	msg := bbcode.Sprintf(format, a...)
	_, err := fmt.Fprint(s.terminal, msg)
	if err != nil {
		return err
	}
	return nil
}

func (s *TelnetSession) Prompt(prompt string) (string, error) {
	s.terminal.SetPrompt(prompt)
	return s.terminal.ReadLine()
}

func (s *TelnetSession) Password(prompt string) (string, error) {
	return s.terminal.ReadPassword(prompt)
}

func (s *TelnetSession) Title(title string) error {
	_, err := s.conn.Write([]byte("\033]0;" + title + "\007"))
	if err != nil {
		return err
	}
	return nil
}
