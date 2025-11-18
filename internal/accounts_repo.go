package audiofeeler

import (
	"github.com/mradmacher/audiofeeler/optiomist"
	"github.com/mradmacher/audiofeeler/sqlbuilder"
)

type accountRecord struct {
	id    uint32
	name  string
	title string
	url   string
}

type AccountsRepo struct {
	Db *DbClient
}

func (repo *AccountsRepo) Create(account Account) (uint32, error) {
	fields := sqlbuilder.Fields{
		"name":  account.Name,
		"title": account.Title,
		"url":   account.Url,
	}
	query, values := fields.BuildInsert("accounts")
	row := repo.Db.Conn.QueryRow(
		query,
		values...,
	)
	var id uint32
	err := row.Scan(&id)
	return id, err
}

func buildAccountParams(record accountRecord) *Account {
	account := Account{
		Id:    optiomist.Optiomize(record.id, true),
		Name:  optiomist.Optiomize(record.name, true),
		Title: optiomist.Optiomize(record.title, true),
		Url:   optiomist.Optiomize(record.url, true),
	}
	return &account
}

func (repo *AccountsRepo) FindByName(name string) (Account, error) {
	row := repo.Db.Conn.QueryRow(
		`
        SELECT id, name, title, url
		FROM accounts
		WHERE name = $1;
		`,
		name,
	)

	var record accountRecord
	err := row.Scan(
		&record.id,
		&record.name,
		&record.title,
		&record.url,
	)

	if err != nil {
		return Account{}, wrapRecordNotFound(err)
	}

	return *buildAccountParams(record), nil
}

func (repo *AccountsRepo) FindAll() ([]Account, error) {
	var accounts []Account

	rows, err := repo.Db.Conn.Query(
		`
        SELECT id, name, title, url FROM accounts;
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
			&record.title,
			&record.url,
		)

		accounts = append(accounts, *buildAccountParams(record))
	}

	if rows.Err() != nil {
		return nil, rows.Err()
	}

	return accounts, nil
}
