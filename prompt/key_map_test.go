package prompt

import (
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestKeyMap_reverse(t *testing.T) {
	k := KeyMapDefault
	kr, err := k.reverse()
	assert.NotNil(t, kr)
	assert.Nil(t, err)
	if kr != nil {
		assert.NotEmpty(t, kr.AutoComplete)
		assert.NotEmpty(t, kr.Insert)
	}

	k.Insert.MoveToEndOfLine = k.Insert.Abort
	kr, err = k.reverse()
	assert.Nil(t, kr)
	assert.NotNil(t, err)
	assert.True(t, errors.Is(err, ErrDuplicateKeyAssignment))
	assert.Contains(t, err.Error(), "- more than one action defined for 'ctrl+c': [Abort, MoveToEndOfLine]")
	assert.Contains(t, err.Error(), "- more than one action defined for 'ctrl+d': [Abort, MoveToEndOfLine]")
	assert.Contains(t, err.Error(), "- more than one action defined for 'escape': [Abort, MoveToEndOfLine]")
	fmt.Println(err.Error())
}
