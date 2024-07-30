package api_test

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/ciphermarco/BOAST/api"
	"github.com/ciphermarco/BOAST/log"
)

func TestMain(m *testing.M) {
	log.SetOutput(ioutil.Discard)
	os.Exit(m.Run())
}

func TestAuthorizeSuccess(t *testing.T) {
	req, err := newEventsRequest()
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Add("Authorization", fmt.Sprintf("Secret %s", tB64TestSecret))

	mockStrg := &mockStorage{}
	rr := httptest.NewRecorder()
	handler := api.NewTestAPI("/test-status", mockStrg)
	handler.ServeHTTP(rr, req)

	checkStatusCode(http.StatusOK, rr.Code, t)
	checkEventsBody(rr.Body, tTest.ID, 0, t)
}

func TestAuthorizeWithoutHeader(t *testing.T) {
	req, err := newEventsRequest()
	if err != nil {
		t.Fatal(err)
	}
	// Two slightly broken headers
	req.Header.Add("Authorizatio", fmt.Sprintf("Secret %s", tB64TestSecret))
	req.Header.Add("Authorization", fmt.Sprintf("Secret%s", tB64TestSecret))

	mockStrg := &mockStorage{}
	rr := httptest.NewRecorder()
	handler := api.NewTestAPI("/test-status", mockStrg)
	handler.ServeHTTP(rr, req)

	checkStatusCode(http.StatusUnauthorized, rr.Code, t)
	checkEventsBody(rr.Body, "", 0, t)
}

func TestAuthorizeWithWrongHeaderFormat(t *testing.T) {
	req, err := newEventsRequest()
	if err != nil {
		t.Fatal(err)
	}

	mockStrg := &mockStorage{}
	rr := httptest.NewRecorder()
	handler := api.NewTestAPI("/test-status", mockStrg)
	handler.ServeHTTP(rr, req)

	checkStatusCode(http.StatusUnauthorized, rr.Code, t)
	checkEventsBody(rr.Body, "", 0, t)
}

func TestAuthorizeWithWrongAuthType(t *testing.T) {
	req, err := newEventsRequest()
	if err != nil {
		t.Fatal(err)
	}
	// Secrt instead of Secret
	req.Header.Add("Authorization", fmt.Sprintf("Secrt %s", tB64TestSecret))

	mockStrg := &mockStorage{}
	rr := httptest.NewRecorder()
	handler := api.NewTestAPI("/test-status", mockStrg)
	handler.ServeHTTP(rr, req)

	checkStatusCode(http.StatusUnauthorized, rr.Code, t)
	checkEventsBody(rr.Body, "", 0, t)
}

func TestAuthorizeWithTooLongSecret(t *testing.T) {
	req, err := newEventsRequest()
	if err != nil {
		t.Fatal(err)
	}
	maxSize := 44
	longSecret := base64.StdEncoding.EncodeToString(randBytes(maxSize + 1))
	req.Header.Add("Authorization", fmt.Sprintf("Secret %s", longSecret))

	mockStrg := &mockStorage{}
	rr := httptest.NewRecorder()
	handler := api.NewTestAPI("/test-status", mockStrg)
	handler.ServeHTTP(rr, req)

	checkStatusCode(http.StatusUnauthorized, rr.Code, t)
	checkEventsBody(rr.Body, "", 0, t)
}

func TestAuthorizeWithInvalidBase64(t *testing.T) {
	req, err := newEventsRequest()
	if err != nil {
		t.Fatal(err)
	}
	invalidB64Secret := tB64TestSecret[1:]
	req.Header.Add("Authorization", fmt.Sprintf("Secret %s", invalidB64Secret))

	mockStrg := &mockStorage{}
	rr := httptest.NewRecorder()
	handler := api.NewTestAPI("/test-status", mockStrg)
	handler.ServeHTTP(rr, req)

	checkStatusCode(http.StatusUnauthorized, rr.Code, t)
	checkEventsBody(rr.Body, "", 0, t)
}

type errMockStorage struct {
	mockStorage
}

func (e *errMockStorage) SetTest(secret []byte) (string, string, error) {
	return "", "", errors.New("fake error")
}

func TestAuthorizeWithSetTestError(t *testing.T) {
	req, err := newEventsRequest()
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Add("Authorization", fmt.Sprintf("Secret %s", tB64TestSecret))

	mockStrg := &errMockStorage{}
	rr := httptest.NewRecorder()
	handler := api.NewTestAPI("/test-status", mockStrg)
	handler.ServeHTTP(rr, req)

	checkStatusCode(http.StatusUnauthorized, rr.Code, t)
	checkEventsBody(rr.Body, "", 0, t)
}

func newEventsRequest(headers ...map[string]string) (req *http.Request, err error) {
	req, err = http.NewRequest("GET", "/events", nil)
	if err != nil {
		return nil, err
	}
	return req, nil
}

func checkStatusCode(want int, got int, t *testing.T) {
	if want != got {
		t.Errorf("handler returned wrong status code: %v (want) != %v (got)",
			want, got)
	}
}

func checkEventsBody(buf *bytes.Buffer, wantID string, wantEventsLen int, t *testing.T) {
	res, err := unmarshalEventsResponse(buf)
	if err != nil {
		t.Fatal(err)
	}

	if res.ID != wantID {
		t.Errorf("wrong ID: %v (want) != %v (got)", wantID, res.ID)
	}
	if len(res.Events) != wantEventsLen {
		t.Errorf("wrong Events length: %v (want) != %v (got)", wantEventsLen, len(res.Events))
	}
}

func unmarshalEventsResponse(buf *bytes.Buffer) (res mockEventsResponse, err error) {
	b, err := ioutil.ReadAll(buf)
	if err != nil {
		return res, err
	}
	if err = json.Unmarshal(b, &res); err != nil {
		return res, err
	}
	return res, err
}

func randBytes(l int) []byte {
	rand.Seed(time.Now().UnixNano())
	b := make([]byte, l)
	for i := range b {
		b[i] = byte(rand.Intn(255))
	}
	return b
}
