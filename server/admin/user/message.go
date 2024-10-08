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
	_, err := fmt.Fprint(s.terminal, msg)
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

			r, g, b := color.RGB255()
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

		if length == 0 {
			continue
		}

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

	data, err := banner.ReadFile("banner.txt")
	if err != nil {
		return err
	}

	//grad, _ := colorgrad.NewGradient().
	//	HtmlColors("#800080", "#0000FF", "#00FFFF").
	//	Build()

	if err := s.Gradientf(
		colorgrad.Spectral(),
		s.Center(string(data)),
	); err != nil {
		return err
	}
	return err
}

func (s *TelnetSession) Prompt(prompt string) (string, error) {
	s.terminal.SetPrompt(prompt)
	return s.terminal.ReadLine()
}

func (s *TelnetSession) Password(prompt string) (string, error) {
	return s.terminal.ReadPassword(prompt)
}

func (s *TelnetSession) Titlef(text string, a ...interface{}) error {
	return s.Printf("\x1b]0;%s\x07", fmt.Sprintf(text, a...))
}
