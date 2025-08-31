package github

import "time"

type PRData struct {
	Number      int      `json:"number"`
	Title       string   `json:"title"`
	Body        string   `json:"body"`
	State       string   `json:"state"`
	BaseRefName string   `json:"baseRefName"`
	MergeCommit *Commit  `json:"mergeCommit"`
	Commits     []Commit `json:"commits"`
	Labels      []Label  `json:"labels"`
}

type Commit struct {
	SHA     string `json:"oid"`
	Message string `json:"messageHeadline"`
	Author  Author `json:"author"`
}

type Author struct {
	Name  string    `json:"name"`
	Email string    `json:"email"`
	Date  time.Time `json:"date"`
}

type Label struct {
	Name  string `json:"name"`
	Color string `json:"color"`
}

type Config struct {
	DryRun          bool
	SkipMergedCheck bool
}
