package testutils

import (
	"context"
	"testing"
)

func (tdb *TestDB) DeleteSessionsByUserID(t *testing.T, id int) {
	t.Helper()
	_, err := tdb.pool.Exec(context.Background(),
		`
	DELETE FROM user_sessions WHERE user_id = $1
	`,
		id)
	if err != nil {
		t.Fatalf("DeleteSessionsByUserID: %v", err)
	}
}
