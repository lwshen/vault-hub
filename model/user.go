package model

import (
	"fmt"
	"regexp"

	"github.com/lwshen/vault-hub/internal/auth"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Email    string  `gorm:"size:255;uniqueIndex"`
	Password *string `gorm:"type:text"`
	Name     *string `gorm:"size:255"`
	Avatar   *string `gorm:"type:text"`
}

func (u *User) GetByEmail() error {
	err := DB.Where("email = ?", u.Email).First(&u).Error
	if err != nil {
		return err
	}
	return nil
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
	if ok, msg := isPasswordValid(params.Password); !ok {
		errors["password"] = msg
	}
	return errors
}

func (params *CreateUserParams) Create() (*User, error) {
	hashedPassword, err := hashPassword(params.Password)
	if err != nil {
		return nil, err
	}

	user := User{
		Email:    params.Email,
		Password: &hashedPassword,
		Name:     &params.Name,
	}

	err = DB.Create(&user).Error
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func isEmailValid(e string) bool {
	emailRegex := regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`)
	return emailRegex.MatchString(e)
}

func isPasswordValid(e string) (bool, string) {
	if len(e) < 8 {
		return false, "password must be at least 8 characters long"
	}
	if len(e) > 64 {
		return false, "password must be less than 64 characters long"
	}
	var (
		hasUpper   = regexp.MustCompile(`[A-Z]`).MatchString
		hasLower   = regexp.MustCompile(`[a-z]`).MatchString
		hasNumber  = regexp.MustCompile(`[0-9]`).MatchString
		hasSpecial = regexp.MustCompile(`[!@#\$%\^&\*\(\)_\+\-=\[\]{};':"\\|,.<>\/?]`).MatchString
	)
	if !hasUpper(e) && !hasLower(e) && !hasNumber(e) && !hasSpecial(e) {
		return false, "password must include at least one uppercase letter, one lowercase letter, one number, and one special character"
	}
	return true, ""
}

func hashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

func (u *User) ComparePassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(*u.Password), []byte(password))
	return err == nil
}

func (u *User) GenerateToken() (string, error) {
	return auth.GenerateToken(u.ID)
}
