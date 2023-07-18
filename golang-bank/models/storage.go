package models

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

const DATABASE_URL = "postgres://postgres:Password@HOSTNAME:PORT/GolangBank?sslmode=disable"

type Storage interface {
	CreateAccount(*Account) error
	DeleteAccount(int) error
	UpdateAccount(*Account) error
	GetAccounts() ([]*Account, error)
	GetAccountById(int) (*Account, error)
	GetAccountByNumber(int) (*Account, error)
}
type PostgresStore struct {
	db *sql.DB
}

func NewPostgresStrore() (*PostgresStore, error) {

	db, err := sql.Open("postgres", DATABASE_URL)

	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return &PostgresStore{
		db: db,
	}, nil

}
func (s *PostgresStore) Init() error {
	return s.createAccountTable()
}

func (s *PostgresStore) createAccountTable() error {
	query := `CREATE TABLE IF NOT EXISTS ACCOUNT (
		id serial primary key,
		first_name varchar(100),
		last_name varchar(100),
		password varchar(255),
		number serial,
		balance serial,
		created_at timestamp
	)`

	_, err := s.db.Exec(query)
	return err
}

func (s *PostgresStore) CreateAccount(account *Account) error {

	query := `INSERT INTO ACCOUNT(first_name,last_name,password,number,balance,created_at) values ($1, $2, $3, $4, $5, $6)`

	_, err := s.db.Query(query,
		account.FirstName,
		account.LastName,
		account.EncryptedPassword,
		account.Number,
		account.Balance,
		account.CreatedAt,
	)

	if err != nil {
		return err

	}
	return nil

}
func (s *PostgresStore) DeleteAccount(id int) error {
	_, err := s.db.Query("DELETE FROM account WHERE id = $1", id)
	return err
}

func (s *PostgresStore) UpdateAccount(*Account) error {
	return nil
}
func (s *PostgresStore) GetAccounts() ([]*Account, error) {
	query := "SELECT * FROM ACCOUNT"

	rows, err := s.db.Query(query)

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

func (s *PostgresStore) GetAccountById(id int) (*Account, error) {
	rows, err := s.db.Query(`SELECT * FROM ACCOUNT WHERE id = $1`, id)

	if err != nil {
		return nil, err
	}

	for rows.Next() {
		return scanIntoAccount(rows)
	}
	return nil, fmt.Errorf("account not found %d", id)
}

func scanIntoAccount(rows *sql.Rows) (*Account, error) {
	account := new(Account)
	err := rows.Scan(
		&account.ID,
		&account.FirstName,
		&account.LastName,
		&account.EncryptedPassword,
		&account.Number,
		&account.Balance,
		&account.CreatedAt)

	return account, err
}

func (s *PostgresStore) GetAccountByNumber(number int) (*Account, error) {
	rows, err := s.db.Query(`SELECT * FROM ACCOUNT WHERE number = $1`, number)

	if err != nil {
		return nil, err
	}

	for rows.Next() {
		return scanIntoAccount(rows)
	}
	return nil, fmt.Errorf("account not found %d", number)
}
