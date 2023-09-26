package prompt

import (
	"fmt"
	"time"
)

const (
	DefaultPrefix = "> "
)

// Prefixer returns the string to precede any new prompt.
type Prefixer func() string

// PrefixNone uses no prompt prefix.
func PrefixNone() Prefixer {
	return PrefixText("")
}

// PrefixSimple uses "> " as the prompt prefix.
func PrefixSimple() Prefixer {
	return PrefixText(DefaultPrefix)
}

// PrefixText uses the given text as the prompt prefix.
func PrefixText(text string) Prefixer {
	return func() string {
		return text
	}
}

// PrefixTimestamp uses a timestamp and a prefix as the prompt prefix.
func PrefixTimestamp(timeLayout string, prefix string) Prefixer {
	return func() string {
		return fmt.Sprintf("%s %s", time.Now().Format(timeLayout), prefix)
	}
}
