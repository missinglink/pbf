package json

import (
	"encoding/json"
	"fmt"
)

// LatLon struct
type LatLon struct {
	Lat float64 `json:"lat"`
	Lon float64 `json:"lon"`
}

// Print json
func (ll LatLon) Print() {
	json, _ := json.Marshal(ll)
	fmt.Println(string(json))
}

// Bytes - return json
func (ll LatLon) Bytes() []byte {
	json, _ := json.Marshal(ll)
	return json
}

// NewLatLon - generate a new JSON struct based off a parse struct
func NewLatLon(lat float64, lon float64) *LatLon {
	return &LatLon{
		Lat: truncate(lat),
		Lon: truncate(lon),
	}
}
