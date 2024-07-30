package log_test

import (
	"bytes"
	"testing"

	"github.com/ciphermarco/BOAST/log"
)

func TestSetLevelWorks(t *testing.T) {
	var buf bytes.Buffer

	log.SetOutput(&buf)
	log.SetLevel(0) // debug
	log.Debug("testing Debug logs")

	got := buf.String()
	if got == "" {
		t.Errorf("empty log line: <log line content> (want) != \"%v\" (got)",
			got)
	}

	buf.Reset()
	log.Info("testing Info logs")

	got = buf.String()
	if got == "" {
		t.Errorf("empty log line: <log line content> (want) != \"%v\" (got)",
			got)
	}

	buf.Reset()
	log.SetLevel(1) // info
	log.Debug("testing debug does not log")

	got = buf.String()
	if got != "" {
		t.Errorf("wrong log line: <empty log line> (want) != \"%v\" (got)",
			got)
	}

	buf.Reset()
	log.Info("testing info still logs")

	got = buf.String()
	if got == "" {
		t.Errorf("empty log line: <log line content> (want) != \"%v\" (got)",
			got)
	}
}

func TestPrintfLogs(t *testing.T) {
	var buf bytes.Buffer

	log.SetOutput(&buf)
	log.SetLevel(0) // debug
	log.Printf("testing Printf logs")

	got := buf.String()
	if got == "" {
		t.Errorf("empty log line: <log line content> (want) != \"%v\" (got)",
			got)
	}

	buf.Reset()
	log.SetLevel(1) // info
	log.Printf("testing Printf still logs")

	got = buf.String()
	if got == "" {
		t.Errorf("empty log line: <log line content> (want) != \"%v\" (got)",
			got)
	}
}

func TestPrintLogs(t *testing.T) {
	var buf bytes.Buffer

	log.SetOutput(&buf)
	log.SetLevel(0) // debug
	log.Print("testing Printf logs")

	got := buf.String()
	if got == "" {
		t.Errorf("empty log line: <log line content> (want) != \"%v\" (got)",
			got)
	}

	buf.Reset()
	log.SetLevel(1) // info
	log.Print("testing Printf still logs")

	got = buf.String()
	if got == "" {
		t.Errorf("empty log line: <log line content> (want) != \"%v\" (got)",
			got)
	}
}

func TestPrintlnLogs(t *testing.T) {
	var buf bytes.Buffer

	log.SetOutput(&buf)
	log.SetLevel(0) // debug
	log.Println("testing Println logs")

	got := buf.String()
	if got == "" {
		t.Errorf("empty log line: <log line content> (want) != \"%v\" (got)",
			got)
	}

	buf.Reset()
	log.SetLevel(1) // info
	log.Println("testing Println still logs")

	got = buf.String()
	if got == "" {
		t.Errorf("empty log line: <log line content> (want) != \"%v\" (got)",
			got)
	}
}
