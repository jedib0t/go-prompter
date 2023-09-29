package prompt

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	p, err := New()
	assert.Nil(t, p)
	assert.NotNil(t, err)
	assert.True(t, errors.Is(err, ErrNonInteractiveShell))
}
