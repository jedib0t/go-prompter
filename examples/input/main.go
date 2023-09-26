package main

import (
	"context"
	"flag"
	"fmt"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/jedib0t/go-prompter/input"
)

var (
	flagWatchMouseAll   = flag.Bool("watch-mouse-all", false, "Watch for Mouse click & movement events?")
	flagWatchMouseClick = flag.Bool("watch-mouse-click", false, "Watch for Mouse click events?")
	flagWatchWindowSize = flag.Bool("watch-window-size", false, "Watch for Window size/resize events?")
)

func getInputsReaderOptions() []input.Option {
	var opts []input.Option
	if *flagWatchMouseAll {
		opts = append(opts, input.WatchMouseAll())
	} else if *flagWatchMouseClick {
		opts = append(opts, input.WatchMouseClick())
	}
	if *flagWatchWindowSize {
		opts = append(opts, input.WatchWindowSize())
	}
	return opts
}

func logWithTime(msg string, a ...any) {
	fmt.Printf("%s | %s\n", time.Now().Format("15:04:05.000000"), fmt.Sprintf(msg, a...))
}

func main() {
	flag.Parse()

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer func() {
		time.Sleep(time.Second)
		cancel()
	}()

	reader := input.NewReader(getInputsReaderOptions()...)
	go reader.Begin(ctx)
	defer func() {
		reader.End()
	}()

	fmt.Printf("Listening for events ... (ctrl+c to quit)\n\n")
	tick := time.Tick(time.Second)
	for {
		select {
		case <-ctx.Done():
			logWithTime("... boom ...")
			return
		case <-tick:
			logWithTime("... tick ...")
		case err := <-reader.Errors():
			logWithTime("> Error: %v", err)
			return
		case ev := <-reader.KeyEvents():
			logWithTime("> KeyEvent: %s", ev)
			if ev.Type == tea.KeyCtrlC || ev.Type == tea.KeyEscape {
				logWithTime("... done ...")
				return
			}
		case ev := <-reader.MouseEvents():
			logWithTime("> MouseEvent: %#v", ev)
		case ev := <-reader.WindowSizeEvents():
			logWithTime("> ResizeEvent: %#v", ev)
		}
	}
}
