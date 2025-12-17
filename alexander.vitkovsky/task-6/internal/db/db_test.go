package db

import (
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/require"
)

type rowTestDB struct {
	names       []string
	errExpected error
}

var testTable = []rowTestDB{
	{
		names: []string{"Bob", "Steve", "Svetozar"},
	},
	{
		names:       nil,
		errExpected: errors.New("no rows in result set"),
	},
}

func mockDBRows(names []string) *sqlmock.Rows {
	rows := sqlmock.NewRows([]string{"name"})
	for _, name := range names {
		rows = rows.AddRow(name)
	}
	return rows
}

func TestNew(t *testing.T) {
	/*
		Насколько я понял, этот тест не имеет смысла.
		Конструктор не имеет вообще никого ветвления -> нечему ломаться.
		Но тест есть, чтобы обеспечить покрытие 100%.
	*/
	mockDB, _, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	dbService := New(mockDB)

	require.Equal(t, mockDB, dbService.DB,
		"expected: %s, got: %s", mockDB, dbService.DB)
}

func TestGetNames(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	dbService := DBService{DB: mockDB}

	for i, row := range testTable {
		mock.ExpectQuery("SELECT name FROM users").
			WillReturnRows(mockDBRows(row.names)).
			WillReturnError(row.errExpected)
		names, err := dbService.GetNames()

		if row.errExpected != nil {
			require.ErrorIs(t, err, row.errExpected,
				"row: %d, expected: %w, got: %w", i, row.errExpected, err)
			require.Nil(t, names, "row: %d, names must be nil", i)
			continue
		}

		require.NoError(t, err, "row: %d, names must be nil", i)
		require.Equal(t, row.names, names,
			"row: %d, expected: %s, got: %s", i, row.names, names)
	}
}

func TestGetNames_RowsError(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	require.NoError(t, err)

	service := DBService{DB: mockDB}

	rows := sqlmock.NewRows([]string{"name"}).
		AddRow("Ivan").
		RowError(0, errors.New("rows error"))

	mock.ExpectQuery("SELECT name FROM users").
		WillReturnRows(rows)

	names, err := service.GetNames()

	require.Error(t, err)
	require.Nil(t, names)
}

func TestGetNames_ScanError(t *testing.T) {
	mockDB, mock, _ := sqlmock.New()
	service := DBService{DB: mockDB}

	rows := sqlmock.NewRows([]string{"name"}).
		AddRow(nil)

	mock.ExpectQuery("SELECT name FROM users").
		WillReturnRows(rows)

	names, err := service.GetNames()

	require.Error(t, err)
	require.Nil(t, names)
}

func TestGetUniqueNames(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	dbService := DBService{DB: mockDB}

	for i, row := range testTable {
		mock.ExpectQuery("SELECT DISTINCT name FROM users").
			WillReturnRows(mockDBRows(row.names)).
			WillReturnError(row.errExpected)
		names, err := dbService.GetUniqueNames()

		if row.errExpected != nil {
			require.ErrorIs(t, err, row.errExpected,
				"row: %d, expected: %w, got: %w", i, row.errExpected, err)
			require.Nil(t, names, "row: %d, names must be nil", i)
			continue
		}

		require.NoError(t, err, "row: %d, names must be nil", i)
		require.Equal(t, row.names, names,
			"row: %d, expected: %s, got: %s", i, row.names, names)
	}
}

func TestGetUniqueNames_RowsError(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	require.NoError(t, err)

	service := DBService{DB: mockDB}

	rows := sqlmock.NewRows([]string{"name"}).
		AddRow("Ivan").
		RowError(0, errors.New("rows error"))

	mock.ExpectQuery("SELECT DISTINCT name FROM users").
		WillReturnRows(rows)

	names, err := service.GetUniqueNames()

	require.Error(t, err)
	require.Nil(t, names)
}

func TestGetUniqueNames_ScanError(t *testing.T) {
	mockDB, mock, _ := sqlmock.New()
	service := DBService{DB: mockDB}

	rows := sqlmock.NewRows([]string{"name"}).
		AddRow(nil)

	mock.ExpectQuery("SELECT DISTINCT name FROM users").
		WillReturnRows(rows)

	names, err := service.GetUniqueNames()

	require.Error(t, err)
	require.Nil(t, names)
}
