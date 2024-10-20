package hivebox

import (
	"encoding/json"
	"fmt"
	"io"
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
	tenHoursAgo := time.Now().UTC().Add(-10 * time.Hour)
	formattedTime := tenHoursAgo.Format(time.RFC3339)
	want := fmt.Sprintf("https://api.opensensemap.org/boxes/data?phenomenon=temperature&bbox=-180,-90,180,90&from-date=%s&format=json", formattedTime)
	got := CreateUrl()
	if want != got {
		valueError(t, got, want)
	}
}

func TestApiAvgTemp(t *testing.T) {
	url := CreateUrl()
	var data []TempData
	tenHoursAgo := time.Now().UTC().Add(-10 * time.Hour)
	resp, err := http.Get(url)
	if err != nil {
		t.Errorf("Got unexpected http error %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Invalid response code %v", resp.StatusCode)
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Errorf("Got unexpected io error %v", err)
	}
	var sensorData []SensorData
	if err := json.Unmarshal(body, &sensorData); err != nil {
		t.Errorf("Got unexpected json error %v", err)
	}
	var i int
	var total float64
	for _, data_sensor := range sensorData {
		if data_sensor.CreatedAt.After(tenHoursAgo) {
			var temp float64
			fmt.Sscanf(data_sensor.Value, "%f", &temp)
			data = append(data, TempData{data_sensor.CreatedAt, temp})
			total += temp
			i++
		} else {
			t.Fatalf("Got fetched older than 10 hours %v, ten hours ago:%v", data_sensor, tenHoursAgo)
		}
	}

	want_average := total / float64(i)
	got_average := AvgTemp(data)

	if got_average != want_average {
		t.Errorf("Average calucation is wrong. Got %.2f want %.2f", got_average, want_average)
	}

	if got_average <= -100.0 || got_average >= 100.0 {
		t.Errorf("Average not in normal range %.2f", got_average)
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
