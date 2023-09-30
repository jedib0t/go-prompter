package input

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/stretchr/testify/assert"
)

var (
	errFoo            = fmt.Errorf("test-error")
	testKeyMsg        = tea.KeyMsg{}
	testMouseMsg      = tea.MouseMsg{}
	testWindowSizeMsg = tea.WindowSizeMsg{}
)

func generateTestReader(ctx context.Context, t *testing.T, opts ...Option) (Reader, *reader) {
	var in bytes.Reader
	opts = append([]Option{WithInput(&in)}, opts...)
	r := NewReader(opts...)
	go r.Begin(ctx)
	<-time.After(time.Second / 4) // time to begin

	rObj, ok := r.(*reader)
	assert.NotNil(t, rObj)
	assert.True(t, ok)
	if !ok {
		t.FailNow()
	}
	assert.NotNil(t, rObj.program)

	return r, rObj
}

func TestReader_Begin(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	r, rObj := generateTestReader(ctx, t)
	defer r.End()

	rObj.program.Kill()
	err, ok := <-r.Errors()
	assert.NotNil(t, err)
	assert.True(t, ok)
	assert.Contains(t, err.Error(), "killed")
}

func TestReader_Errors(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	r, rObj := generateTestReader(ctx, t)
	defer r.End()

	rObj.program.Send(errFoo)
	received, ok := <-r.Errors()
	assert.NotNil(t, received)
	assert.True(t, ok)
	assert.Equal(t, errFoo, received)
}

func TestReader_KeyEvents(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	r, rObj := generateTestReader(ctx, t)
	defer r.End()

	rObj.program.Send(testKeyMsg)
	received, ok := <-r.KeyEvents()
	assert.NotNil(t, received)
	assert.True(t, ok)
	assert.Equal(t, testKeyMsg, received)
}

func TestReader_MouseEvents(t *testing.T) {
	runTest := func(opts ...Option) {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
		defer cancel()

		r, rObj := generateTestReader(ctx, t, opts...)
		defer r.End()

		rObj.program.Send(testMouseMsg)
		received, ok := <-r.MouseEvents()
		assert.NotNil(t, received)
		assert.True(t, ok)
		assert.Equal(t, testMouseMsg, received)
	}

	runTest(WatchMouseClick())
	runTest(WatchMouseAll())
}

func TestReader_Reset(t *testing.T) {
	var in bytes.Reader
	r := NewReader(WithInput(&in), WatchMouseAll())
	rObj, ok := r.(*reader)
	assert.NotNil(t, rObj)
	assert.True(t, ok)
	if !ok {
		t.FailNow()
	}

	originalChDone := rObj.chDone
	originalChErrors := rObj.chErrors
	originalChKeyEvents := rObj.chKeyEvents
	originalChMouseEvents := rObj.chMouseEvents
	originalChWindowSizeEvents := rObj.chWindowSizeEvents
	err := r.Reset()
	assert.Nil(t, err)
	assert.NotEqual(t, originalChDone, rObj.chDone)
	assert.NotEqual(t, originalChErrors, rObj.chErrors)
	assert.NotEqual(t, originalChKeyEvents, rObj.chKeyEvents)
	assert.NotEqual(t, originalChMouseEvents, rObj.chMouseEvents)
	assert.NotEqual(t, originalChWindowSizeEvents, rObj.chWindowSizeEvents)
}

func TestReader_Send(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	var in bytes.Reader
	r := NewReader(WithInput(&in), WatchMouseClick(), WatchWindowSize())
	go r.Begin(ctx)
	defer r.End()
	<-time.After(time.Second / 4) // time to begin

	t.Run("errors", func(t *testing.T) {
		err := r.Send(errFoo)
		assert.Nil(t, err)
		received, ok := <-r.Errors()
		assert.NotNil(t, received)
		assert.True(t, ok)
		assert.Equal(t, errFoo, received)
	})

	t.Run("KeyMsg", func(t *testing.T) {
		err := r.Send(testKeyMsg)
		assert.Nil(t, err)
		received, ok := <-r.KeyEvents()
		assert.NotNil(t, received)
		assert.True(t, ok)
		assert.Equal(t, testKeyMsg, received)
	})

	t.Run("MouseMsg", func(t *testing.T) {
		err := r.Send(testMouseMsg)
		assert.Nil(t, err)
		received, ok := <-r.MouseEvents()
		assert.NotNil(t, received)
		assert.True(t, ok)
		assert.Equal(t, testMouseMsg, received)
	})

	t.Run("WindowSizeMsg", func(t *testing.T) {
		err := r.Send(testWindowSizeMsg)
		assert.Nil(t, err)
		received, ok := <-r.WindowSizeEvents()
		assert.NotNil(t, received)
		assert.True(t, ok)
		assert.Equal(t, testWindowSizeMsg, received)
	})

	t.Run("Unsupported", func(t *testing.T) {
		err := r.Send(r)
		assert.NotNil(t, err)
		assert.True(t, errors.Is(err, ErrUnsupportedMessage))
	})
}

func TestReader_WindowSizeEvents(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	r, rObj := generateTestReader(ctx, t, WatchWindowSize())
	defer r.End()

	rObj.program.Send(testWindowSizeMsg)
	received, ok := <-r.WindowSizeEvents()
	assert.NotNil(t, received)
	assert.True(t, ok)
	assert.Equal(t, testWindowSizeMsg, received)
}
