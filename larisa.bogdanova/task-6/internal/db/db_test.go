package db_test

import (
	"errors"
	"testing"

	"task-6/internal/db"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	ErrDBConnectionFailed      = errors.New("db connection failed")
	ErrNetworkFailureIteration = errors.New("network failure during iteration")
	ErrDBTimeout               = errors.New("db timeout")
	ErrConnectionReset         = errors.New("connection reset")
)

func TestNew(t *testing.T) {
	t.Parallel()

	mockDB, _, err := sqlmock.New()
	require.NoError(t, err)
	defer mockDB.Close()

	service := db.New(mockDB)

	assert.NotNil(t, service.DB, "DBService должен быть инициализирован с DB")
}

func TestGetNames_Success(t *testing.T) {
	t.Parallel()

	mockDB, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer mockDB.Close()

	service := db.New(mockDB)

	expectedNames := []string{"Alice", "Bob", "Alice"}
	rows := sqlmock.NewRows([]string{"name"}).
		AddRow("Alice").
		AddRow("Bob").
		AddRow("Alice")

	mock.ExpectQuery("SELECT name FROM users").WillReturnRows(rows)

	names, err := service.GetNames()

	require.NoError(t, err)
	assert.Equal(t, expectedNames, names)

	require.NoError(t, mock.ExpectationsWereMet())
}

func TestGetNames_DBQueryError(t *testing.T) {
	t.Parallel()

	mockDB, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer mockDB.Close()

	service := db.New(mockDB)
	expectedErr := ErrDBConnectionFailed

	mock.ExpectQuery("SELECT name FROM users").WillReturnError(expectedErr)

	names, err := service.GetNames()

	require.Error(t, err)
	assert.Contains(t, err.Error(), "db query: db connection failed")
	assert.Nil(t, names)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetNames_RowsScanError(t *testing.T) {
	t.Parallel()

	mockDB, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer mockDB.Close()

	service := db.New(mockDB)

	rows := sqlmock.NewRows([]string{"name"}).
		AddRow(nil)

	mock.ExpectQuery("SELECT name FROM users").WillReturnRows(rows)

	names, err := service.GetNames()

	require.Error(t, err)
	assert.Contains(t, err.Error(), "rows scanning:")
	assert.Contains(t, err.Error(), "NULL to string")
	assert.Nil(t, names)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetNames_RowsIterationError(t *testing.T) {
	t.Parallel()

	mockDB, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer mockDB.Close()

	service := db.New(mockDB)
	expectedErr := ErrNetworkFailureIteration

	rows := sqlmock.NewRows([]string{"name"}).
		AddRow("Alice").
		RowError(0, expectedErr)

	mock.ExpectQuery("SELECT name FROM users").WillReturnRows(rows)

	names, err := service.GetNames()

	require.Error(t, err)
	assert.Contains(t, err.Error(), "rows error: network failure during iteration")
	assert.Nil(t, names)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetUniqueNames_Success(t *testing.T) {
	t.Parallel()

	mockDB, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer mockDB.Close()

	service := db.New(mockDB)

	expectedNames := []string{"Alice", "Bob"}
	rows := sqlmock.NewRows([]string{"name"}).
		AddRow("Alice").
		AddRow("Bob")

	mock.ExpectQuery("SELECT DISTINCT name FROM users").WillReturnRows(rows)

	names, err := service.GetUniqueNames()

	require.NoError(t, err)
	assert.Equal(t, expectedNames, names)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetUniqueNames_DBQueryError(t *testing.T) {
	t.Parallel()

	mockDB, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer mockDB.Close()

	service := db.New(mockDB)
	expectedErr := ErrDBTimeout

	mock.ExpectQuery("SELECT DISTINCT name FROM users").WillReturnError(expectedErr)

	names, err := service.GetUniqueNames()

	require.Error(t, err)
	assert.Contains(t, err.Error(), "db query: db timeout")
	assert.Nil(t, names)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetUniqueNames_RowsScanError(t *testing.T) {
	t.Parallel()

	mockDB, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer mockDB.Close()

	service := db.New(mockDB)

	rows := sqlmock.NewRows([]string{"name"}).
		AddRow(nil)

	mock.ExpectQuery("SELECT DISTINCT name FROM users").WillReturnRows(rows)

	names, err := service.GetUniqueNames()

	require.Error(t, err)
	assert.Contains(t, err.Error(), "rows scanning:")
	assert.Contains(t, err.Error(), "NULL to string")
	assert.Nil(t, names)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetUniqueNames_RowsIterationError(t *testing.T) {
	t.Parallel()

	mockDB, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer mockDB.Close()

	service := db.New(mockDB)
	expectedErr := ErrConnectionReset

	rows := sqlmock.NewRows([]string{"name"}).
		AddRow("Jane").
		RowError(0, expectedErr)

	mock.ExpectQuery("SELECT DISTINCT name FROM users").WillReturnRows(rows)

	names, err := service.GetUniqueNames()

	require.Error(t, err)
	assert.Contains(t, err.Error(), "rows error: connection reset")
	assert.Nil(t, names)

	assert.NoError(t, mock.ExpectationsWereMet())
}
