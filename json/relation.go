package json

import (
	"encoding/json"
	"fmt"

	"github.com/missinglink/gosmparse"
)

// Relation struct
type Relation struct {
	ID      int64             `json:"id"`
	Type    string            `json:"type"`
	Hash    string            `json:"hash,omitempty"`
	Tags    map[string]string `json:"tags,omitempty"`
	Members []Member          `json:"members"`
}

// Print json
func (relation Relation) Print() {
	json, _ := json.Marshal(relation)
	fmt.Println(string(json))
}

// Bytes - return json
func (relation Relation) Bytes() []byte {
	json, _ := json.Marshal(relation)
	return json
}

// RelationFromParser - generate a new JSON struct based off a parse struct
func RelationFromParser(item gosmparse.Relation) *Relation {

	// members
	members := make([]Member, 0, len(item.Members))
	for _, member := range item.Members {

		// detect type
		var typ = "node"
		if member.Type == gosmparse.WayType {
			typ = "way"
		} else if member.Type == gosmparse.RelationType {
			typ = "relation"
		}

		members = append(members, Member{
			ID:   member.ID,
			Type: typ,
			Role: member.Role,
		})
	}

	return &Relation{
		ID:      item.ID,
		Type:    "relation",
		Tags:    item.Tags,
		Members: members,
	}
}
