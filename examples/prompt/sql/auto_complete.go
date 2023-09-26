package main

import "github.com/jedib0t/go-prompter/prompt"

var (
	tableAndColumnNames = []prompt.Suggestion{
		{Value: "employees", Hint: "Table of Employees"},
		{Value: "id", Hint: "ID (primary key)"},
		{Value: "username", Hint: "User ID"},
		{Value: "first_name", Hint: "First Name"},
		{Value: "last_name", Hint: "Last Name"},
		{Value: "salary", Hint: "Salary Per Month in USD"},
		{Value: "notes", Hint: "Notes"},
		{Value: "date_applied", Hint: "Application Date"},
		{Value: "date_onboarded", Hint: "Onboarding Date"},
	}
)
