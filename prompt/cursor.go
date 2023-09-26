package prompt

import "fmt"

// CursorLocation contains the current cursor position in a 2d-wall-of-text; the
// values are 0-indexed to keep it simple to manipulate the wall of text
type CursorLocation struct {
	Line   int `json:"line"`
	Column int `json:"column"`
}

// String returns the co-ordinates in human-readable format.
func (cl CursorLocation) String() string {
	return fmt.Sprintf("[%d, %d]", cl.Line+1, cl.Column+1)
}
