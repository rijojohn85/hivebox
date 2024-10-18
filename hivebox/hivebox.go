package hivebox

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

const version = "0.0.1"

func GetVersion(w http.ResponseWriter, req *http.Request) {
	json.NewEncoder(w).Encode(version)
}

func CreateUrl() string {
	tenHoursAgo := time.Now().Add(-10 * time.Hour)
	formattedTime := tenHoursAgo.Format(time.RFC3339)
	url := fmt.Sprintf("https://api.opensensemap.org/boxes/data?phenomenon=temperature&bbox=-180,-90,180,90&date=%s&format=json", formattedTime)
	return url
}
