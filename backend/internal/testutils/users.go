package testutils

import (
	"context"
	"github.com/Tanzor-Disco/skaM/models"
	"testing"
)

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

func (tdb *TestDB) GetUserIDByEmail(t *testing.T, email string) models.UserID {
	t.Helper()
	var id models.UserID
	err := tdb.pool.QueryRow(context.Background(),
		`
	SELECT id FROM users WHERE email = $1
	`,
		email).Scan(&id)
	if err != nil {
		t.Fatalf("getUserIDByEmail: %v", err)
	}
	return id
}
