package handler

import (
	"fmt"
	"strings"
)

// DeleteTags - delete blacklisted tags from map
func DeleteTags(elemTags map[string]string, blacklist map[string]bool) {
	for key, isPrefix := range blacklist {
		if isPrefix {
			for tag := range elemTags {
				if strings.HasPrefix(tag, key) {
					delete(elemTags, tag)
				}
			}
		} else {
			delete(elemTags, key)
		}
	}
}

func encode(str string) string {
	var encoded = ""
	for i, rune := range str {
		switch true {
		case rune == 10:
			encoded += fmt.Sprintf("%%%x%%", rune)
		case rune == 32:
			encoded += fmt.Sprintf("%%%x%%", rune)
		case rune == 37:
			encoded += fmt.Sprintf("%%%x%%", rune)
		case rune == 44:
			encoded += fmt.Sprintf("%%%x%%", rune)
		case rune == 61:
			encoded += fmt.Sprintf("%%%x%%", rune)
		case rune == 64:
			encoded += fmt.Sprintf("%%%x%%", rune)
		case rune == 127:
			encoded += fmt.Sprintf("%%%x%%", rune)
		case rune == 160:
			encoded += fmt.Sprintf("%%%x%%", rune)
		case rune > 1535:
			encoded += fmt.Sprintf("%%%x%%", rune)
		default:
			encoded += string(str[i])
		}
	}
	return encoded
}
