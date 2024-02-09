package handler

import (
	"regexp"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

type UserRegistrationRequest struct {
	PhoneNumber string `json:"phone_number"`
	Name        string `json:"name"`
	Password    string `json:"password"`
}

func (r *UserRegistrationRequest) Validate() string {
	var err []string

	if len(r.PhoneNumber) < 10 || len(r.PhoneNumber) > 13 {
		err = append(err, "phone numbers must be at minimum 10 characters and maximum 13 characters")
	}

	if r.PhoneNumber[0:3] != "+62" {
		err = append(err, "phone numbers must start with the Indonesia country code '+62'")
	}

	if len(r.Name) < 3 || len(r.Name) > 30 {
		err = append(err, "name must be at minimum 3 characters and maximum 60 characters")
	}

	if len(r.Password) < 6 || len(r.Password) > 64 {
		err = append(err, "passwords must be minimum 6 characters and maximum 64 characters")
	}

	regexPassword := "/^(?=.*[0-9])(?=.*[a-zA-Z])([a-zA-Z0-9]+)(?=.*[@#$%^&+=]).*$/"
	if match, _ := regexp.MatchString(regexPassword, r.Password); match {
		err = append(err, "password containing at least 1 capital characters AND 1 number AND 1 special (non alpha-numeric) characters")
	}

	return strings.Join(err, ", ")
}

func (r *UserRegistrationRequest) HashAndSaltPassword() (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(r.Password), bcrypt.MinCost)
	if err != nil {
		return "", err
	}

	return string(hash), nil
}

type UserRegistrationResponse struct {
	ID    int    `json:"id"`
	Error string `json:"error"`
}

type UserLoginRequest struct {
	PhoneNumber string `json:"phone_number"`
	Password    string `json:"password"`
}

type UserLoginResponse struct {
	ID    int    `json:"id"`
	Token string `json:"token"`
	Error string `json:"error"`
}

type GetUserResponse struct {
	PhoneNumber string `json:"phone_number"`
	Name        string `json:"name"`
	Error       string `json:"error"`
}

type UpdateUserRequest struct {
	PhoneNumber string `json:"phone_number"`
	Name        string `json:"name"`
}

func (r *UpdateUserRequest) Validate() string {
	var err []string

	if len(r.PhoneNumber) < 10 || len(r.PhoneNumber) > 13 {
		err = append(err, "phone numbers must be at minimum 10 characters and maximum 13 characters")
	}

	if r.PhoneNumber[0:3] != "+62" {
		err = append(err, "phone numbers must start with the Indonesia country code '+62'")
	}

	if len(r.Name) < 3 || len(r.Name) > 30 {
		err = append(err, "name must be at minimum 3 characters and maximum 60 characters")
	}

	return strings.Join(err, ", ")
}
