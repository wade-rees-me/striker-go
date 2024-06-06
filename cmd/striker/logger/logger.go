package logger

import (
	"io"
	"io/ioutil"
	"log"
	"os"
)

type Logger struct {
	infoLogger    *log.Logger
	warningLogger *log.Logger
	errorLogger   *log.Logger
	debugLogger   *log.Logger
	debugFile     *os.File
}

var Log = NewLogger(os.Stdout, os.Stdout, os.Stderr, ioutil.Discard)

func (l *Logger) OpenDebugFile(fileName string) {
	//debugFile, err := os.OpenFile(fileName, os.O_CREATE|os.O_WRONLY, 0666)
	debugFile, err := os.Create(fileName)
	if err != nil {
		log.Fatalf("Failed to open debug log file: %v", err)
	}
	l.debugLogger.SetOutput(debugFile)
	l.debugFile = debugFile
}

func (l *Logger) CloseDebugFile() {
	l.debugFile.Close()
}

// NewLogger creates a new instance of Logger
func NewLogger(infoHandle io.Writer, warningHandle io.Writer, errorHandle io.Writer, debugHandle io.Writer) *Logger {
	return &Logger{
		infoLogger:    log.New(infoHandle, "INFO: ", log.Ldate|log.Ltime /*|log.Lshortfile*/),
		warningLogger: log.New(warningHandle, "WARNING: ", log.Ldate|log.Ltime /*|log.Lshortfile*/),
		errorLogger:   log.New(errorHandle, "ERROR: ", log.Ldate|log.Ltime /*|log.Lshortfile*/),
		debugLogger:   log.New(debugHandle, "DEBUG: ", log.Ldate|log.Ltime /*|log.Lshortfile*/),
	}
}

// Info logs an info message
func (l *Logger) Info(message string) {
	l.infoLogger.Println(message)
	l.debugLogger.Println(message)
}

// Warning logs a warning message
func (l *Logger) Warning(message string) {
	l.warningLogger.Println(message)
	l.debugLogger.Println(message)
}

// Error logs an error message
func (l *Logger) Error(message string) {
	l.errorLogger.Println(message)
	l.debugLogger.Println(message)
}

// Debug logs a debug message
func (l *Logger) Debug(message string) {
	l.debugLogger.Println(message)
}
