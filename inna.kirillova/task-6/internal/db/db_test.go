package db

import (
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNew(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()
	s := New(db)
	assert.NotNil(t, s)
	assert.Equal(t, db, s.DB)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestGetNamesSuccess(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()
	mock.ExpectQuery(`SELECT name FROM users`).
		WillReturnRows(sqlmock.NewRows([]string{"name"}).AddRow("Alex").AddRow("Sam"))
	s := New(db)
	names, err := s.GetNames()
	require.NoError(t, err)
	assert.Equal(t, []string{"Alex", "Sam"}, names)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestGetNamesEmpty(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()
	mock.ExpectQuery(`SELECT name FROM users`).
		WillReturnRows(sqlmock.NewRows([]string{"name"}))
	s := New(db)
	names, err := s.GetNames()
	require.NoError(t, err)
	assert.Nil(t, names)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestGetNamesError(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()
	mock.ExpectQuery(`SELECT name FROM users`).
		WillReturnError(errors.New("db error"))
	s := New(db)
	names, err := s.GetNames()
	require.Error(t, err)
	assert.Contains(t, err.Error(), "db query:")
	assert.Nil(t, names)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestGetNamesScanError(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()
	mock.ExpectQuery(`SELECT name FROM users`).
		WillReturnRows(sqlmock.NewRows([]string{"name"}).AddRow(nil))
	s := New(db)
	names, err := s.GetNames()
	require.Error(t, err)
	assert.Contains(t, err.Error(), "rows scanning:")
	assert.Nil(t, names)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestGetNamesRowsError(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()
	rows := sqlmock.NewRows([]string{"name"}).AddRow("Alex")
	rows.RowError(0, errors.New("row error"))
	mock.ExpectQuery(`SELECT name FROM users`).WillReturnRows(rows)
	s := New(db)
	names, err := s.GetNames()
	require.Error(t, err)
	assert.Contains(t, err.Error(), "rows error:")
	assert.Nil(t, names)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestGetUniqueNamesSuccess(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()
	mock.ExpectQuery(`SELECT DISTINCT name FROM users`).
		WillReturnRows(sqlmock.NewRows([]string{"name"}).AddRow("Tom").AddRow("Anna"))
	s := New(db)
	names, err := s.GetUniqueNames()
	require.NoError(t, err)
	assert.Equal(t, []string{"Tom", "Anna"}, names)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestGetUniqueNamesEmpty(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()
	mock.ExpectQuery(`SELECT DISTINCT name FROM users`).
		WillReturnRows(sqlmock.NewRows([]string{"name"}))
	s := New(db)
	names, err := s.GetUniqueNames()
	require.NoError(t, err)
	assert.Nil(t, names)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestGetUniqueNamesError(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()
	mock.ExpectQuery(`SELECT DISTINCT name FROM users`).
		WillReturnError(errors.New("timeout"))
	s := New(db)
	names, err := s.GetUniqueNames()
	require.Error(t, err)
	assert.Contains(t, err.Error(), "db query:")
	assert.Nil(t, names)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestGetUniqueNamesScanError(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()
	mock.ExpectQuery(`SELECT DISTINCT name FROM users`).
		WillReturnRows(sqlmock.NewRows([]string{"name"}).AddRow(nil))
	s := New(db)
	names, err := s.GetUniqueNames()
	require.Error(t, err)
	assert.Contains(t, err.Error(), "rows scanning:")
	assert.Nil(t, names)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestGetUniqueNamesRowsError(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()
	rows := sqlmock.NewRows([]string{"name"}).AddRow("Tom")
	rows.RowError(0, errors.New("row error"))
	mock.ExpectQuery(`SELECT DISTINCT name FROM users`).WillReturnRows(rows)
	s := New(db)
	names, err := s.GetUniqueNames()
	require.Error(t, err)
	assert.Contains(t, err.Error(), "rows error:")
	assert.Nil(t, names)
	require.NoError(t, mock.ExpectationsWereMet())
}
