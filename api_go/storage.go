package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type Storage interface {
	CreateAccount(*Account) error
	DeleteAccount(int) error
	UpdateAccount(*Account) error
	GetAccounts() ([]*Account, error)
	GetAccountById(int) (*Account, error)
	GetAccountByEmail(string) (*Account, error)
}

type PostgresStore struct {
	db *sql.DB
}

func NewPostgresStore() (*PostgresStore, error) {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	
	connStr := os.Getenv("CONN_STR")
	db, err := sql.Open("postgres", connStr)
	// se o db não existe fica eternamente carregando a api
	if err != nil {
		return nil, err
	}

	return &PostgresStore{
		db: db,
	}, nil
}

func (s *PostgresStore) Init() error {
	return s.CreateAccountTable()
}

func (s *PostgresStore) CreateAccountTable() error {
	query := `create table if not exists account (
		id serial primary key,
		email CITEXT unique not null,
		fullname varchar(254) not null,
		password varchar(254) not null,
		admin boolean not null,
		sex text not null,
		country text not null,
		language text not null,
		created_at timestamp,
		updated_at timestamp
	)`

	_, err := s.db.Exec(query)
	return err
}
 
func (s *PostgresStore) CreateAccount(acc *Account) error {
	query := `
		insert into account
		(email, fullname, password, admin, sex, country, language, created_at, updated_at)
		values ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	`

	_, err := s.db.Query(
		query,
		acc.Email,
		acc.FullName,
		acc.Password,
		acc.Admin,
		acc.Sex,
		acc.Country,
		acc.Language,
		acc.CreatedAt,
		acc.UpdatedAt,
	)

	if err != nil {
		return err
	}

	// fmt.Printf("%v+\n", resp)
	return nil
}

func (s *PostgresStore) DeleteAccount(id int) error {
	query := `delete from account where id = $1`
	_, err := s.db.Query(
		query,
		id,
	)

	return err
}

func (s *PostgresStore) UpdateAccount(*Account) error {
	return nil
}

func (s *PostgresStore) GetAccounts() ([]*Account, error) {
	query := `select id, email, fullName, admin, sex, country, language, created_at, updated_at from account`
	rows, err := s.db.Query(query)
	if err != nil {
		return nil, err
	}

	fmt.Printf("%v\n", rows)
	accounts := []*Account{}
	for rows.Next() {
		account, err := ScanIntoAccount(rows)
		if err != nil {
			return nil, err
		}

		accounts = append(accounts, account)
	}

	return accounts, nil
}

func (s *PostgresStore) GetAccountById(id int) (*Account, error) {
	query := `select id, email, fullName, admin, sex, country, language, created_at, updated_at from account where id = $1`
	row, err := s.db.Query(
		query,
		id,
	)
	if err != nil {
		return nil, err
	}

	for row.Next() {
		return ScanIntoAccount(row)
	}

	return nil, fmt.Errorf("Unable to get user from DB where id: %d", id)
}

func (s *PostgresStore) GetAccountByEmail(email string) (*Account, error) {
	query := `select id, email, fullName, password, admin, sex, country, language, created_at, updated_at from account where email = $1`
	row, err := s.db.Query(
		query,
		email,
	)
	if err != nil {
		return nil, err
	}

	for row.Next() {
		return ScanIntoAccount(row)
	}

	return nil, fmt.Errorf("Unable to get user from DB where email: %d", email)
}

func ScanIntoAccount(rows *sql.Rows) (*Account, error) {
	account := &Account{}
	err := rows.Scan(
		&account.ID,
		&account.Email,
		&account.FullName,
		&account.Password,
		&account.Admin,
		&account.Sex,
		&account.Country,
		&account.Language,
		&account.CreatedAt,
		&account.UpdatedAt,
	)

	return account, err
}
