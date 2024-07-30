package httprcv_test

import app "github.com/ciphermarco/BOAST"

var tID = "mpqhomfbxab55m5de32mywvfoy"
var tCanary = "k2b27meg7dfifvxuxmnfnm24oa"

type mockStorage struct{}

func (s *mockStorage) SetTest(secret []byte) (id string, canary string, err error) {
	return tID, tCanary, nil
}

func (s *mockStorage) LoadEvents(id string) (evts []app.Event, loaded bool) {
	return evts, false
}

func (s *mockStorage) SearchTest(f func(k, v string) bool) (id string, canary string) {
	return tID, tCanary
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
