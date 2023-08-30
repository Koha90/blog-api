package storage

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/koha90/blog-api/internal/types"
)

type Storage interface {
	CreateAccount(*types.Account) error
	DeleteAccount(int) error
	// UpdateAccount(*types.Account) error
	GetAccounts() ([]*types.Account, error)
	GetAccountByID(int) (*types.Account, error)
}

type PostgresStore struct {
	maxPoolSize  int
	connAttempts int
	connTimeout  time.Duration

	pool *pgxpool.Pool
}

func NewPostgresStore(url string) (*PostgresStore, error) {
	pool, err := pgxpool.New(context.Background(), url)
	if err != nil {
		return nil, err
	}

	err = pool.Ping(context.Background())
	if err != nil {
		return nil, err
	}

	return &PostgresStore{
		pool: pool,
	}, nil
}

func (ps *PostgresStore) Init() error {
	return ps.createAccountTable()
}

func (ps *PostgresStore) createAccountTable() error {
	query := `create table if not exists users (
    id serial primary key,
    first_name varchar(50),
    last_name varchar(50),
    number serial,
    encrypted_password varchar(100),
    balance decimal(10, 2),
    created_at timestamp
    )`

	_, err := ps.pool.Exec(context.Background(), query)
	return err
}

func (ps *PostgresStore) CreateAccount(account *types.Account) error {
	_, err := ps.pool.Exec(
		context.Background(),
		`insert into users (first_name, last_name, number, encrypted_password, balance, created_at) values($1, $2, $3, $4, $5, $6)`,
		account.FirstName,
		account.LastName,
		account.Number,
		account.EncryptedPassword,
		account.Balance,
		account.CreatedAt,
	)
	return err
}

func (ps *PostgresStore) GetAccounts() ([]*types.Account, error) {
	rows, err := ps.pool.Query(context.Background(), "select * from users")
	if err != nil {
		return nil, err
	}

	accounts := []*types.Account{}
	for rows.Next() {
		account, err := scanIntoAccount(rows)
		if err != nil {
			return nil, err
		}
		accounts = append(accounts, account)
	}

	return accounts, nil
}

func (ps *PostgresStore) GetAccountByID(id int) (*types.Account, error) {
	rows, err := ps.pool.Query(context.Background(), "select * from users where id = $1", id)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		return scanIntoAccount(rows)
	}

	return nil, fmt.Errorf("account %d not found", id)
}

func (ps *PostgresStore) DeleteAccount(id int) error {
	_, err := ps.pool.Query(context.Background(), "delete from users where id = $1", id)

	return err
}

func (ps *PostgresStore) Close() {
	ps.pool.Close()
}

func scanIntoAccount(rows pgx.Rows) (*types.Account, error) {
	account := new(types.Account)
	err := rows.Scan(
		&account.ID,
		&account.FirstName,
		&account.LastName,
		&account.Number,
		&account.EncryptedPassword,
		&account.Balance,
		&account.CreatedAt,
	)

	return account, err
}
