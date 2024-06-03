package logger

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

type Logger struct {
	Level       LogLevel
	writer      io.Writer
	indentation int
}

type LogLevel int

const (
	TraceLevel LogLevel = iota
	DebugLevel
	InfoLevel
	WarnLevel
	ErrorLevel
)

var traceStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("8")).Background(lipgloss.Color("15"))
var errorStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("15")).Background(lipgloss.Color("168"))
var infoStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("8")).Background(lipgloss.Color("111"))
var warnStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("8")).Background(lipgloss.Color("222"))
var debugStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("8")).Background(lipgloss.Color("149"))

var DefaultLogger = &Logger{writer: os.Stderr}

func NewLogger(writer io.Writer, level LogLevel) *Logger {
	if writer == nil {
		panic("Logger received nil writter")
	}

	return &Logger{
		writer: writer,
		Level:  level,
	}
}

func (l *Logger) SetIndentation(indLevel int) {
	l.indentation = indLevel
}

func (l *Logger) ClearIndentation() {
	l.indentation = 0
}

func (l *Logger) Printf(format string, v ...any) {
	l.doPrintf(-1, format, v...)
}

func (l *Logger) Errorf(format string, v ...any) {
	if l.Level > ErrorLevel {
		return
	}

	l.doPrintf(ErrorLevel, format, v...)
}

func (l *Logger) Warnf(format string, v ...any) {
	if l.Level > WarnLevel {
		return
	}

	l.doPrintf(WarnLevel, format, v...)
}

func (l *Logger) Infof(format string, v ...any) {
	if l.Level > InfoLevel {
		return
	}

	l.doPrintf(InfoLevel, format, v...)
}

func (l *Logger) Debugf(format string, v ...any) {
	if l.Level > DebugLevel {
		return
	}

	l.doPrintf(DebugLevel, format, v...)
}

func (l *Logger) Tracef(format string, v ...any) {
	if l.Level > TraceLevel {
		return
	}

	l.doPrintf(TraceLevel, format, v...)
}

func (l *Logger) Error(str string) {
	if l.Level > ErrorLevel {
		return
	}

	l.doPrint(ErrorLevel, str)
}

func (l *Logger) Warn(str string) {
	if l.Level > WarnLevel {
		return
	}

	l.doPrint(WarnLevel, str)
}

func (l *Logger) Info(str string) {
	if l.Level > InfoLevel {
		return
	}

	l.doPrint(InfoLevel, str)
}

func (l *Logger) Debug(str string) {
	if l.Level > DebugLevel {
		return
	}

	l.doPrint(DebugLevel, str)
}

func (l *Logger) Trace(str string) {
	if l.Level > TraceLevel {
		return
	}

	l.doPrint(TraceLevel, str)
}

func (l *Logger) Print(str string) {
	l.doPrint(-1, str)
}

func (l *Logger) getWriter() io.Writer {
	if l.writer == nil {
		return os.Stderr
	}

	return l.writer
}

func (l *Logger) getPrefix(level LogLevel) string {
	var prefix string
	switch level {
	case TraceLevel:
		prefix = traceStyle.Render(" TRACE ") + " "
	case DebugLevel:
		prefix = debugStyle.Render(" DEBUG ") + " "
	case InfoLevel:
		prefix = infoStyle.Render(" INFO ") + " "
	case WarnLevel:
		prefix = warnStyle.Render(" WARN ") + " "
	case ErrorLevel:
		prefix = errorStyle.Render(" ERROR ") + " "
	}

	return prefix
}

func (l *Logger) doPrintf(level LogLevel, format string, v ...any) {
	w := l.getWriter()
	str := fmt.Sprintf(format, v...)

	if l.indentation == 0 {
		w.Write([]byte(l.getPrefix(level)))
		w.Write([]byte(str))
	} else {
		s := lipgloss.JoinHorizontal(0, strings.Repeat("    ", l.indentation), str)
		w.Write([]byte(s))
	}
	w.Write([]byte("\n"))
}

func (l *Logger) doPrint(level LogLevel, str string) {
	w := l.getWriter()

	if l.indentation == 0 {
		w.Write([]byte(l.getPrefix(level)))
		w.Write([]byte(str))
	} else {
		s := lipgloss.JoinHorizontal(0, strings.Repeat("    ", l.indentation), str)
		w.Write([]byte(s))
	}
	w.Write([]byte("\n"))
}
