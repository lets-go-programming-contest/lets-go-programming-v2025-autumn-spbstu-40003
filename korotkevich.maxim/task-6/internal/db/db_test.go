package db_test

import (
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	dbpkg "github.com/KrrMaxim/task-6/internal/db"
)

var (
	errQuery = errors.New("query error")
	errRow   = errors.New("row iteration error")
)

func TestNew(t *testing.T) {
	t.Parallel()

	mockDB, _, err := sqlmock.New()
	require.NoError(t, err)
	defer mockDB.Close()

	service := dbpkg.New(mockDB)
	require.NotNil(t, service)
	require.Equal(t, mockDB, service.DB)
}

func TestGetNames(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name        string
		setupMock   func(sqlmock.Sqlmock)
		expected    []string
		errContains string
	}{
		{
			name: "success multiple rows",
			setupMock: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"name"}).
					AddRow("Maxim").
					AddRow("Artem")
				mock.ExpectQuery(`SELECT name FROM users`).WillReturnRows(rows)
			},
			expected: []string{"Maxim", "Artem"},
		},
		{
			name: "success empty",
			setupMock: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"name"})
				mock.ExpectQuery(`SELECT name FROM users`).WillReturnRows(rows)
			},
			expected: nil,
		},
		{
			name: "query error",
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(`SELECT name FROM users`).WillReturnError(errQuery)
			},
			errContains: "db query",
		},
		{
			name: "scan error",
			setupMock: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"name"}).AddRow(nil)
				mock.ExpectQuery(`SELECT name FROM users`).WillReturnRows(rows)
			},
			errContains: "rows scanning",
		},
		{
			name: "rows iteration error",
			setupMock: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"name"}).
					AddRow("first").
					AddRow("second")
				rows.RowError(1, errRow)
				mock.ExpectQuery(`SELECT name FROM users`).WillReturnRows(rows)
			},
			errContains: "rows error",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			mockDB, mock, err := sqlmock.New()
			require.NoError(t, err)
			defer mockDB.Close()

			tc.setupMock(mock)

			service := dbpkg.New(mockDB)
			result, err := service.GetNames()

			require.NoError(t, mock.ExpectationsWereMet())

			if tc.errContains != "" {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tc.errContains)
				assert.Nil(t, result)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tc.expected, result)
			}
		})
	}
}

func TestGetUniqueNames(t *testing.T) {
	t.Parallel()

	runCase := func(
		t *testing.T,
		prepare func(sqlmock.Sqlmock),
		expected []string,
		errContains string,
	) {
		t.Helper()

		mockDB, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer mockDB.Close()

		prepare(mock)

		service := dbpkg.New(mockDB)
		result, err := service.GetUniqueNames()

		require.NoError(t, mock.ExpectationsWereMet())

		if errContains != "" {
			require.Error(t, err)
			assert.Contains(t, err.Error(), errContains)
			assert.Nil(t, result)

			return
		}

		require.NoError(t, err)
		assert.Equal(t, expected, result)
	}

	t.Run("success distinct", func(t *testing.T) {
		t.Parallel()

		prepare := func(mock sqlmock.Sqlmock) {
			rows := sqlmock.NewRows([]string{"name"}).
				AddRow("u1").
				AddRow("u2")
			mock.ExpectQuery("SELECT DISTINCT name FROM users").WillReturnRows(rows)
		}

		runCase(t, prepare, []string{"u1", "u2"}, "")
	})

	t.Run("query error", func(t *testing.T) {
		t.Parallel()

		prepare := func(mock sqlmock.Sqlmock) {
			mock.ExpectQuery("SELECT DISTINCT name FROM users").
				WillReturnError(errQuery)
		}

		runCase(t, prepare, nil, "db query")
	})

	t.Run("scan error", func(t *testing.T) {
		t.Parallel()

		prepare := func(mock sqlmock.Sqlmock) {
			rows := sqlmock.NewRows([]string{"name"}).AddRow(nil)
			mock.ExpectQuery("SELECT DISTINCT name FROM users").WillReturnRows(rows)
		}

		runCase(t, prepare, nil, "rows scanning")
	})

	t.Run("rows error", func(t *testing.T) {
		t.Parallel()

		prepare := func(mock sqlmock.Sqlmock) {
			rows := sqlmock.NewRows([]string{"name"}).
				AddRow("ok").
				AddRow("bad")
			rows.RowError(1, errRow)
			mock.ExpectQuery("SELECT DISTINCT name FROM users").WillReturnRows(rows)
		}

		runCase(t, prepare, nil, "rows error")
	})
}
