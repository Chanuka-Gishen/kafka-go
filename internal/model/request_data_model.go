package model

type RequestData struct {
	Data DataValues `json:"data" validate:"required"`
}

type DataValues struct {
	Email      string `json:"email" validate:"required"`
	FirstName  string `json:"first_name" validate:"required"`
	LastName   string `json:"last_name" validate:"required"`
	TimeZoneID string `json:"time_zone_id" validate:"required"`
}
