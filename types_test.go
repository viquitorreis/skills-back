package main

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewAcc(t *testing.T) {
	acc, err := NewAccount("a", "b", "123", false, "male", "brazil", "br")
	assert.Nil(t, err)

	fmt.Printf("%+v\n", acc)
}

//{
//		Email: email,
//		FullName: fullName,
//		Password: string(encryptedPassword),
//		Admin: admin,
//		Sex: givenSex,
//		Country: country,
//		Language: givenLanguage,
//		CreatedAt: time.Now().In(location),
//		UpdatedAt: time.Now().In(location),
//	}
