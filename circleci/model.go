package circleci

import (
	"fmt"
	"time"
)

// Error represents an error from CircleCI
type Error struct {
	HTTPStatusCode int
	Message        string
}

// Error returns the error message
func (e *Error) Error() string {
	return fmt.Sprintf("%d: %s", e.HTTPStatusCode, e.Message)
}

// Query definition
type Query struct {
	Username string `form:"-"`
	Project  string `form:"-"`
	Branch   string `form:"branch"`
	Status   string `form:"filter"`
	Limit    int    `form:"limit"`
	Offset   int    `form:"offset"`
}

// Build represents a build
type Build struct {
	AuthorDate      time.Time     `json:"author_date"`
	AuthorEmail     string        `json:"author_email"`
	AuthorName      string        `json:"author_name"`
	Reponame        string        `json:"reponame" header:"repository"`
	Branch          string        `json:"branch" header:"branch"`
	BuildNum        int           `json:"build_num" header:"build_num,text"`
	Body            string        `json:"body"`
	BuildTimeMillis int           `json:"build_time_millis"`
	BuildURL        string        `json:"build_url"`
	CommitterDate   time.Time     `json:"committer_date"`
	CommitterEmail  string        `json:"committer_email"`
	CommitterName   string        `json:"committer_name"`
	DontBuild       interface{}   `json:"dont_build"`
	Fleet           string        `json:"fleet"`
	Lifecycle       string        `json:"lifecycle"`
	Outcome         string        `json:"outcome"`
	Parallel        int           `json:"parallel"`
	Platform        string        `json:"platform"`
	PullRequests    []interface{} `json:"pull_requests"`
	QueuedAt        time.Time     `json:"queued_at"`
	StartTime       time.Time     `json:"start_time"`
	Status          string        `json:"status" header:"status"`
	StopTime        time.Time     `json:"stop_time"`
	Subject         string        `json:"subject"`
	UsageQueuedAt   time.Time     `json:"usage_queued_at"`
	User            User          `json:"user"`
	Username        string        `json:"username"`
	VcsRevision     string        `json:"vcs_revision"`
	VcsTag          interface{}   `json:"vcs_tag"`
	VcsURL          string        `json:"vcs_url"`
	Why             string        `json:"why"`
	Workflows       Workflow      `json:"workflows" header:"inline"`
}

// User rerpesent a user that caused that job execution
type User struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	Login     string `json:"login"`
	VcsType   string `json:"vcs_type"`
	IsUser    bool   `json:"is_user"`
	AvatarURL string `json:"avatar_url"`
}

// Workflow represents the workflow information
type Workflow struct {
	JobID                  string                 `json:"job_id"`
	JobName                string                 `json:"job_name" header:"job_name"`
	WorkflowName           string                 `json:"workflow_name" header:"workflow"`
	WorkflowID             string                 `json:"workflow_id"`
	WorkspaceID            string                 `json:"workspace_id"`
	UpstreamJobIds         []string               `json:"upstream_job_ids"`
	UpstreamConcurrencyMap map[string]interface{} `json:"upstream_concurrency_map"`
}
