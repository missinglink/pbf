package json

import (
	"encoding/json"
	"fmt"
)

// Coords struct
type Coords struct {
	ID   int64   `json:"id"`
	Type string  `json:"type"`
	Lat  float64 `json:"lat"`
	Lon  float64 `json:"lon"`
}

// Print json
func (coords Coords) Print() {
	json, _ := json.Marshal(coords)
	fmt.Println(string(json))
}

// Bytes - return json
func (coords Coords) Bytes() []byte {
	json, _ := json.Marshal(coords)
	return json
}
