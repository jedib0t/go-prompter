package main

import "github.com/jedib0t/go-prompter/prompt"

var (
	tableAndColumnNames = []prompt.Suggestion{
		{Value: "employees", Hint: "Table: Table of Employees"},
		{Value: "id", Hint: "Column: ID (primary key)"},
		{Value: "username", Hint: "Column: User ID"},
		{Value: "first_name", Hint: "Column: First Name"},
		{Value: "last_name", Hint: "Column: Last Name"},
		{Value: "salary", Hint: "Column: Salary Per Month in USD"},
		{Value: "notes", Hint: "Column: Notes"},
		{Value: "date_applied", Hint: "Column: Application Date"},
		{Value: "date_onboarded", Hint: "Column: Onboarding Date"},
	}
)
