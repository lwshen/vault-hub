package model

import (
	"fmt"
	"regexp"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Email    string  `gorm:"size:255;uniqueIndex"`
	Password *string `gorm:"type:text"`
	Name     *string `gorm:"size:255"`
	Avatar   *string `gorm:"type:text"`
}

type CreateUserParams struct {
	Email    string
	Password string
	Name     string
}

func (params *CreateUserParams) Validate() map[string]string {
	errors := map[string]string{}
	if !isEmailValid(params.Email) {
		errors["email"] = fmt.Sprintf("email %s is invalid", params.Email)
	}
	return errors
}

func isEmailValid(e string) bool {
	emailRegex := regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`)
	return emailRegex.MatchString(e)
}
