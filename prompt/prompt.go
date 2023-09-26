package prompt

import (
	"context"
	"fmt"
	"os"
	"strings"
	"sync"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/jedib0t/go-prompter/input"
	"github.com/muesli/termenv"
	"golang.org/x/term"
)

var (
	defaultRefreshInterval = time.Millisecond * 100
)

type prompt struct {
	autoCompleter           AutoCompleter
	autoCompleterContextual AutoCompleter
	debug                   bool
	footerGenerator         LineGenerator
	isInAutoComplete        bool
	headerGenerator         LineGenerator
	history                 History
	historyExecPrefix       string
	historyListPrefix       string
	keyMap                  KeyMap
	keyMapReversed          *keyMapReversed
	prefixer                Prefixer
	promptMutex             sync.Mutex
	refreshInterval         time.Duration
	shortcuts               map[KeySequence]string
	style                   *Style
	syntaxHighlighter       SyntaxHighlighter
	terminationChecker      TerminationChecker
	widthEnforcer           WidthEnforcer

	// render state
	active                      bool
	activeMutex                 sync.RWMutex
	buffer                      *buffer
	cursorColor                 Color
	cursorColorMutex            sync.RWMutex
	debugData                   map[string]string
	debugDataMutex              sync.RWMutex
	displayWidth                int
	displayWidthMutex           sync.RWMutex
	footer                      string
	footerMutex                 sync.RWMutex
	header                      string
	headerMutex                 sync.RWMutex
	linesMutex                  sync.Mutex
	linesRendered               []string
	linesToRender               []string
	reader                      input.Reader
	renderingPaused             bool
	renderingPausedMutex        sync.RWMutex
	suggestions                 []Suggestion
	suggestionsIdx              int
	suggestionsMutex            sync.RWMutex
	syntaxHighlighterCache      map[string][]string
	syntaxHighlighterCacheMutex sync.RWMutex
	timeGen                     time.Duration
	timeSyntaxGen               time.Duration
	timeBufferGen               time.Duration
	timeAutoComplete            time.Duration
}

// CursorLocation returns the current location of the cursor on the prompt.
func (p *prompt) CursorLocation() CursorLocation {
	if p.buffer != nil {
		return p.buffer.Cursor()
	}
	return CursorLocation{}
}

// History returns all the commands executed until now.
func (p *prompt) History() []HistoryCommand {
	return p.history.Commands
}

// IsActive returns true if there is an active prompt at the moment.
func (p *prompt) IsActive() bool {
	p.activeMutex.Lock()
	defer p.activeMutex.Unlock()

	return p.active
}

// KeyMap returns the current KeyMap for inspection/customization.
func (p *prompt) KeyMap() KeyMap {
	return p.keyMap
}

// NumLines returns the number of lines of text in the current active prompt.
func (p *prompt) NumLines() int {
	if p.buffer != nil {
		return p.buffer.NumLines()
	}
	return 0
}

// Prompt prompts. It also watches for cancel KeyEvents on the Context to abort the
// prompt and return control to client.
func (p *prompt) Prompt(ctx context.Context) (string, error) {
	p.promptMutex.Lock()
	defer p.promptMutex.Unlock()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	p.markActive()
	defer p.markInactive()

	// sanity checks
	if err := p.style.validate(); err != nil {
		return "", err
	}

	// init output
	output := termenv.NewOutput(os.Stdout)
	defer func() {
		output.Reset()
	}()

	userInput, err := p.render(ctx, output)
	if err == nil {
		p.history.Append(userInput)
	}
	return userInput, err
}

// SendInput lets you send strings/runes/KeySequence to the currently active
// prompt.
func (p *prompt) SendInput(a []any, delayBetweenRunes ...time.Duration) error {
	delay := time.Duration(0)
	if len(delayBetweenRunes) > 0 {
		delay = delayBetweenRunes[0]
	}

	for idx, item := range a {
		if obj, ok := item.([]rune); ok {
			item = string(obj)
		}

		switch obj := item.(type) {
		case KeySequence:
			if km, ok := keySequenceKeyMsgMap[obj]; ok {
				err := p.reader.Send(km)
				if err != nil {
					return err
				}
			}
		case time.Duration:
			time.Sleep(obj)
		case rune:
			err := p.reader.Send(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{obj}})
			if err != nil {
				return err
			}
			time.Sleep(delay)
		case string:
			for _, r := range obj {
				err := p.reader.Send(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{r}})
				if err != nil {
					return err
				}
				time.Sleep(delay)
			}
		default:
			return fmt.Errorf("%w: [#%d] %#v (allowed: prompt.KeySequence, time.Duration, rune, string)",
				ErrUnsupportedInput, idx, item)
		}
	}
	return nil
}

