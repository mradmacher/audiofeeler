package audiofeeler

type Account struct {
	Id   DatabaseId
	Name string
}

type accountRecord struct {
	id   int64
	name string
}

type AccountsRepo struct {
	Db *DbClient
}

func (repo *AccountsRepo) Save(account Account) (DatabaseId, error) {
	var query string
	if account.Id == 0 {
		query = "INSERT INTO accounts (name) " +
			"VALUES ($1) " +
			"RETURNING id;"
	}

	result, err := repo.Db.Conn.Exec(
		query,
		account.Name,
	)
	if err != nil {
		return DatabaseId(0), err
	}
	id, _ := result.LastInsertId()
	return DatabaseId(id), err
}

func buildAccount(record accountRecord) *Account {
	account := Account{
		Id:   DatabaseId(record.id),
		Name: record.name,
	}
	return &account
}

func (repo *AccountsRepo) FindByName(name string) (Account, error) {
	row := repo.Db.Conn.QueryRow(
		`
        SELECT id, name
		FROM accounts
		WHERE name = $1;
		`,
		name,
	)

	record := accountRecord{}
	err := row.Scan(
		&record.id,
		&record.name,
	)

	if err != nil {
		return Account{}, wrapRecordNotFound(err)
	}

	return *buildAccount(record), nil
}

func (repo *AccountsRepo) FindAll() ([]Account, error) {
	var accounts []Account

	rows, err := repo.Db.Conn.Query(
		`
        SELECT id, name FROM accounts;
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
		)

		accounts = append(accounts, *buildAccount(record))
	}

	if rows.Err() != nil {
		return nil, rows.Err()
	}

	return accounts, nil
}
