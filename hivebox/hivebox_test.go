package hivebox

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestVersion(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/version", nil)
	w := httptest.NewRecorder()
	GetVersion(w, req)
	statusCheck(t, *w)
	var response string
	err := json.NewDecoder(w.Body).Decode(&response)
	if err != nil {
		jsonError(t, err)
	}
	if response != version {
		valueError(t, response, version)
	}
}

func TestUrl(t *testing.T) {
	tenHoursAgo := time.Now().Add(-10 * time.Hour)
	formattedTime := tenHoursAgo.Format(time.RFC3339)
	want := fmt.Sprintf("https://api.opensensemap.org/boxes/data?phenomenon=temperature&bbox=-180,-90,180,90&date=%s&format=json", formattedTime)
	got := CreateUrl()
	if want != got {
		valueError(t, got, want)
	}
}

func statusCheck(t *testing.T, w httptest.ResponseRecorder) {
	t.Helper()
	if w.Code != http.StatusOK {
		t.Fatalf("Expcted status code %v, for got %v with error: %v", http.StatusOK, w.Code, w.Body)
	}
}

func jsonError(t *testing.T, err error) {
	t.Helper()
	t.Errorf("Expcted no error, but got: %v. from decoding Json", err.Error())
}

func valueError(t *testing.T, response, item any) {
	t.Errorf("Expcted: %v, got: %v", item, response)
}
