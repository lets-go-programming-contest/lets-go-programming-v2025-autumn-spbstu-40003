package db_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/Vurvaa/task-6/internal/db"
	"github.com/stretchr/testify/require"
)

var (
	ErrQueryError = errors.New("ExpectedError")
	ErrRowsError  = errors.New("RowError")
)

type testCase struct {
	rows     *sqlmock.Rows
	queryErr error

	expected      []string
	expectedError bool
}

type callFunc func(db.DBService) ([]string, error)

func TestNew(t *testing.T) {
	t.Parallel()

	mockDB, _, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	dbService := db.New(mockDB)
	require.Equal(t, mockDB, dbService.DB)
}

func TestGetNames(t *testing.T) {
	t.Parallel()

	testCases := []testCase{
		{
			rows:     getMockDBRows([]string{"Egor", "Egor", "Sasha", "Sasha"}),
			expected: []string{"Egor", "Egor", "Sasha", "Sasha"},
		},
		{
			queryErr:      ErrQueryError,
			expectedError: true,
		},
		{
			rows:          sqlmock.NewRows([]string{"name"}).AddRow(nil),
			expectedError: true,
		},
		{
			rows: sqlmock.NewRows([]string{"name"}).
				AddRow("Aboba").
				AddRow("Alex").
				RowError(1, ErrRowsError),
			expectedError: true,
		},
		{
			rows:     sqlmock.NewRows([]string{"name"}),
			expected: nil,
		},
	}

	runTestTable(t, "SELECT name FROM users", testCases, db.DBService.GetNames)
}

func TestGetUniqueNames(t *testing.T) {
	t.Parallel()

	testCases := []testCase{
		{
			rows:     getMockDBRows([]string{"Egor", "Sasha"}),
			expected: []string{"Egor", "Sasha"},
		},
		{
			queryErr:      ErrQueryError,
			expectedError: true,
		},
		{
			rows:          sqlmock.NewRows([]string{"name"}).AddRow(nil),
			expectedError: true,
		},
		{
			rows: sqlmock.NewRows([]string{"name"}).
				AddRow("Aboba").
				AddRow("Alex").
				RowError(1, ErrRowsError),
			expectedError: true,
		},
		{
			rows:     sqlmock.NewRows([]string{"name"}),
			expected: nil,
		},
	}

	runTestTable(t, "SELECT DISTINCT name FROM users", testCases, db.DBService.GetUniqueNames)
}

func runTestTable(t *testing.T, query string, testCases []testCase, call callFunc) {
	t.Helper()

	for i, tc := range testCases {
		t.Run(fmt.Sprintf("case_%d", i), func(t *testing.T) {
			runTestCase(t, query, tc, call)
		})
	}
}

func runTestCase(t *testing.T, query string, tc testCase, call callFunc) {
	t.Helper()
	t.Parallel()

	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	service := db.DBService{DB: mockDB}

	if tc.queryErr != nil {
		mock.ExpectQuery(query).WillReturnError(tc.queryErr)
	} else {
		mock.ExpectQuery(query).WillReturnRows(tc.rows)
	}

	got, err := call(service)

	if tc.expectedError {
		require.Error(t, err)
		require.Nil(t, got)

		return
	}

	require.NoError(t, err)
	require.Equal(t, tc.expected, got)
}

func getMockDBRows(names []string) *sqlmock.Rows {
	rows := sqlmock.NewRows([]string{"name"})
	for _, name := range names {
		rows.AddRow(name)
	}

	return rows
}
