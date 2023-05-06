package model

type GitHubFollowers []*Followers

type Followers struct {
	UserName string `json:"login"`
}
