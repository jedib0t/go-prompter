package main

import (
	"time"

	"github.com/go-openapi/strfmt"
	"github.com/jedib0t/go-prompter/prompt"
)

var (
	history = []prompt.HistoryCommand{
		{
			Command:   "/* Employees.LoadByID */ select * from employees where id = 1;",
			Timestamp: strfmt.DateTime(timeStamp1),
		}, {
			Command: "/* Employees.Insert */ insert into employees (first_name, last_name, salary, notes) values\n" +
				"  ('Arya', 'Stark', 3000, 'Not today.'),\n" +
				"  ('Jon', 'Snow', 2000, 'Knows nothing.'),\n" +
				"  ('Tyrion', 'Lannister', 5000, 'Pays his debts.');",
			Timestamp: strfmt.DateTime(timeStamp2),
		}, {
			Command:   "/* Employees.LoadBySalaryRange */ select * from employees where salary between 1000 and 6000 order by id;",
			Timestamp: strfmt.DateTime(timeStamp3),
		}, {
			Command:   "/* Employees.Delete */ delete from employees where salary < 10000;",
			Timestamp: strfmt.DateTime(timeStamp4),
		},
	}
	timeStamp1, _ = time.Parse(time.DateTime, "2023-09-01 13:00:00")
	timeStamp2, _ = time.Parse(time.DateTime, "2023-09-02 15:00:00")
	timeStamp3, _ = time.Parse(time.DateTime, "2023-09-03 17:00:00")
	timeStamp4, _ = time.Parse(time.DateTime, "2023-09-03 19:00:00")
)
