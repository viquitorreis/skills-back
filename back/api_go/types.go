package main

import (
	"fmt"
	"time"
)

type Sex string
type Language string

const (
	Male Sex = "male"
	Female Sex = "female"
	Other Sex = "other"
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
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}
 
 //id | email | fullname | password | admin | sex | language | created_at | updated_at 
type CreateAccountRequest struct {
	Email string `json:"email"`
	FullName string `json:"fullName"`
	Password string `json:"password"`
	Sex *Sex `json:"sex"`
	Language *Language `json:"language"`
}

func validateSex(sex string) (*Sex, error) {
	validSex := map[string]Sex{
		"male":   Male,
		"female": Female,
		"other":  Other,
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

	location, err := getBrazilCurrentTimeHelper()
	if err != nil {
		return nil, err
	}

	return &Account{
		Email: email,
		FullName: fullName,
		Password: password,
		Admin: admin,
		Sex: givenSex,
		Language: givenLanguage,
		CreatedAt: time.Now().In(location),
		UpdatedAt: time.Now().In(location),
	}, nil
}
