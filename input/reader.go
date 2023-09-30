package input

import (
	"context"
	"fmt"
	"io"
	"sync"

	tea "github.com/charmbracelet/bubbletea"
)

// Reader channels Key, Mouse and Window Resize events to the caller through the
// publicly defined channels below. The caller needs to drain all the channels
// to prevent blocking.
//
// It uses BubbleTea framework underneath the hood to make this happen, but
// provides a cleaner interface to just input handling.
type Reader interface {
	// Begin begins reading inputs.
	Begin(ctx context.Context)
	// End terminates reading.
	End()
	// Errors returns the channel with errors encountered while reading.
	Errors() <-chan error
	// KeyEvents returns the channel with events from the Keyboard.
	KeyEvents() <-chan tea.KeyMsg
	// MouseEvents returns the channel with events from the Mouse.
	MouseEvents() <-chan tea.MouseMsg
	// Send sends the message over the appropriate channel back to client. This
	// can be useful for testing and automated inputs.
	Send(msg any) error
	// WindowSizeEvents returns the channel with resize events from the Terminal
	// window.
	WindowSizeEvents() <-chan tea.WindowSizeMsg
	// Reset resets everything and prepares for new input, and will block if
	// called when reading is in progress and has not been End-ed yet.
	Reset() error
}

// NewReader returns a new Reader with the provided options applied on it.
func NewReader(opts ...Option) Reader {
	r := &reader{}
	for _, opt := range opts {
		opt(r)
	}
	r.init()
	return r
}

type reader struct {
	chDone             chan bool
	chErrors           chan error
	chKeyEvents        chan tea.KeyMsg
	chMouseEvents      chan tea.MouseMsg
	chWindowSizeEvents chan tea.WindowSizeMsg
	done               bool
	input              io.Reader
	mutex              sync.Mutex
	program            *tea.Program
	programMutex       sync.Mutex
	teaBag             *teaBag
	watchMouseAll      bool
	watchMouseClick    bool
	watchWindowSize    bool
}

// Begin sets up the channels for reading and begins capturing inputs.
func (r *reader) Begin(ctx context.Context) {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	r.done = false

	r.programMutex.Lock()
	r.program = tea.NewProgram(r.teaBag, r.progOpts(ctx)...)
	r.programMutex.Unlock()

	_, err := r.program.Run()
	if err != nil {
		r.chErrors <- err
	}
	r.chDone <- true
}

// End stops the input handling and cleans up.
func (r *reader) End() {
	r.programMutex.Lock()
	defer r.programMutex.Unlock()

	if r.program != nil {
		r.program.Quit()
		r.program = nil
		<-r.chDone
	}
}

// Errors returns the channel that passes on error events.
func (r *reader) Errors() <-chan error {
	return r.chErrors
}

// KeyEvents returns the channel that passes on keyboard events.
func (r *reader) KeyEvents() <-chan tea.KeyMsg {
	return r.chKeyEvents
}

// MouseEvents returns the channel that passes on mouse events.
func (r *reader) MouseEvents() <-chan tea.MouseMsg {
	return r.chMouseEvents
}

// Send sends the message over the appropriate channel back to client. This can
// be useful for testing and automated inputs.
func (r *reader) Send(msg any) error {
	r.programMutex.Lock()
	defer r.programMutex.Unlock()

	switch obj := msg.(type) {
	case error:
		r.chErrors <- obj
	case tea.KeyMsg:
		r.chKeyEvents <- obj
	case tea.MouseMsg:
		r.chMouseEvents <- obj
	case tea.WindowSizeMsg:
		r.chWindowSizeEvents <- obj
	default:
		return fmt.Errorf("%w: %#v", ErrUnsupportedMessage, msg)
	}
	return nil
}

// WindowSizeEvents returns the channel that passes on window size events.
func (r *reader) WindowSizeEvents() <-chan tea.WindowSizeMsg {
	return r.chWindowSizeEvents
}

// Reset resets the state as long as the reading has not begun. Works before
// calling Begin or after calling End.
func (r *reader) Reset() error {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	r.programMutex.Lock()
	defer r.programMutex.Unlock()

	r.init()
	return nil
}

func (r *reader) init() {
	r.chDone = make(chan bool, 1)
	r.chErrors = make(chan error, 5)
	r.chKeyEvents = make(chan tea.KeyMsg, 5)
	r.chMouseEvents = make(chan tea.MouseMsg, 5)
	r.chWindowSizeEvents = make(chan tea.WindowSizeMsg, 5)
	r.teaBag = &teaBag{
		ErrorEvents:     r.chErrors,
		KeyEvents:       r.chKeyEvents,
		MouseEvents:     r.chMouseEvents,
		ResizeEvents:    r.chWindowSizeEvents,
		watchMouse:      r.watchMouseAll || r.watchMouseClick,
		watchWindowSize: r.watchWindowSize,
	}
}

func (r *reader) progOpts(ctx context.Context) []tea.ProgramOption {
	opts := []tea.ProgramOption{
		tea.WithContext(ctx),
	}
	if r.input != nil {
		opts = append(opts, tea.WithInput(r.input))
	}
	if r.watchMouseAll {
		opts = append(opts, tea.WithMouseAllMotion())
	} else if r.watchMouseClick {
		opts = append(opts, tea.WithMouseCellMotion())
	}
	return opts
}

// teaBag wraps a bubbletea model. Get it?
type teaBag struct {
	ErrorEvents     chan error
	KeyEvents       chan tea.KeyMsg
	MouseEvents     chan tea.MouseMsg
	ResizeEvents    chan tea.WindowSizeMsg
	watchMouse      bool
	watchWindowSize bool
}

func (tb *teaBag) Init() tea.Cmd {
	return nil
}

func (tb *teaBag) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case error:
		tb.ErrorEvents <- msg
	case tea.KeyMsg:
		tb.KeyEvents <- msg
	case tea.MouseMsg:
		if tb.watchMouse {
			tb.MouseEvents <- msg
		}
	case tea.WindowSizeMsg:
		if tb.watchWindowSize {
			tb.ResizeEvents <- msg
		}
	}
	return tb, nil
}

func (tb *teaBag) View() string {
	return ""
}
