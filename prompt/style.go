package prompt

import (
	"fmt"
	"time"

	"github.com/muesli/termenv"
)

// Style is used to customize the look and feel of everything about the prompt.
type Style struct {
	AutoComplete StyleAutoComplete `json:"auto_complete"`
	Colors       StyleColors       `json:"colors"`
	Cursor       StyleCursor       `json:"cursor"`
	Dimensions   StyleDimensions   `json:"dimensions"`
	LineNumbers  StyleLineNumbers  `json:"line_numbers"`
	Scrollbar    StyleScrollbar    `json:"scrollbar"`
	TabString    string            `json:"tab_string"`
}

// Validate ensures that the Style can be used without issues.
func (s Style) Validate() error {
	if err := s.Dimensions.validate(); err != nil {
		return err
	}
	return nil
}

// StyleDefault - default Style when none provided.
var StyleDefault = Style{
	AutoComplete: StyleAutoCompleteDefault,
	Colors:       StyleColorsDefault,
	Cursor:       StyleCursorDefault,
	Dimensions:   StyleDimensionsDefault,
	LineNumbers:  StyleLineNumbersNone,
	Scrollbar:    StyleScrollbarDefault,
	TabString:    "    ",
}

// StyleAutoComplete is used to customize the look and feel of the auto-complete
// dropdown.
type StyleAutoComplete struct {
	HintColor          Color          `json:"hint_color"`
	HintSelectedColor  Color          `json:"hint_selected_color"`
	HintLengthMin      int            `json:"hint_length_min"`
	HintLengthMax      int            `json:"hint_length_max"`
	MinChars           int            `json:"min_chars"`
	NumItems           int            `json:"num_items"`
	Scrollbar          StyleScrollbar `json:"scrollbar"`
	ValueColor         Color          `json:"value_color"`
	ValueSelectedColor Color          `json:"value_selected_color"`
	ValueLengthMin     int            `json:"value_length_min"`
	ValueLengthMax     int            `json:"value_length_max"`
	WordDelimiters     map[byte]bool  `json:"word_delimiters"`
}

// StyleAutoCompleteDefault - default Style when none provided.
var StyleAutoCompleteDefault = StyleAutoComplete{
	HintColor: Color{
		Foreground: termenv.ANSI256Color(0),
		Background: termenv.ANSI256Color(39),
	},
	HintSelectedColor: Color{
		Foreground: termenv.ANSI256Color(16),
		Background: termenv.ANSI256Color(208),
	},
	HintLengthMin: 8,
	HintLengthMax: 32,
	MinChars:      0,
	NumItems:      4,
	Scrollbar:     StyleScrollbarAutoComplete,
	ValueColor: Color{
		Foreground: termenv.ANSI256Color(16),
		Background: termenv.ANSI256Color(45),
	},
	ValueSelectedColor: Color{
		Foreground: termenv.ANSI256Color(16),
		Background: termenv.ANSI256Color(214),
	},
	ValueLengthMin: 8,
	ValueLengthMax: 32,
	WordDelimiters: map[byte]bool{
		' ':  true,
		'(':  true,
		')':  true,
		',':  true,
		';':  true,
		'[':  true,
		'\n': true,
		'\t': true,
		']':  true,
		'{':  true,
		'}':  true,
	},
}

// StyleColors is used to customize the colors used on the prompt.
type StyleColors struct {
	Debug Color `json:"debug"`
	Error Color `json:"error"`
}

// StyleColorsDefault - default style when none provided.
var StyleColorsDefault = StyleColors{
	Debug: Color{
		Foreground: termenv.ANSI256Color(22),
		Background: termenv.ANSI256Color(232),
	},
	Error: Color{
		Foreground: termenv.ANSI256Color(9),
		Background: termenv.BackgroundColor(),
	},
}

// StyleCursor is used to customize the look and feel of the cursor.
type StyleCursor struct {
	Blink         bool          `json:"blink"`
	BlinkInterval time.Duration `json:"blink_interval"`
	Color         Color         `json:"color"`
	ColorAlt      Color         `json:"color_alt"`
	Enabled       bool          `json:"enabled"`
}

// StyleCursorDefault - default style when none provided.
var StyleCursorDefault = StyleCursor{
	Blink:         true,
	BlinkInterval: time.Millisecond * 500,
	Color: Color{
		Foreground: termenv.ANSI256Color(232),
		Background: termenv.ANSI256Color(6),
	},
	ColorAlt: Color{
		Foreground: termenv.ANSI256Color(232),
		Background: termenv.ANSI256Color(14),
	},
	Enabled: true,
}

