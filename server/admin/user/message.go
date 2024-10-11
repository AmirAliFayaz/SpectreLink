package user

import (
	"embed"
	"fmt"
	"github.com/iskaa02/qalam/bbcode"
	"github.com/mazznoer/colorgrad"
	"strings"
	"unicode/utf8"
)

func (s *TelnetSession) Printf(format string, a ...any) error {
	msg := bbcode.Sprintf(format, a...)
	msg = strings.ReplaceAll(msg, "\r\n", "\n")
	msg = strings.ReplaceAll(msg, "\n", "\r\n")
	_, err := fmt.Fprint(s.terminal, msg)
	if err != nil {
		return err
	}
	return nil
}

func (s *TelnetSession) Rprintf(format string, a ...any) error {
	msg := strings.NewReplacer("\r\n", "\n", "\r", "\n").Replace(fmt.Sprintf(format, a...))
	_, err := fmt.Fprintf(s.terminal, msg)
	if err != nil {
		return err
	}
	return nil
}

func (s *TelnetSession) Gradientf(grad colorgrad.Gradient, format string, a ...interface{}) error {
	formatted := fmt.Sprintf(format, a...)

	replace := strings.NewReplacer("\r\n", "\n", "\r", "\n").Replace(formatted)
	lines := strings.Split(replace, "\n")

	for idx, line := range lines {
		count := utf8.RuneCountInString(line)

		for i, run := range line {
			color := grad.At(float64(utf8.RuneCountInString(line[:i])) / float64(count))

			r, g, b, _ := color.RGBA255()
			if _, err := fmt.Fprintf(s.terminal, "\x1b[38;2;%d;%d;%dm%c\x1b[0m", r, g, b, run); err != nil {
				return err
			}
		}

		if idx != len(lines)-1 || strings.HasSuffix(formatted, "\n") {
			if _, err := fmt.Fprint(s.terminal, "\r\n"); err != nil {
				return err
			}
		}
	}

	return nil
}

func (s *TelnetSession) Center(text string) string {
	replace := strings.NewReplacer("\r\n", "\n", "\r", "\n").Replace(text)
	lines := strings.Split(replace, "\n")

	builder := new(strings.Builder)

	width := int(s.Width)

	for idx, line := range lines {
		length := utf8.RuneCountInString(line)

		spaces := (width - length) / 2
		if spaces < 0 {
			spaces = 0
		}

		builder.WriteString(strings.Repeat(" ", spaces))
		builder.WriteString(line)

		if idx != len(lines)-1 || strings.HasSuffix(text, "\n") {
			builder.WriteString("\r\n")
		}
	}

	return builder.String()
}

func (s *TelnetSession) Clear() error {
	return s.Printf("\x1b[2J\x1b[0;0H")
}

//go:embed banner.txt
var banner embed.FS

func (s *TelnetSession) SendBanner() error {
	if err := s.Clear(); err != nil {
		return err
	}

	bannerData, err := banner.ReadFile("banner.txt")
	if err != nil {
		return err
	}

	lines := strings.Split(string(bannerData), "\n")
	maxWidth := 0
	for _, line := range lines {
		width := utf8.RuneCountInString(line)
		if width > maxWidth {
			maxWidth = width
		}
	}

	if maxWidth > int(s.Width) || len(lines) > int(s.Height) {
		return nil
	}

	gradient := colorgrad.Spectral()
	centered := s.Center(string(bannerData))
	if err := s.Gradientf(gradient, centered); err != nil {
		return err
	}

	return nil
}

func (s *TelnetSession) Promptf(prompt string, args ...any) (string, error) {
	s.terminal.SetPrompt(bbcode.Sprintf(prompt, args...))
	return s.terminal.ReadLine()
}

func (s *TelnetSession) Password(prompt string) (string, error) {
	return s.terminal.ReadPassword(prompt)
}

func (s *TelnetSession) Titlef(text string, a ...interface{}) error {
	_, err := s.conn.Write([]byte(fmt.Sprintf("\x1b]0;%v\x07", fmt.Sprintf(text, a...))))
	return err
}

func (s *TelnetSession) Error(err string) error {
	return s.Messagef("[#ff0000][!][/#ff0000] [#cc4949]%s[/#cc4949]\n", err)
}

func (s *TelnetSession) Messagef(msg string, a ...interface{}) error {
	lines := strings.Split("\n"+bbcode.Sprintf(msg, a...), "\n")

	for _, line := range lines {
		line = " " + line

		if err := s.Rprintf("%s\n", line); err != nil {
			return err
		}
	}

	return nil
}
