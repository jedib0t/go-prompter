package prompt

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func BenchmarkAutoComplete(b *testing.B) {
	ac := AutoCompleteGoLangKeywords()

	// one-time check to make sure it returns right amount of suggestions
	matches := ac("var foo in", "in", 9)
	assert.Len(b, matches, 6)

	// bench
	for idx := 0; idx < b.N; idx++ {
		ac("var foo in", "in", 9)
	}
}

func TestAutoCompleteGoLangKeywords(t *testing.T) {
	ac := AutoCompleteGoLangKeywords()

	matches := ac("foo", "foo", 0)
	assert.Empty(t, matches)

	matches = ac("ca", "ca", 0)
	assert.Len(t, matches, 2)
	assert.Equal(t, "cap", matches[0].Value)
	assert.Equal(t, "case", matches[1].Value)
	matches = ac("cas", "cas", 0)
	assert.Len(t, matches, 1)
	assert.Equal(t, "case", matches[0].Value)
	matches = ac("case", "case", 0)
	assert.Empty(t, matches)
}

func TestAutoCompletePythonKeywords(t *testing.T) {
	ac := AutoCompletePythonKeywords()

	matches := ac("foo", "foo", 0)
	assert.Empty(t, matches)

	matches = ac("ba", "ba", 0)
	assert.Len(t, matches, 1)
	assert.Equal(t, "basestring", matches[0].Value)
	matches = ac("bas", "bas", 0)
	assert.Len(t, matches, 1)
	assert.Equal(t, "basestring", matches[0].Value)

	matches = ac("Bas", "Bas", 0)
	assert.Len(t, matches, 1)
	assert.Equal(t, "BaseException", matches[0].Value)

	matches = ac("ass", "ass", 0)
	assert.Len(t, matches, 1)
	assert.Equal(t, "assert", matches[0].Value)
	matches = ac("assert", "assert", 0)
	assert.Empty(t, matches)
}

func TestAutoCompleteSQLKeywords(t *testing.T) {
	ac := AutoCompleteSQLKeywords()

	matches := ac("foo", "foo", 0)
	assert.Empty(t, matches)

	matches = ac("se", "sel", 0)
	assert.Len(t, matches, 2)
	assert.Equal(t, "select", matches[0].Value)
	assert.Equal(t, "self", matches[1].Value)
	matches = ac("sele", "sele", 0)
	assert.Len(t, matches, 1)
	assert.Equal(t, "select", matches[0].Value)
	matches = ac("select", "select", 0)
	assert.Empty(t, matches)

	matches = ac("SELECT * fr", "fr", 9)
	assert.Len(t, matches, 3)
	assert.Equal(t, "free", matches[0].Value)
	assert.Equal(t, "freeze", matches[1].Value)
	assert.Equal(t, "from", matches[2].Value)
	matches = ac("SELECT * fro", "fro", 9)
	assert.Len(t, matches, 1)
	assert.Equal(t, "from", matches[0].Value)
	matches = ac("SELECT * from", "from", 9)
	assert.Empty(t, matches)
}

func TestAutoCompleteWords(t *testing.T) {
	t.Run("case insensitive", func(t *testing.T) {
		suggestions := []Suggestion{{Value: "foo"}, {Value: "BAZ"}, {Value: "bar"}}
		ac := AutoCompleteSimple(suggestions, 2, true)

		matches := ac("A Big Croc", "Croc", 6)
		assert.Empty(t, matches)

		matches = ac("fo", "f", 0)
		assert.Empty(t, matches)
		matches = ac("fo", "fo", 0)
		assert.Len(t, matches, 1)
		assert.Equal(t, "foo", matches[0].Value)
		matches = ac("foo", "foo", 0)
		assert.Empty(t, matches)

		matches = ac("foo BA", "BA", 4)
		assert.Len(t, matches, 2)
		assert.Equal(t, "bar", matches[0].Value)
		assert.Equal(t, "baz", matches[1].Value)
	})

	t.Run("case sensitive", func(t *testing.T) {
		suggestions := []Suggestion{{Value: "foo"}, {Value: "BAZ"}, {Value: "bar"}}
		ac := AutoCompleteSimple(suggestions, 2, false)

		matches := ac("A Big Croc", "Croc", 6)
		assert.Empty(t, matches)

		matches = ac("fo", "f", 0)
		assert.Empty(t, matches)
		matches = ac("fo", "fo", 0)
		assert.Len(t, matches, 1)
		assert.Equal(t, "foo", matches[0].Value)
		matches = ac("foo", "foo", 0)
		assert.Empty(t, matches)

		matches = ac("foo ba", "ba", 4)
		assert.Len(t, matches, 1)
		assert.Equal(t, "bar", matches[0].Value)

		matches = ac("foo BA", "BA", 4)
		assert.Len(t, matches, 1)
		assert.Equal(t, "BAZ", matches[0].Value)
	})
}
