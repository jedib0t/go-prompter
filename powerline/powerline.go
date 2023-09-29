package powerline

import (
	"fmt"
	"os"
	"strings"
	"sync"

	"github.com/jedib0t/go-pretty/v6/text"
	"github.com/jedib0t/go-prompter/prompt"
	"golang.org/x/term"
)

// Powerline helps construct a powerline like prompt with segmented contents.
type Powerline struct {
	autoAdjustWidth  bool
	hasChanges       bool
	left             []*Segment
	leftRendered     string
	right            []*Segment
	rightRendered    string
	renderedMaxWidth int
	style            *Style
	mutex            sync.Mutex
}

// Append appends the given segment to the Left side of the prompt.
func (p *Powerline) Append(s *Segment) {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	p.hasChanges = true
	p.left = append(p.left, s)
}

// AppendRight appends the given segment to the Right side of the prompt.
func (p *Powerline) AppendRight(s *Segment) {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	p.hasChanges = true
	p.right = append(p.right, s)
}

// AutoAdjustWidth turns on automatically reducing number of segments to fit the
// provided width.
func (p *Powerline) AutoAdjustWidth(v bool) {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	if v != p.autoAdjustWidth {
		p.autoAdjustWidth = v
		p.hasChanges = true
	}
}

// Render renders the Powerline prompt with all the segments.
func (p *Powerline) Render(maxWidth int) string {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	defer func() {
		p.hasChanges = false
	}()

	if p.style == nil {
		p.style = &StyleDefault
	}
	if maxWidth == 0 {
		maxWidth, _, _ = term.GetSize(int(os.Stdout.Fd()))
	}
	defer func() {
		p.renderedMaxWidth = maxWidth
	}()
	nsLeft, nsRight, paddingSpace := p.autoAdjustNumSegments(maxWidth)

	left := p.renderLeft(maxWidth, nsLeft, nsRight, paddingSpace)
	if len(p.right) > 0 {
		leftLen := text.RuneWidthWithoutEscSequences(left)
		right := p.renderRight(maxWidth, nsLeft, nsRight, paddingSpace)
		rightLen := text.RuneWidthWithoutEscSequences(right)
		padding := p.renderPadding(maxWidth - (leftLen + rightLen))
		return fmt.Sprintf("%s%s%s", left, padding, right)
	}
	return left
}

func (p *Powerline) SetStyle(style Style) {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	p.hasChanges = true
	p.style = &style
}

func (p *Powerline) autoAdjustNumSegments(maxWidth int) (int, int, int) {
	if !p.autoAdjustWidth {
		return len(p.left), len(p.right), 0
	}

	idxLeft, idxRight := 0, 0
	isDone, usedWidth := false, 0
	for idx := 0; !isDone && idx < 1024; /* 512 segments per side */ idx++ {
		if idxLeft == len(p.left) && idxRight == len(p.right) { // appended everything
			break
		}
		if idx%2 == 0 { // append a left segment
			idxLeft, usedWidth, isDone = appendSegmentIfUnderWidth(p.left, idxLeft, usedWidth, maxWidth)
		} else { // append a right segment
			idxRight, usedWidth, isDone = appendSegmentIfUnderWidth(p.right, idxRight, usedWidth, maxWidth)
		}
	}
	return idxLeft, idxRight, maxWidth - usedWidth
}

func (p *Powerline) hasChangesLeft() bool {
	if p.hasChanges {
		return true
	}
	for _, segment := range p.left {
		if segment.HasChanges() {
			return true
		}
	}
	return false
}

func (p *Powerline) hasChangesRight() bool {
	if p.hasChanges {
		return true
	}
	for _, segment := range p.right {
		if segment.HasChanges() {
			return true
		}
	}
	return false
}

func (p *Powerline) renderLeft(maxWidth, nsLeft, nsRight, paddingSpace int) string {
	if p.leftRendered != "" && maxWidth == p.renderedMaxWidth && !p.hasChangesLeft() {
		return p.leftRendered
	}
	var left []string

	// restrict the segments that get rendered
	segments := p.left
	if nsLeft <= len(segments) {
		segments = segments[:nsLeft]
	}

	// inject all the segments now
	sep := p.style.SeparatorLeft
	for idx, segment := range segments {
		// segment: set up margin and append
		if idx == 0 {
			segment.setPaddingLeft(p.style.MarginLeft)
		}
		left = append(left, segment.Render())

		// separator
		if len(sep) > 0 {
			sepColor := p.separatorColor(segments, idx, true)
			left = append(left, sepColor.Sprint(sep))
		}
	}

	p.leftRendered = strings.Join(left, "")
	return p.leftRendered
}

func (p *Powerline) renderPadding(length int) string {
	if length < 0 {
		length = 0
	}
	return p.style.Color.Sprint(strings.Repeat(" ", length))
}

func (p *Powerline) renderRight(maxWidth, nsLeft, nsRight, paddingSpace int) string {
	if p.leftRendered != "" && maxWidth == p.renderedMaxWidth && !p.hasChangesRight() {
		return p.rightRendered
	}
	var right []string

	// restrict the segments that get rendered
	segments := p.right
	if nsRight <= len(segments) {
		segments = segments[len(segments)-nsRight:]
	}

	// inject all the segments now
	sep := p.style.SeparatorRight
	for idx, segment := range segments {
		// separator
		if len(sep) > 0 {
			sepColor := p.separatorColor(segments, idx, false)
			right = append(right, sepColor.Sprint(sep))
		}

		// segment: set up the margin and append
		if idx == len(segments)-1 {
			segment.setPaddingLeft(p.style.MarginLeft)
		}
		right = append(right, segment.Render())
	}

	p.rightRendered = strings.Join(right, "")
	return p.rightRendered
}

func appendSegmentIfUnderWidth(segments []*Segment, idx int, usedWidth int, maxWidth int) (int, int, bool) {
	if idx < len(segments) {
		addlWidth := segments[idx].Width()
		if usedWidth+addlWidth > maxWidth {
			return idx, usedWidth, true
		}
		usedWidth += addlWidth
		idx++
	}
	return idx, usedWidth, false
}

func (p *Powerline) separatorColor(segments []*Segment, idx int, leftSide bool) prompt.Color {
	c := prompt.Color{
		Background: p.style.Color.Background,
		Foreground: segments[idx].Color().Background,
	}

	if leftSide {
		if idx+1 < len(segments) {
			c.Background = segments[idx+1].Color().Background
		} else if idx == len(segments)-1 { // last segment, use powerline's background
			c.Background = p.style.Color.Background
		}
	} else { // right side
		if idx > 0 {
			c.Background = segments[idx-1].Color().Background
		}
	}

	if p.style.InvertSeparatorColors {
		c = c.Invert()
	}

	return c
}
