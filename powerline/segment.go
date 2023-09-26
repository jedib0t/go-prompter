package powerline

import (
	"fmt"
	"hash/crc32"
	"strings"
	"sync"

	"github.com/jedib0t/go-pretty/v6/text"
	"github.com/jedib0t/go-prompter/prompt"
	"github.com/muesli/termenv"
)

// Segment contains the contents for a "segment" of data in the Powerline
// prompt.
type Segment struct {
	color        *prompt.Color
	content      string
	contentColor *prompt.Color
	hasChanges   bool
	icon         string
	mutex        sync.Mutex
	paddingLeft  *string
	paddingRight *string
	rendered     string
	width        int
}

// Color returns the color that will be used for the segment.
func (s *Segment) Color() prompt.Color {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if s.color != nil {
		return *s.color
	}
	if s.contentColor != nil {
		return *s.contentColor
	}
	return prompt.Color{
		Foreground: termenv.ANSI256Color(7),
		Background: termenv.ANSI256Color(16),
	}
}

// HasChanges returns true if Render is going to return a different result
// compared to its last invocation.
func (s *Segment) HasChanges() bool {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	return s.hasChanges
}

// SetIcon sets the optional Icon/Emoji to be rendered before the text.
func (s *Segment) SetIcon(icon string) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	if s.icon == icon {
		return
	}
	defer func() {
		s.width = s.calculateWidth()
	}()

	s.hasChanges = true
	s.icon = icon
}

// SetColor sets the colors to be used for the segment. If not set, the hash of
// the content determines the colors.
func (s *Segment) SetColor(color prompt.Color) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	if s.color != nil && *s.color == color {
		return
	}

	s.hasChanges = true
	s.color = &color
}

// SetContent sets the text of the segment.
//
// Normally, the client can set the color for the content using SetColor.
// However, in case the client doesn't do it, the color is determined
// automatically by hashing one of the following:
//   - the "tags" values
//   - the "content" value if no "tags" provided
func (s *Segment) SetContent(content string, tags ...string) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	if s.content == content {
		return // if there are no changes, ignore request
	}
	defer func() {
		s.width = s.calculateWidth()
	}()

	s.hasChanges = true
	s.content = content

	// determine the content color based on hash value of content
	if s.color == nil {
		h := crc32.NewIEEE()
		if len(tags) > 0 {
			_, _ = h.Write([]byte(fmt.Sprint(tags)))
		} else {
			_, _ = h.Write([]byte(s.content))
		}
		hash := h.Sum32()
		bg := (hash % (231 - 16)) + 16
		fg := 16 // black
		if (bg-16)%36 < 18 {
			fg = 7 // white
		}
		s.contentColor = &prompt.Color{
			Foreground: termenv.ANSI256Color(fg),
			Background: termenv.ANSI256Color(bg),
		}
	}
}

// Render returns the segment rendered in appropriate colors.
func (s *Segment) Render() string {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if s.content == "" {
		return ""
	}

	if s.hasChanges {
		color := s.color
		if color == nil {
			color = s.contentColor
		}

		out := strings.Builder{}
		if s.paddingLeft != nil {
			out.WriteString(*s.paddingLeft)
		} else {
			out.WriteRune(' ')
		}
		if s.icon != "" {
			out.WriteString(s.icon)
			out.WriteRune(' ')
		}
		out.WriteString(s.content)
		if s.paddingRight != nil {
			out.WriteString(*s.paddingRight)
		} else {
			out.WriteRune(' ')
		}

		s.rendered = color.Sprint(out.String())
	}
	s.hasChanges = false
	return s.rendered
}

// ResetColor resets the color of the content to defaults.
func (s *Segment) ResetColor() {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	s.color = nil
}

// Width returns the width of the segment when printed on screen.
func (s *Segment) Width() int {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	return s.width
}

func (s *Segment) calculateWidth() int {
	w := 1 + // separator
		1 + // left margin
		text.RuneWidthWithoutEscSequences(s.content) +
		1 // right margin

	if s.icon != "" {
		w += 1 + // left margin
			text.RuneWidthWithoutEscSequences(s.icon)
	}
	return w
}

func (s *Segment) setPaddingLeft(p string) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	s.hasChanges = true
	s.paddingLeft = &p
}

func (s *Segment) setPaddingRight(p string) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	s.hasChanges = true
	s.paddingRight = &p
}
