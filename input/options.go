package input

import "io"

// Option tells the Reader how to behave.
type Option func(r *reader)

// WithInput makes Reader listen for events on the given source. Use this if you
// want to use something other than os.Stdin as the source of events.
func WithInput(input io.Reader) Option {
	return func(r *reader) {
		r.input = input
	}
}

// WatchMouseAll makes Reader watch and pipe Mouse click and move events.
func WatchMouseAll() Option {
	return func(r *reader) {
		r.watchMouseAll = true
	}
}

// WatchMouseClick makes Reader watch and pipe Mouse click events.
func WatchMouseClick() Option {
	return func(r *reader) {
		r.watchMouseClick = true
	}
}

// WatchWindowSize makes Reader watch and pipe Window size/resize events.
func WatchWindowSize() Option {
	return func(r *reader) {
		r.watchWindowSize = true
	}
}
