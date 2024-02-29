package repo

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
)

type DbClient struct {
	Conn *pgxpool.Pool
}

func Connect(dbUrl string) (*DbClient, error) {
	conn, err := pgxpool.New(context.Background(), dbUrl)
	if err != nil {
		return nil, err
	}
	return &DbClient{conn}, nil
}

func (db *DbClient) Close() {
	db.Conn.Close()
}

func (db *DbClient) Ping() bool {
	err := db.Conn.Ping(context.Background())
	if err != nil {
		return false
	} else {
		return true
	}
}

func (db *DbClient) CreateStructure() error {
	_, err := db.Conn.Exec(
		context.Background(),
		`
        CREATE TABLE IF NOT EXISTS events (
            id SERIAL PRIMARY KEY,
            date date,
            hour time,
            venue VARCHAR(255),
            address VARCHAR(255),
            town VARCHAR(255)
        );
        `,
	)

	if err != nil {
		return err
	}
	return nil
}

func (db *DbClient) RemoveStructure() error {
	_, err := db.Conn.Exec(
		context.Background(),
		`
		DROP TABLE IF EXISTS events;
		`,
	)
	return err
}
