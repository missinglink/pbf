package lib

import (
	"log"

	"github.com/missinglink/gosmparse"
)

// MemberType - map numeric value to string
func MemberType(t gosmparse.MemberType) string {
	switch t {
	case gosmparse.NodeType:
		return "node"
	case gosmparse.WayType:
		return "way"
	case gosmparse.RelationType:
		return "relation"
	default:
		log.Println("unknown member type", t)
		return "node"
	}
}
