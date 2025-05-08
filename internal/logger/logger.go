package logger

import (
	"io"
	"log"
	"os"
	"strings"
)

type LogLevel int

const (
	DEBUG LogLevel = iota
	INFO
	WARN
	ERROR
)

var (
	debugLogger  *log.Logger
	infoLogger   *log.Logger
	warnLogger   *log.Logger
	errorLogger  *log.Logger
	currentLevel LogLevel
)

type MultiWriter struct {
	writers []io.Writer
}

func (t *MultiWriter) Write(p []byte) (n int, err error) {
	for _, w := range t.writers {
		n, err = w.Write(p)
		if err != nil {
			return
		}
		if n != len(p) {
			err = io.ErrShortWrite
			return
		}
	}
	return len(p), nil
}

func Init(level string, writers ...io.Writer) {
	if len(writers) == 0 {
		writers = []io.Writer{os.Stdout}
	}

	multiWriter := &MultiWriter{writers: writers}

	debugLogger = log.New(multiWriter, "DEBUG: ", log.Ldate|log.Ltime)
	infoLogger = log.New(multiWriter, "INFO:  ", log.Ldate|log.Ltime)
	warnLogger = log.New(multiWriter, "WARN:  ", log.Ldate|log.Ltime)
	errorLogger = log.New(multiWriter, "ERROR: ", log.Ldate|log.Ltime)

	switch strings.ToLower(level) {
	case "debug":
		currentLevel = DEBUG
	case "info":
		currentLevel = INFO
	case "warn":
		currentLevel = WARN
	case "error":
		currentLevel = ERROR
	default:
		warnLogger.Printf("Invalid log level \"%s\" using default level \"info\"", level)
		currentLevel = INFO
	}
}

func Debug(format string, v ...interface{}) {
	if currentLevel <= DEBUG {
		debugLogger.Printf(format, v...)
	}
}

func Info(format string, v ...interface{}) {
	if currentLevel <= INFO {
		infoLogger.Printf(format, v...)
	}
}

func Warn(format string, v ...interface{}) {
	if currentLevel <= WARN {
		warnLogger.Printf(format, v...)
	}
}

func Error(format string, v ...interface{}) {
	if currentLevel <= ERROR {
		errorLogger.Printf(format, v...)
	}
}
