package db_test

import (
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/d1mene/task-6/internal/db"
	"github.com/stretchr/testify/require"
)

type dbTestestase struct {
	name          string
	query         string
	mockBehavior  func(mock sqlmock.Sqlmock)
	expectedNames []string
	expectedError string
}

func TestNew(t *testing.T) {
	mockDB, _, _ := sqlmock.New()
	service := db.New(mockDB)
	require.Equal(t, mockDB, service.DB)
}

func TestGetNames(t *testing.T) {
	query := "SELECT name FROM users"

	testTable := []dbTestestase{
		{
			name:  "Success Query",
			query: query,
			mockBehavior: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"name"}).AddRow("d1mene").AddRow("pupsik")
				mock.ExpectQuery(query).WillReturnRows(rows)
			},
			expectedNames: []string{"d1mene", "pupsik"},
		},
		{
			name:  "Query Error",
			query: query,
			mockBehavior: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(query).WillReturnError(errors.New("db down"))
			},
			expectedError: "db query: db down",
		},
		{
			name:  "Scan Error",
			query: query,
			mockBehavior: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"name"}).AddRow(nil)
				mock.ExpectQuery(query).WillReturnRows(rows)
			},
			expectedError: "rows scanning",
		},
		{
			name:  "Rows Iteration Error",
			query: query,
			mockBehavior: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"name"}).
					RowError(0, errors.New("iteration error")).
					AddRow("d1mene")
				mock.ExpectQuery(query).WillReturnRows(rows)
			},
			expectedError: "rows error: iteration error",
		},
	}

	for _, test := range testTable {
		t.Run(test.name, func(t *testing.T) {
			mockDB, mock, err := sqlmock.New()
			require.NoError(t, err)
			defer mockDB.Close()

			dbService := db.New(mockDB)

			test.mockBehavior(mock)

			names, err := dbService.GetNames()

			if test.expectedError != "" {
				require.Error(t, err)
				require.Contains(t, err.Error(), test.expectedError)
				require.Nil(t, names)
			} else {
				require.NoError(t, err)
				require.Equal(t, test.expectedNames, names)
			}

			require.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestGetUniqueNames(t *testing.T) {
	query := "SELECT DISTINCT name FROM users"

	testTable := []dbTestestase{
		{
			name:  "Success Query",
			query: query,
			mockBehavior: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"name"}).AddRow("d1mene")
				mock.ExpectQuery(query).WillReturnRows(rows)
			},
			expectedNames: []string{"d1mene"},
		},
		{
			name:  "Query Error",
			query: query,
			mockBehavior: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(query).WillReturnError(errors.New("fatal error"))
			},
			expectedError: "db query: fatal error",
		},
		{
			name:  "Scan Error",
			query: query,
			mockBehavior: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"name"}).AddRow(nil)
				mock.ExpectQuery(query).WillReturnRows(rows)
			},
			expectedError: "rows scanning",
		},
		{
			name:  "Rows Iteration Error",
			query: query,
			mockBehavior: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"name"}).
					RowError(0, errors.New("iteration error")).
					AddRow("d1mene")
				mock.ExpectQuery(query).WillReturnRows(rows)
			},
			expectedError: "rows error: iteration error",
		},
	}

	for _, test := range testTable {
		t.Run(test.name, func(t *testing.T) {
			mockDB, mock, err := sqlmock.New()
			require.NoError(t, err)
			defer mockDB.Close()

			dbService := db.New(mockDB)

			test.mockBehavior(mock)

			names, err := dbService.GetUniqueNames()

			if test.expectedError != "" {
				require.Error(t, err)
				require.Contains(t, err.Error(), test.expectedError)
				require.Nil(t, names)
			} else {
				require.NoError(t, err)
				require.Equal(t, test.expectedNames, names)
			}

			require.NoError(t, mock.ExpectationsWereMet())
		})
	}
}
