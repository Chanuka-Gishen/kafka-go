package model

type GitHubUser struct {
	UserName string  `json:"login"`
	Email    *string `json:"email"`
}
