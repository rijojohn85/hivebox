package hivebox

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"
)

const version = "0.0.1"

type SensorData struct {
	CreatedAt time.Time `json:"createdAt"`
	Value     string    `json:"value"`
}

type TempData struct {
	CreatedAt time.Time
	Value     float64
}

func GetVersion(w http.ResponseWriter, req *http.Request) {
	json.NewEncoder(w).Encode(version)
}

func GetAvgTemp(w http.ResponseWriter, req *http.Request) {
	url := CreateUrl()
	resp, err := http.Get(url)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	var sensorData []SensorData
	if err := json.Unmarshal(body, &sensorData); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	var i int
	var total float64
	for _, data_sensor := range sensorData {
		var temp float64
		fmt.Sscanf(data_sensor.Value, "%f", &temp)
		total += temp
		i++
	}

	avg := total / float64(i)
	json.NewEncoder(w).Encode(avg)
}

func CreateUrl() string {
	tenHoursAgo := time.Now().UTC().Add(-10 * time.Hour)
	formattedTime := tenHoursAgo.Format(time.RFC3339)
	url := fmt.Sprintf("https://api.opensensemap.org/boxes/data?phenomenon=temperature&bbox=-180,-90,180,90&from-date=%s&format=json", formattedTime)
	return url
}

func FetchData(apiUrl string) ([]TempData, error) {
	var data []TempData
	tenHoursAgo := time.Now().Add(-10 * time.Hour)
	resp, err := http.Get(apiUrl)
	if err != nil {
		return data, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return data, errors.New("Bad status code returned")
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return data, err
	}
	var sensorData []SensorData
	if err := json.Unmarshal(body, &sensorData); err != nil {
		return data, err
	}
	for _, data_sensor := range sensorData {
		if data_sensor.CreatedAt.After(tenHoursAgo) {
			var temp float64
			fmt.Sscanf(data_sensor.Value, "%f", &temp)
			data = append(data, TempData{data_sensor.CreatedAt, temp})
		}
	}
	return data, nil
}

func AvgTemp(data []TempData) float64 {
	var total float64
	for _, each := range data {
		total += each.Value
	}
	return total / float64(len(data))
}
