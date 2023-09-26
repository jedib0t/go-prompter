package prompt

import (
	tea "github.com/charmbracelet/bubbletea"
)

// KeySequence defines a special key-sequence that the user presses.
type KeySequence string

// KeySequences are a slice of KeySequence(s).
type KeySequences []KeySequence

// Supported Keys
const (
	AltA            KeySequence = "alt+a"
	AltB            KeySequence = "alt+b"
	AltC            KeySequence = "alt+c"
	AltD            KeySequence = "alt+d"
	AltE            KeySequence = "alt+e"
	AltF            KeySequence = "alt+f"
	AltG            KeySequence = "alt+g"
	AltH            KeySequence = "alt+h"
	AltI            KeySequence = "alt+i"
	AltJ            KeySequence = "alt+j"
	AltK            KeySequence = "alt+k"
	AltL            KeySequence = "alt+l"
	AltM            KeySequence = "alt+m"
	AltN            KeySequence = "alt+n"
	AltO            KeySequence = "alt+o"
	AltP            KeySequence = "alt+p"
	AltQ            KeySequence = "alt+q"
	AltR            KeySequence = "alt+r"
	AltS            KeySequence = "alt+s"
	AltT            KeySequence = "alt+t"
	AltU            KeySequence = "alt+u"
	AltV            KeySequence = "alt+v"
	AltW            KeySequence = "alt+w"
	AltX            KeySequence = "alt+x"
	AltY            KeySequence = "alt+y"
	AltZ            KeySequence = "alt+z"
	ArrowDown       KeySequence = "arrow-down"
	ArrowLeft       KeySequence = "arrow-left"
	ArrowRight      KeySequence = "arrow-right"
	ArrowUp         KeySequence = "arrow-up"
	Backspace       KeySequence = "backspace"
	CtrlA           KeySequence = "ctrl+a"
	CtrlArrowDown   KeySequence = "ctrl+down"
	CtrlArrowLeft   KeySequence = "ctrl+left"
	CtrlArrowRight  KeySequence = "ctrl+right"
	CtrlArrowUp     KeySequence = "ctrl+up"
	CtrlB           KeySequence = "ctrl+b"
	CtrlC           KeySequence = "ctrl+c"
	CtrlD           KeySequence = "ctrl+d"
	CtrlE           KeySequence = "ctrl+e"
	CtrlEnd         KeySequence = "ctrl+end"
	CtrlF           KeySequence = "ctrl+f"
	CtrlG           KeySequence = "ctrl+g"
	CtrlH           KeySequence = "ctrl+h"
	CtrlHome        KeySequence = "ctrl+home"
	CtrlI           KeySequence = "ctrl+i"
	CtrlJ           KeySequence = "ctrl+j"
	CtrlK           KeySequence = "ctrl+k"
	CtrlL           KeySequence = "ctrl+l"
	CtrlM           KeySequence = "ctrl+m"
	CtrlN           KeySequence = "ctrl+n"
	CtrlO           KeySequence = "ctrl+o"
	CtrlP           KeySequence = "ctrl+p"
	CtrlQ           KeySequence = "ctrl+q"
	CtrlR           KeySequence = "ctrl+r"
	CtrlS           KeySequence = "ctrl+s"
	CtrlSpace       KeySequence = "ctrl+space"
	CtrlT           KeySequence = "ctrl+t"
	CtrlU           KeySequence = "ctrl+u"
	CtrlV           KeySequence = "ctrl+v"
	CtrlW           KeySequence = "ctrl+w"
	CtrlX           KeySequence = "ctrl+x"
	CtrlY           KeySequence = "ctrl+y"
	CtrlZ           KeySequence = "ctrl+z"
	Delete          KeySequence = "delete"
	End             KeySequence = "end"
	Enter           KeySequence = "enter"
	Escape          KeySequence = "escape"
	F1              KeySequence = "f1"
	F10             KeySequence = "f10"
	F11             KeySequence = "f11"
	F12             KeySequence = "f12"
	F2              KeySequence = "f2"
	F3              KeySequence = "f3"
	F4              KeySequence = "f4"
	F5              KeySequence = "f5"
	F6              KeySequence = "f6"
	F7              KeySequence = "f7"
	F8              KeySequence = "f8"
	F9              KeySequence = "f9"
	Home            KeySequence = "home"
	Insert          KeySequence = "insert"
	PageDown        KeySequence = "page-down"
	PageUp          KeySequence = "page-up"
	ShiftArrowDown  KeySequence = "shift+down"
	ShiftArrowLeft  KeySequence = "shift+left"
	ShiftArrowRight KeySequence = "shift+right"
	ShiftArrowUp    KeySequence = "shift+up"
	ShiftEnd        KeySequence = "shift+end"
	ShiftHome       KeySequence = "shift+home"
	ShiftTab        KeySequence = "shift-tab"
	Space           KeySequence = "space"
	Tab             KeySequence = "tab"
)

