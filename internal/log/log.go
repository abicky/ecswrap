package log

import (
	"io"
	"log"
)

type Level uint

const (
	TraceLevel Level = iota
	DebugLevel
	InfoLevel
	WarnLevel
	ErrorLevel
	FatalLevel
)

func (l Level) String() string {
	switch l {
	case TraceLevel:
		return "TRACE"
	case DebugLevel:
		return "DEBUG"
	case InfoLevel:
		return "INFO"
	case WarnLevel:
		return "WARN"
	case ErrorLevel:
		return "ERROR"
	case FatalLevel:
		return "FATAL"
	}
	return "UNKNOWN"
}

var currentLevel = InfoLevel

func SetOutput(w io.Writer) {
	log.SetOutput(w)
}

func SetLevel(level Level) {
	currentLevel = level
}

func SetPrefix(prefix string) {
	log.SetPrefix(prefix)
}

func SetVerbosity(verbosity int) {
	level := Level(int(InfoLevel) - verbosity)
	if level < TraceLevel {
		level = TraceLevel
	} else if level > FatalLevel {
		level = FatalLevel
	}
	currentLevel = level
}

func logf(level Level, format string, v ...interface{}) {
	if level >= currentLevel {
		log.Printf("["+level.String()+"] "+format, v...)
	}
}

func Tracef(format string, v ...interface{}) {
	logf(TraceLevel, format, v...)
}

func Debugf(format string, v ...interface{}) {
	logf(DebugLevel, format, v...)
}

func Infof(format string, v ...interface{}) {
	logf(InfoLevel, format, v...)
}

func Warnf(format string, v ...interface{}) {
	logf(WarnLevel, format, v...)
}

func Errorf(format string, v ...interface{}) {
	logf(ErrorLevel, format, v...)
}

func Fatalf(format string, v ...interface{}) {
	logf(FatalLevel, format, v...)
}

func logln(level Level, v ...interface{}) {
	if level >= currentLevel {
		prefix := []interface{}{"[" + level.String() + "]"}
		log.Println(append(prefix, v...)...)
	}
}

func Traceln(v ...interface{}) {
	logln(TraceLevel, v...)
}

func Debugln(v ...interface{}) {
	logln(DebugLevel, v...)
}

func Infoln(v ...interface{}) {
	logln(InfoLevel, v...)
}

func Warnln(v ...interface{}) {
	logln(WarnLevel, v...)
}

func Errorln(v ...interface{}) {
	logln(ErrorLevel, v...)
}

func Fatalln(v ...interface{}) {
	logln(FatalLevel, v...)
}
