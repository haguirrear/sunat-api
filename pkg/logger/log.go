package logger

import (
	"fmt"
	"io"
	"os"
)

type Logger struct {
	Level  LogLevel
	writer io.Writer
}

type LogLevel int

const (
	TraceLevel LogLevel = iota
	DebugLevel
	InfoLevel
	WarnLevel
	ErrorLevel
)

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

func (l *Logger) Errorf(format string, v ...any) {
	if l.Level > ErrorLevel {
		return
	}

	l.doPrintf(format, v...)
}

func (l *Logger) Warnf(format string, v ...any) {
	if l.Level > WarnLevel {
		return
	}

	l.doPrintf(format, v...)
}

func (l *Logger) Infof(format string, v ...any) {
	if l.Level > InfoLevel {
		return
	}

	l.doPrintf(format, v...)
}

func (l *Logger) Debugf(format string, v ...any) {
	if l.Level > DebugLevel {
		return
	}

	l.doPrintf(format, v...)
}

func (l *Logger) Tracef(format string, v ...any) {
	if l.Level > TraceLevel {
		return
	}

	l.doPrintf(format, v...)
}

func (l *Logger) Error(str string) {
	if l.Level > ErrorLevel {
		return
	}

	l.doPrint(str)
}

func (l *Logger) Warn(str string) {
	if l.Level > WarnLevel {
		return
	}

	l.doPrint(str)
}

func (l *Logger) Info(str string) {
	if l.Level > InfoLevel {
		return
	}

	l.doPrint(str)
}

func (l *Logger) Debug(str string) {
	if l.Level > DebugLevel {
		return
	}

	l.doPrint(str)
}

func (l *Logger) Trace(str string) {
	if l.Level > TraceLevel {
		return
	}

	l.doPrint(str)
}

func (l *Logger) getWriter() io.Writer {
	if l.writer == nil {
		return os.Stderr
	}

	return l.writer
}

func (l *Logger) getPrefix() string {
	var prefix string
	switch l.Level {
	case TraceLevel:
		prefix = "[TRACE] "
	case DebugLevel:
		prefix = "[DEBUG] "
	case InfoLevel:
		prefix = "[INFO] "
	case WarnLevel:
		prefix = "[WARN] "
	case ErrorLevel:
		prefix = "[ERROR] "
	}

	return prefix
}

func (l *Logger) doPrintf(format string, v ...any) {
	w := l.getWriter()

	w.Write([]byte(l.getPrefix()))
	w.Write([]byte(fmt.Sprintf(format, v...)))
	w.Write([]byte("\n"))
}

func (l *Logger) doPrint(str string) {
	w := l.getWriter()

	w.Write([]byte(l.getPrefix()))
	w.Write([]byte(str))
	w.Write([]byte("\n"))
}
