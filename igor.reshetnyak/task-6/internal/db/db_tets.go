package db_test

import (
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	db "github.com/uchuaip/task-6/internal/db"
)

const (
	getAllQuery    = "SELECT name FROM users"
	getUniqueQuery = "SELECT DISTINCT name FROM users"
)

var testErr = errors.New("test error")

func TestNewDBService(t *testing.T) {
	t.Parallel()

	mockDB, _, err := sqlmock.New()
	require.NoError(t, err)
	defer mockDB.Close()

	service := db.New(mockDB)
	assert.Equal(t, mockDB, service.DB)
}

func TestGetNames_Success(t *testing.T) {
	t.Parallel()

	mockDB, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer mockDB.Close()

	service := db.DBService{DB: mockDB}

	rows := sqlmock.NewRows([]string{"name"}).
		AddRow("John").
		AddRow("Jane").
		AddRow("Bob")
	mock.ExpectQuery(getAllQuery).WillReturnRows(rows)

	names, err := service.GetNames()
	require.NoError(t, err)
	assert.Equal(t, []string{"John", "Jane", "Bob"}, names)
}

func TestGetNames_EmptyResult(t *testing.T) {
	t.Parallel()

	mockDB, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer mockDB.Close()

	service := db.DBService{DB: mockDB}

	rows := sqlmock.NewRows([]string{"name"})
	mock.ExpectQuery(getAllQuery).WillReturnRows(rows)

	names, err := service.GetNames()
	require.NoError(t, err)
	assert.Empty(t, names)
}

func TestGetNames_QueryError(t *testing.T) {
	t.Parallel()

	mockDB, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer mockDB.Close()

	service := db.DBService{DB: mockDB}

	mock.ExpectQuery(getAllQuery).WillReturnError(testErr)

	names, err := service.GetNames()
	require.ErrorContains(t, err, "db query")
	assert.Nil(t, names)
}

func TestGetNames_ScanError(t *testing.T) {
	t.Parallel()

	mockDB, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer mockDB.Close()

	service := db.DBService{DB: mockDB}

	rows := sqlmock.NewRows([]string{"name"}).AddRow(nil)
	mock.ExpectQuery(getAllQuery).WillReturnRows(rows)

	names, err := service.GetNames()
	require.ErrorContains(t, err, "rows scanning")
	assert.Nil(t, names)
}

func TestGetNames_RowsError(t *testing.T) {
	t.Parallel()

	mockDB, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer mockDB.Close()

	service := db.DBService{DB: mockDB}

	rows := sqlmock.NewRows([]string{"name"}).AddRow("Test")
	rows.RowError(0, testErr)
	mock.ExpectQuery(getAllQuery).WillReturnRows(rows)

	names, err := service.GetNames()
	require.ErrorContains(t, err, "rows error")
	assert.Nil(t, names)
}

func TestGetUniqueNames_Success(t *testing.T) {
	t.Parallel()

	mockDB, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer mockDB.Close()

	service := db.DBService{DB: mockDB}

	rows := sqlmock.NewRows([]string{"name"}).
		AddRow("Alice").
		AddRow("Bob").
		AddRow("Alice") // Дубликат, но DISTINCT
	mock.ExpectQuery(getUniqueQuery).WillReturnRows(rows)

	names, err := service.GetUniqueNames()
	require.NoError(t, err)
	assert.Equal(t, []string{"Alice", "Bob", "Alice"}, names)
}

func TestGetUniqueNames_NoRows(t *testing.T) {
	t.Parallel()

	mockDB, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer mockDB.Close()

	service := db.DBService{DB: mockDB}

	rows := sqlmock.NewRows([]string{"name"})
	mock.ExpectQuery(getUniqueQuery).WillReturnRows(rows)

	names, err := service.GetUniqueNames()
	require.NoError(t, err)
	assert.Empty(t, names)
}

func TestGetUniqueNames_QueryFailed(t *testing.T) {
	t.Parallel()

	mockDB, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer mockDB.Close()

	service := db.DBService{DB: mockDB}

	mock.ExpectQuery(getUniqueQuery).WillReturnError(testErr)

	names, err := service.GetUniqueNames()
	require.ErrorContains(t, err, "db query")
	assert.Nil(t, names)
}

func TestGetUniqueNames_ScanFailed(t *testing.T) {
	t.Parallel()

	mockDB, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer mockDB.Close()

	service := db.DBService{DB: mockDB}

	rows := sqlmock.NewRows([]string{"name"}).AddRow(nil)
	mock.ExpectQuery(getUniqueQuery).WillReturnRows(rows)

	names, err := service.GetUniqueNames()
	require.ErrorContains(t, err, "rows scanning")
	assert.Nil(t, names)
}
