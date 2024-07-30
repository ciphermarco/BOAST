package log

import (
	"fmt"
	"io"
	stdLog "log"
	"os"
	"syscall"
	"time"
)

type level int

const (
	debug level = iota
	info
)

var (
	// Logger represents a custom logging object.
	// It's exported so it can be used by api.httplogger until it's changed.
	Logger   = stdLog.New(&logWriter{out: os.Stdout}, "", stdLog.Lshortfile)
	curLevel = info
	labels   = map[level]string{
		debug: "DEBUG",
		info:  "INFO",
	}
)

type logWriter struct {
	out io.Writer
}

func (w logWriter) Write(b []byte) (int, error) {
	tid := fmt.Sprintf(" %d ", syscall.Gettid())
	return fmt.Fprint(w.out, time.Now().UTC().Format("2006-01-02T15:04:05.999Z")+tid+string(b))
}

func log(lvl level, format string, v ...interface{}) {
	if lvl < curLevel || curLevel > info || curLevel < debug {
		return
	}
	format = fmt.Sprintf("[%s] %s", labels[lvl], format)
	Logger.Output(3, fmt.Sprintf(format, v...))
}

// SetLevel sets the logging level.
func SetLevel(lvl int) {
	curLevel = level(lvl)
}

// SetOutput sets the output to a new io.Writer object.
func SetOutput(w io.Writer) {
	Logger.SetOutput(&logWriter{out: w})
}

// Info logs an INFO logging line.
func Info(format string, v ...interface{}) {
	log(info, format, v...)
}

// Debug logs a DEBUG logging line.
func Debug(format string, v ...interface{}) {
	log(debug, format, v...)
}

// Printf calls logger.Output to print to the logger without any labels.
// Arguments are handled in the manner of fmt.Printf.
func Printf(format string, v ...interface{}) {
	Logger.Output(2, fmt.Sprintf(format, v...))
}

// Print calls logger.Output to print to the logger without any labels.
// Arguments are handled in the manner of fmt.Print.
func Print(v ...interface{}) {
	Logger.Output(2, fmt.Sprint(v...))
}

// Println calls logger.Output to print to the logger without any labels.
// Arguments are handled in the manner of fmt.Println.
func Println(v ...interface{}) {
	Logger.Output(2, fmt.Sprintln(v...))
}

// Fatalf is equivalent to Printf() followed by a call to os.Exit(1).
func Fatalf(format string, v ...interface{}) {
	Logger.Output(2, fmt.Sprintf(format, v...))
	os.Exit(1)
}

// Fatalln is equivalent to Println() followed by a call to os.Exit(1).
func Fatalln(v ...interface{}) {
	Logger.Output(2, fmt.Sprintln(v...))
	os.Exit(1)
}
