package main

import (
	"log"
	"math/rand"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type LoginRequest struct {
	Number   int64  `json:"number"`
	Password string `json:"password"`
}

type TransferRequest struct {
	ToAccount int64 `json:"toAccount"`
	Amount    int64 `json:"amount"`
}

type CreateAccountRequest struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Password string `json:"password"`
}

type Account struct {
	ID        int       `json:"id"`
	FirstName string    `json:"firstName"`
	LastName  string    `json:"secondName"`
	Number    int64     `json:"number"`
	EncryptedPassword string `json:"-"`
	Balance   int64     `json:"balance"`
	CreatedAt time.Time `json:"createdAt"`
}

func NewAccount(firstname string, lastname string, password string) *Account {

	encryptedPass, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Println(err)
		return nil
	}

	return &Account{
		FirstName: firstname,
		LastName:  lastname,
		Number:    int64(rand.Intn(10000000)),
		EncryptedPassword: string(encryptedPass),
		Balance:   0,
		CreatedAt: time.Now().UTC(),
	}
}
