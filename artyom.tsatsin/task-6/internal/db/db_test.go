package db

import (
	"database/sql"
	"errors"
	"reflect"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

func setupDB(t *testing.T) (*sql.DB, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("sqlmock error: %v", err)
	}

	t.Cleanup(func() {
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Fatalf("unmet sqlmock expectations: %v", err)
		}
	})

	return db, mock
}

func TestGetNames_OK(t *testing.T) {
	dbConn, mock := setupDB(t)
	defer dbConn.Close()

	rows := sqlmock.NewRows([]string{"name"}).
		AddRow("Alice").
		AddRow("Bob")

	mock.ExpectQuery("^SELECT name FROM users$").
		WillReturnRows(rows)

	service := New(dbConn)

	got, err := service.GetNames()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	want := []string{"Alice", "Bob"}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v, want %v", got, want)
	}
}

func TestGetNames_QueryError(t *testing.T) {
	dbConn, mock := setupDB(t)
	defer dbConn.Close()

	mock.ExpectQuery("^SELECT name FROM users$").
		WillReturnError(errors.New("query error"))

	service := New(dbConn)

	_, err := service.GetNames()
	if err == nil {
		t.Fatal("expected query error")
	}
}

func TestGetNames_ScanError(t *testing.T) {
	dbConn, mock := setupDB(t)
	defer dbConn.Close()

	rows := sqlmock.NewRows([]string{"name"}).
		AddRow(123)

	mock.ExpectQuery("^SELECT name FROM users$").
		WillReturnRows(rows)

	service := New(dbConn)

	_, err := service.GetNames()
	if err == nil {
		t.Fatal("expected scan error")
	}
}

func TestGetNames_RowsError(t *testing.T) {
	dbConn, mock := setupDB(t)
	defer dbConn.Close()

	rows := sqlmock.NewRows([]string{"name"}).
		AddRow("Alice").
		RowError(0, errors.New("rows error"))

	mock.ExpectQuery("^SELECT name FROM users$").
		WillReturnRows(rows)

	service := New(dbConn)

	_, err := service.GetNames()
	if err == nil {
		t.Fatal("expected rows error")
	}
}

func TestGetUniqueNames_OK(t *testing.T) {
	dbConn, mock := setupDB(t)
	defer dbConn.Close()

	rows := sqlmock.NewRows([]string{"name"}).
		AddRow("Alice")

	mock.ExpectQuery("^SELECT DISTINCT name FROM users$").
		WillReturnRows(rows)

	service := New(dbConn)

	got, err := service.GetUniqueNames()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	want := []string{"Alice"}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v, want %v", got, want)
	}
}

func TestGetUniqueNames_QueryError(t *testing.T) {
	dbConn, mock := setupDB(t)
	defer dbConn.Close()

	mock.ExpectQuery("^SELECT DISTINCT name FROM users$").
		WillReturnError(errors.New("query error"))

	service := New(dbConn)

	_, err := service.GetUniqueNames()
	if err == nil {
		t.Fatal("expected query error")
	}
}

func TestGetUniqueNames_ScanError(t *testing.T) {
	dbConn, mock := setupDB(t)
	defer dbConn.Close()

	rows := sqlmock.NewRows([]string{"name"}).
		AddRow(123)

	mock.ExpectQuery("^SELECT DISTINCT name FROM users$").
		WillReturnRows(rows)

	service := New(dbConn)

	_, err := service.GetUniqueNames()
	if err == nil {
		t.Fatal("expected scan error")
	}
}

func TestGetUniqueNames_RowsError(t *testing.T) {
	dbConn, mock := setupDB(t)
	defer dbConn.Close()

	rows := sqlmock.NewRows([]string{"name"}).
		AddRow("Alice").
		RowError(0, errors.New("rows error"))

	mock.ExpectQuery("^SELECT DISTINCT name FROM users$").
		WillReturnRows(rows)

	service := New(dbConn)

	_, err := service.GetUniqueNames()
	if err == nil {
		t.Fatal("expected rows error")
	}
}
