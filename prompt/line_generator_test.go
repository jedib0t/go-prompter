package prompt

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLineCentered(t *testing.T) {
	header := LineCentered("<foo>")
	assert.Equal(t, "", header(0))
	assert.Equal(t, "<fo", header(3))
	assert.Equal(t, "   <foo>  ", header(10))
	assert.Equal(t,
		"                                      <foo>                                     ",
		header(80),
	)

	color := StyleLineNumbersEnabled.Color
	header = LineCentered("<foo>", color)
	assert.Equal(t, "", header(0))
	assert.Equal(t, color.Sprint("<fo"), header(3))
	assert.Equal(t, color.Sprint("   <foo>  "), header(10))
	assert.Equal(t,
		color.Sprint("                                      <foo>                                     "),
		header(80),
	)
}

func TestLineRuler(t *testing.T) {
	color := StyleLineNumbersEnabled.Color
	header := LineRuler()
	assert.Equal(t, "", header(0))
	assert.Equal(t, color.Sprint("----+"), header(5))
	assert.Equal(t, color.Sprint("----+----1"), header(10))
	assert.Equal(t, color.Sprint("----+----1----+----2"), header(20))
	assert.Equal(t, color.Sprint("----+----1----+----2----+----3----+----4"), header(40))
	assert.Equal(t,
		color.Sprint("----+----1----+----2----+----3----+----4----+----5----+----6----+----7----+----8"),
		header(80),
	)

	color = StyleAutoCompleteDefault.ValueColor
	header = LineRuler(color)
	assert.Equal(t, "", header(0))
	assert.Equal(t, color.Sprint("----+"), header(5))
	assert.Equal(t, color.Sprint("----+----1"), header(10))
	assert.Equal(t, color.Sprint("----+----1----+----2"), header(20))
	assert.Equal(t, color.Sprint("----+----1----+----2----+----3----+----4"), header(40))
	assert.Equal(t,
		color.Sprint("----+----1----+----2----+----3----+----4----+----5----+----6----+----7----+----8"),
		header(80),
	)
}

func TestLineSimple(t *testing.T) {
	header := LineSimple("<foo>")
	assert.Equal(t, "", header(0))
	assert.Equal(t, "<fo", header(3))
	assert.Equal(t, "<foo>", header(80))

	color := StyleLineNumbersEnabled.Color
	header = LineSimple("<foo>", color)
	assert.Equal(t, "", header(0))
	assert.Equal(t, color.Sprint("<fo"), header(3))
	assert.Equal(t, color.Sprint("<foo>"), header(80))
}
