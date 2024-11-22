package validator

import (
	"regexp"
	"strings"
)

type Validator struct {
	Errors map[string]string
}

func New() *Validator {
	return &Validator{
		Errors: make(map[string]string),
	}
}

func (v *Validator) ValidateEmail(email string) bool {
	emailRegex := regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`)
	if !emailRegex.MatchString(strings.ToLower(email)) {
		v.Errors["email"] = "Invalid email format"
		return false
	}
	return true
}

func (v *Validator) ValidatePassword(password string) bool {
	if len(password) < 8 {
		v.Errors["password"] = "Password must be at least 8 characters"
		return false
	}

	hasNumber := regexp.MustCompile(`[0-9]`).MatchString(password)
	hasUpper := regexp.MustCompile(`[A-Z]`).MatchString(password)
	hasLower := regexp.MustCompile(`[a-z]`).MatchString(password)
	hasSpecial := regexp.MustCompile(`[!@#$%^&*]`).MatchString(password)

	if !hasNumber || !hasUpper || !hasLower || !hasSpecial {
		v.Errors["password"] = "Password must contain at least one number, uppercase letter, lowercase letter, and special character"
		return false
	}

	return true
}

func (v *Validator) ValidatePhone(phone string) bool {
	phoneRegex := regexp.MustCompile(`^\+?[1-9]\d{1,14}$`)
	if !phoneRegex.MatchString(phone) {
		v.Errors["phone"] = "Invalid phone number format"
		return false
	}
	return true
}

func (v *Validator) HasErrors() bool {
	return len(v.Errors) > 0
}
