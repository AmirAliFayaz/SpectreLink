package log

import (
	"github.com/charmbracelet/lipgloss"
	glog "github.com/charmbracelet/log"
	"github.com/iskaa02/qalam/bbcode"
	"os"
	"strings"
	"time"
)

var logger = glog.NewWithOptions(os.Stderr, glog.Options{
	ReportTimestamp: true,
	TimeFormat:      time.Kitchen,
	Prefix:          "SpectreLink",
})

var errlogger = glog.NewWithOptions(os.Stderr, glog.Options{
	ReportCaller:    true,
	ReportTimestamp: true,
	TimeFormat:      time.Kitchen,
	Prefix:          "SpectreLink",
})

func init() {
	styles := glog.DefaultStyles()

	styles.Timestamp = lipgloss.NewStyle().
		Padding(0, 1, 0, 1).
		Background(lipgloss.Color("#54aeff")).
		Foreground(lipgloss.Color("#000000"))

	styles.Caller = lipgloss.NewStyle().
		Faint(true).
		Padding(0, 1, 0, 1).
		Background(lipgloss.Color("#bb0000")).
		Foreground(lipgloss.Color("#000000"))

	styles.Prefix = lipgloss.NewStyle().
		Faint(true).
		Bold(true).
		Padding(0, 1, 0, 1).
		Background(lipgloss.Color("#d29922")).
		Foreground(lipgloss.Color("#000000"))

	styles.Levels[glog.ErrorLevel] = lipgloss.NewStyle().
		Bold(true).
		SetString("ERROR").
		Padding(0, 1, 0, 1).
		Background(lipgloss.Color("#FF0000")).
		Foreground(lipgloss.Color("#FFFFFF"))

	styles.Levels[glog.WarnLevel] = lipgloss.NewStyle().
		Bold(true).
		SetString("WARN").
		Padding(0, 1, 0, 1).
		Background(lipgloss.Color("#FFFF00")).
		Foreground(lipgloss.Color("#000000"))

	styles.Levels[glog.InfoLevel] = lipgloss.NewStyle().
		Bold(true).
		SetString("INFO").
		Padding(0, 1, 0, 1).
		Background(lipgloss.Color("#00FF00")).
		Foreground(lipgloss.Color("#000000"))

	styles.Levels[glog.DebugLevel] = lipgloss.NewStyle().
		Bold(true).
		SetString("DEBUG").
		Padding(0, 1, 0, 1).
		Background(lipgloss.Color("#00FFFF")).
		Foreground(lipgloss.Color("#000000"))

	errlogger.SetStyles(styles)
	logger.SetStyles(styles)
}

func Exception(err error, message ...string) {
	if err == nil {
		return
	}

	errlogger.Helper()

	if len(message) != 0 {
		errlogger.Error(bbcode.Sprintf(strings.Join(message, " ")), "err", err)
	} else {
		errlogger.Error(err)
	}

}

func Infof(format string, a ...any) {
	logger.Infof(bbcode.Sprintf(format, a...))
}

func Errorf(format string, a ...any) {
	logger.Errorf(bbcode.Sprintf(format, a...))
}

func Debugf(format string, a ...any) {
	logger.Debugf(bbcode.Sprintf(format, a...))
}

func Warnf(format string, a ...any) {
	logger.Warnf(bbcode.Sprintf(format, a...))
}
