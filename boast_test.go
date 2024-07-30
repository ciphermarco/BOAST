package boast_test

import (
	"reflect"
	"testing"
	"time"

	app "github.com/ciphermarco/BOAST"
)

func TestNewEvent(t *testing.T) {
	want := app.Event{
		ID:         "TEST ID",
		Time:       time.Now(),
		TestID:     "TEST TestID",
		Receiver:   "TEST Receiver",
		RemoteAddr: "TEST RemoteAddr",
		Dump:       "TEST Dump",
	}
	got, err := app.NewEvent(
		"TEST TestID",
		"TEST Receiver",
		"TEST RemoteAddr",
		"TEST Dump",
	)

	if err != nil {
		t.Errorf("unexpected error: %v (want) != %v (got)", nil, err)
	}
	if want.Time.UnixNano() >= got.Time.UnixNano() {
		t.Errorf("wrong Time: %v (want) != %v (got)", want.Time, got.Time)
	}

	want.Time = got.Time

	wantLenID := 26
	gotLenID := len(got.ID)
	if wantLenID != gotLenID {
		t.Errorf("wrong ID length: %v (want) != %v (got)",
			wantLenID, gotLenID)
	}

	want.ID = got.ID

	if !reflect.DeepEqual(want, got) {
		t.Errorf("wrong event")
		t.Errorf("Want:")
		t.Errorf("%+v", want)
		t.Errorf("Got:")
		t.Errorf("%+v", got)
	}
}

func TestNewDNSEvent(t *testing.T) {
	want := app.Event{
		ID:         "TEST ID",
		Time:       time.Now(),
		TestID:     "TEST TestID",
		Receiver:   "TEST Receiver",
		RemoteAddr: "TEST RemoteAddr",
		Dump:       "TEST Dump",
		QueryType:  "TEST QueryType",
	}
	got, err := app.NewDNSEvent(
		"TEST TestID",
		"TEST Receiver",
		"TEST RemoteAddr",
		"TEST Dump",
		"TEST QueryType",
	)

	if err != nil {
		t.Errorf("unexpected error: %v (want) != %v (got)", nil, err)
	}
	if want.Time.UnixNano() >= got.Time.UnixNano() {
		t.Errorf("wrong Time: %v (want) != %v (got)", want.Time, got.Time)
	}

	want.Time = got.Time

	wantLenID := 26
	gotLenID := len(got.ID)
	if wantLenID != gotLenID {
		t.Errorf("wrong ID length: %v (want) != %v (got)",
			wantLenID, gotLenID)
	}

	want.ID = got.ID

	if !reflect.DeepEqual(want, got) {
		t.Errorf("wrong event")
		t.Errorf("Want:")
		t.Errorf("%+v", want)
		t.Errorf("Got:")
		t.Errorf("%+v", got)
	}
}

func TestToBase32(t *testing.T) {
	want := map[string]string{
		"smcdpjlzu":    "onwwgzdqnjwhu5i",
		"bdwbpjiic":    "mjshoytqnjuwsyy",
		"epizdfjvnjt":  "mvygs6temzvhm3tkoq",
		"fahkwl":       "mzqwq23xnq",
		"xfzc":         "pbthuyy",
		"xvvdji":       "pb3hmzdkne",
		"ufmuhvebnxr":  "ovtg25liozswe3tyoi",
		"whelamjko":    "o5ugk3dbnvvgw3y",
		"amcaabgp":     "mfwwgylbmjtxa",
		"zkqsgvnkhs":   "pjvxc43hozxgw2dt",
		"dcbf":         "mrrwezq",
		"ceehfulcs":    "mnswk2dgovwgg4y",
		"rstfckezdp":   "ojzxiztdnnsxuzdq",
		"apbwwsoetv":   "mfyge53xonxwk5dw",
		"dpmvb":        "mryg25tc",
		"pgsl":         "obtxg3a",
		"tajzyt":       "orqwu6tzoq",
		"nufpbgwlqpxx": "nz2wm4dcm53wy4lqpb4a",
		"vprqw":        "ozyhe4lx",
		"yhejthchj":    "pfugk2tunbrwq2q",
	}

	for k, v := range want {
		got := app.ToBase32([]byte(k))
		if v != got {
			t.Errorf("wrong base32: %v (want) != %v (got)", want, got)
		}
	}
}
