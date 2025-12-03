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
	Password *string
	Name     string
}

func (params *CreateUserParams) Validate() map[string]string {
	errors := map[string]string{}
	if !isEmailValid(params.Email) {
		errors["email"] = fmt.Sprintf("email %s is invalid", params.Email)
	}
	// Only validate password if it's provided (not nil)
	// OIDC users don't need passwords
	if params.Password != nil {
		if ok, msg := isPasswordValid(*params.Password); !ok {
			errors["password"] = msg
		}
	}
	return errors
}

func (params *CreateUserParams) Create() (*User, error) {
	user := User{
		Email: params.Email,
		Name:  &params.Name,
	}

	// Only hash and set password if it's provided
	// OIDC users will have nil password
	if params.Password != nil {
		hashedPassword, err := HashPassword(*params.Password)
		if err != nil {
			return nil, err
		}
		user.Password = &hashedPassword
	}

	err := DB.Create(&user).Error
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
	if !hasUpper(e) || !hasLower(e) || !hasNumber(e) || !hasSpecial(e) {
		return false, "password must include at least one uppercase letter, one lowercase letter, one number, and one special character"
	}
	return true, ""
}

func HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

func (u *User) ComparePassword(password string) bool {
	// OIDC users have nil passwords and cannot login with password
	if u.Password == nil {
		return false
	}
	err := bcrypt.CompareHashAndPassword([]byte(*u.Password), []byte(password))
	return err == nil
}

func (u *User) GenerateToken() (string, error) {
	return auth.GenerateToken(u.ID)
}

// Demo user constants
const (
	DemoUserEmail    = "mock@demo.com"
	demoUserPassword = "Test1234!"
	demoUserName     = "demo"
)

// EnsureDemoUser checks and creates or verifies the demo user when demo mode is enabled
func EnsureDemoUser() error {
	// Check if user exists
	var user User
	err := DB.Where("email = ?", DemoUserEmail).First(&user).Error
	if err == gorm.ErrRecordNotFound {
		// Create demo user
		password := demoUserPassword
		params := CreateUserParams{
			Email:    DemoUserEmail,
			Password: &password,
			Name:     demoUserName,
		}

		if errs := params.Validate(); len(errs) > 0 {
			return fmt.Errorf("demo user validation failed: %v", errs)
		}

		_, err = params.Create()
		if err != nil {
			return fmt.Errorf("failed to create demo user: %w", err)
		}

		return nil
	}

	if err != nil {
		return fmt.Errorf("failed to check for demo user: %w", err)
	}

	// User exists, verify password and name
	if !user.ComparePassword(demoUserPassword) {
		return fmt.Errorf("demo user exists but password does not match expected value")
	}

	if user.Name == nil || *user.Name != demoUserName {
		return fmt.Errorf("demo user exists but name does not match expected value")
	}

	return nil
}
