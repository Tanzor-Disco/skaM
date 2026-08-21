package testutils

import (
	"context"
	"github.com/Tanzor-Disco/skaM/models"
	"testing"
)

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

func (tdb *TestDB) GetAllPendingUsers(t *testing.T) (users []models.PendingUser) {
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
