package prompt

import (
	"context"
	"strings"
	"testing"
	"time"

	"github.com/muesli/termenv"
	"github.com/stretchr/testify/assert"
)

var (
	testReasonDebug = "test-debug"
)

func TestPrompt_renderView(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	p := generateTestPrompt(t, ctx)
	p.SetFooterGenerator(LineRuler(StyleLineNumbersEnabled.Color))
	p.SetHeaderGenerator(LineRuler(StyleLineNumbersEnabled.Color))
	p.SetPrefixer(PrefixNone())
	p.SetTerminationChecker(TerminationCheckerSQL()) // enable multi-line
	p.Style().LineNumbers = StyleLineNumbersEnabled
	p.init(ctx)

	testSubtitle := "The first render"
	output := strings.Builder{}
	p.buffer.InsertString("This is a test")
	p.updateModel(true)
	p.renderView(termenv.NewOutput(&output), "test")
	expectedLines := []string{
		"\x1b[2K\x1b[38;5;239;48;5;235m----+----1----+----2----+----3----+----4----+----5----+----6----+----7----+----8----+----9----+----0----+----1----+----2\x1b[0m",
		"\x1b[2K\x1b[38;5;239;48;5;235m 1 \x1b[0m This is a test\x1b[38;5;232;48;5;6m \x1b[0m",
		"\x1b[2K\x1b[38;5;239;48;5;235m----+----1----+----2----+----3----+----4----+----5----+----6----+----7----+----8----+----9----+----0----+----1----+----2\x1b[0m",
		"",
	}
	compareModelLines(t, expectedLines, strings.Split(output.String(), "\n"), testSubtitle)

	testSubtitle = "Add a new line of text"
	output.Reset()
	p.buffer.InsertString("\nand this is not a test")
	p.updateModel(true)
	p.renderView(termenv.NewOutput(&output), "test")
	expectedLines = []string{
		"\x1b[3A\x1b[1B\x1b[2K\x1b[38;5;239;48;5;235m 1 \x1b[0m This is a test",
		"\x1b[2K\x1b[38;5;239;48;5;235m 2 \x1b[0m and this is not a test\x1b[38;5;232;48;5;6m \x1b[0m",
		"\x1b[2K\x1b[38;5;239;48;5;235m----+----1----+----2----+----3----+----4----+----5----+----6----+----7----+----8----+----9----+----0----+----1----+----2\x1b[0m",
		"",
	}
	compareModelLines(t, expectedLines, strings.Split(output.String(), "\n"), testSubtitle)

	testSubtitle = "Add one more new line of text"
	output.Reset()
	p.buffer.InsertString("\nand no idea what this is about.")
	p.updateModel(true)
	p.renderView(termenv.NewOutput(&output), "test", true)
	expectedLines = []string{
		"\x1b[4A\x1b[1B\x1b[1B\x1b[2K\x1b[38;5;239;48;5;235m 2 \x1b[0m and this is not a test",
		"\x1b[2K\x1b[38;5;239;48;5;235m 3 \x1b[0m and no idea what this is about.\x1b[38;5;232;48;5;6m \x1b[0m",
		"\x1b[2K\x1b[38;5;239;48;5;235m----+----1----+----2----+----3----+----4----+----5----+----6----+----7----+----8----+----9----+----0----+----1----+----2\x1b[0m",
		"",
	}
	compareModelLines(t, expectedLines, strings.Split(output.String(), "\n"), testSubtitle)

	testSubtitle = "Render the whole thing again"
	output.Reset()
	p.linesRendered = make([]string, 0)
	p.updateModel(true)
	p.renderView(termenv.NewOutput(&output), "test")
	expectedLines = []string{
		"\x1b[2K\x1b[38;5;239;48;5;235m----+----1----+----2----+----3----+----4----+----5----+----6----+----7----+----8----+----9----+----0----+----1----+----2\x1b[0m",
		"\x1b[2K\x1b[38;5;239;48;5;235m 1 \x1b[0m This is a test",
		"\x1b[2K\x1b[38;5;239;48;5;235m 2 \x1b[0m and this is not a test",
		"\x1b[2K\x1b[38;5;239;48;5;235m 3 \x1b[0m and no idea what this is about.\x1b[38;5;232;48;5;6m \x1b[0m",
		"\x1b[2K\x1b[38;5;239;48;5;235m----+----1----+----2----+----3----+----4----+----5----+----6----+----7----+----8----+----9----+----0----+----1----+----2\x1b[0m",
		"",
	}
	compareModelLines(t, expectedLines, strings.Split(output.String(), "\n"), testSubtitle)

	testSubtitle = "Move cursor to the very top"
	output.Reset()
	p.buffer.MoveToBeginning()
	p.updateModel(true)
	p.renderView(termenv.NewOutput(&output), "test")
	expectedLines = []string{
		"\x1b[5A\x1b[1B\x1b[2K\x1b[38;5;239;48;5;235m 1 \x1b[0m \x1b[38;5;232;48;5;6mT\x1b[0mhis is a test",
		"\x1b[1B\x1b[2K\x1b[38;5;239;48;5;235m 3 \x1b[0m and no idea what this is about.",
		"\x1b[1B",
	}
	compareModelLines(t, expectedLines, strings.Split(output.String(), "\n"), testSubtitle)

	testSubtitle = "Hide the cursor"
	output.Reset()
	p.updateModel(false)
	p.renderView(termenv.NewOutput(&output), "test")
	expectedLines = []string{
		"\x1b[5A\x1b[1B\x1b[2K\x1b[38;5;239;48;5;235m 1 \x1b[0m This is a test",
		"\x1b[1B\x1b[1B\x1b[1B",
	}
	compareModelLines(t, expectedLines, strings.Split(output.String(), "\n"), testSubtitle)

	testSubtitle = "Hide the rulers"
	output.Reset()
	p.SetFooterGenerator(nil)
	p.SetHeaderGenerator(nil)
	p.footer = ""
	p.header = ""
	p.updateModel(false)
	p.renderView(termenv.NewOutput(&output), "test")
	expectedLines = []string{
		"\x1b[1A\x1b[2K\x1b[1A\x1b[2K\x1b[3A\x1b[2K\x1b[38;5;239;48;5;235m 1 \x1b[0m This is a test",
		"\x1b[2K\x1b[38;5;239;48;5;235m 2 \x1b[0m and this is not a test",
		"\x1b[2K\x1b[38;5;239;48;5;235m 3 \x1b[0m and no idea what this is about.",
		"",
	}
	compareModelLines(t, expectedLines, strings.Split(output.String(), "\n"), testSubtitle)

	testSubtitle = "Render the whole thing again with debug mode"
	output.Reset()
	p.SetDebug(true)
	p.linesRendered = make([]string, 0)
	p.setDebugData("foo", "bar")
	p.updateModel(true)
	p.renderView(termenv.NewOutput(&output), testReasonDebug)
	actualLines := strings.Split(output.String(), "\n")
	assert.Len(t, actualLines, 5, testSubtitle)
	assert.Contains(t, actualLines[3], p.debugDataAsString(), testSubtitle)
	assert.Contains(t, actualLines[3], "time=", testSubtitle)

	testSubtitle = "Render the diff with debug mode"
	output.Reset()
	p.setDebugData("bar", "baz")
	p.updateModel(true)
	p.renderView(termenv.NewOutput(&output), testReasonDebug)
	actualLines = strings.Split(output.String(), "\n")
	assert.Len(t, actualLines, 2, testSubtitle)
	assert.Contains(t, actualLines[0], p.debugDataAsString(), testSubtitle)
	assert.Contains(t, actualLines[0], "time=", testSubtitle)

	testSubtitle = "Paused"
	output.Reset()
	p.pauseRender()
	p.setDebugData("baz", "foo")
	p.updateModel(true)
	p.renderView(termenv.NewOutput(&output), testReasonDebug)
	assert.Equal(t, "", output.String(), testSubtitle)

	testSubtitle = "Resumed"
	output.Reset()
	p.resumeRender()
	p.updateModel(true)
	p.renderView(termenv.NewOutput(&output), testReasonDebug)
	actualLines = strings.Split(output.String(), "\n")
	assert.Len(t, actualLines, 2, testSubtitle)
	assert.Contains(t, actualLines[0], p.debugDataAsString(), testSubtitle)
	assert.Contains(t, actualLines[0], "time=", testSubtitle)
}
