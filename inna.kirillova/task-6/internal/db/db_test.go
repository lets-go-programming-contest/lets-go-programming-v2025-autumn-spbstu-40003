package db

import (
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNew(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	service := New(db)

	assert.NotNil(t, service)
	assert.Equal(t, db, service.DB)

	require.NoError(t, mock.ExpectationsWereMet())
}

func TestDBService_GetNames(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer db.Close()

		mock.ExpectQuery(`SELECT name FROM users`).
			WillReturnRows(sqlmock.NewRows([]string{"name"}).AddRow("a").AddRow("b"))

		svc := New(db)
		names, err := svc.GetNames()

		require.NoError(t, err)
		assert.Equal(t, []string{"a", "b"}, names)
		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("empty", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer db.Close()

		mock.ExpectQuery(`SELECT name FROM users`).
			WillReturnRows(sqlmock.NewRows([]string{"name"}))

		svc := New(db)
		names, err := svc.GetNames()

		require.NoError(t, err)
		assert.Nil(t, names)
		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("query error", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer db.Close()

		mock.ExpectQuery(`SELECT name FROM users`).
			WillReturnError(errors.New("err"))

		svc := New(db)
		names, err := svc.GetNames()

		require.Error(t, err)
		assert.Contains(t, err.Error(), "db query:")
		assert.Nil(t, names)
		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("scan error", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer db.Close()

		mock.ExpectQuery(`SELECT name FROM users`).
			WillReturnRows(sqlmock.NewRows([]string{"name"}).AddRow(nil))

		svc := New(db)
		names, err := svc.GetNames()

		require.Error(t, err)
		assert.Contains(t, err.Error(), "rows scanning:")
		assert.Nil(t, names)
		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("rows error", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer db.Close()

		rows := sqlmock.NewRows([]string{"name"}).AddRow("x")
		rows.RowError(0, errors.New("row err"))
		mock.ExpectQuery(`SELECT name FROM users`).WillReturnRows(rows)

		svc := New(db)
		names, err := svc.GetNames()

		require.Error(t, err)
		assert.Contains(t, err.Error(), "rows error:")
		assert.Nil(t, names)
		require.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestDBService_GetUniqueNames(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer db.Close()

		mock.ExpectQuery(`SELECT DISTINCT name FROM users`).
			WillReturnRows(sqlmock.NewRows([]string{"name"}).AddRow("x").AddRow("y"))

		svc := New(db)
		names, err := svc.GetUniqueNames()

		require.NoError(t, err)
		assert.Equal(t, []string{"x", "y"}, names)
		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("empty", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer db.Close()

		mock.ExpectQuery(`SELECT DISTINCT name FROM users`).
			WillReturnRows(sqlmock.NewRows([]string{"name"}))

		svc := New(db)
		names, err := svc.GetUniqueNames()

		require.NoError(t, err)
		assert.Nil(t, names)
		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("query error", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer db.Close()

		mock.ExpectQuery(`SELECT DISTINCT name FROM users`).
			WillReturnError(errors.New("err"))

		svc := New(db)
		names, err := svc.GetUniqueNames()

		require.Error(t, err)
		assert.Contains(t, err.Error(), "db query:")
		assert.Nil(t, names)
		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("scan error", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer db.Close()

		mock.ExpectQuery(`SELECT DISTINCT name FROM users`).
			WillReturnRows(sqlmock.NewRows([]string{"name"}).AddRow(nil))

		svc := New(db)
		names, err := svc.GetUniqueNames()

		require.Error(t, err)
		assert.Contains(t, err.Error(), "rows scanning:")
		assert.Nil(t, names)
		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("rows error", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer db.Close()

		rows := sqlmock.NewRows([]string{"name"}).AddRow("z")
		rows.RowError(0, errors.New("row err"))
		mock.ExpectQuery(`SELECT DISTINCT name FROM users`).WillReturnRows(rows)

		svc := New(db)
		names, err := svc.GetUniqueNames()

		require.Error(t, err)
		assert.Contains(t, err.Error(), "rows error:")
		assert.Nil(t, names)
		require.NoError(t, mock.ExpectationsWereMet())
	})
}
