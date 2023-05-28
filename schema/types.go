package schema

import (
	"fmt"
)

// Issue represents a JIRA issue
type Issue struct {
	ID     string `json:"id"`
	Key    string `json:"key"`
	Self   string `json:"self"`
	Fields Fields `json:"fields"`
}

// String returns a string representation of an issue
func (i Issue) String() string {
	var s string
	s = fmt.Sprintln(s, fmt.Sprintf("https://nwty.atlassian.net/browse/%s", i.Key))
	s = fmt.Sprintln(s, fmt.Sprintf("*%s*", i.Fields.Summary))
	s = fmt.Sprintln(s, "Type:", i.Fields.IssueType.Name)
	s = fmt.Sprintln(s, "Status:", i.Fields.Status.Name)
	s = fmt.Sprintln(s, "Assignee", i.Fields.Assignee.DisplayName)
	s = fmt.Sprintln(s, "Reporter:", i.Fields.Reporter.DisplayName)
	s = fmt.Sprintln(s, "Created:", i.Fields.Created)
	s = fmt.Sprintln(s, "Updated:", i.Fields.Updated)
	return s
}

// Fields represents a JIRA issue field
type Fields struct {
	Summary     string    `json:"summary"`
	Description string    `json:"description"`
	Project     Project   `json:"project"`
	IssueType   IssueType `json:"issuetype"`
	Assignee    User      `json:"assignee"`
	Reporter    User      `json:"reporter"`
	Created     string    `json:"created"`
	Updated     string    `json:"updated"`
	Status      Status    `json:"status"`
	Priority    Priority  `json:"priority"`
}

// Project represents a JIRA project
type Project struct {
	ID   string `json:"id"`
	Key  string `json:"key"`
	Name string `json:"name"`
	Self string `json:"self"`
}

// IssueType represents a JIRA issue type
type IssueType struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Self string `json:"self"`
}

// User represents a JIRA user
type User struct {
	Name         string `json:"name"`
	DisplayName  string `json:"displayName"`
	EmailAddress string `json:"emailAddress"`
}

// Status represents a JIRA status
type Status struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Self string `json:"self"`
}

// Priority represents a JIRA priority
type Priority struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Self string `json:"self"`
}
