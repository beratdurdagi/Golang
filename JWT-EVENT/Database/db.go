package database

import (
	"context"
	"database/sql"
	"fmt"
	"os"

	"github.com/karalakrepp/Golang/JWT/Types"
	_ "github.com/lib/pq"
)

type PostgresStore struct {
	db *sql.DB
}
type Storer interface {
	CreateUser(context.Context, *Types.User) error
	DeleteAccount(context.Context, int) error
	GetAccountByMail(context.Context, string) (*Types.User, error)
	GetUserById(context.Context, int) (*Types.User, error)
	CreateEvent(context.Context, *Types.Event) error
	GetEvents(context.Context, int) ([]*Types.Event, error)
}

func (s *PostgresStore) CreateEvent(ctx context.Context, event *Types.Event) error {

	query := `INSERT INTO EVENT(title,descr,date,user_id,start_time,end_time,created_at,updated_at) values ($1, $2, $3, $4, $5, $6,$7,$8)`

	_, err := s.db.Query(query,
		event.Title,
		event.Description,
		event.Date,
		event.User_id,
		event.StartTime,
		event.EndTime,
		event.Created_At,
		event.Updated_At,
	)

	if err != nil {
		return err

	}
	return nil

}

func (s *PostgresStore) GetEvents(ctx context.Context, id int) ([]*Types.Event, error) {

	rows, err := s.db.Query(`SELECT * FROM EVENT WHERE user_id = $1`, id)

	if err != nil {
		return nil, err
	}

	events := []*Types.Event{}

	for rows.Next() {
		event, err := scanIntoEvent(rows)
		if err != nil {
			return nil, err
		}
		events = append(events, event)
	}

	return events, nil

}

func (s *PostgresStore) CreateUser(ctx context.Context, user *Types.User) error {

	query := `INSERT INTO ACCOUNT(first_name,last_name,email,password,created_at,updated_at) values ($1, $2, $3, $4, $5, $6)`

	_, err := s.db.Query(query,
		user.Name,
		user.Surname,
		user.Email,
		user.EncPassword,
		user.Created_At,
		user.Updated_At,
	)

	if err != nil {
		return err

	}
	return nil

}

func (s *PostgresStore) GetUserById(ctx context.Context, id int) (*Types.User, error) {
	rows, err := s.db.Query(`SELECT * FROM ACCOUNT WHERE id = $1`, id)

	if err != nil {
		return nil, err
	}

	for rows.Next() {
		return scanIntoAccount(rows)
	}
	return nil, fmt.Errorf("account not found %d", id)
}

func (s *PostgresStore) createAccountTable() error {
	query := `CREATE TABLE IF NOT EXISTS ACCOUNT (
		id serial primary key,
		first_name varchar(100),
		last_name varchar(100),
		email  varchar(255),
		password varchar(255),
		created_at timestamp,
		updated_at  timestamp
	)`

	_, err := s.db.Exec(query)
	return err
}

func (s *PostgresStore) createEventTable() error {
	query := `CREATE TABLE IF NOT EXISTS EVENT (
		id SERIAL PRIMARY KEY,
		title VARCHAR(255) NOT NULL,
		descr TEXT,
		date DATE,
		user_id INT REFERENCES ACCOUNT(id),
		start_time TIME,
		end_time TIME,
		created_at TIMESTAMP,
		updated_at TIMESTAMP
	);
	`
	_, err := s.db.Exec(query)
	return err

}

func NewPostgresStrore() (*PostgresStore, error) {

	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))

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
	if err := s.createAccountTable(); err != nil {
		return err
	}

	if err := s.createEventTable(); err != nil {
		return err
	}

	return nil
}

func (d *PostgresStore) GetAccountByMail(ctx context.Context, mail string) (*Types.User, error) {

	{
		rows, err := d.db.Query(`SELECT * FROM ACCOUNT WHERE email = $1`, mail)

		if err != nil {
			return nil, err
		}

		for rows.Next() {
			return scanIntoAccount(rows)
		}
		return nil, fmt.Errorf("account not found %d", mail)
	}

}

func (s *PostgresStore) DeleteAccount(ctx context.Context, id int) error {
	_, err := s.db.Query("DELETE FROM account WHERE id = $1", id)
	return err
}

func scanIntoAccount(rows *sql.Rows) (*Types.User, error) {
	account := new(Types.User)
	err := rows.Scan(
		&account.ID,
		&account.Name,
		&account.Surname,
		&account.Email,
		&account.EncPassword,
		&account.Created_At,
		&account.Updated_At)

	return account, err
}
func scanIntoEvent(rows *sql.Rows) (*Types.Event, error) {
	event := new(Types.Event)
	err := rows.Scan(
		&event.Id,
		&event.Title,
		&event.Description,
		&event.Date,
		&event.User_id,
		&event.StartTime,
		&event.EndTime,
		&event.Created_At,
		&event.Updated_At)

	return event, err
}