// SetAutoCompleter sets up the AutoCompleter that will be used to provide
// suggestions. Consider this as an auto completer which will be useful for
// suggesting global stuff like language keywords, or global variables.
//
// This should ideally be set ONCE for the lifetime of a Prompter object.
func (p *prompt) SetAutoCompleter(autoCompleter AutoCompleter) {
	p.autoCompleter = autoCompleter
}

// SetAutoCompleterContextual sets up the AutoCompleter that will be used to
// provide suggestions with more priority than the regular AutoCompleter. This
// is supposed to play the role of providing suggestions like local variables or
// table names when connected to a database because of a previous command.
func (p *prompt) SetAutoCompleterContextual(autoCompleter AutoCompleter) {
	p.autoCompleterContextual = autoCompleter
}

// SetCommandShortcuts sets up command shortcuts. For example, if you want to
// get the prompt input as "/help" when the user presses F1, you'd call this
// function with the argument:
//
//	map[KeySequence]string{
//	   F1: "/help",
//	}
//
// These shortcuts will take precedence over anything in the prompt and
// overwrite the contents of the prompt and return control to the caller.
func (p *prompt) SetCommandShortcuts(shortcuts map[KeySequence]string) {
	p.shortcuts = shortcuts
}

// SetDebug enables/disables debug logs/messages in the prompt.
func (p *prompt) SetDebug(debug bool) {
	p.debug = debug
}

// SetFooter sets up the footer line above the prompt line for each render
// cycle.
func (p *prompt) SetFooter(header string) {
	p.footerGenerator = LineSimple(header)
}

// SetFooterGenerator sets up the FooterGenerator to be used to generate the
// footer line below the prompt line for each render cycle.
func (p *prompt) SetFooterGenerator(footer LineGenerator) {
	p.footerGenerator = footer
}

// SetHeader sets up the header line above the prompt line for each render
// cycle.
func (p *prompt) SetHeader(header string) {
	p.headerGenerator = LineSimple(header)
}

// SetHeaderGenerator sets up the HeaderGenerator to be used to generate the
// header line above the prompt line for each render cycle.
func (p *prompt) SetHeaderGenerator(header LineGenerator) {
	p.headerGenerator = header
}

// SetHistory sets up the history (past inputs) for use when the user wants to
// move backwards/forwards through the command list
func (p *prompt) SetHistory(commands []HistoryCommand) {
	p.history = History{
		Commands: commands,
		Index:    len(commands),
	}
}

// SetHistoryExecPrefix sets up the pattern used to exec command from history.
// Example (prefix="!"):
//   - !10 == execute 10th command
func (p *prompt) SetHistoryExecPrefix(prefix string) {
	p.historyExecPrefix = prefix
}

// SetHistoryListPrefix sets up the prefix used to list commands from history.
// Example (prefix="!!"):
//   - !! == list all commands in history;
//   - !! 10 == list last 10 commands
func (p *prompt) SetHistoryListPrefix(prefix string) {
	p.historyListPrefix = prefix
}

// SetKeyMap sets up the KeyMap used for interacting with the user's input.
func (p *prompt) SetKeyMap(keyMap KeyMap) error {
	kmr, err := keyMap.reverse()
	if err != nil {
		return err
	}

	p.keyMap = keyMap
	p.keyMapReversed = kmr
	return nil
}

// SetPrefix sets up the prefix to use before the prompt.
//
// SetPrefix and SetPrefixer override the same property and the last function
// called take priority.
func (p *prompt) SetPrefix(prefix string) {
	p.prefixer = PrefixText(prefix)
}

// SetPrefixer sets up the prefixer to be called to generate the prefix before
// the prompt for each render cycle.
//
// SetPrefix and SetPrefixer override the same property and the last function
// called take priority.
func (p *prompt) SetPrefixer(prefixer Prefixer) {
	p.prefixer = prefixer
}

// SetRefreshInterval sets up the minimum interval between consecutive renders.
// Note that this can be overridden in case of a mandatory override event like
// a cursor blink.
func (p *prompt) SetRefreshInterval(interval time.Duration) {
	if interval > 0 {
		p.refreshInterval = interval
	} else {
		p.refreshInterval = defaultRefreshInterval
	}
}

// SetStyle sets up the Style sheet to be followed for the render.
func (p *prompt) SetStyle(s Style) {
	p.style = &s
}

