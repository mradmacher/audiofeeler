package audiofeeler

type Account struct {
	Id   DatabaseId
	Name string
}

type accountRecord struct {
	id   string
	name string
}

type AccountsRepo struct {
	Db *DbClient
}

func (repo *AccountsRepo) Save(account Account) (DatabaseId, error) {
	var query string
	var id DatabaseId
	var err error
	if IsDatabaseIdSet(account.Id) {
		id = account.Id
	} else {
		id = account.Id
		query = "INSERT INTO account (id, name) VALUES ($1, $2)"
		if id, err = NewDatabaseId(); err != nil {
			return id, err
		}
	}
	if _, err = repo.Db.Conn.Exec(query, id, account.Name); err != nil {
		return NewUnsetDatabaseId(), err
	}
	return id, nil
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
		FROM account
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
        SELECT id, name FROM account;
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
