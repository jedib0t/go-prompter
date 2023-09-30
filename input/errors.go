package input

import "fmt"

// ErrUnsupportedMessage is returned on an attempt to Send an unsupported object
// or message.
var ErrUnsupportedMessage = fmt.Errorf("unsupported message")