// SetSyntaxHighlighter sets up the function that will colorize and highlight
// keywords in the user-input
func (p *prompt) SetSyntaxHighlighter(highlighter SyntaxHighlighter) {
	p.syntaxHighlighter = highlighter
}

// SetTerminationChecker sets up the termination checker to check if the user
// input is done and can be returned to caller on "Terminate" action.
func (p *prompt) SetTerminationChecker(checker TerminationChecker) {
	p.terminationChecker = checker
}

// SetWidthEnforcer sets up the function to wrap lines longer than the prompt
// width.
func (p *prompt) SetWidthEnforcer(enforcer WidthEnforcer) {
	p.widthEnforcer = enforcer
}

// Style returns the current Style in use, so it can be modified on the fly in
// between two prompts.
func (p *prompt) Style() *Style {
	return p.style
}

func (p *prompt) changeSuggestionsIdx(v int) bool {
	p.suggestionsMutex.Lock()
	defer p.suggestionsMutex.Unlock()

	if v > 0 && p.suggestionsIdx < len(p.suggestions)-1 {
		p.suggestionsIdx++
		return true
	} else if v < 0 && p.suggestionsIdx > 0 {
		p.suggestionsIdx--
		return true
	}
	return false
}

func (p *prompt) clearDebugData(prefixes ...string) {
	p.debugDataMutex.Lock()
	defer p.debugDataMutex.Unlock()

	if len(prefixes) == 0 {
		p.debugData = make(map[string]string)
	} else {
		for _, prefix := range prefixes {
			for k := range p.debugData {
				if strings.HasPrefix(k, prefix) {
					delete(p.debugData, k)
				}
			}
		}
	}
}

func (p *prompt) clearSyntaxHighlighterCache() {
	p.syntaxHighlighterCacheMutex.Lock()
	p.syntaxHighlighterCache = make(map[string][]string)
	p.syntaxHighlighterCacheMutex.Unlock()
}

func (p *prompt) debugDataAsString() string {
	p.debugDataMutex.RLock()
	defer p.debugDataMutex.RUnlock()

	return fmt.Sprint(p.debugData)
}

func (p *prompt) doSyntaxHighlighting(lines []string) []string {
	linesStr := strings.Join(lines, "\n")

	p.syntaxHighlighterCacheMutex.RLock()
	cacheVal, ok := p.syntaxHighlighterCache[linesStr]
	p.syntaxHighlighterCacheMutex.RUnlock()

	if !ok {
		cacheVal = strings.Split(p.syntaxHighlighter(linesStr), "\n")

		p.syntaxHighlighterCacheMutex.Lock()
		p.syntaxHighlighterCache[linesStr] = cacheVal
		p.syntaxHighlighterCacheMutex.Unlock()
	}

	return cacheVal
}

func (p *prompt) getCursorColor() Color {
	p.cursorColorMutex.RLock()
	defer p.cursorColorMutex.RUnlock()

	return p.cursorColor
}

func (p *prompt) getDisplayWidth() int {
	p.displayWidthMutex.RLock()
	defer p.displayWidthMutex.RUnlock()

	return p.displayWidth
}

func (p *prompt) getSuggestionsAndIdx() ([]Suggestion, int) {
	p.suggestionsMutex.RLock()
	defer p.suggestionsMutex.RUnlock()

	return append([]Suggestion{}, p.suggestions...), p.suggestionsIdx
}

func (p *prompt) getFooter() string {
	p.footerMutex.RLock()
	defer p.footerMutex.RUnlock()

	return p.footer
}

func (p *prompt) getHeader() string {
	p.headerMutex.RLock()
	defer p.headerMutex.RUnlock()

	return p.header
}

func (p *prompt) init(ctx context.Context) {
	p.initSync(ctx)

	go p.switchCursorColors(ctx)
	go p.updateSuggestions(ctx)
	go p.updateHeaderAndFooterAsync(ctx)
}

func (p *prompt) initSync(ctx context.Context) {
	termWidth, _, _ := term.GetSize(int(os.Stdout.Fd()))
	p.updateDisplayWidth(termWidth)
	p.updateHeaderAndFooter()

	// in the buffer or reset it to previous state
	if p.buffer == nil {
		p.buffer = newBuffer()
	} else {
		p.buffer.Reset()
	}
	p.buffer.SetTab(p.style.TabString)

	// clear the rendering state model
	p.linesMutex.Lock()
	p.linesRendered = make([]string, 0)
	p.linesToRender = make([]string, 0)
	p.linesMutex.Unlock()

	// clear/reset the reader
	if p.reader == nil {
		p.reader = input.NewReader(input.WatchWindowSize())
	} else {
		_ = p.reader.Reset()
	}

	// clear other things
	p.clearDebugData()
	p.clearSyntaxHighlighterCache()
	p.history.syntaxHighlighter = p.syntaxHighlighter
	p.resumeRender()
	p.setCursorColor(p.style.Cursor.Color)

	// reset timestamps
	p.timeGen = time.Duration(0)
	p.timeSyntaxGen = time.Duration(0)
	p.timeBufferGen = time.Duration(0)
	p.timeAutoComplete = time.Duration(0)
}

