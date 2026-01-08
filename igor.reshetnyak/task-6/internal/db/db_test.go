package db_test

import (
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	db "github.com/ReshetnyakIgor/task-6/internal/db"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	getAllQuery    = "SELECT name FROM users"
	getUniqueQuery = "SELECT DISTINCT name FROM users"
)

var (
	errTest          = errors.New("test error")
	errRowsIteration = errors.New("rows iteration error")
)

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

	service := db.New(mockDB)

	rows := sqlmock.NewRows([]string{"name"}).
		AddRow("John").
		AddRow("Jane").
		AddRow("Bob")
	mock.ExpectQuery(getAllQuery).WillReturnRows(rows)

	names, err := service.GetNames()
	require.NoError(t, err)
	assert.Equal(t, []string{"John", "Jane", "Bob"}, names)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetNames_EmptyResult(t *testing.T) {
	t.Parallel()

	mockDB, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer mockDB.Close()

	service := db.New(mockDB)

	rows := sqlmock.NewRows([]string{"name"})
	mock.ExpectQuery(getAllQuery).WillReturnRows(rows)

	names, err := service.GetNames()
	require.NoError(t, err)
	assert.Empty(t, names)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetNames_QueryError(t *testing.T) {
	t.Parallel()

	mockDB, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer mockDB.Close()

	service := db.New(mockDB)

	mock.ExpectQuery(getAllQuery).WillReturnError(errTest)

	names, err := service.GetNames()
	require.ErrorContains(t, err, "db query")
	assert.Nil(t, names)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetNames_ScanError(t *testing.T) {
	t.Parallel()

	mockDB, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer mockDB.Close()

	service := db.New(mockDB)

	rows := sqlmock.NewRows([]string{"name"}).AddRow(nil)
	mock.ExpectQuery(getAllQuery).WillReturnRows(rows)

	names, err := service.GetNames()
	require.ErrorContains(t, err, "rows scanning")
	assert.Nil(t, names)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetNames_RowsNextError(t *testing.T) {
	t.Parallel()

	mockDB, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer mockDB.Close()

	service := db.New(mockDB)

	rows := sqlmock.NewRows([]string{"name"}).AddRow("Test")
	rows.RowError(0, errTest)
	mock.ExpectQuery(getAllQuery).WillReturnRows(rows)

	names, err := service.GetNames()
	require.ErrorContains(t, err, "rows error")
	assert.Nil(t, names)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetNames_RowsErrAfterLoop(t *testing.T) {
	t.Parallel()

	mockDB, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer mockDB.Close()

	service := db.New(mockDB)

	rows := sqlmock.NewRows([]string{"name"}).
		AddRow("Test1").
		AddRow("Test2").
		CloseError(errRowsIteration)

	mock.ExpectQuery(getAllQuery).WillReturnRows(rows)

	names, err := service.GetNames()
	require.Error(t, err)
	require.ErrorContains(t, err, "rows error")
	assert.Nil(t, names)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetUniqueNames_Success(t *testing.T) {
	t.Parallel()

	mockDB, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer mockDB.Close()

	service := db.New(mockDB)

	rows := sqlmock.NewRows([]string{"name"}).
		AddRow("Alice").
		AddRow("Bob").
		AddRow("Alice")
	mock.ExpectQuery(getUniqueQuery).WillReturnRows(rows)

	names, err := service.GetUniqueNames()
	require.NoError(t, err)
	assert.Equal(t, []string{"Alice", "Bob", "Alice"}, names)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetUniqueNames_NoRows(t *testing.T) {
	t.Parallel()

	mockDB, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer mockDB.Close()

	service := db.New(mockDB)

	rows := sqlmock.NewRows([]string{"name"})
	mock.ExpectQuery(getUniqueQuery).WillReturnRows(rows)

	names, err := service.GetUniqueNames()
	require.NoError(t, err)
	assert.Empty(t, names)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetUniqueNames_QueryFailed(t *testing.T) {
	t.Parallel()

	mockDB, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer mockDB.Close()

	service := db.New(mockDB)

	mock.ExpectQuery(getUniqueQuery).WillReturnError(errTest)

	names, err := service.GetUniqueNames()
	require.ErrorContains(t, err, "db query")
	assert.Nil(t, names)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetUniqueNames_ScanFailed(t *testing.T) {
	t.Parallel()

	mockDB, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer mockDB.Close()

	service := db.New(mockDB)

	rows := sqlmock.NewRows([]string{"name"}).AddRow(nil)
	mock.ExpectQuery(getUniqueQuery).WillReturnRows(rows)

	names, err := service.GetUniqueNames()
	require.ErrorContains(t, err, "rows scanning")
	assert.Nil(t, names)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetUniqueNames_RowsNextError(t *testing.T) {
	t.Parallel()

	mockDB, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer mockDB.Close()

	service := db.New(mockDB)

	rows := sqlmock.NewRows([]string{"name"}).AddRow("Test")
	rows.RowError(0, errTest)
	mock.ExpectQuery(getUniqueQuery).WillReturnRows(rows)

	names, err := service.GetUniqueNames()
	require.ErrorContains(t, err, "rows error")
	assert.Nil(t, names)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetUniqueNames_RowsErrAfterLoop(t *testing.T) {
	t.Parallel()

	mockDB, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer mockDB.Close()

	service := db.New(mockDB)

	rows := sqlmock.NewRows([]string{"name"}).
		AddRow("Test1").
		AddRow("Test2").
		CloseError(errRowsIteration)

	mock.ExpectQuery(getUniqueQuery).WillReturnRows(rows)

	names, err := service.GetUniqueNames()
	require.Error(t, err)
	require.ErrorContains(t, err, "rows error")
	assert.Nil(t, names)

	assert.NoError(t, mock.ExpectationsWereMet())
}
