package prompt

import (
	"strconv"
	"strings"
	"time"

	"github.com/go-openapi/strfmt"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
)

// History contains the past commands executed by the user using Prompt.
type History struct {
	Commands          []HistoryCommand
	Index             int
	syntaxHighlighter SyntaxHighlighter
}

// HistoryCommand contains the command and associated timestamp.
type HistoryCommand struct {
	Command   string          `json:"command"`
	Timestamp strfmt.DateTime `json:"timestamp"`
}

// Append appends a command to the history.
func (h *History) Append(cmd string, optionalTimestamp ...time.Time) {
	timeStamp := time.Now()
	if len(optionalTimestamp) > 0 {
		timeStamp = optionalTimestamp[0]
	}
	h.Commands = append(h.Commands, HistoryCommand{
		Command:   cmd,
		Timestamp: strfmt.DateTime(timeStamp),
	})
	h.Index = len(h.Commands)
}

// Get returns the Nth command from history. N is zero-indexed.
func (h *History) Get(n int) string {
	if n >= 0 && n < len(h.Commands) {
		return h.Commands[n].Command
	}
	return ""
}

// GetNext returns the next command in history.
func (h *History) GetNext() string {
	h.Index++
	if h.Index > len(h.Commands) {
		h.Index = len(h.Commands)
	}
	return h.Get(h.Index)
}

// GetPrev returns the previous command in history.
func (h *History) GetPrev() string {
	h.Index--
	if h.Index < 0 {
		h.Index = 0
	}
	return h.Get(h.Index)
}

// Render renders the list of historic commands to a table and returns the
// string to the client for printing.
func (h *History) Render(numItems int, dispWidth int) string {
	startIdx := 0
	if numItems > 0 && len(h.Commands) > numItems {
		startIdx = len(h.Commands) - numItems
	}
	cmdColWidth := 120
	if dispWidth > 0 {
		cmdColWidth = dispWidth -
			4 /* id */ -
			22 /* timestamp */ -
			2 /* command margin */
	}

	tw := table.NewWriter()
	tw.AppendHeader(table.Row{"#", "Timestamp", "Command"})
	for idx := startIdx; idx < len(h.Commands); idx++ {
		hc := h.Commands[idx]
		timeStamp := "(unknown)"
		if !time.Time(hc.Timestamp).IsZero() {
			timeStamp = time.Time(hc.Timestamp).Format(time.DateTime)
		}
		command := hc.Command
		if h.syntaxHighlighter != nil {
			command = h.syntaxHighlighter(command)
		}
		tw.AppendRow(table.Row{idx + 1, timeStamp, command})
	}
	tw.SetStyle(table.StyleLight)
	tw.SetColumnConfigs([]table.ColumnConfig{
		{Number: 3, WidthMax: cmdColWidth, WidthMaxEnforcer: text.WrapText},
	})
	tw.Style().Options.DrawBorder = false

	if tw.Length() == 0 {
		return "History is empty; nothing to list.\n"
	}
	return tw.Render() + "\n"
}

type historyCommandType string

const (
	historyCommandNone historyCommandType = ""
	historyCommandList historyCommandType = "list"
	historyCommandExec historyCommandType = "exec"
)

type historyCommand struct {
	Type  historyCommandType
	Value int
}

func (p *prompt) processHistoryCommand(input string) *historyCommand {
	input = strings.TrimSpace(input)

	// list?
	if p.historyListPrefix != "" && strings.HasPrefix(input, p.historyListPrefix) {
		input = strings.Replace(input, p.historyListPrefix, "", 1)
		input = strings.TrimSpace(input)
		numItems, err := strconv.Atoi(input)
		if err != nil {
			numItems = 0
		}
		return &historyCommand{Type: historyCommandList, Value: numItems}
	}

	// exec?
	if p.historyExecPrefix != "" && strings.HasPrefix(input, p.historyExecPrefix) {
		input = strings.Replace(input, p.historyExecPrefix, "", 1)
		input = strings.TrimSpace(input)
		itemNum, err := strconv.Atoi(input)
		if err != nil {
			itemNum = 0
		}
		return &historyCommand{Type: historyCommandExec, Value: itemNum}
	}

	return &historyCommand{Type: historyCommandNone}
}
