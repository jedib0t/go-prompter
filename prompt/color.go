package prompt

import (
	"fmt"
	"strings"
	"sync"

	"github.com/muesli/termenv"
)

const (
	escSeqStart = '\x1b'
	escSeqSep   = ';'
	escSeqStop  = 'm'
	escSeqReset = string(escSeqStart) + "[0" + string(escSeqStop)
)

// Color contain the foreground and background colors to use to format text.
type Color struct {
	Foreground termenv.Color
	Background termenv.Color

	escSeq string
}

var (
	escSeqCache      = make(map[string]string)
	escSeqCacheMutex = sync.RWMutex{}
)

// Invert flips the background and foreground colors and returns a new Color
// object.
func (c Color) Invert() Color {
	return Color{
		Foreground: c.Background,
		Background: c.Foreground,
	}
}

// Sprint behaves like fmt.Sprint but with color sequence wrapping the output.
func (c Color) Sprint(a ...any) string {
	if c.Foreground == nil {
		c.Foreground = termenv.ForegroundColor()
	}
	if c.Background == nil {
		c.Background = termenv.BackgroundColor()
	}
	cacheKey := fmt.Sprintf("%s%s", c.Foreground, c.Background)

	// find value in cache
	escSeqCacheMutex.RLock()
	cacheVal, ok := escSeqCache[cacheKey]
	escSeqCacheMutex.RUnlock()

	// if none found, generate and store in cache
	if !ok {
		sb := strings.Builder{}
		sb.WriteRune(escSeqStart)
		sb.WriteRune('[')
		sb.WriteString(c.Foreground.Sequence(false))
		sb.WriteRune(escSeqSep)
		sb.WriteString(c.Background.Sequence(true))
		sb.WriteRune(escSeqStop)
		if sb.String() == "\x1b[;m" {
			cacheVal = ""
		} else {
			cacheVal = sb.String()
		}

		escSeqCacheMutex.Lock()
		escSeqCache[cacheKey] = cacheVal
		escSeqCacheMutex.Unlock()
	}

	// if no escape sequence was generated, just return input
	out := fmt.Sprint(a...)
	if len(out) == 0 || len(cacheVal) == 0 {
		return out
	}

	// generate the colored string
	sb := strings.Builder{}
	sb.WriteString(cacheVal)
	sb.WriteString(out)
	sb.WriteString(escSeqReset)
	return sb.String()
}

// Sprintf behaves like fmt.Sprintf but with color sequence wrapping the output.
func (c Color) Sprintf(msg string, a ...any) string {
	return c.Sprint(fmt.Sprintf(msg, a...))
}

// Ref.: https://talyian.github.io/ansicolors/
