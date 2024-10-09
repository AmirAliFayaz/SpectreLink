package user

import (
	"fmt"
	"github.com/google/shlex"
	"strings"
)

type versionCommand struct {
	session  *TelnetSession
	Username string `short:"u" long:"username" description:"username"`
}

func (c *versionCommand) Execute(args []string) error {
	return c.session.Messagef("[green]Version: 1.0-SNAPSHOT[/green]\n")
}

type clearCommand struct {
	session *TelnetSession
}

func (c *clearCommand) Execute(args []string) error {
	return c.session.SendBanner()
}

func (s *TelnetSession) RegisterCommands() {
	s.manager.AddCommand("version", "show the version of the application", "show the version of the application", &versionCommand{s, ""})
	s.manager.AddCommand("clear", "clear the screen", "clear the screen", &clearCommand{s})
	
	s.registerAutoComplete()
}

func (s *TelnetSession) registerAutoComplete() {
	s.terminal.SetBracketedPasteMode(true)
	s.terminal.AutoCompleteCallback = func(line string, pos int, key rune) (newLine string, newPos int, ok bool) {
		if key == '\t' && line != "" {
			split, err := shlex.Split(line)
			if err != nil || len(split) == 0 {
				return line, pos, false
			}
			
			if strings.HasSuffix(line, " ") || strings.HasSuffix(line, "\t") || len(split) > 1 {
				return s.completeFlag(line, split, pos)
			}
			
			return s.completeCommand(line, pos)
		}
		
		return line, pos, false
	}
}

func (s *TelnetSession) completeCommand(line string, pos int) (newLine string, newPos int, ok bool) {
	for _, cmd := range s.manager.Commands() {
		if strings.HasPrefix(cmd.Name, line) {
			return cmd.Name, len(cmd.Name), true
		}
		
		for _, alias := range cmd.Aliases {
			if strings.HasPrefix(alias, line) {
				return alias, len(alias), true
			}
		}
	}
	
	return line, pos, false
}

func (s *TelnetSession) completeFlag(line string, parts []string, pos int) (newLine string, newPos int, ok bool) {
	cmd := s.manager.Find(parts[0])
	if cmd == nil {
		return line, pos, false
	}
	
	options := cmd.Options()
	
	if len(parts) == 1 {
		if len(options) == 0 {
			return line, pos, false
		}
		
		long := options[0].LongName
		newLine = fmt.Sprintf("%s--%s", line, long)
		return newLine, len(newLine), true
	}
	
	knownTarget := parts[len(parts)-1]
	target := knownTarget
	target = strings.TrimPrefix(target, "--")
	target = strings.TrimPrefix(target, "-")
	target = strings.TrimPrefix(target, "/")
	
	for _, option := range options {
		line = strings.ReplaceAll(line, knownTarget, "")
		
		short := string(option.ShortName)
		if strings.HasPrefix(short, target) {
			newLine = fmt.Sprintf("%s-%s", line, short)
			return newLine, len(newLine), true
		}
		
		long := option.LongName
		if strings.HasPrefix(long, target) {
			newLine = fmt.Sprintf("%s--%s", line, long)
			return newLine, len(newLine), true
		}
		
	}
	
	return line, pos, true
}
