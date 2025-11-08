package api

import (
	"net/http"

	"github.com/lwshen/vault-hub/model"
)

// BuildCurrentUserResponse converts the authenticated user to API shape.
func BuildCurrentUserResponse(user *model.User) (GetUserResponse, *APIError) {
	if user == nil {
		return GetUserResponse{}, newAPIError(http.StatusUnauthorized, "user not found in context")
	}

	resp := GetUserResponse{
		Email:  user.Email,
		Avatar: user.Avatar,
		Name:   user.Name,
	}
	return resp, nil
}
