package json

import (
	"encoding/json"
	"fmt"
)

// Member struct
type Member struct {
	ID   int64  `json:"ref"`
	Type string `json:"type"`
	Role string `json:"role"`
}

// Print json
func (member Member) Print() {
	json, _ := json.Marshal(member)
	fmt.Println(string(json))
}

// Bytes - return json
func (member Member) Bytes() []byte {
	json, _ := json.Marshal(member)
	return json
}
