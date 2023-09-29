package prompt

import (
	_ "embed"
	"fmt"
	"sort"
	"strings"
)

// Suggestion is what is returned by the auto-completer.
type Suggestion struct {
	Value string
	Hint  string
}

// String appeases the stringer interface.
func (s Suggestion) String() string {
	return fmt.Sprintf("%s: %s", s.Value, s.Hint)
}

// AutoCompleter defines a function that takes the entire user input, the word
// the user is specifically on, and the location of the cursor on the entire
// sentence. It is expected to return zero or more strings that match what the
// user may type.
type AutoCompleter func(sentence string, word string, location uint) []Suggestion

//go:embed suggestions/golang.txt
var suggestionsFileGoLang string

// AutoCompleteGoLangKeywords is a simple auto-completer that helps
// auto-complete most of the known GoLang suggestions.
func AutoCompleteGoLangKeywords() AutoCompleter {
	return AutoCompleteSimple(suggestionsFromFile(suggestionsFileGoLang), false)
}

//go:embed suggestions/python.txt
var suggestionsFilePython string

// AutoCompletePythonKeywords is a simple auto-completer that helps
// auto-complete most of the known Python suggestions.
func AutoCompletePythonKeywords() AutoCompleter {
	return AutoCompleteSimple(suggestionsFromFile(suggestionsFilePython), false)
}

//go:embed suggestions/sql.txt
var suggestionsFileSQL string

// AutoCompleteSQLKeywords is a simple auto-completer that helps
// auto-complete most of the known SQL suggestions.
func AutoCompleteSQLKeywords() AutoCompleter {
	return AutoCompleteSimple(suggestionsFromFile(suggestionsFileSQL), true)
}

// AutoCompleteSimple returns an AutoCompleter which will use the given list of
// suggestions in an optimized fashion for look-ups.
func AutoCompleteSimple(suggestions []Suggestion, caseInsensitive bool) AutoCompleter {
	// sort ahead and avoid sorting while searching in the loop below
	if caseInsensitive {
		for idx := range suggestions {
			suggestions[idx].Value = strings.ToLower(suggestions[idx].Value)
		}
	}
	sort.SliceStable(suggestions, func(i, j int) bool {
		return suggestions[i].Value < suggestions[j].Value
	})

	// build a map of the first character to the list of words to make look-up
	// reasonably fast
	lookupMap, firstRunes := processSuggestionsForLookup(suggestions)

	// return the auto-completer
	return func(sentence string, word string, location uint) []Suggestion {
		var matches []Suggestion
		if word == "" { // recommend everything
			for _, firstRune := range firstRunes {
				matches = append(matches, lookupMap[firstRune]...)
			}
		} else {
			if caseInsensitive {
				word = strings.ToLower(word)
			}
			for _, possibleMatch := range lookupMap[word[0:1]] {
				if strings.HasPrefix(possibleMatch.Value, word) && len(possibleMatch.Value) >= len(word) {
					matches = append(matches, possibleMatch)
				}
			}
		}
		return matches
	}
}

func processSuggestionsForLookup(suggestions []Suggestion) (map[string][]Suggestion, []string) {
	lookupMap := make(map[string][]Suggestion)
	possibleMatchesFirstRuneMap := make(map[string]bool)
	for _, suggestion := range suggestions {
		firstRune := suggestion.Value[0:1]
		lookupMap[firstRune] = append(lookupMap[firstRune], suggestion)
		possibleMatchesFirstRuneMap[firstRune] = true
	}

	firstRunes := make([]string, 0)
	for k := range possibleMatchesFirstRuneMap {
		firstRunes = append(firstRunes, k)
	}
	sort.Strings(firstRunes)

	return lookupMap, firstRunes
}

func suggestionsFromFile(contents string) []Suggestion {
	var suggestions []Suggestion
	for _, line := range strings.Split(contents, "\n") {
		line = strings.TrimSpace(line)
		if line != "" {
			tokens := strings.SplitN(line, ":", 2)
			value, hint := strings.TrimSpace(tokens[0]), ""
			if len(tokens) > 1 {
				hint = strings.TrimSpace(tokens[1])
			}
			suggestions = append(suggestions, Suggestion{Value: value, Hint: hint})
		}
	}

	sort.SliceStable(suggestions, func(i, j int) bool {
		return suggestions[i].Value < suggestions[j].Value
	})

	return suggestions
}
