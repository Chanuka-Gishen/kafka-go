package model

type User struct {
	ID         int      `json:"id"`
	Username   string   `json:"userName"`
	Followers  []string `json:"followers,omitempty"`
	Repos      []string `json:"repos,omitempty"`
	Email      string   `json:"userEmail"`
	FirstName  string   `json:"userFirstName"`
	LastName   string   `json:"userLastName"`
	TimeZoneID string   `json:"userTimeZoneId"`
}
