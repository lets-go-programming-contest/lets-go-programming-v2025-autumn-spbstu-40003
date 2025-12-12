package db_test

import (
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/kryjkaqq/task-6/internal/db"
	"github.com/stretchr/testify/require"
)

var (
	errDatabase     = errors.New("database error")
	errRowIteration = errors.New("row iteration error")
)

func TestGetNames(t *testing.T) {
	t.Parallel()

	database, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create mock: %v", err)
	}
	defer database.Close()

	service := db.New(database)

	t.Run("successful get names", func(t *testing.T) {
		t.Parallel()

		rows := sqlmock.NewRows([]string{"name"}).
			AddRow("Ivan").
			AddRow("Maria").
			AddRow("Petr")

		mock.ExpectQuery("SELECT name FROM users").WillReturnRows(rows)

		names, err := service.GetNames()

		require.NoError(t, err)
		require.Equal(t, []string{"Ivan", "Maria", "Petr"}, names)
		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("empty result", func(t *testing.T) {
		t.Parallel()

		rows := sqlmock.NewRows([]string{"name"})

		mock.ExpectQuery("SELECT name FROM users").WillReturnRows(rows)

		names, err := service.GetNames()

		require.NoError(t, err)
		require.Empty(t, names)
		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("query error", func(t *testing.T) {
		t.Parallel()

		mock.ExpectQuery("SELECT name FROM users").
			WillReturnError(errDatabase)

		names, err := service.GetNames()

		require.Error(t, err)
		require.Nil(t, names)
		require.Contains(t, err.Error(), "db query")
		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("scan error", func(t *testing.T) {
		t.Parallel()

		rows := sqlmock.NewRows([]string{"name", "extra"}).
			AddRow("Ivan", "extra")

		mock.ExpectQuery("SELECT name FROM users").WillReturnRows(rows)

		names, err := service.GetNames()

		require.Error(t, err)
		require.Nil(t, names)
		require.Contains(t, err.Error(), "rows scanning")
		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("rows error", func(t *testing.T) {
		t.Parallel()

		rows := sqlmock.NewRows([]string{"name"}).
			AddRow("Ivan").
			AddRow("Maria").
			RowError(1, errRowIteration)

		mock.ExpectQuery("SELECT name FROM users").WillReturnRows(rows)

		names, err := service.GetNames()

		require.Error(t, err)
		require.Nil(t, names)
		require.Contains(t, err.Error(), "rows error")
		require.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestGetUniqueNames(t *testing.T) {
	t.Parallel()

	database, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create mock: %v", err)
	}
	defer database.Close()

	service := db.New(database)

	t.Run("successful get unique names", func(t *testing.T) {
		t.Parallel()

		rows := sqlmock.NewRows([]string{"name"}).
			AddRow("Ivan").
			AddRow("Maria")

		mock.ExpectQuery("SELECT DISTINCT name FROM users").WillReturnRows(rows)

		names, err := service.GetUniqueNames()

		require.NoError(t, err)
		require.Equal(t, []string{"Ivan", "Maria"}, names)
		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("empty result", func(t *testing.T) {
		t.Parallel()

		rows := sqlmock.NewRows([]string{"name"})

		mock.ExpectQuery("SELECT DISTINCT name FROM users").WillReturnRows(rows)

		names, err := service.GetUniqueNames()

		require.NoError(t, err)
		require.Empty(t, names)
		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("query error", func(t *testing.T) {
		t.Parallel()

		mock.ExpectQuery("SELECT DISTINCT name FROM users").
			WillReturnError(errDatabase)

		names, err := service.GetUniqueNames()

		require.Error(t, err)
		require.Nil(t, names)
		require.Contains(t, err.Error(), "db query")
		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("scan error", func(t *testing.T) {
		t.Parallel()

		rows := sqlmock.NewRows([]string{"name", "extra"}).
			AddRow("Ivan", "extra")

		mock.ExpectQuery("SELECT DISTINCT name FROM users").WillReturnRows(rows)

		names, err := service.GetUniqueNames()

		require.Error(t, err)
		require.Nil(t, names)
		require.Contains(t, err.Error(), "rows scanning")
		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("rows error", func(t *testing.T) {
		t.Parallel()

		rows := sqlmock.NewRows([]string{"name"}).
			AddRow("Ivan").
			AddRow("Maria").
			RowError(1, errRowIteration)

		mock.ExpectQuery("SELECT DISTINCT name FROM users").WillReturnRows(rows)

		names, err := service.GetUniqueNames()

		require.Error(t, err)
		require.Nil(t, names)
		require.Contains(t, err.Error(), "rows error")
		require.NoError(t, mock.ExpectationsWereMet())
	})
}
