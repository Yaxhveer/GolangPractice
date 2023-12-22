package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

type Storage interface {
	LoginAccount(*LoginRequest) (string, error)
	CreateAccount(*Account) (*Account, error)
	DeleteAccount(int) error
	GetAccountByID(int) (*Account, error)
	GetAccountByAccNumber(int64) (*Account, error)
	GetAccounts() ([]*Account, error)
}

type PostgresStore struct {
	db *sql.DB
}

func NewPostgresStore() (*PostgresStore, error) {

	user := os.Getenv("USER")
	pass := os.Getenv("PASS")

	connStr := "postgres://"+user+":"+pass+"@localhost/goproj?sslmode=disable"

	db, err := sql.Open("postgres", connStr)

	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return &PostgresStore{db: db}, nil
}

func (s *PostgresStore) Init() error {
	return s.createAccountTable()
}

func (s *PostgresStore) LoginAccount(req *LoginRequest) (string, error) {

	account, err := s.GetAccountByAccNumber(req.Number)

	if err != nil {
		return "", err
	}

	err = bcrypt.CompareHashAndPassword([]byte(account.EncryptedPassword), []byte(req.Password))

	if err != nil {
		return "", err
	}

	tokenString, err := createJWT(account)

	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func (s *PostgresStore) createAccountTable() error {
	query := `
		CREATE TABLE if not exists account(
		id serial primary key,
		first_name text,
		last_name text,
		number text,
		password text,
		balance serial,
		created_at timestamp
	)`

	_, err := s.db.Exec(query)
	return err
}

func (s *PostgresStore) CreateAccount(account *Account) (*Account, error) {

	query := `insert into account 
		(first_name, last_name, number, password, balance, created_at)
		values ($1, $2, $3, $4, $5, $6)
		returning *`

	rows, err := s.db.Query(query, account.FirstName, account.LastName, account.Number, account.EncryptedPassword, account.Balance, account.CreatedAt)

	if err != nil {
		return nil, err
	}
	
	log.Println("Account Created")

	for rows.Next() {
		return scanIntoAccount(rows)
	}
	return nil, fmt.Errorf("Account couldnot be created")
}

func (s *PostgresStore) DeleteAccount(id int) error {

	_, err := s.db.Query(`delete from account where id=$1`, id)
	if err != nil {
		return err
	}

	log.Println("Account deleted ID: ", id)

	return nil
}

func (s *PostgresStore) GetAccountByAccNumber(number int64) (*Account, error) {
	rows, err := s.db.Query(`select * from account where number=$1`, number)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		return scanIntoAccount(rows)
	}
	return nil, fmt.Errorf("Account not found")
}

func (s *PostgresStore) GetAccountByID(id int) (*Account, error) {
	rows, err := s.db.Query(`select * from account where id=$1`, id)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		return scanIntoAccount(rows)
	}
	return nil, fmt.Errorf("Account not found")
}

func (s *PostgresStore) GetAccounts() ([]*Account, error) {

	rows, err := s.db.Query(`select * from account`)

	if err != nil {
		return nil, err
	}

	accounts := []*Account{}

	for rows.Next() {
		account, err := scanIntoAccount(rows)

		if err != nil {
			return nil, err
		}

		accounts = append(accounts, account)
	}

	return accounts, nil
}

func scanIntoAccount(rows *sql.Rows) (*Account, error) {
	account := new(Account)
	err := rows.Scan(
		&account.ID,
		&account.FirstName,
		&account.LastName,
		&account.Number,
		&account.EncryptedPassword,
		&account.Balance,
		&account.CreatedAt)

	return account, err
}
