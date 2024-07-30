package api_test

import (
	"bytes"
	"encoding/base64"
	"errors"
	"time"

	app "github.com/ciphermarco/BOAST"
)

var tB64TestSecret = "872k5eD/lGRbMZ3GqIPB0bUzqRjBlt1lhLH4+/42sKa="

var tTestSecret, _ = base64.StdEncoding.DecodeString(tB64TestSecret)

var tTest = mockTest{
	ID:     "mpqhomfbxab55m5de32mywvfoy",
	Canary: "k2b27meg7dfifvxuxmnfnm24oa",
	Events: []app.Event{},
}

type mockTest struct {
	ID     string
	Canary string
	Events []app.Event
}

type mockEventsResponse struct {
	ID     string      `json:"id"`
	Canary string      `json:"canary"`
	Events []app.Event `json:"events"`
}

type mockStorage struct{}

func (s *mockStorage) SetTest(secret []byte) (id string, canary string, err error) {
	if eq := bytes.Compare(secret, tTestSecret); eq != 0 {
		return "", "", errors.New("mock test not found")
	}
	return tTest.ID, tTest.Canary, nil
}

func (s *mockStorage) LoadEvents(id string) (evts []app.Event, loaded bool) {
	mockEvt := app.Event{
		ID:         "TEST ID",
		Time:       time.Now(),
		TestID:     "TEST TestID",
		Receiver:   "TEST Receiver",
		RemoteAddr: "203.0.113.113",
		Dump:       "TEST DUMP",
		QueryType:  "TEST QueryType",
	}
	evts = []app.Event{mockEvt}
	return evts, false
}

func (s *mockStorage) SearchTest(f func(k, v string) bool) (id string, canary string) {
	return "", ""
}

func (s *mockStorage) StoreEvent(evt app.Event) error {
	return nil
}

func (s *mockStorage) TotalTests() int {
	return 0
}

func (s *mockStorage) TotalEvents() int {
	return 0
}

func (s *mockStorage) StartExpire(err chan error) {}