// StyleDimensions is used to customize the sizing of the prompt
type StyleDimensions struct {
	HeightMin uint `json:"height_min"`
	HeightMax uint `json:"height_max"`
	WidthMin  uint `json:"width_min"`
	WidthMax  uint `json:"width_max"`
}

// StyleDimensionsDefault - default style when none provided.
var StyleDimensionsDefault = StyleDimensions{
	HeightMin: 0, // no minimum height (uses as much as possible)
	HeightMax: 0, // no maximum height (uses as much as possible)
	WidthMin:  0, // no minimum width (uses as much as possible)
	WidthMax:  0, // no maximum width (uses as much as possible)
}

func (sd StyleDimensions) validate() error {
	if sd.HeightMin > sd.HeightMax && sd.HeightMax > 0 {
		return fmt.Errorf("%w: height-min [%d] cannot be greater than height-max [%d]",
			ErrInvalidDimensions, sd.HeightMin, sd.HeightMax)
	}
	if sd.WidthMin > sd.WidthMax && sd.WidthMax > 0 {
		return fmt.Errorf("%w: width-min [%d] cannot be greater than width-max [%d]",
			ErrInvalidDimensions, sd.WidthMin, sd.WidthMax)
	}
	if sd.WidthMax < 0 {
		sd.WidthMax = 0
	}
	return nil
}

// StyleLineNumbers is used to customize the look and feel of the line numbers
// in the prompt.
type StyleLineNumbers struct {
	Enabled      bool  `json:"enabled"`
	Color        Color `json:"color"`
	ZeroPrefixed bool  `json:"zero_prefixed"`
}

var (
	// StyleLineNumbersNone - line numbers not enabled.
	StyleLineNumbersNone = StyleLineNumbers{
		Enabled: false,
	}

	// StyleLineNumbersEnabled - enabled with sane defaults.
	StyleLineNumbersEnabled = StyleLineNumbers{
		Enabled: true,
		Color: Color{
			Foreground: termenv.ANSI256Color(239),
			Background: termenv.ANSI256Color(235),
		},
		ZeroPrefixed: false,
	}
)

// StyleScrollbar is used to customize the look and feel of the scrollbar.
type StyleScrollbar struct {
	Color          Color
	Indicator      rune
	IndicatorEmpty rune
}

var (
	// StyleScrollbarDefault - default style when none provided.
	StyleScrollbarDefault = StyleScrollbar{
		Color: Color{
			Foreground: termenv.ANSI256Color(237),
			Background: termenv.ANSI256Color(233),
		},
		Indicator:      '█',
		IndicatorEmpty: '░',
	}

	// StyleScrollbarAutoComplete - default style for the auto-complete
	// drop-down.
	StyleScrollbarAutoComplete = StyleScrollbar{
		Color: Color{
			Foreground: termenv.ANSI256Color(27),
			Background: termenv.ANSI256Color(39),
		},
		Indicator:      '█',
		IndicatorEmpty: '░',
	}
)

// Generate generates the scroll bar strings to be used as suffixes for the
// content lines.
func (s StyleScrollbar) Generate(contentHeight int, cursorLine int, scrollbarHeight int) ([]string, bool) {
	//fmt.Println(contentHeight, cursorLine, scrollbarHeight)

	rsp := make([]string, scrollbarHeight)
	if scrollbarHeight == 0 || scrollbarHeight >= contentHeight {
		return rsp, false
	}

	cursorLocation := (cursorLine * 100) / contentHeight
	scrollWeight, scrollWeightIdx := 100, 0
	for idx := range rsp {
		scrollLocation := ((idx + 1) * 100) / scrollbarHeight
		scrollOffset := scrollLocation - cursorLocation
		if scrollOffset >= 0 && scrollOffset < scrollWeight {
			scrollWeight = scrollOffset
			scrollWeightIdx = idx
		}
		//fmt.Println(idx, cursorLocation, scrollLocation, scrollLocation-cursorLocation, scrollWeight, scrollWeightIdx)
	}
	for idx := range rsp {
		indicator, color := s.IndicatorEmpty, s.Color
		if scrollWeightIdx == idx {
			indicator = s.Indicator
			if indicator == ' ' {
				color = color.Invert()
			}
		}
		rsp[idx] = color.Sprintf("%c", indicator)
	}
	return rsp, true
}
