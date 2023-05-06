package model

type GitHubRepos []*Repos

type Repos struct {
	Name string `json:"name"`
}
