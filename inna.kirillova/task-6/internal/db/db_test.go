package db

import (
	"errors"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestNew(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()
	s := New(db)
	assert.NotNil(t, s)
	assert.Equal(t, db, s.DB)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestGetNames_Success(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()
	mock.ExpectQuery(`SELECT name FROM users`).
		WillReturnRows(sqlmock.NewRows([]string{"name"}).AddRow("A"))
	s := New(db)
	names, err := s.GetNames()
	require.NoError(t, err)
	assert.Equal(t, []string{"A"}, names)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestGetNames_Empty(t *testing.T) {
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

func TestGetNames_QueryError(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()
	mock.ExpectQuery(`SELECT name FROM users`).
		WillReturnError(errors.New("err"))
	s := New(db)
	names, err := s.GetNames()
	require.Error(t, err)
	assert.Contains(t, err.Error(), "db query:")
	assert.Nil(t, names)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestGetNames_ScanError(t *testing.T) {
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

func TestGetNames_RowsError(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()
	rows := sqlmock.NewRows([]string{"name"}).AddRow("X")
	rows.RowError(0, errors.New("row err"))
	mock.ExpectQuery(`SELECT name FROM users`).WillReturnRows(rows)
	s := New(db)
	names, err := s.GetNames()
	require.Error(t, err)
	assert.Contains(t, err.Error(), "rows error:")
	assert.Nil(t, names)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestGetUniqueNames_Success(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()
	mock.ExpectQuery(`SELECT DISTINCT name FROM users`).
		WillReturnRows(sqlmock.NewRows([]string{"name"}).AddRow("B"))
	s := New(db)
	names, err := s.GetUniqueNames()
	require.NoError(t, err)
	assert.Equal(t, []string{"B"}, names)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestGetUniqueNames_Empty(t *testing.T) {
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

func TestGetUniqueNames_QueryError(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()
	mock.ExpectQuery(`SELECT DISTINCT name FROM users`).
		WillReturnError(errors.New("err"))
	s := New(db)
	names, err := s.GetUniqueNames()
	require.Error(t, err)
	assert.Contains(t, err.Error(), "db query:")
	assert.Nil(t, names)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestGetUniqueNames_ScanError(t *testing.T) {
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

func TestGetUniqueNames_RowsError(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()
	rows := sqlmock.NewRows([]string{"name"}).AddRow("Y")
	rows.RowError(0, errors.New("row err"))
	mock.ExpectQuery(`SELECT DISTINCT name FROM users`).WillReturnRows(rows)
	s := New(db)
	names, err := s.GetUniqueNames()
	require.Error(t, err)
	assert.Contains(t, err.Error(), "rows error:")
	assert.Nil(t, names)
	require.NoError(t, mock.ExpectationsWereMet())
}
