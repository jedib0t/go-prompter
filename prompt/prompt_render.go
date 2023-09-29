package prompt

import (
	"context"
	"fmt"
	"time"

	"github.com/jedib0t/go-pretty/v6/text"
	"github.com/muesli/termenv"
)

//gocyclo:ignore
func (p *prompt) render(ctx context.Context, output *termenv.Output) (rsp string, err error) {
	p.init(ctx)

	// set up cleanup
	defer func() {
		p.pauseRender()
		p.updateModel(false)
		p.renderView(output, "done", true)
		p.buffer.Reset()
	}()

	// instantiate an input reader and begin looking for inputs
	p.readerMutex.Lock()
	go p.reader.Begin(ctx)
	defer func() {
		p.reader.End()
		p.readerMutex.Unlock()
	}()

	// first time render
	p.updateModel(true)

	// start handling input events and rendering to screen
	tick := time.Tick(p.refreshInterval)
	tickCursor := time.Tick(p.style.Cursor.BlinkInterval)
	for {
		select {
		case <-ctx.Done():
			return "", ctx.Err()
		case <-tick:
			if p.buffer.HasChanges() {
				p.updateModel(true)
			}
			p.renderView(output, "tick")
		case <-tickCursor:
			p.updateModel(true)
		case err = <-p.reader.Errors():
			return "", err
		case key := <-p.reader.KeyEvents():
			if err = p.handleKey(output, key); err != nil {
				return "", err
			}
			if p.buffer.IsDone() {
				return p.buffer.String(), nil
			}
		case resize := <-p.reader.WindowSizeEvents():
			p.updateDisplayWidth(resize.Width)
		}
	}
}

func (p *prompt) renderView(output *termenv.Output, reason string, forced ...bool) {
	p.linesMutex.Lock()
	defer p.linesMutex.Unlock()

	// if paused, don't do anything
	if p.isRenderPaused() && len(forced) == 0 {
		return
	}
	p.setDebugData("reason", reason)

	timeStart := time.Now()
	defer func() {
		p.linesRendered = p.linesToRender
	}()

	// calculate movement
	numLinesToRender, numLinesRendered := len(p.linesToRender), len(p.linesRendered)

	// move cursor up and clear printed lines
	if numLinesRendered > 0 {
		if numLinesToRender < numLinesRendered {
			for idx := numLinesRendered; idx > numLinesToRender; idx-- {
				output.CursorUp(1)
				output.ClearLine()
			}
			output.CursorUp(numLinesToRender)
		} else {
			output.CursorUp(numLinesRendered)
		}
		if p.debug { // for the final debug footer
			output.CursorUp(1)
		}
	}

	// print all the changed lines
	for idx, line := range p.linesToRender {
		if idx < len(p.linesRendered) && p.linesRendered[idx] == line { // nothing changed
			output.CursorDown(1)
			continue
		}

		// something is different, clear and reprint
		output.ClearLine()
		if p.debug { // render the "second" this line was rendered to screen
			_, _ = output.WriteString(p.style.Colors.Debug.Sprintf(" %02d ", time.Now().Second()))
		}
		_, _ = output.WriteString(fmt.Sprintf("%s\n", line))
	}

	if p.debug {
		stats := fmt.Sprintf("%s; time=%v [gen=sh:%v/bf:%v/ac:%v/%v]",
			p.debugDataAsString(), time.Since(timeStart).Round(time.Microsecond),
			p.timeSyntaxGen, p.timeBufferGen, p.timeAutoComplete, p.timeGen,
		)

		output.ClearLine()
		_, _ = output.WriteString(p.style.Colors.Debug.Sprintf(" %02d ", time.Now().Second()))
		_, _ = output.WriteString(p.style.Colors.Debug.Sprintf(
			text.AlignCenter.Apply(stats, p.getDisplayWidth())) + "\n",
		)
	}
}
