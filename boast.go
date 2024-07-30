package boast

import (
	"crypto/rand"
	"encoding/base32"
	"encoding/json"
	"strings"
	"time"

	"github.com/ciphermarco/BOAST/log"
)

// Storage represents the BOAST's storage implementation.
// It's implemented by any type that provides these methods so it can be easily swapped
// by a DB or other kind of storage if needed.
type Storage interface {
	SetTest(secret []byte) (id string, canary string, err error)
	SearchTest(f func(k, v string) bool) (id string, canary string)
	StoreEvent(evt Event) error
	LoadEvents(id string) (evts []Event, loaded bool)
	TotalTests() int
	TotalEvents() int
	StartExpire(err chan error)
}

// Event represents an interaction event.
type Event struct {
	ID         string    `json:"id"`
	Time       time.Time `json:"time"`
	TestID     string    `json:"testID"`
	Receiver   string    `json:"receiver"`
	RemoteAddr string    `json:"remoteAddress,omitempty"`
	Dump       string    `json:"dump,omitempty"`
	QueryType  string    `json:"queryType,omitempty"`
}

// String satisfies the Stringer interface for pretty-printing Event.
// This should only be used for debugging.
func (e *Event) String() string {
	s, err := json.MarshalIndent(e, "", "\t")
	if err != nil {
		log.Debug("Event's String method error: %v", err)
		return ""
	}
	return string(s)
}

// NewEvent allocates a new Event struct and returns its copy.
// The raison d'être of this function is to provide an easy interface to generate an
// event with a standard ID without the caller having to deal with it.
func NewEvent(testID, receiver, addr, dump string) (Event, error) {
	id, err := genEventID()
	if err != nil {
		return Event{}, err
	}

	return Event{
		ID:         id,
		Time:       time.Now(),
		TestID:     testID,
		Receiver:   receiver,
		RemoteAddr: addr,
		Dump:       dump,
	}, nil
}

// NewDNSEvent allocates a new Event using NewEvent but with the difference of recording
// the passed DNS query type to keep more information for DNS queries.
func NewDNSEvent(testID, receiver, addr, dump, qType string) (Event, error) {
	evt, err := NewEvent(testID, receiver, addr, dump)
	if err != nil {
		return evt, err
	}
	evt.QueryType = qType
	return evt, nil
}

// genEventID generates a random event ID.
// The returned event ID is a 16 bytes long value base32 encoded by ToBase32.
func genEventID() (string, error) {
	c := 16
	random := make([]byte, c)

	if _, err := rand.Read(random); err != nil {
		return "", err
	}

	res := ToBase32(random)
	return res, nil
}

// ToBase32 encodes b to the base32 format used by BOAST's components.
func ToBase32(b []byte) string {
	enc := base32.StdEncoding.WithPadding(-1)
	res := enc.EncodeToString(b)
	res = strings.ToLower(res)
	return res
}
