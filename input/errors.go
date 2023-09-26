package input

import "fmt"

// ErrDisallowedAction is returned when an action is attempted that is not
// allowed at this time.
var ErrDisallowedAction = fmt.Errorf("action not allowed at this time")
