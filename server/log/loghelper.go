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
	TimeFormat:      time.DateTime,
	Prefix:          "SpectreLink",
})

var errlogger = glog.NewWithOptions(os.Stderr, glog.Options{
	ReportCaller:    true,
	ReportTimestamp: true,
	TimeFormat:      time.DateTime,
	Prefix:          "SpectreLink",
})

func init() {
	styles := glog.DefaultStyles()

	styles.Timestamp = lipgloss.NewStyle().
		Padding(0, 1, 0, 1).
		SetString("🕒").
		Background(lipgloss.Color("#AEEEEE")).
		Foreground(lipgloss.Color("#000000"))

	styles.Caller = lipgloss.NewStyle().
		Faint(true).
		Padding(0, 1, 0, 1).
		Background(lipgloss.Color("#FFDAB9")).
		Foreground(lipgloss.Color("#000000"))

	styles.Prefix = lipgloss.NewStyle().
		Faint(true).
		Bold(true).
		Padding(0, 1, 0, 1).
		Background(lipgloss.Color("#90EE90")).
		Foreground(lipgloss.Color("#000000"))

	styles.Levels[glog.ErrorLevel] = lipgloss.NewStyle().
		Bold(true).
		SetString("❌ ERROR").
		Padding(0, 1, 0, 1).
		Background(lipgloss.Color("#FFB3B3")).
		Foreground(lipgloss.Color("#000000"))

	styles.Levels[glog.WarnLevel] = lipgloss.NewStyle().
		Bold(true).
		SetString("⚠️ WARN").
		Padding(0, 1, 0, 1).
		Background(lipgloss.Color("#FFFFCC")).
		Foreground(lipgloss.Color("#000000"))

	styles.Levels[glog.InfoLevel] = lipgloss.NewStyle().
		Bold(true).
		SetString("ℹ️ INFO").
		Padding(0, 1, 0, 1).
		Background(lipgloss.Color("#CCFFFF")).
		Foreground(lipgloss.Color("#000000"))

	styles.Levels[glog.DebugLevel] = lipgloss.NewStyle().
		Bold(true).
		SetString("🐞 DEBUG").
		Padding(0, 1, 0, 1).
		Background(lipgloss.Color("#FFA5CC")).
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
