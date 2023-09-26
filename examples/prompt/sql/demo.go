package main

import (
	"time"

	"github.com/jedib0t/go-prompter/prompt"
)

var (
	pauseAutoComplete     = time.Second / 4
	pauseAutoCompletePost = time.Second / 4
	pauseBeforeCommand    = time.Second * 2
	pauseBeforeExec       = time.Second / 2
)

var (
	demo01SQLQueries = [][]any{
		{
			"/* Employees.LoadByID */ select * from ",
			"emp", pauseAutoComplete, prompt.Tab, pauseAutoCompletePost,
			"where id = 1;", pauseBeforeExec, prompt.Enter,
		}, {
			"/* Employees.Insert */ insert into ",
			"emp", pauseAutoComplete, prompt.Tab, pauseAutoCompletePost,
			"(fir", pauseAutoComplete, prompt.Tab, prompt.Backspace, pauseAutoCompletePost,
			", las", pauseAutoComplete, prompt.Tab, prompt.Backspace, pauseAutoCompletePost,
			", sal", pauseAutoComplete, prompt.Tab, prompt.Backspace, pauseAutoCompletePost,
			", not", pauseAutoComplete, prompt.Tab, prompt.Backspace, pauseAutoCompletePost,
			") values\n",
			"  ('Arya', 'Stark', 3000, 'Not today.'),\n",
			"  ('Jon', 'Snow', 2000, 'Knows nothing.'),\n",
			"  ('Tyrion', 'Lannister', 5000, 'Pays his debts.');", pauseBeforeExec, prompt.Enter,
		}, {
			"/* Employees.LoadBySalaryRange */ select * from emp", pauseAutoComplete, prompt.Tab, pauseAutoCompletePost,
			"where sal", pauseAutoComplete, prompt.Tab, pauseAutoCompletePost,
			"between 1000 and 6000 order by id;", pauseBeforeExec, prompt.Enter,
		}, {
			"/* Employees.Delete */ delete from emp", pauseAutoComplete, prompt.Tab, pauseAutoCompletePost,
			"where salary < 10000;", pauseBeforeExec, prompt.Enter,
		},
	}
	demo99HistoryAndEnd = [][]any{
		{
			"!!", pauseBeforeExec, prompt.Enter,
		}, {
			"/* demo done */", pauseBeforeExec, " /quit", pauseBeforeExec, prompt.Enter,
		},
	}
)

func runDemo(p prompt.Prompter) {
	var cmdBlocks [][]any
	cmdBlocks = append(cmdBlocks, demo01SQLQueries...)
	cmdBlocks = append(cmdBlocks, demo99HistoryAndEnd...)

	for _, cmdBlock := range cmdBlocks {
		time.Sleep(pauseBeforeCommand)
		_ = p.SendInput(cmdBlock, time.Second/10)
	}
}