func (p *prompt) isRenderPaused() bool {
	p.renderingPausedMutex.RLock()
	defer p.renderingPausedMutex.RUnlock()

	return p.renderingPaused
}

func (p *prompt) markActive() {
	p.activeMutex.Lock()
	defer p.activeMutex.Unlock()

	p.active = true
}

func (p *prompt) markInactive() {
	p.activeMutex.Lock()
	defer p.activeMutex.Unlock()

	p.active = false
}

func (p *prompt) pauseRender() {
	p.renderingPausedMutex.Lock()
	defer p.renderingPausedMutex.Unlock()

	p.renderingPaused = true
}

func (p *prompt) resumeRender() {
	p.renderingPausedMutex.Lock()
	defer p.renderingPausedMutex.Unlock()

	p.renderingPaused = false
}

func (p *prompt) setCursorColor(c Color) {
	p.cursorColorMutex.Lock()
	defer p.cursorColorMutex.Unlock()

	p.cursorColor = c
}

func (p *prompt) setDebugData(k, v string) {
	p.debugDataMutex.Lock()
	defer p.debugDataMutex.Unlock()

	if v == "" {
		delete(p.debugData, k)
	} else {
		p.debugData[k] = v
	}
}

func (p *prompt) setDisplayWidth(w int) {
	p.displayWidthMutex.Lock()
	defer p.displayWidthMutex.Unlock()

	p.displayWidth = w
}

func (p *prompt) setSuggestions(s []Suggestion) {
	p.suggestionsMutex.Lock()
	defer p.suggestionsMutex.Unlock()

	p.suggestions = s
}

func (p *prompt) setSuggestionsIdx(idx int) {
	p.suggestionsMutex.Lock()
	defer p.suggestionsMutex.Unlock()

	p.suggestionsIdx = idx
}

func (p *prompt) switchCursorColors(ctx context.Context) {
	if p.style.Cursor.Blink {
		isLow := true
		tick := time.Tick(p.style.Cursor.BlinkInterval)
		for {
			select {
			case <-ctx.Done():
				return
			case <-tick:
				cursorColor := p.style.Cursor.Color
				if isLow {
					cursorColor = p.style.Cursor.ColorAlt
				}
				isLow = !isLow

				p.setCursorColor(cursorColor)
			}
		}
	}
}

func (p *prompt) translateKeyToKeySequence(key tea.KeyMsg) KeySequence {
	var ks KeySequence
	if key.Alt == true && len(key.Runes) > 0 {
		ks = altKeySequenceMap[key.Runes[0]]
	} else {
		ks = keyTypeKeySequenceMap[key.Type]
	}
	return ks
}

func (p *prompt) translateKeyToAutoCompleteAction(key tea.KeyMsg) Action {
	return p.keyMapReversed.AutoComplete[p.translateKeyToKeySequence(key)]
}

func (p *prompt) translateKeyToInsertAction(key tea.KeyMsg) Action {
	return p.keyMapReversed.Insert[p.translateKeyToKeySequence(key)]
}

func (p *prompt) updateDisplayWidth(termWidth int) {
	termWidth = clampValue(
		termWidth,
		int(p.style.Dimensions.WidthMin),
		int(p.style.Dimensions.WidthMax),
	)
	if p.debug {
		termWidth -= 4 // account for debug margin
	}
	p.setDisplayWidth(termWidth)
}

func (p *prompt) updateHeaderAndFooter() {
	if p.headerGenerator != nil {
		header := p.headerGenerator(p.getDisplayWidth())
		p.headerMutex.Lock()
		p.header = header
		p.headerMutex.Unlock()
	}

	if p.footerGenerator != nil {
		footer := p.footerGenerator(p.getDisplayWidth())
		p.footerMutex.Lock()
		p.footer = footer
		p.footerMutex.Unlock()
	}
}

func (p *prompt) updateHeaderAndFooterAsync(ctx context.Context) {
	if p.headerGenerator != nil {
		tick := time.Tick(p.refreshInterval)
		for {
			select {
			case <-ctx.Done():
				return
			case <-tick:
				p.updateHeaderAndFooter()
			}
		}
	}
}
