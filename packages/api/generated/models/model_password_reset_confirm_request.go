package models

type PasswordResetConfirmRequest struct {

	Token string `json:"token"`

	NewPassword string `json:"newPassword"`
}
