package testutils

import (
	"context"
	"github.com/Tanzor-Disco/skaM/models"
	"github.com/jackc/pgx/v4/pgxpool"
	"testing"
)

type TestDB struct {
	pool *pgxpool.Pool
}

func Connect(URI string) (*TestDB, error) {
	pool, err := pgxpool.Connect(context.Background(), URI)
	if err != nil {
		return &TestDB{}, err
	}
	return &TestDB{
		pool: pool,
	}, err
}

func (tdb *TestDB) Close() {
	tdb.pool.Close()
}

func (tdb *TestDB) GetUsersByEmail(t *testing.T, email string) []models.User {
	t.Helper()
	rows, err := tdb.pool.Query(context.Background(),
		`
		SELECT * FROM users WHERE email = $1
	`,
		email)
	if err != nil {
		t.Fatal(err)
	}

	var users []models.User

	for rows.Next() {
		var user models.User
		err := rows.Scan(&user.Id, &user.Email, &user.Username, &user.PasswordHash, &user.LastYearActive)
		if err != nil {
			t.Fatal(err)
		}
		users = append(users, user)
	}
	return users
}

func (tdb *TestDB) GetPendingUserByEmail(t *testing.T, email string) models.PendingUser {
	t.Helper()
	rows, err := tdb.pool.Query(context.Background(),
		`
		SELECT * FROM pending_users WHERE email = $1
	`,
		email)
	if err != nil {
		t.Fatal(err)
	}

	var users []models.PendingUser

	for rows.Next() {
		var user models.PendingUser
		err := rows.Scan(&user.Id, &user.Email, &user.Username, &user.PasswordHash, &user.TokenHash, &user.ExpiresAt)
		if err != nil {
			t.Fatal(err)
		}
		users = append(users, user)
	}
	if len(users) != 1 {
		t.Fatalf("the amount of found users doesn't equal to one: %+v", users)
	}
	return users[0]

}

func (tdb *TestDB) DeleteUsersByEmail(t *testing.T, email string) {
	t.Helper()
	_, err := tdb.pool.Exec(context.Background(),
		`
	DELETE FROM users WHERE email = $1
	`, email)
	if err != nil {
		t.Fatal(err)
	}
}

func (tdb *TestDB) DeletePendingUsersByEmail(t *testing.T, email string) {
	t.Helper()
	_, err := tdb.pool.Exec(context.Background(),
		`
	DELETE FROM pending_users WHERE email = $1
	`, email)
	if err != nil {
		t.Fatal(err)
	}
}

func (tdb *TestDB) GetAllUsers(t *testing.T) (userRows []models.User) {
	rows, err := tdb.pool.Query(context.Background(),
		`
	SELECT * FROM users
	`)
	if err != nil {
		t.Fatal(err)
		return
	}
	for rows.Next() {
		var user models.User
		err := rows.Scan(&user.Id, &user.Email, &user.Username, &user.PasswordHash, &user.LastYearActive)
		if err != nil {
			t.Fatal(err)
		}
		userRows = append(userRows, user)
	}
	return
}

func (tdb *TestDB) GetAllRequests(t *testing.T) (users []models.PendingUser) {
	rows, err := tdb.pool.Query(context.Background(),
		`
	SELECT * FROM pending_users
	`)
	if err != nil {
		t.Fatal(err)
		return
	}
	for rows.Next() {
		var user models.PendingUser
		err := rows.Scan(&user.Id, &user.Email, &user.Username, &user.PasswordHash, &user.TokenHash, &user.ExpiresAt)
		if err != nil {
			t.Fatal(err)
		}
		users = append(users, user)
	}
	return
}
