package database

import (
	"os"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMySQLConnQueryTimesOut(t *testing.T) {
	t.Setenv("DB_QUERY_TIMEOUT_MS", "1")

	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	conn := NewMySQLConn(db)
	mock.ExpectQuery("SELECT 1").WillDelayFor(50 * time.Millisecond).WillReturnRows(sqlmock.NewRows([]string{"v"}).AddRow(1))

	_, err = conn.Query("SELECT 1")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "failed to execute query")
}

func TestMySQLConnQuerySuccess(t *testing.T) {
	t.Setenv("DB_QUERY_TIMEOUT_MS", "1000")

	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	conn := NewMySQLConn(db)
	mock.ExpectQuery("SELECT 1").WillReturnRows(sqlmock.NewRows([]string{"v"}).AddRow(1))

	rows, err := conn.Query("SELECT 1")
	require.NoError(t, err)
	require.NotNil(t, rows)
	defer rows.Close()

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestReadDBTimeoutFallsBack(t *testing.T) {
	_ = os.Setenv("DB_QUERY_TIMEOUT_MS", "bad")
	defer os.Unsetenv("DB_QUERY_TIMEOUT_MS")

	value := readDBTimeout()
	assert.Equal(t, defaultDBTimeout, value)
}

func TestMySQLConnQueryRowScanSuccess(t *testing.T) {
	t.Setenv("DB_QUERY_TIMEOUT_MS", "1000")

	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	conn := NewMySQLConn(db)
	mock.ExpectQuery("SELECT name FROM workspaces WHERE id = \\?").
		WithArgs(1).
		WillReturnRows(sqlmock.NewRows([]string{"name"}).AddRow("acme"))

	row := conn.QueryRow("SELECT name FROM workspaces WHERE id = ?", 1)
	var name string
	err = row.Scan(&name)
	require.NoError(t, err)
	assert.Equal(t, "acme", name)
	assert.NoError(t, mock.ExpectationsWereMet())
}
