package model

type UserInfoChanged struct {
	Meta struct {
		Type      string `json:"type"`
		EventID   string `json:"event_id"`
		CreatedAt int64  `json:"created_at"`
		TraceID   string `json:"trace_id"`
		ServiceID string `json:"service_id"`
	} `json:"meta"`
	Payload struct {
		ID         int      `json:"id"`
		Username   string   `json:"userName"`
		Followers  []string `json:"followers,omitempty"`
		Repos      []string `json:"repos,omitempty"`
		Email      string   `json:"userEmail"`
		FirstName  string   `json:"userFirstName"`
		LastName   string   `json:"userLastName"`
		TimeZoneID string   `json:"userTimeZoneId"`
	} `json:"payload"`
}
