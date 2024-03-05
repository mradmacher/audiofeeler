package repo

import (
	"context"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/mradmacher/audiofeeler/internal"
	"github.com/mradmacher/audiofeeler/optiomist"
)

type accountRecord struct {
	id   pgtype.Uint32
	name pgtype.Text
	url  pgtype.Text
}

type AccountsRepo struct {
	Db *DbClient
}

func (repo *AccountsRepo) Create(account audiofeeler.Account) (uint32, error) {
	fields := Fields{
		"name": account.Name,
		"url":  account.Url,
	}
	query, values := fields.BuildInsert("accounts")
	row := repo.Db.Conn.QueryRow(context.Background(),
		query,
		values...,
	)
	var id uint32
	err := row.Scan(&id)
	return id, err
}

func buildAccountParams(record accountRecord) *audiofeeler.Account {
	account := audiofeeler.Account{
		Id:   optiomist.Optiomize(record.id.Uint32, record.id.Valid),
		Name: optiomist.Optiomize(record.name.String, record.name.Valid),
		Url:  optiomist.Optiomize(record.url.String, record.url.Valid),
	}
	return &account
}

func (repo *AccountsRepo) FindAll() ([]audiofeeler.Account, error) {
	var accounts []audiofeeler.Account

	rows, err := repo.Db.Conn.Query(context.Background(),
		`
        SELECT id, name, url FROM accounts;
        `,
	)
	defer rows.Close()

	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var record accountRecord
		err = rows.Scan(
			&record.id,
			&record.name,
			&record.url,
		)

		accounts = append(accounts, *buildAccountParams(record))
	}

	if rows.Err() != nil {
		return nil, rows.Err()
	}

	return accounts, nil
}
