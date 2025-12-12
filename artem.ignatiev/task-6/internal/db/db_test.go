package db

import (
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestGetNames(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create mock: %v", err)
	}
	defer db.Close()

	service := New(db)

	t.Run("successful get names", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"name"}).
			AddRow("Ivan").
			AddRow("Maria").
			AddRow("Petr")

		mock.ExpectQuery("SELECT name FROM users").WillReturnRows(rows)

		names, err := service.GetNames()

		assert.NoError(t, err)
		assert.Equal(t, []string{"Ivan", "Maria", "Petr"}, names)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("empty result", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"name"})

		mock.ExpectQuery("SELECT name FROM users").WillReturnRows(rows)

		names, err := service.GetNames()

		assert.NoError(t, err)
		assert.Empty(t, names)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("query error", func(t *testing.T) {
		mock.ExpectQuery("SELECT name FROM users").
			WillReturnError(errors.New("database error"))

		names, err := service.GetNames()

		assert.Error(t, err)
		assert.Nil(t, names)
		assert.Contains(t, err.Error(), "db query")
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("scan error", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"name"}).
			AddRow("Ivan").
			AddRow(nil).
			RowError(1, errors.New("scan error"))

		mock.ExpectQuery("SELECT name FROM users").WillReturnRows(rows)

		names, err := service.GetNames()

		assert.Error(t, err)
		assert.Nil(t, names)
		assert.Contains(t, err.Error(), "rows scanning")
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestGetUniqueNames(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create mock: %v", err)
	}
	defer db.Close()

	service := New(db)

	t.Run("successful get unique names", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"name"}).
			AddRow("Ivan").
			AddRow("Maria")

		mock.ExpectQuery("SELECT DISTINCT name FROM users").WillReturnRows(rows)

		names, err := service.GetUniqueNames()

		assert.NoError(t, err)
		assert.Equal(t, []string{"Ivan", "Maria"}, names)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("empty result", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"name"})

		mock.ExpectQuery("SELECT DISTINCT name FROM users").WillReturnRows(rows)

		names, err := service.GetUniqueNames()

		assert.NoError(t, err)
		assert.Empty(t, names)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("query error", func(t *testing.T) {
		mock.ExpectQuery("SELECT DISTINCT name FROM users").
			WillReturnError(errors.New("database error"))

		names, err := service.GetUniqueNames()

		assert.Error(t, err)
		assert.Nil(t, names)
		assert.Contains(t, err.Error(), "db query")
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}
