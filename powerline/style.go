package powerline

import (
	"github.com/jedib0t/go-prompter/prompt"
	"github.com/muesli/termenv"
)

// Style is used to customize the look and feel of the Powerline.
type Style struct {
	Color                 prompt.Color
	InvertSeparatorColors bool
	MarginLeft            string
	MarginRight           string
	SeparatorLeft         string
	SeparatorRight        string
}

var (
	// StyleDefault - default Style when none provided.
	StyleDefault = Style{
		Color: prompt.Color{
			Foreground: termenv.ForegroundColor(),
			Background: termenv.BackgroundColor(),
		},
		InvertSeparatorColors: false,
		MarginLeft:            " ",
		MarginRight:           " ",
		SeparatorLeft:         "",
		SeparatorRight:        "",
	}

	// StyleNonPatched assumes use of regular non-patched fonts.
	StyleNonPatched = Style{
		Color: prompt.Color{
			Foreground: termenv.ForegroundColor(),
			Background: termenv.BackgroundColor(),
		},
		InvertSeparatorColors: true,
		MarginLeft:            " ",
		MarginRight:           " ",
		SeparatorLeft:         "◢",
		SeparatorRight:        "◣",
	}

	// StylePatched assumes use of patched fonts with the separator characters:
	// https://github.com/powerline/fonts
	StylePatched = Style{
		Color: prompt.Color{
			Foreground: termenv.ForegroundColor(),
			Background: termenv.BackgroundColor(),
		},
		InvertSeparatorColors: false,
		MarginLeft:            " ",
		MarginRight:           " ",
		SeparatorLeft:         "\uE0B0",
		SeparatorRight:        "\uE0B2",
	}
)
