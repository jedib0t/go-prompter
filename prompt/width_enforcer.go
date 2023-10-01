package prompt

import (
	"strings"
)

// WidthEnforcer is a function that will enforce a "max-length" condition on the
// given text.
type WidthEnforcer func(input string, maxLen int) string

// WidthEnforcerDefault -
func WidthEnforcerDefault(str string, maxLen int) string {
	if maxLen == 0 {
		return str
	}

	sLen := len(str)
	out := strings.Builder{}
	out.Grow(sLen + (sLen / maxLen))
	lineLen, escSeq, inEscSeq := 0, make([]rune, 0), false
	for _, r := range str {
		if lineLen == maxLen {
			if len(escSeq) > 0 {
				out.WriteString(escSeqReset)    // reset before end of line
				out.WriteRune('\n')             // end of line
				out.WriteString(string(escSeq)) // restart on next line
			} else {
				out.WriteRune('\n')
			}
			lineLen = 0
		}

		if r == escSeqStart {
			inEscSeq = true
		}
		if inEscSeq {
			escSeq = append(escSeq, r)
		}

		out.WriteRune(r)
		if !inEscSeq {
			lineLen++
		}

		if r == escSeqStop {
			inEscSeq = false
			if strings.HasSuffix(string(escSeq), escSeqReset) {
				escSeq = make([]rune, 0)
			}
		}
	}
	return out.String()
}