var (
	altKeySequenceMap = map[rune]KeySequence{
		'A': AltA,
		'B': AltB,
		'C': AltC,
		'D': AltD,
		'E': AltE,
		'F': AltF,
		'G': AltG,
		'H': AltH,
		'I': AltI,
		'J': AltJ,
		'K': AltK,
		'L': AltL,
		'M': AltM,
		'N': AltN,
		'O': AltO,
		'P': AltP,
		'Q': AltQ,
		'R': AltR,
		'S': AltS,
		'T': AltT,
		'U': AltU,
		'V': AltV,
		'W': AltW,
		'X': AltX,
		'Y': AltY,
		'Z': AltZ,
		'a': AltA,
		'b': AltB,
		'c': AltC,
		'd': AltD,
		'e': AltE,
		'f': AltF,
		'g': AltG,
		'h': AltH,
		'i': AltI,
		'j': AltJ,
		'k': AltK,
		'l': AltL,
		'm': AltM,
		'n': AltN,
		'o': AltO,
		'p': AltP,
		'q': AltQ,
		'r': AltR,
		's': AltS,
		't': AltT,
		'u': AltU,
		'v': AltV,
		'w': AltW,
		'x': AltX,
		'y': AltY,
		'z': AltZ,
	}
	keyTypeKeySequenceMap = map[tea.KeyType]KeySequence{
		tea.KeyDown:      ArrowDown,
		tea.KeyLeft:      ArrowLeft,
		tea.KeyRight:     ArrowRight,
		tea.KeyUp:        ArrowUp,
		tea.KeyBackspace: Backspace,
		tea.KeyCtrlA:     CtrlA,
		tea.KeyCtrlDown:  CtrlArrowDown,
		tea.KeyCtrlLeft:  CtrlArrowLeft,
		tea.KeyCtrlRight: CtrlArrowRight,
		tea.KeyCtrlUp:    CtrlArrowUp,
		tea.KeyCtrlB:     CtrlB,
		tea.KeyCtrlC:     CtrlC,
		tea.KeyCtrlD:     CtrlD,
		tea.KeyCtrlE:     CtrlE,
		tea.KeyCtrlEnd:   CtrlEnd,
		tea.KeyCtrlF:     CtrlF,
		tea.KeyCtrlG:     CtrlG,
		tea.KeyCtrlH:     CtrlH,
		tea.KeyCtrlHome:  CtrlHome,
		//tea.KeyCtrlI:   CtrlI, // same as tea.KeyTab
		tea.KeyCtrlJ: CtrlJ,
		tea.KeyCtrlK: CtrlK,
		tea.KeyCtrlL: CtrlL,
		//tea.KeyCtrlM:    CtrlM, // same as tea.Enter
		tea.KeyCtrlN:      CtrlN,
		tea.KeyCtrlO:      CtrlO,
		tea.KeyCtrlP:      CtrlP,
		tea.KeyCtrlQ:      CtrlQ,
		tea.KeyCtrlR:      CtrlR,
		tea.KeyCtrlS:      CtrlS,
		tea.KeyCtrlAt:     CtrlSpace,
		tea.KeyCtrlT:      CtrlT,
		tea.KeyCtrlU:      CtrlU,
		tea.KeyCtrlV:      CtrlV,
		tea.KeyCtrlW:      CtrlW,
		tea.KeyCtrlX:      CtrlX,
		tea.KeyCtrlY:      CtrlY,
		tea.KeyCtrlZ:      CtrlZ,
		tea.KeyDelete:     Delete,
		tea.KeyEnd:        End,
		tea.KeyEnter:      Enter,
		tea.KeyEscape:     Escape,
		tea.KeyF10:        F10,
		tea.KeyF11:        F11,
		tea.KeyF12:        F12,
		tea.KeyF1:         F1,
		tea.KeyF2:         F2,
		tea.KeyF3:         F3,
		tea.KeyF4:         F4,
		tea.KeyF5:         F5,
		tea.KeyF6:         F6,
		tea.KeyF7:         F7,
		tea.KeyF8:         F8,
		tea.KeyF9:         F9,
		tea.KeyHome:       Home,
		tea.KeyInsert:     Insert,
		tea.KeyPgDown:     PageDown,
		tea.KeyPgUp:       PageUp,
		tea.KeyShiftDown:  ShiftArrowDown,
		tea.KeyShiftLeft:  ShiftArrowLeft,
		tea.KeyShiftRight: ShiftArrowRight,
		tea.KeyShiftUp:    ShiftArrowUp,
		tea.KeyShiftEnd:   ShiftEnd,
		tea.KeyShiftHome:  ShiftHome,
		tea.KeyShiftTab:   ShiftTab,
		tea.KeySpace:      Space,
		tea.KeyTab:        Tab,
	}
	keySequenceKeyMsgMap = map[KeySequence]tea.KeyMsg{}
)

func init() {
	for kt, ks := range keyTypeKeySequenceMap {
		keySequenceKeyMsgMap[ks] = tea.KeyMsg{Type: kt}
	}
	keySequenceKeyMsgMap[CtrlI] = tea.KeyMsg{Type: tea.KeyCtrlI}
	keySequenceKeyMsgMap[CtrlM] = tea.KeyMsg{Type: tea.KeyCtrlM}

	for r, ks := range altKeySequenceMap {
		keySequenceKeyMsgMap[ks] = tea.KeyMsg{
			Type:  tea.KeyRunes,
			Runes: []rune{r},
			Alt:   true,
		}
	}
}
