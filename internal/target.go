package internal

type Target struct {
	Id            DatabaseId
	AccountId     DatabaseId
	Server        string
	Username      string
	UsernameNonce string
	Password      string
	PasswordNonce string
	Dir           string
}

type TargetRepo struct {
	Db DbEngine
}

/*
func (repo *TargetRepo) Save(target Target) (DatabaseId, error) {
	var query string
	var id DatabaseId
	var err error
	query = "INSERT INTO target (id, account_id, server, dir, username, username_nonce, password, password_nonce) " +
		"VALUES ($1, $2, $3, $4, $5, $6, $7, $8);"
	if id, err = NewDatabaseId(); err != nil {
		return id, err
	}
	_, err = repo.Db.Conn.Exec(
		query,
		id, target.AccountId, target.Server, target.Dir, username, username_nonce, password, password_nonce,
	)
	if err != nil {
		return NewUnsetDatabaseId(), err
	}
	return id, err
}
*/
