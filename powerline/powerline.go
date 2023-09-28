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

	l, r := 0, 0
	currWidth, handledRight := 0, true
	for {
		if l == len(p.left) && r == len(p.right) { // appended everything
			break
		}
		if handledRight { // append a left segment
			if l < len(p.left) {
				addlWidth := p.left[l].Width()
				if currWidth+addlWidth > maxWidth {
					break
				}
				currWidth += addlWidth
				l++
			}
			handledRight = false
		} else { // append a right segment
			if r < len(p.right) {
				addlWidth := p.right[r].Width()
				if currWidth+addlWidth > maxWidth {
					break
				}
				currWidth += addlWidth
				r++
			}
			handledRight = true
		}
	}
	return l, r, maxWidth - currWidth
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
		if idx == 0 {
			segment.setPaddingLeft(p.style.MarginLeft)
		}

		// segment
		left = append(left, segment.Render())

		// separator
		if len(sep) > 0 {
			sepColor := prompt.Color{
				Background: p.style.Color.Background,
				Foreground: segment.Color().Background,
			}
			if idx < len(segments)-1 {
				sepColor.Background = segments[idx+1].Color().Background
			} else if idx == len(segments)-1 {
				sepColor.Background = p.style.Color.Background
			}
			if p.style.InvertSeparatorColors {
				sepColor = sepColor.Invert()
			}
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
			sepColor := prompt.Color{
				Foreground: segment.Color().Background,
				Background: p.style.Color.Background,
			}
			if idx > 0 {
				sepColor.Background = segments[idx-1].Color().Background
			}
			if p.style.InvertSeparatorColors {
				sepColor = sepColor.Invert()
			}
			right = append(right, sepColor.Sprint(sep))
		}

		// set up the margin
		if idx == len(segments)-1 {
			segment.setPaddingLeft(p.style.MarginLeft)
		}

		// segment
		right = append(right, segment.Render())
	}

	p.rightRendered = strings.Join(right, "")
	return p.rightRendered
}
