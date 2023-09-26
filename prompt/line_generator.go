package prompt

import (
	"fmt"

	"github.com/jedib0t/go-pretty/v6/text"
)

// LineGenerator is a function that takes the width of the terminal as input and
// generates a "line" of content to be display on the terminal above or below
// the prompt.
type LineGenerator func(width int) string

// LineCentered uses the given text as the "header" text and centers it.
func LineCentered(title string, optionalColor ...Color) LineGenerator {
	return func(width int) string {
		rsp := title
		if text.RuneWidthWithoutEscSequences(rsp) > width {
			rsp = text.Trim(rsp, width)
		}
		rsp = text.AlignCenter.Apply(rsp, width)
		if len(rsp) > 0 && len(optionalColor) > 0 {
			rsp = optionalColor[0].Sprint(rsp)
		}
		return rsp
	}
}

// LineRuler prints a ruler above the prompt.
func LineRuler(optionalColor ...Color) LineGenerator {
	rulerMap := make(map[int]string) // cache
	color := StyleLineNumbersEnabled.Color
	if len(optionalColor) > 0 {
		color = optionalColor[0]
	}

	return func(width int) string {
		if width <= 0 {
			return ""
		}
		if _, ok := rulerMap[width]; !ok {
			ruler := make([]rune, width)
			for idx := 1; idx <= width; idx++ {
				if idx%10 == 0 {
					tenths := []rune(fmt.Sprint(idx / 10))
					ruler[idx-1] = tenths[len(tenths)-1]
				} else if idx%5 == 0 {
					ruler[idx-1] = '+'
				} else {
					ruler[idx-1] = '-'
				}
			}
			rulerMap[width] = color.Sprint(string(ruler))
		}
		return rulerMap[width]
	}
}

// LineSimple uses the given text as the "header" text.
func LineSimple(title string, optionalColor ...Color) LineGenerator {
	return func(width int) string {
		rsp := title
		if text.RuneWidthWithoutEscSequences(rsp) > width {
			rsp = text.Trim(rsp, width)
		}
		if len(rsp) > 0 && len(optionalColor) > 0 {
			rsp = optionalColor[0].Sprint(rsp)
		}
		return rsp
	}
}
