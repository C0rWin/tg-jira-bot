package service

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/c0rwin/jira/schema"
)

// Issues is a list of Jira issues
type Issues []schema.Issue

func (i Issues) String() string {
	var s string
	for _, issue := range i {
		s = fmt.Sprintln(s, issue.String())
	}
	return s
}

// Jira service interface to extract data from Jira
type Jira interface {
	// GetRecentTasks returns list of recent tasks created during last week
	GetRecentOpenTasks() (Issues, error)
	// GetAllOpenTasks returns list of all open tasks
	GetAllOpenTasks() (Issues, error)
	// Query returns list of tasks by given query
	Query(query string) (Issues, error)
}

type jira struct {
	projectKey string
	url        string
	username   string
	token      string
	client     *http.Client
}

// GetRecentTasks returns list of recent tasks created in past 24 hours
func (j *jira) GetRecentOpenTasks() (Issues, error) {
	req, err := j.request(fmt.Sprintf("project = %s AND status = Open AND created >= startOfDay(-7d)", j.projectKey))
	if err != nil {
		return nil, err
	}
	// Send the request
	resp, err := j.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return j.issues(resp)
}

// GetAllOpenTasks returns list of all open tasks
func (j *jira) GetAllOpenTasks() (Issues, error) {
	req, err := j.request(fmt.Sprintf("project = %s AND status = Open", j.projectKey))
	if err != nil {
		return nil, err
	}
	// Send the request
	resp, err := j.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return j.issues(resp)
}

// Query returns list of tasks by given query
func (j *jira) Query(query string) (Issues, error) {
	req, err := j.request(query)
	if err != nil {
		return nil, err
	}
	// Send the request
	resp, err := j.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return j.issues(resp)
}

func (j *jira) request(query string) (*http.Request, error) {
	vals := url.Values{}
	vals.Set("jql", query)
	urlStr, err := url.Parse(j.url + "/search?" + vals.Encode())
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("GET", urlStr.String(), nil)
	if err != nil {
		return nil, err
	}

	req.SetBasicAuth(j.username, j.token) // Set basic authentication header
	return req, nil
}

func (j *jira) issues(resp *http.Response) (Issues, error) {
	// Read the response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// Parse the JSON response
	var list struct {
		Issues []schema.Issue `json:"issues"`
	}
	err = json.Unmarshal(body, &list)
	if err != nil {
		return nil, err
	}

	return list.Issues, nil
}

// NewJira returns new instance of Jira service
func NewJira(url, username, token, projectKey string) Jira {
	fmt.Println("Jira URL:", url)
	return &jira{
		projectKey: projectKey,
		url:        url,
		username:   username,
		token:      token,
		client:     &http.Client{},
	}
}
