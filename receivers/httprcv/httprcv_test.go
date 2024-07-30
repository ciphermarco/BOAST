package httprcv_test

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/ciphermarco/BOAST/log"
	"github.com/ciphermarco/BOAST/receivers/httprcv"
)

func TestMain(m *testing.M) {
	log.SetOutput(ioutil.Discard)
	os.Exit(m.Run())
}

func TestExpectedCanary(t *testing.T) {
	req, err := http.NewRequest("GET", "/mpqhomfbxab55m5de32mywvfoy", nil)
	if err != nil {
		t.Fatal(err)
	}

	mockStrg := &mockStorage{}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(httprcv.CatchAll(mockStrg, ""))
	handler.ServeHTTP(rr, req)

	checkStatusCode(http.StatusOK, rr.Code, t)

	want := fmt.Sprintf("<html><body>%s</body></html>", tCanary)
	got := rr.Body.String()

	if want != got {
		t.Errorf("canary not found: %v (want) != %v (got)", want, got)
	}
}

func checkStatusCode(want int, got int, t *testing.T) {
	if want != got {
		t.Errorf("handler returned wrong status code: %v (want) != %v (got)",
			want, got)
	}
}
