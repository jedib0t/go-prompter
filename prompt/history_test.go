package prompt

import (
	"context"
	"strings"
	"testing"
	"time"

	"github.com/go-openapi/strfmt"
	"github.com/stretchr/testify/assert"
)

var (
	testTimestamp1, _   = time.Parse(time.DateTime, "2023-09-01 13:14:15")
	testTimestamp2, _   = time.Parse(time.DateTime, "2023-09-02 14:15:16")
	testHistoryCommands = []HistoryCommand{
		{
			Command:   "foo",
			Timestamp: strfmt.DateTime(testTimestamp1),
		}, {
			Command:   "bar",
			Timestamp: strfmt.DateTime(testTimestamp2),
		},
	}
)

func TestHistoryAppend(t *testing.T) {
	h := History{}
	assert.Len(t, h.Commands, 0)

	h.Append("foo")
	assert.Len(t, h.Commands, 1)

	h.Append("bar")
	assert.Len(t, h.Commands, 2)
}

func TestHistoryGet(t *testing.T) {
	h := History{}
	for _, cmd := range testHistoryCommands {
		h.Append(cmd.Command)
	}

	assert.Equal(t, testHistoryCommands[0].Command, h.Get(0))
	assert.Equal(t, testHistoryCommands[1].Command, h.Get(1))
	assert.Equal(t, "", h.Get(2))
}

func TestHistoryGetNext(t *testing.T) {
	h := History{}
	for _, cmd := range testHistoryCommands {
		h.Append(cmd.Command)
	}
	h.Index = -1

	assert.Equal(t, testHistoryCommands[0].Command, h.GetNext())
	assert.Equal(t, testHistoryCommands[1].Command, h.GetNext())
	assert.Equal(t, "", h.GetNext())
}

func TestHistoryGetPrev(t *testing.T) {
	h := History{}
	for _, cmd := range testHistoryCommands {
		h.Append(cmd.Command)
	}

	assert.Equal(t, testHistoryCommands[1].Command, h.GetPrev())
	assert.Equal(t, testHistoryCommands[0].Command, h.GetPrev())
	assert.Equal(t, testHistoryCommands[0].Command, h.GetPrev())
}

func TestHistoryRender(t *testing.T) {
	h := History{}
	for _, cmd := range testHistoryCommands {
		h.Append(cmd.Command, time.Time(cmd.Timestamp))
	}

	expected := ` # │ TIMESTAMP           │ COMMAND 
───┼─────────────────────┼─────────
 1 │ 2023-09-01 13:14:15 │ foo     
 2 │ 2023-09-02 14:15:16 │ bar     
`
	assert.Equal(t, expected, h.Render(0, 0))

	expected = ` # │ TIMESTAMP           │ COMMAND 
───┼─────────────────────┼─────────
 2 │ 2023-09-02 14:15:16 │ bar     
`
	assert.Equal(t, expected, h.Render(1, 0))

	expected = ` # │ TIMESTAMP           │ COMMAND 
───┼─────────────────────┼─────────
 1 │ 2023-09-01 13:14:15 │ foo     
 2 │ 2023-09-02 14:15:16 │ bar     
`
	assert.Equal(t, expected, h.Render(2, 0))

	h.syntaxHighlighter = strings.ToUpper
	expected = ` # │ TIMESTAMP           │ COMMAND 
───┼─────────────────────┼─────────
 1 │ 2023-09-01 13:14:15 │ FOO     
 2 │ 2023-09-02 14:15:16 │ BAR     
`
	assert.Equal(t, expected, h.Render(3, 0))
}

func TestPromptprocessHistoryCommand(t *testing.T) {
	p := generateTestPrompt(t, context.Background())
	p.SetHistory([]HistoryCommand{
		{Command: "foo"},
		{Command: "bar"},
		{Command: "baz"},
	})
	p.SetHistoryExecPrefix("!")
	p.SetHistoryListPrefix("!!")

	hc := p.processHistoryCommand("!")
	assert.NotNil(t, hc)
	assert.Equal(t, historyCommandExec, hc.Type)
	assert.Equal(t, 0, hc.Value)

	hc = p.processHistoryCommand("!10")
	assert.NotNil(t, hc)
	assert.Equal(t, historyCommandExec, hc.Type)
	assert.Equal(t, 10, hc.Value)

	hc = p.processHistoryCommand("!!")
	assert.NotNil(t, hc)
	assert.Equal(t, historyCommandList, hc.Type)
	assert.Equal(t, 0, hc.Value)

	hc = p.processHistoryCommand("!! 10")
	assert.NotNil(t, hc)
	assert.Equal(t, historyCommandList, hc.Type)
	assert.Equal(t, 10, hc.Value)

	p.SetHistoryExecPrefix("!")
	p.SetHistoryListPrefix("/!")
	hc = p.processHistoryCommand("!!")
	assert.NotNil(t, hc)
	assert.Equal(t, historyCommandExec, hc.Type)
	assert.Equal(t, 3, hc.Value)
}
