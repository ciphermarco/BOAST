package httplogger

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/ciphermarco/BOAST/log"
)

var (
	// LogEntryCtxKey is the context.Context key to store the request log entry.
	LogEntryCtxKey = &contextKey{"LogEntry"}

	// DefaultLogger is called by the Logger middleware handler to log each request.
	// Its made a package-level variable so that it can be reconfigured for custom
	// logging configurations.
	DefaultLogger = RequestLogger(&DefaultLogFormatter{})
)

// contextKey is a value for use with context.WithValue. It's used as a pointer so it
// fits in an interface{} without allocation. This technique for defining context keys
// was copied from Go 1.7's new use of context in net/http.
type contextKey struct {
	name string
}

func (k *contextKey) String() string {
	return "httplogger context value " + k.name

}

// Logger is a middleware that logs the start and end of each request, along with some
// useful data about what was requested, what the response status was, and how long it
// took to return. When standard output is a TTY, Logger will print in color, otherwise
// it will print in black and white. Logger prints a request ID if one is provided.
//
// Alternatively, look at https://github.com/goware/httplog for a more in-depth http
// logger with structured logging support.
func Logger(next http.Handler) http.Handler {
	return DefaultLogger(next)
}

// RequestLogger returns a logger handler using a custom LogFormatter.
func RequestLogger(f LogFormatter) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			entry := f.NewLogEntry(r)
			ww := NewWrapResponseWriter(w, r.ProtoMajor)

			t1 := time.Now()
			defer func() {
				entry.Write(ww.Status(), ww.BytesWritten(), ww.Header(), time.Since(t1), nil)
			}()

			next.ServeHTTP(ww, WithLogEntry(r, entry))
		}
		return http.HandlerFunc(fn)
	}
}

// LogFormatter initiates the beginning of a new LogEntry per request.
// See DefaultLogFormatter for an example implementation.
type LogFormatter interface {
	NewLogEntry(r *http.Request) LogEntry
}

// LogEntry records the final log when a request completes.
// See defaultLogEntry for an example implementation.
type LogEntry interface {
	Write(status, bytes int, header http.Header, elapsed time.Duration, extra interface{})
}

// GetLogEntry returns the in-context LogEntry for a request.
// func GetLogEntry(r *http.Request) LogEntry {
// 	entry, _ := r.Context().Value(LogEntryCtxKey).(LogEntry)
// 	return entry
// }

// WithLogEntry sets the in-context LogEntry for a request.
func WithLogEntry(r *http.Request, entry LogEntry) *http.Request {
	r = r.WithContext(context.WithValue(r.Context(), LogEntryCtxKey, entry))
	return r
}

// DefaultLogFormatter is a simple logger that implements a LogFormatter.
type DefaultLogFormatter struct{}

// NewLogEntry creates a new LogEntry for the request.
func (l *DefaultLogFormatter) NewLogEntry(r *http.Request) LogEntry {
	entry := &defaultLogEntry{
		request: r,
		buf:     &bytes.Buffer{},
	}

	entry.buf.WriteString("\"")
	fmt.Fprintf(entry.buf, "%s ", r.Method)

	scheme := "http"
	if r.TLS != nil {
		scheme = "https"
	}
	fmt.Fprintf(entry.buf, "%s://%s%s %s\" ", scheme, r.Host, r.RequestURI, r.Proto)

	entry.buf.WriteString("from ")
	entry.buf.WriteString(r.RemoteAddr)
	entry.buf.WriteString(" - ")

	return entry
}

type defaultLogEntry struct {
	request *http.Request
	buf     *bytes.Buffer
}

func (l *defaultLogEntry) Write(status, bytes int, header http.Header, elapsed time.Duration, extra interface{}) {
	fmt.Fprintf(l.buf, "%03d", status)
	fmt.Fprintf(l.buf, " %dB", bytes)
	fmt.Fprintf(l.buf, "%s", elapsed)
	log.Info(l.buf.String())
}
