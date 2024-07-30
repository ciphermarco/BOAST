package storage_test

import (
	"container/heap"
	"io/ioutil"

	"math/rand"
	"os"
	"reflect"
	"testing"
	"time"

	app "github.com/ciphermarco/BOAST"
	"github.com/ciphermarco/BOAST/log"
	"github.com/ciphermarco/BOAST/storage"
)

type testEnv struct {
	strg *storage.ExportStorage
	cfg  *storage.Config
}

func newTestEnv() *testEnv {
	cfg := storage.NewTestConfig()
	return &testEnv{
		strg: storage.NewTestStorage(cfg),
		cfg:  cfg,
	}
}

func TestMain(m *testing.M) {
	log.SetOutput(ioutil.Discard)
	os.Exit(m.Run())
}

func TestNew(t *testing.T) {
	env := newTestEnv()

	want := storage.NewMockStorage(env.cfg)
	got, err := storage.New(env.cfg)

	if err != nil {
		t.Errorf("unexpected error: %v (want) != %v (got)", nil, err)
	}

	if !reflect.DeepEqual(want, got) {
		t.Errorf("wrong storage")
		t.Errorf("Want:")
		t.Errorf("%+v", want)
		t.Errorf("Got:")
		t.Errorf("%+v", got)
	}

	tCfg := storage.NewTestConfig()
	tCfg.HMACKey = storage.RandBytes(65)
	_, gotErr := storage.New(tCfg)

	if gotErr == nil {
		t.Errorf("did not fail: error (want) != %v (got)", gotErr)
	}
}

func TestSetTestBasic(t *testing.T) {
	env := newTestEnv()

	wantID := storage.TTest.ID()
	wantCanary := storage.TTest.Canary()
	gotID, gotCanary, err := env.strg.SetTest(storage.TTest.Secret)

	if err != nil {
		t.Errorf("unexpected error: %v (want) != %v (got)", nil, err)
	}
	if wantID != gotID {
		t.Errorf("wrong ID: %v (want) != %v (got)", wantID, gotID)
	}
	if wantCanary != gotCanary {
		t.Errorf("wrong Canary: %v (want) != %v (got)", wantCanary, gotCanary)
	}

	wantTotal := 1
	gotTotal := env.strg.TotalTests()

	if wantTotal != gotTotal {
		t.Errorf("wrong total: %v (want) != %v (got)", wantTotal, gotTotal)
	}

	// wantID does not change
	// wantCanary does not change
	gotID, gotCanary, err = env.strg.SetTest(storage.TTest.Secret)

	if err != nil {
		t.Errorf("unexpected error: %v (want) != %v (got)", nil, err)
	}
	if wantID != gotID {
		t.Errorf("wrong ID: %v (want) != %v (got)", wantID, gotID)
	}
	if wantCanary != gotCanary {
		t.Errorf("wrong Canary: %v (want) != %v (got)", wantCanary, gotCanary)
	}

	// wantTotal does not change
	gotTotal = env.strg.TotalTests()

	if wantTotal != gotTotal {
		t.Errorf("wrong total: %v (want) != %v (got)", wantTotal, gotTotal)
	}

}

func TestCreateSameTestMultipleTimes(t *testing.T) {
	env := newTestEnv()

	wantID := storage.TTest.ID()
	wantCanary := storage.TTest.Canary()
	wantTotal := 1
	rand.Seed(time.Now().UnixNano())
	for i := 0; i < rand.Intn(env.strg.MaxTests()); i++ {
		gotID, gotCanary, err := env.strg.SetTest(storage.TTest.Secret)

		if err != nil {
			t.Errorf("unexpected error: %v (want) != %v (got)", nil, err)
		}
		if wantID != gotID {
			t.Errorf("wrong ID: %v (want) != %v (got)", wantID, gotID)
			break
		}
		if wantCanary != gotCanary {
			t.Errorf("wrong Canary: %v (want) != %v (got)", wantCanary, gotCanary)
		}

		gotTotal := env.strg.TotalTests()

		if wantTotal != gotTotal {
			t.Errorf("wrong total: %v (want) != %v (got)", wantTotal, gotTotal)
			break
		}
	}
}

func TestCreateDifferengTests(t *testing.T) {
	env := newTestEnv()

	rand.Seed(time.Now().UnixNano())
	wantTotal := rand.Intn(env.strg.MaxTests())

	for i := 0; i < wantTotal; i++ {
		_, _, err := env.strg.SetTest(storage.RandBytes(32))
		if err != nil {
			t.Errorf("unexpected error: %v (want) != %v (got)", nil, err)
			break
		}
	}

	gotTotal := env.strg.TotalTests()

	if wantTotal != gotTotal {
		t.Errorf("wrong total: %v (want) != %v (got)", wantTotal, gotTotal)
	}
}

