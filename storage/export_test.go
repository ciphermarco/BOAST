package storage

import (
	"encoding/base64"
	"hash"
	"math/rand"
	"time"

	app "github.com/ciphermarco/BOAST"
	"github.com/ciphermarco/BOAST/log"

	"golang.org/x/crypto/blake2b"
)

type ExportTest struct {
	test
	Secret []byte
}

var tTestSecret, _ = base64.StdEncoding.DecodeString(
	"872k5eD/lGRbMZ3GqIPB0bUzqRjBlt1lhLH4+/42sKa=")

var TTest = &ExportTest{
	test: test{
		id:     "mpqhomfbxab55m5de32mywvfoy",
		canary: "k2b27meg7dfifvxuxmnfnm24oa",
		events: &eventHeap{},
	},
	Secret: tTestSecret,
}

func (u *ExportTest) ID() string {
	return u.id
}

func (u *ExportTest) Canary() string {
	return u.canary
}

type ExportStorage struct {
	*Storage
}

func (s *ExportStorage) TotalEvents() int {
	return s.Storage.TotalEvents()
}

func (s *ExportStorage) MaxTests() int {
	return s.maxTests
}

func (s *ExportStorage) MaxEvents() int {
	return s.cfg.MaxEvents
}

func (s *ExportStorage) MaxEventsByTest() int {
	return s.cfg.MaxEventsByTest
}

func (s *ExportStorage) TTL() time.Duration {
	return s.cfg.TTL
}

func (s *ExportStorage) CheckInterval() time.Duration {
	return s.cfg.CheckInterval
}

func NewTestConfig() *Config {
	return &Config{
		TTL:             100 * time.Minute,
		CheckInterval:   1 * time.Second,
		MaxRestarts:     10,
		MaxEvents:       1000,
		MaxEventsByTest: 10,
		HMACKey:         []byte("testing"),
	}
}

func NewTestStorage(cfg *Config) *ExportStorage {
	strg, err := New(cfg)
	if err != nil {
		log.Fatalln("NewTestStorage:", err)
	}
	return &ExportStorage{
		Storage: strg,
	}
}

func NewMockStorage(cfg *Config) *Storage {
	hmac := NewTestHMAC(cfg.HMACKey)
	return &Storage{
		tests:    make(map[string]test),
		maxTests: cfg.MaxEvents / cfg.MaxEventsByTest,
		hmac:     hmac,
		cfg:      *cfg,
	}
}

func NewTestHMAC(key []byte) hash.Hash {
	hmac, err := blake2b.New256(key)
	if err != nil {
		log.Fatalln("NewTestHMAC:", err)
	}
	return hmac
}

func NewTestEvent() app.Event {
	return app.Event{
		ID:         string(RandBytes(16)),
		Time:       time.Now(),
		TestID:     TTest.id,
		Receiver:   "TEST Receiver",
		RemoteAddr: "203.0.113.113",
		Dump:       "TEST Dump",
		QueryType:  "TEST QueryType",
	}
}

type ExportEventHeap struct {
	*eventHeap
}

func NewEmptyEventsHeap() *ExportEventHeap {
	return &ExportEventHeap{eventHeap: &eventHeap{}}
}

func RandBytes(l int) []byte {
	rand.Seed(time.Now().UnixNano())
	b := make([]byte, l)
	for i := range b {
		b[i] = byte(rand.Intn(255))
	}
	return b
}
