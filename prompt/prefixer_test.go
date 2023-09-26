package prompt

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestPrefixNone(t *testing.T) {
	prefixer := PrefixNone()

	assert.Equal(t, "", prefixer())
}

func TestPrefixSimple(t *testing.T) {
	prefixer := PrefixSimple()

	assert.Equal(t, "> ", prefixer())
}

func TestPrefixText(t *testing.T) {
	prefixer := PrefixText("foo> ")

	assert.Equal(t, "foo> ", prefixer())
}

func TestPrefixTimestamp(t *testing.T) {
	prefixer := PrefixTimestamp(time.DateTime, "> ")
	timestamp := time.Now().Format(time.DateTime)

	assert.Equal(t, fmt.Sprintf("%s %s", timestamp, "> "), prefixer())
}