func TestSetTestLimit(t *testing.T) {
	env := newTestEnv()

	wantmaxTests := env.strg.MaxEvents() / env.strg.MaxEventsByTest()
	gotmaxTests := env.strg.MaxTests()

	if wantmaxTests != gotmaxTests {
		t.Errorf("wrong max tests: %v (want) != %v (got)", wantmaxTests, gotmaxTests)
	}

	for i := 0; i < gotmaxTests; i++ {
		_, _, err := env.strg.SetTest(storage.RandBytes(8))
		if err != nil {
			t.Errorf("unexpected error: %v (want) != %v (got)", nil, err)
			break
		}
	}

	_, _, err := env.strg.SetTest(storage.RandBytes(8))
	if err == nil {
		t.Errorf("did not fail: error (want) != %v (got)", err)
	}
}

func TestSearchTest(t *testing.T) {
	env := newTestEnv()
	f := func(searchID string) func(key, value string) bool {
		return func(key, value string) bool {
			return searchID == key
		}
	}

	wantID, wantCanary := "", ""
	gotID, gotCanary := env.strg.SearchTest(f(wantID))

	if wantID != gotID {
		t.Errorf("wrong ID: %v (want) != %v (got)", wantID, gotID)
	}
	if wantCanary != gotCanary {
		t.Errorf("wrong canary: %v (want) != %v (got)", wantCanary, gotCanary)
	}

	env.strg.SetTest(storage.TTest.Secret)

	rand.Seed(time.Now().UnixNano())
	for i := 0; i < rand.Intn(11); i++ {
		env.strg.SetTest(storage.RandBytes(8))
	}

	wantID, wantCanary = storage.TTest.ID(), storage.TTest.Canary()
	gotID, gotCanary = env.strg.SearchTest(f(wantID))

	if wantID != gotID {
		t.Errorf("wrong ID: %v (want) != %v (got)", wantID, gotID)
	}
	if wantCanary != gotCanary {
		t.Errorf("wrong canary: %v (want) != %v (got)", wantCanary, gotCanary)
	}
}

func TestStoreEvent(t *testing.T) {
	env := newTestEnv()
	evt := storage.NewTestEvent()

	err := env.strg.StoreEvent(evt)
	if err == nil {
		t.Errorf("did not fail: error (want) != %v (got)", err)
	}

	env.strg.SetTest(storage.TTest.Secret)
	err = env.strg.StoreEvent(evt)
	if err != nil {
		t.Errorf("unexpected error: %v (want) != %v (got)", nil, err)
	}
}

func TestStoreEventLimit(t *testing.T) {
	env := newTestEnv()
	evt := storage.NewTestEvent()
	env.strg.SetTest(storage.TTest.Secret)

	totalEvts := env.strg.MaxEventsByTest() + 10
	for i := 0; i < totalEvts; i++ {
		err := env.strg.StoreEvent(evt)
		if err != nil {
			t.Errorf("unexpected error: %v (want) != %v (got)", nil, err)
			break
		}
	}

	wantTotal := env.strg.MaxEventsByTest()
	gotTotal := env.strg.TotalEvents()

	if wantTotal != gotTotal {
		t.Errorf("wrong total: %v (want) != %v (got)", wantTotal, gotTotal)
	}
}

func TestLoadEvents(t *testing.T) {
	env := newTestEnv()
	evt := storage.NewTestEvent()
	env.strg.SetTest(storage.TTest.Secret)

	totalEvts := env.strg.MaxEventsByTest() + 10
	for i := 0; i < totalEvts; i++ {
		err := env.strg.StoreEvent(evt)
		if err != nil {
			t.Errorf("unexpected error: %v (want) != %v (got)", nil, err)
			break
		}
	}

	wantLoaded := true
	wantTotal := env.strg.MaxEventsByTest()
	gotEvts, gotLoaded := env.strg.LoadEvents(storage.TTest.ID())
	gotTotal := len(gotEvts)

	if wantLoaded != gotLoaded {
		t.Errorf("response was not set as loaded: %v (want) != %v (got)", wantLoaded, gotLoaded)
	}
	if wantTotal != gotTotal {
		t.Errorf("wrong length: %v (want) != %v (got)", wantTotal, gotTotal)
	}

	var wantEvts []app.Event
	wantLoaded = false
	gotEvts, gotLoaded = env.strg.LoadEvents(string(storage.RandBytes(8)))

	if wantLoaded != gotLoaded {
		t.Errorf("response was set as loaded: %v (want) != %v (got)", wantLoaded, gotLoaded)
	}
	if !reflect.DeepEqual(wantEvts, gotEvts) {
		t.Errorf("wrong slice: %v (want) != %v (got)", wantEvts, gotEvts)
	}
}

