package db_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/wedwincode/task-6/internal/db"
	"gopkg.in/DATA-DOG/go-sqlmock.v1"
)

var (
	errRows     = errors.New("RowsError")
	errScan     = errors.New("ScanError")
	errExpected = errors.New("ExpectedError")
)

type rowTestDB struct {
	names       []string
	errExpected error
	scanErrIdx  int
	scanErr     error
	rowsErr     error
}

func getTestTable() []rowTestDB {
	return []rowTestDB{
		{
			names:       []string{"name1", "name2"},
			errExpected: nil,
			scanErrIdx:  0,
			scanErr:     nil,
			rowsErr:     nil,
		},
		{
			names:       []string{"name1", "name2"},
			errExpected: nil,
			scanErrIdx:  0,
			scanErr:     nil,
			rowsErr:     errRows,
		},
		{
			names:       []string{"name1", "name2"},
			errExpected: nil,
			scanErrIdx:  1,
			scanErr:     errScan,
			rowsErr:     nil,
		},
		{
			names:       nil,
			errExpected: errExpected,
			scanErrIdx:  0,
			scanErr:     nil,
			rowsErr:     nil,
		},
	}
}

func TestNew(t *testing.T) {
	t.Parallel()

	mockDB, _, err := sqlmock.New()
	require.NoError(t, err, "failed to create mock database")

	dbService := db.New(mockDB)

	require.NotNil(t, dbService, "dbService should not be nil")
	require.NotNil(t, dbService.DB, "dbService.DB should not be nil")
	require.Equal(t, mockDB, dbService.DB, "dbService.DB should equal the provided database")
}

func TestDBService_GetNames(t *testing.T) {
	t.Parallel()

	testQuery(t,
		"SELECT name FROM users",
		func(service db.DBService) ([]string, error) {
			return service.GetNames()
		},
	)
}

func TestDBService_GetUniqueNames(t *testing.T) {
	t.Parallel()

	testQuery(t,
		"SELECT DISTINCT name FROM users",
		func(service db.DBService) ([]string, error) {
			return service.GetUniqueNames()
		},
	)
}

func testQuery(t *testing.T, query string, call func(service db.DBService) ([]string, error)) {
	t.Helper()

	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	dbService := db.DBService{DB: mockDB}

	for idx, row := range getTestTable() {
		mock.ExpectQuery(query).WillReturnRows(
			mockDBRows(row)).WillReturnError(row.errExpected)

		names, err := call(dbService)

		if row.rowsErr != nil {
			require.ErrorIs(t, err, row.rowsErr,
				"row: %d, expected error: %w, actual error: %w", idx, row.rowsErr, err)

			continue
		}

		if row.scanErr != nil {
			require.Error(t, err, "row: %d, error: %w", idx, err)

			continue
		}

		if row.errExpected != nil {
			require.ErrorIs(t, err, row.errExpected,
				"row: %d, expected error: %w, actual error: %w", idx, row.errExpected, err)
			require.Nil(t, names, "row: %d, names must be nil", idx)

			continue
		}

		require.NoError(t, err, "row: %d, error must be nil", idx)
		require.Equal(t, names, row.names,
			"row: %d, expected names: %s, actual names: %s", idx, row.names, names)
	}
}

func mockDBRows(row rowTestDB) *sqlmock.Rows {
	rows := sqlmock.NewRows([]string{"name"})

	for i, name := range row.names {
		if row.scanErr != nil && row.scanErrIdx == i {
			rows = rows.AddRow(nil)

			continue
		}

		rows = rows.AddRow(name)
	}

	if row.rowsErr != nil {
		last := len(row.names) - 1
		rows.RowError(last, row.rowsErr)
	}

	return rows
}
