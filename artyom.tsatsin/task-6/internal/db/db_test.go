package db_test

import (
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/require"
)

func newMockDB(t *testing.T) (*DBService, sqlmock.Sqlmock, func()) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)

	cleanup := func() {
		require.NoError(t, mock.ExpectationsWereMet())
		db.Close()
	}

	return &DBService{DB: db}, mock, cleanup
}

func TestGetNames_Success(t *testing.T) {
	service, mock, done := newMockDB(t)
	defer done()

	rows := sqlmock.NewRows([]string{"name"}).
		AddRow("Alice").
		AddRow("Bob")

	mock.ExpectQuery(`^SELECT name FROM users$`).
		WillReturnRows(rows)

	names, err := service.GetNames()
	require.NoError(t, err)
	require.Equal(t, []string{"Alice", "Bob"}, names)
}

func TestGetNames_EmptyResult(t *testing.T) {
	service, mock, done := newMockDB(t)
	defer done()

	rows := sqlmock.NewRows([]string{"name"})

	mock.ExpectQuery(`^SELECT name FROM users$`).
		WillReturnRows(rows)

	names, err := service.GetNames()
	require.NoError(t, err)
	require.Nil(t, names)
}

func TestGetNames_QueryFail(t *testing.T) {
	service, mock, done := newMockDB(t)
	defer done()

	mock.ExpectQuery(`^SELECT name FROM users$`).
		WillReturnError(errors.New("db down"))

	result, err := service.GetNames()
	require.Error(t, err)
	require.Nil(t, result)
}

func TestGetNames_ScanFail(t *testing.T) {
	service, mock, done := newMockDB(t)
	defer done()

	rows := sqlmock.NewRows([]string{"name"}).
		AddRow(123)

	mock.ExpectQuery(`^SELECT name FROM users$`).
		WillReturnRows(rows)

	result, err := service.GetNames()
	require.Error(t, err)
	require.Nil(t, result)
}

func TestGetNames_RowsFail(t *testing.T) {
	service, mock, done := newMockDB(t)
	defer done()

	rows := sqlmock.NewRows([]string{"name"}).
		AddRow("ok").
		AddRow("bad")
	rows.RowError(1, errors.New("iteration failed"))

	mock.ExpectQuery(`^SELECT name FROM users$`).
		WillReturnRows(rows)

	result, err := service.GetNames()
	require.Error(t, err)
	require.Nil(t, result)
}

func TestGetUniqueNames_Success(t *testing.T) {
	service, mock, done := newMockDB(t)
	defer done()

	rows := sqlmock.NewRows([]string{"name"}).
		AddRow("u1").
		AddRow("u2")

	mock.ExpectQuery(`^SELECT DISTINCT name FROM users$`).
		WillReturnRows(rows)

	names, err := service.GetUniqueNames()
	require.NoError(t, err)
	require.Equal(t, []string{"u1", "u2"}, names)
}

func TestGetUniqueNames_QueryFail(t *testing.T) {
	service, mock, done := newMockDB(t)
	defer done()

	mock.ExpectQuery(`^SELECT DISTINCT name FROM users$`).
		WillReturnError(errors.New("query error"))

	result, err := service.GetUniqueNames()
	require.Error(t, err)
	require.Nil(t, result)
}

func TestGetUniqueNames_ScanFail(t *testing.T) {
	service, mock, done := newMockDB(t)
	defer done()

	rows := sqlmock.NewRows([]string{"name"}).
		AddRow(nil)

	mock.ExpectQuery(`^SELECT DISTINCT name FROM users$`).
		WillReturnRows(rows)

	result, err := service.GetUniqueNames()
	require.Error(t, err)
	require.Nil(t, result)
}

func TestGetUniqueNames_RowsFail(t *testing.T) {
	service, mock, done := newMockDB(t)
	defer done()

	rows := sqlmock.NewRows([]string{"name"}).
		AddRow("ok").
		AddRow("bad")
	rows.RowError(1, errors.New("rows error"))

	mock.ExpectQuery(`^SELECT DISTINCT name FROM users$`).
		WillReturnRows(rows)

	result, err := service.GetUniqueNames()
	require.Error(t, err)
	require.Nil(t, result)
}
