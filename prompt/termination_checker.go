package prompt

import (
	"regexp"
	"strings"
)

// TerminationChecker returns true if the command is terminated and is ready to
// be processed. This is called once the "Terminate" action is invoked by the
// appropriate key sequence - this function can skip the "terminate" action if
// it returns false.
type TerminationChecker func(input string) bool

// TerminationCheckerNone returns true for all inputs and does not abort the
// "terminate" action. Basically an "Enter" key would terminate input, and
// return the accumulated command back to caller.
func TerminationCheckerNone() TerminationChecker {
	return func(input string) bool {
		return true
	}
}

var (
	reSqlComments = regexp.MustCompile(`(/\*.*\*/|--[^\n]*\n|--[^\n]*$)`)
)

// TerminationCheckerSQL returns true if the input is supposed to be a SQL
// statement, and is terminated properly with a semicolon, or is a command
// starting with "/".
func TerminationCheckerSQL() TerminationChecker {
	return func(input string) bool {
		input = reSqlComments.ReplaceAllString(input, "")
		input = strings.TrimSpace(input)

		// SQLs end with a ';'
		if strings.HasSuffix(input, ";") {
			return true
		}
		// SQL command can begin with a '/'
		if strings.HasPrefix(input, "/") {
			return true
		}
		return false
	}
}
