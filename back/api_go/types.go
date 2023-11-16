package main

import (
	"fmt"
	"math/rand"
)

type Sex string
type Language string

const (
	Male Sex = "Male"
	Female Sex = "Female"
	Other Sex = "Other"
)

const (
	En Language = "en"
	Br Language = "br"
)

type Account struct {
	ID int `json:"id"`
	Email string `json:"email"`
	FullName string `json:"fullName"`
	Password string `json:"password"`
	Admin bool `json:"admin"`
	Sex *Sex `json:"sex"`
	Language *Language `json:"language"`
}

func validateSex(sex string) (*Sex, error) {
	validSex := map[string]Sex{
		"Male":   Male,
		"Female": Female,
		"Other":  Other,
	}

	if value, ok := validSex[sex]; ok {
		return &value, nil
	}

	return nil, fmt.Errorf("Invalid value for sex: %s", sex)
}

func validateLanguage(language string) (*Language, error) {
	validLanguage := map[string]Language{
		"en": En,
		"br": Br,
	}

	if value, ok := validLanguage[language]; ok {
		return &value, nil
	}
	return nil, fmt.Errorf("Invalid value for lanague: %s", language)
}

func NewAccount(email, fullName, password string, admin bool, sex, language string) (*Account, error) {

	givenSex, err := validateSex(sex)
	if err != nil {
		return nil, err
	}

	givenLanguage, err := validateLanguage(language)
	if err != nil {
		return nil, err
	}

	return &Account{
		ID: rand.Intn(100000), // passar para UUID ( ou algo melhor )
		Email: email,
		FullName: fullName,
		Password: password,
		Admin: true,
		Sex: givenSex,
		Language: givenLanguage,
	}, nil
}
