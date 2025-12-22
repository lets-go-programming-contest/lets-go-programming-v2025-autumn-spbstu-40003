package db_test

import (
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/KiRy6A/task-6/internal/db"
	"github.com/stretchr/testify/require"
)

var (
	errNameExpected = errors.New("expected name")
	errRows         = errors.New("rows error")
	errScan         = errors.New("scan error")
)

type rowTestDB struct {
	names           []string
	errExpected     error
	errScanExpected error
	errRowsExpected error
	errIndex        int
}

func getTestTable() []rowTestDB {
	return []rowTestDB{
		{
			names: []string{"Petr", "Kirill"},
		},
		{
			names: []string{"Ivan", "Oleg"},
		},
		{
			names:       nil,
			errExpected: errNameExpected,
		},
		{
			names:           []string{"Ivan", "Oleg"},
			errScanExpected: errScan,
			errIndex:        1,
		},
		{
			names:           []string{"Ivan", "Oleg"},
			errRowsExpected: errRows,
			errIndex:        1,
		},
	}
}

func TestNew(t *testing.T) {
	t.Parallel()

	mockDB, _, err := sqlmock.New()
	if err != nil {
		t.Fatalf("sqlmock.New() error: %s", err)
	}
	defer mockDB.Close()

	dbService := db.New(mockDB)

	require.NotNil(t, dbService, "dbService should not be nil")
	require.NotNil(t, dbService.DB, "dbService.DB should not be nil")
	require.Equal(t, mockDB, dbService.DB, "dbService.DB should equal the provided database")
}

func TestGetNames(t *testing.T) {
	t.Parallel()

	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("sqlmock.New() error: %s", err)
	}
	defer mockDB.Close()

	dbService := db.DBService{DB: mockDB}

	for i, test := range getTestTable() {
		mock.ExpectQuery("SELECT name FROM users").WillReturnRows(mockDBRows(test)).
			WillReturnError(test.errExpected)

		names, err := dbService.GetNames()

		if test.errScanExpected != nil {
			require.Error(t, err, "row: %d, error: %w", i, err)
			require.Nil(t, names, "row: %d, names must be nil", i)

			continue
		}

		if test.errRowsExpected != nil {
			require.ErrorIs(t, err, test.errRowsExpected, "row: %d, expected error: %w, actual error: %w", i,
				test.errRowsExpected, err)
			require.Nil(t, names, "row: %d, names must be nil", i)

			continue
		}

		if test.errExpected != nil {
			require.ErrorIs(t, err, test.errExpected, "row: %d, expected error: %w, actual error: %w", i,
				test.errExpected, err)
			require.Nil(t, names, "row: %d, names must be nil", i)

			continue
		}

		require.NoError(t, err, "row: %d, error must be nil", i)
		require.Equal(t, test.names, names, "row: %d, expected names: %s, actual names: %s", i, test.names, names)
	}
}

func TestGetUniqueNames(t *testing.T) {
	t.Parallel()

	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("sqlmock.New() error: %s", err)
	}
	defer mockDB.Close()

	dbService := db.DBService{DB: mockDB}

	for i, test := range getTestTable() {
		mock.ExpectQuery("SELECT DISTINCT name FROM users").WillReturnRows(mockDBRows(test)).
			WillReturnError(test.errExpected)

		names, err := dbService.GetUniqueNames()

		if test.errScanExpected != nil {
			require.Error(t, err, "row: %d, error: %w", i, err)
			require.Nil(t, names, "row: %d, names must be nil", i)

			continue
		}

		if test.errRowsExpected != nil {
			require.ErrorIs(t, err, test.errRowsExpected, "row: %d, expected error: %w, actual error: %w", i,
				test.errRowsExpected, err)
			require.Nil(t, names, "row: %d, names must be nil", i)

			continue
		}

		if test.errExpected != nil {
			require.ErrorIs(t, err, test.errExpected, "row: %d, expected error: %w, actual error: %w", i,
				test.errExpected, err)
			require.Nil(t, names, "row: %d, names must be nil", i)

			continue
		}

		require.NoError(t, err, "row: %d, error must be nil", i)
		require.Equal(t, test.names, names, "row: %d, expected names: %s, actual names: %s", i, test.names, names)
	}
}

func mockDBRows(test rowTestDB) *sqlmock.Rows {
	rows := sqlmock.NewRows([]string{"name"})

	for i, name := range test.names {
		if i == test.errIndex && test.errScanExpected != nil {
			rows = rows.AddRow(nil)
		} else {
			rows = rows.AddRow(name)
		}
	}

	if test.errRowsExpected != nil {
		rows.RowError(test.errIndex, errRows)
	}

	return rows
}
