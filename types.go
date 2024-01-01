package main

import (
	"time"

	"golang.org/x/crypto/bcrypt"
)

type Sex string
type Language string

const (
	Male   Sex = "male"
	Female Sex = "female"
	Other  Sex = "other"
)

const (
	En Language = "en"
	Br Language = "pt-BR"
)

type Account struct {
	ID        int       `json:"id"`
	Email     string    `json:"email"`
	FullName  string    `json:"fullName"`
	Password  string    `json:"-"`
	Admin     bool      `json:"admin"`
	Sex       *Sex      `json:"sex"`
	Country   string    `json:"country"`
	Language  *Language `json:"language"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type NewbieCreateAccountRequest struct {
	Email    string `json:"email" validate:"email,required"`
	FullName string `json:"fullName" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type CreateAccountRequest struct {
	Email    string    `json:"email" validate:"email,required"`
	FullName string    `json:"fullName" validate:"required"`
	Password string    `json:"password" validate:"required"`
	Sex      *Sex      `json:"sex" validate:"required"`
	Country  string    `json:"country" validate:"required"`
	Language *Language `json:"language" validate:"required"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Email string `json:"email"`
	Token string `json:"token"`
}

func NewAccount(email, fullName, password string, admin bool, sex, country, language string) (*Account, error) {

	givenSex, err := validateSexHelper(sex)
	if err != nil {
		return nil, err
	}

	givenLanguage, err := validateLanguageHelper(language)
	if err != nil {
		return nil, err
	}

	location, err := GetBrazilCurrentTimeHelper()
	if err != nil {
		return nil, err
	}

	encryptedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	return &Account{
		Email:     email,
		FullName:  fullName,
		Password:  string(encryptedPassword),
		Admin:     admin,
		Sex:       givenSex,
		Country:   country,
		Language:  givenLanguage,
		CreatedAt: time.Now().In(location),
		UpdatedAt: time.Now().In(location),
	}, nil
}
