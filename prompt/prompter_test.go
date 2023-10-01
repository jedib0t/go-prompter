package prompt

import (
	"errors"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	t.Run("error", func(t *testing.T) {
		KeyMapDefault.AutoComplete.ChoosePrevious = KeySequences{ArrowUp, ArrowDown}
		defer func() {
			KeyMapDefault.AutoComplete.ChoosePrevious = KeySequences{ArrowUp}
		}()

		p, err := New()
		assert.Nil(t, p)
		assert.NotNil(t, err)
		assert.True(t, errors.Is(err, ErrDuplicateKeyAssignment))
	})

	t.Run("success", func(t *testing.T) {
		p, err := New()
		assert.NotNil(t, p)
		assert.Nil(t, err)
		if err != nil {
			t.FailNow()
		}

		pObj, ok := p.(*prompt)
		assert.NotNil(t, pObj)
		assert.True(t, ok)
		if !ok {
			t.FailNow()
		}

		assert.Equal(t, DefaultHistoryExecPrefix, pObj.historyExecPrefix)
		assert.Equal(t, DefaultHistoryListPrefix, pObj.historyListPrefix)
		assert.Equal(t, os.Stdin, pObj.input)
		assert.Equal(t, os.Stdout, pObj.output)
		assert.Equal(t, PrefixSimple()(), pObj.prefixer())
		assert.Equal(t, DefaultRefreshInterval, pObj.refreshInterval)
		assert.NotNil(t, pObj.style)
		if pObj.style != nil {
			assert.Equal(t, StyleDefault, *pObj.style)
		}
		assert.NotNil(t, pObj.terminationChecker)
		if pObj.terminationChecker != nil {
			assert.True(t, pObj.terminationChecker("foo"))
		}
		assert.NotNil(t, pObj.widthEnforcer)
		if pObj.widthEnforcer != nil {
			assert.Equal(t, WidthEnforcerDefault("foo", 2), pObj.widthEnforcer("foo", 2))
		}
	})
}
