package db

import (
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/require"
)

func TestDBService_GetNames_OK(t *testing.T) {
	conn, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer conn.Close()

	rows := sqlmock.NewRows([]string{"name"}).
		AddRow("alice").
		AddRow("bob")

	mock.ExpectQuery("^SELECT name FROM users$").WillReturnRows(rows)

	service := New(conn)
	got, err := service.GetNames()
	require.NoError(t, err)
	require.Equal(t, []string{"alice", "bob"}, got)

	require.NoError(t, mock.ExpectationsWereMet())
}

func TestDBService_GetNames_QueryError(t *testing.T) {
	conn, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer conn.Close()

	qErr := errors.New("db down")
	mock.ExpectQuery("^SELECT name FROM users$").WillReturnError(qErr)

	service := New(conn)
	got, err := service.GetNames()
	require.Nil(t, got)
	require.Error(t, err)
	require.ErrorContains(t, err, "db query:")
	require.ErrorIs(t, err, qErr)

	require.NoError(t, mock.ExpectationsWereMet())
}

func TestDBService_GetNames_ScanError(t *testing.T) {
	conn, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer conn.Close()

	rows := sqlmock.NewRows([]string{"name"}).AddRow(nil)
	mock.ExpectQuery("^SELECT name FROM users$").WillReturnRows(rows)

	service := New(conn)
	got, err := service.GetNames()
	require.Nil(t, got)
	require.Error(t, err)
	require.ErrorContains(t, err, "rows scanning:")

	require.NoError(t, mock.ExpectationsWereMet())
}

func TestDBService_GetNames_RowsError(t *testing.T) {
	conn, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer conn.Close()

	rowErr := errors.New("row error")
	rows := sqlmock.NewRows([]string{"name"}).
		AddRow("alice").
		RowError(0, rowErr)

	mock.ExpectQuery("^SELECT name FROM users$").WillReturnRows(rows)

	service := New(conn)
	got, err := service.GetNames()
	require.Nil(t, got)
	require.Error(t, err)
	require.ErrorContains(t, err, "rows error:")
	require.ErrorIs(t, err, rowErr)

	require.NoError(t, mock.ExpectationsWereMet())
}

func TestDBService_GetUniqueNames_OK(t *testing.T) {
	conn, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer conn.Close()

	rows := sqlmock.NewRows([]string{"name"}).
		AddRow("alice").
		AddRow("bob")

	mock.ExpectQuery("^SELECT DISTINCT name FROM users$").WillReturnRows(rows)

	service := New(conn)
	got, err := service.GetUniqueNames()
	require.NoError(t, err)
	require.Equal(t, []string{"alice", "bob"}, got)

	require.NoError(t, mock.ExpectationsWereMet())
}

func TestDBService_GetUniqueNames_QueryError(t *testing.T) {
	conn, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer conn.Close()

	qErr := errors.New("db down")
	mock.ExpectQuery("^SELECT DISTINCT name FROM users$").WillReturnError(qErr)

	service := New(conn)
	got, err := service.GetUniqueNames()
	require.Nil(t, got)
	require.Error(t, err)
	require.ErrorContains(t, err, "db query:")
	require.ErrorIs(t, err, qErr)

	require.NoError(t, mock.ExpectationsWereMet())
}

func TestDBService_GetUniqueNames_ScanError(t *testing.T) {
	conn, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer conn.Close()

	rows := sqlmock.NewRows([]string{"name"}).AddRow(nil)
	mock.ExpectQuery("^SELECT DISTINCT name FROM users$").WillReturnRows(rows)

	service := New(conn)
	got, err := service.GetUniqueNames()
	require.Nil(t, got)
	require.Error(t, err)
	require.ErrorContains(t, err, "rows scanning:")

	require.NoError(t, mock.ExpectationsWereMet())
}

func TestDBService_GetUniqueNames_RowsError(t *testing.T) {
	conn, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer conn.Close()

	rowErr := errors.New("row error")
	rows := sqlmock.NewRows([]string{"name"}).
		AddRow("alice").
		RowError(0, rowErr)

	mock.ExpectQuery("^SELECT DISTINCT name FROM users$").WillReturnRows(rows)

	service := New(conn)
	got, err := service.GetUniqueNames()
	require.Nil(t, got)
	require.Error(t, err)
	require.ErrorContains(t, err, "rows error:")
	require.ErrorIs(t, err, rowErr)

	require.NoError(t, mock.ExpectationsWereMet())
}
