package audiofeeler

import (
	"database/sql"
)

type accountRecord struct {
	id         int64
	name       string
	source_dir string
}
type nullableAccountRecord struct {
	id         sql.Null[int64]
	name       sql.Null[string]
	source_dir sql.Null[string]
}

type AccountsRepo struct {
	Db *DbClient
}

func (repo *AccountsRepo) Save(account Account) (DatabaseId, error) {
	var query string
	if account.Id == 0 {
		query = "INSERT INTO accounts (name, source_dir) " +
			"VALUES ($1, $2) " +
			"RETURNING id;"
	}

	result, err := repo.Db.Conn.Exec(
		query,
		account.Name, account.SourceDir,
	)
	if err != nil {
		return DatabaseId(0), err
	}
	id, _ := result.LastInsertId()
	return DatabaseId(id), err
}

func buildAccountParams(record accountRecord) *Account {
	account := Account{
		Id:        DatabaseId(record.id),
		Name:      record.name,
		SourceDir: record.source_dir,
	}
	return &account
}
func buildNullableAccountParams(record nullableAccountRecord) *Account {
	account := Account{
		Id:        DatabaseId(record.id.V),
		Name:      record.name.V,
		SourceDir: record.source_dir.V,
	}
	return &account
}

func (repo *AccountsRepo) FindByName(name string) (Account, error) {
	row := repo.Db.Conn.QueryRow(
		`
        SELECT id, name, source_dir
		FROM accounts
		WHERE name = $1;
		`,
		name,
	)

	record := nullableAccountRecord{}
	err := row.Scan(
		&record.id,
		&record.name,
		&record.source_dir,
	)

	if err != nil {
		return Account{}, wrapRecordNotFound(err)
	}

	return *buildNullableAccountParams(record), nil
}

func (repo *AccountsRepo) FindAll() ([]Account, error) {
	var accounts []Account

	rows, err := repo.Db.Conn.Query(
		`
        SELECT id, name, source_dir FROM accounts;
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
			&record.source_dir,
		)

		accounts = append(accounts, *buildAccountParams(record))
	}

	if rows.Err() != nil {
		return nil, rows.Err()
	}

	return accounts, nil
}