func TestTotalTests(t *testing.T) {
	env := newTestEnv()

	wantTotal := env.strg.MaxTests() - 1
	for i := 0; i < wantTotal; i++ {
		env.strg.SetTest(storage.RandBytes(8))
	}
	gotTotal := env.strg.TotalTests()

	if wantTotal != gotTotal {
		t.Errorf("wrong total: %v (want) != %v (got)", wantTotal, gotTotal)
	}
}

func TestTotalEvents(t *testing.T) {
	env := newTestEnv()
	evt := storage.NewTestEvent()
	env.strg.SetTest(storage.TTest.Secret)

	totalEvts := env.strg.MaxEventsByTest() - 2
	for i := 0; i < totalEvts; i++ {
		err := env.strg.StoreEvent(evt)
		if err != nil {
			t.Errorf("unexpected error: %v (want) != %v (got)", nil, err)
		}
	}

	wantTotal := totalEvts
	gotTotal := env.strg.TotalEvents()

	if wantTotal != gotTotal {
		t.Errorf("wrong total: %v (want) != %v (got)", wantTotal, gotTotal)
	}
}

func TestExpiration(t *testing.T) {
	check := func(want, got int) {
		if want != got {
			t.Errorf("wrong total: %v (want) != %v (got)", want, got)
		}
	}

	tCfg := storage.NewTestConfig()
	tCfg.TTL = 500 * time.Millisecond
	tCfg.CheckInterval = 1 * time.Millisecond
	tStrg := storage.NewTestStorage(tCfg)
	tStrg.SetTest(storage.TTest.Secret)
	// sets a test meant to keep without events
	tStrg.SetTest([]byte("2sqGqj4FQubefsqqiEksJg=="))

	totalEvts := 3
	for i := 0; i < totalEvts; i++ {
		err := tStrg.StoreEvent(storage.NewTestEvent())
		if err != nil {
			t.Fatal(err)
		}
		wantTotal := i + 1
		check(wantTotal, tStrg.TotalEvents())
	}

	tErr := make(chan error, 1)
	go tStrg.StartExpire(tErr)
	time.Sleep(tStrg.TTL())
	time.Sleep(tStrg.CheckInterval())

	wantTotal := 0
	check(wantTotal, tStrg.TotalEvents())
	check(wantTotal, tStrg.TotalTests())
}

func TestHeapPushWithWrongType(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("function panicked: r == %v (want) r == \"%v\" (got)", nil, r)
		}
	}()
	evts := storage.NewEmptyEventsHeap()
	heap.Init(evts)
	wrongType := "not an event"
	heap.Push(evts, wrongType)

	wantTotal := 0
	gotTotal := evts.Len()

	if wantTotal != gotTotal {
		t.Errorf("wrong length: %v (want) != %v (got)", wantTotal, gotTotal)
	}
}

var evtsBench []app.Event
var loadedBench bool

func BenchmarkLoadEvents(b *testing.B) {
	tCfg := storage.NewTestConfig()
	tCfg.MaxEvents = 10_000
	tCfg.MaxEventsByTest = 10_000
	tStrg := storage.NewTestStorage(tCfg)
	id, _, err := tStrg.SetTest(storage.TTest.Secret)
	if err != nil {
		b.Fatal(err)
	}

	evt := storage.NewTestEvent()
	for i := 0; i < tCfg.MaxEvents; i++ {
		tStrg.StoreEvent(evt)
	}

	var evts []app.Event
	var loaded bool
	for n := 0; n < b.N; n++ {
		evts, loaded = tStrg.LoadEvents(id)
	}

	// To avoid compiler optimisations that could eliminate the function
	// call if the value is not used.
	evtsBench, loadedBench = evts, loaded
}

var idBench string
var canaryBench string

func BenchmarkSearchTest(b *testing.B) {
	tCfg := storage.NewTestConfig()
	tCfg.MaxEvents = 10_000
	tCfg.MaxEventsByTest = 1
	tStrg := storage.NewTestStorage(tCfg)
	_, _, err := tStrg.SetTest(storage.TTest.Secret)
	if err != nil {
		b.Fatal(err)
	}
	f := func(searchID string) func(key, value string) bool {
		return func(key, value string) bool {
			return searchID == key
		}
	}

	for i := 0; i < tStrg.MaxTests(); i++ {
		tStrg.SetTest(storage.RandBytes(8))
	}

	var id string
	var canary string
	for n := 0; n < b.N; n++ {
		id, canary = tStrg.SearchTest(f("unexistent test id"))
	}

	// To avoid compiler optimisations that could eliminate the function
	// call if the value is not used.
	idBench, canaryBench = id, canary
}
