package tags

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTrim(t *testing.T) {

	var tagMap = make(map[string]string)
	tagMap["  A"] = "a"
	tagMap["B  "] = "b"
	tagMap["C"] = "  c"
	tagMap["D"] = "d  "

	var trimmed = Trim(tagMap)
	assert.Equal(t, 4, len(trimmed))
	assert.Equal(t, trimmed["A"], "a")
	assert.Equal(t, trimmed["B"], "b")
	assert.Equal(t, trimmed["C"], "c")
	assert.Equal(t, trimmed["D"], "d")
}
