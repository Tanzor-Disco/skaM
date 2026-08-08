package testutils

import (
	"github.com/jackc/pgx/v4/pgxpool"
	"context"
	"testing"
)

type UserRow struct {
	Id int
	Email string
	Username string
	PasswordHash string
	LastYearActive int
}

type TestDB struct {
	pool *pgxpool.Pool
}

func Connect(URI string) (*TestDB,error) {
	pool,err := pgxpool.Connect(context.Background(),URI)
	if err != nil {
		return &TestDB{},err
	}
	return &TestDB {
		pool:pool,
	},err
}

func (tdb *TestDB) Close() {
	tdb.pool.Close()
}

func (tdb *TestDB) GetUsersByEmail(t *testing.T, email string) []UserRow {
	t.Helper()
	rows,err := tdb.pool.Query(context.Background(),
	`
		SELECT * FROM users WHERE email = $1
	`,
	email)
	if err != nil {
		t.Fatal(err)
	}

	var users []UserRow

	for rows.Next() {
		var user UserRow
		err := rows.Scan(&user.Id,&user.Email,&user.Username,&user.PasswordHash,&user.LastYearActive)
		if err != nil {
			t.Fatal(err)
		}
		users = append(users,user)
	}
	return users

}

func (tdb *TestDB) DeleteUsersByEmail(t *testing.T, email string) {
	t.Helper()
	_,err := tdb.pool.Exec(context.Background(),
	`
	DELETE FROM users WHERE email = $1
	`,email)
	if err != nil {
		t.Fatal(err)
	}
}


func (tdb *TestDB) GetAllUsers(t *testing.T) (userRows []UserRow) {
	rows,err := tdb.pool.Query(context.Background(),
	`
	SELECT * FROM users
	`)
	if err != nil {
		t.Fatal(err)
		return 
	}
	for rows.Next() {
		var user UserRow
		err := rows.Scan(&user.Id,&user.Email,&user.Username,&user.PasswordHash,&user.LastYearActive)
		if err != nil {
			t.Fatal(err)
		}
		userRows = append(userRows,user)
	}
	return 
}
