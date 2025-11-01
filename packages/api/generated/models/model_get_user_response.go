package models

type GetUserResponse struct {

	Email string `json:"email"`

	Name string `json:"name,omitempty"`

	Avatar string `json:"avatar,omitempty"`
}
