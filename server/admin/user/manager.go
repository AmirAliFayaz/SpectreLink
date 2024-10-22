package user

import (
	"SpectreLink/log"
	"fmt"
	"github.com/google/shlex"
	"github.com/jessevdk/go-flags"
	"strings"
)

func (s *TelnetSession) RegisterCommands() {
	_, _ = s.manager.AddCommand("version", "show the version of the application", "show the version of the application", &versionCommand{s})
	_, _ = s.manager.AddCommand("clear", "clear the screen", "clear the screen", &clearCommand{s})
	_, _ = s.manager.AddCommand("exit", "exit the application", "exit the application", &exitCommand{s})
	_, _ = s.manager.AddCommand("help", "show the help", "show the help", &helpCommand{s})
	
	s.registerAutoComplete()
	s.handleAttacks()
}

func (s *TelnetSession) handleAttacks() {
	s.manager.CommandHandler = func(command flags.Commander, args []string) error {
		log.Infof("Args: %v", args)
		
		if command == nil {
			return nil
		}
		
		return command.Execute(args)
	}
}

func (s *TelnetSession) registerAutoComplete() {
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
	target := strings.TrimLeft(knownTarget, "-/")
	
	for _, option := range options {
		line = strings.ReplaceAll(line, knownTarget, "")
		
		short := string(option.ShortName)
		long := option.LongName
		
		for _, l := range parts {
			l = strings.TrimLeft(l, "-/")
			
			if strings.HasPrefix(l, long) || strings.HasPrefix(l, short) {
				goto next
			}
		}
		
		if strings.HasPrefix(short, target) {
			newLine = fmt.Sprintf("%s-%s", line, short)
			return newLine, len(newLine), true
		}
		
		if strings.HasPrefix(long, target) {
			newLine = fmt.Sprintf("%s--%s", line, long)
			return newLine, len(newLine), true
		}
	
	next:
		continue
	}
	
	return line, pos, false
}
