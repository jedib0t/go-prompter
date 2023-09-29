package prompt

import "errors"

// ErrAborted is returned when the execution is terminated for any unexpected
// reason
var ErrAborted = errors.New("aborted")

// ErrDuplicateKeyAssignment is returned when the key map has the same keys
// defined for multiple incompatible actions.
var ErrDuplicateKeyAssignment = errors.New("duplicate key assignment")

// ErrInvalidDimensions is returned when the style sheet has dimensions that
// does not make sense.
var ErrInvalidDimensions = errors.New("invalid dimensions")

// ErrNonInteractiveShell is returned when the Prompter is being invoked on a
// shell which is not-interactive (ex.: script run in headless mode).
var ErrNonInteractiveShell = errors.New("non-interactive shell")

// ErrUnsupportedChromaLanguage is returned when Syntax-Highlighting is
// requested with Chroma library with a language that it does not understand.
var ErrUnsupportedChromaLanguage = errors.New("unsupported language for chroma")

// ErrUnsupportedInput is returned when the input given to a function is not
// supported or handled.
var ErrUnsupportedInput = errors.New("unsupported input")
