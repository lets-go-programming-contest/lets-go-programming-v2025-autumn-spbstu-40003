package db

import (
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func newTestDB(t *testing.T) (DBService, sqlmock.Sqlmock, func()) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	service := New(db)

	return service, mock, func() {
		db.Close()
	}
}

func TestDBService_GetNames(t *testing.T) {
	query := "SELECT name FROM users"

	t.Run("success", func(t *testing.T) {
		service, mock, teardown := newTestDB(t)
		defer teardown()

		rows := sqlmock.NewRows([]string{"name"}).
			AddRow("Alice").
			AddRow("Bob")

		mock.ExpectQuery(query).WillReturnRows(rows)

		names, err := service.GetNames()

		assert.NoError(t, err)
		assert.Equal(t, []string{"Alice", "Bob"}, names)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("query error", func(t *testing.T) {
		service, mock, teardown := newTestDB(t)
		defer teardown()

		mock.ExpectQuery(query).WillReturnError(errors.New("connection refused"))

		names, err := service.GetNames()

		assert.Error(t, err)
		assert.Nil(t, names)
		assert.Contains(t, err.Error(), "db query")
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("scan error", func(t *testing.T) {
		service, mock, teardown := newTestDB(t)
		defer teardown()

		rows := sqlmock.NewRows([]string{"name", "age"}).
			AddRow("Alice", 30)

		mock.ExpectQuery(query).WillReturnRows(rows)

		names, err := service.GetNames()

		assert.Error(t, err)
		assert.Nil(t, names)
		assert.Contains(t, err.Error(), "rows scanning")
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("rows iteration error", func(t *testing.T) {
		service, mock, teardown := newTestDB(t)
		defer teardown()

		rows := sqlmock.NewRows([]string{"name"}).
			AddRow("Alice").
			RowError(0, errors.New("corrupted data"))

		mock.ExpectQuery(query).WillReturnRows(rows)

		names, err := service.GetNames()

		assert.Error(t, err)
		assert.Nil(t, names)
		assert.Contains(t, err.Error(), "rows error")
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestDBService_GetUniqueNames(t *testing.T) {
	query := "SELECT DISTINCT name FROM users"

	t.Run("success", func(t *testing.T) {
		service, mock, teardown := newTestDB(t)
		defer teardown()

		rows := sqlmock.NewRows([]string{"name"}).AddRow("Charlie")
		mock.ExpectQuery(query).WillReturnRows(rows)

		names, err := service.GetUniqueNames()

		assert.NoError(t, err)
		assert.Equal(t, []string{"Charlie"}, names)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("query error", func(t *testing.T) {
		service, mock, teardown := newTestDB(t)
		defer teardown()

		mock.ExpectQuery(query).WillReturnError(errors.New("syntax error"))

		names, err := service.GetUniqueNames()

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "db query")
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("scan error", func(t *testing.T) {
		service, mock, teardown := newTestDB(t)
		defer teardown()

		rows := sqlmock.NewRows([]string{"name"}).AddRow(nil)

		mock.ExpectQuery(query).WillReturnRows(rows)

		names, err := service.GetUniqueNames()

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "rows scanning")
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("rows error", func(t *testing.T) {
		service, mock, teardown := newTestDB(t)
		defer teardown()

		rows := sqlmock.NewRows([]string{"name"}).
			AddRow("Dave").
			CloseError(errors.New("close error"))

		mock.ExpectQuery(query).WillReturnRows(rows)

		names, err := service.GetUniqueNames()

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "rows error")
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}
