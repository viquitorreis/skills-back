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
	GetAccountById(int) (*Account, error)
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
		(email, fullname, password, admin, sex, language, created_at, updated_at)
		values ($1, $2, $3, $4, $5, $6, $7, $8)
	`

	resp, err := s.db.Query(
		query,
		acc.Email,
		acc.FullName,
		acc.Password,
		acc.Admin,
		acc.Sex,
		acc.Language,
		acc.CreatedAt,
		acc.UpdatedAt,
	)

	if err != nil {
		return err
	}

	fmt.Printf("%v+\n", resp)
	return nil
}

func (s *PostgresStore) DeleteAccount(id int) error {
	return nil
}

func (s *PostgresStore) UpdateAccount(*Account) error {
	return nil
}

func (s *PostgresStore) GetAccountById(id int) (*Account, error) {
	return nil, nil
}
