package db

import (
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNew(t *testing.T) {
	t.Parallel()

	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	service := New(db)

	assert.NotNil(t, service)
	assert.Equal(t, db, service.DB)

	require.NoError(t, mock.ExpectationsWereMet())
}

func TestDBService_GetNames(t *testing.T) {
	tests := []struct {
		name         string
		mockSetup    func(sqlmock.Sqlmock)
		expectErr    bool
		expectErrMsg string
		expectNames  []string
	}{
		{
			name: "success with two names",
			mockSetup: func(m sqlmock.Sqlmock) {
				m.ExpectQuery(`SELECT name FROM users`).
					WillReturnRows(sqlmock.NewRows([]string{"name"}).AddRow("Alex").AddRow("Sam"))
			},
			expectNames: []string{"Alex", "Sam"},
		},
		{
			name: "empty result",
			mockSetup: func(m sqlmock.Sqlmock) {
				m.ExpectQuery(`SELECT name FROM users`).
					WillReturnRows(sqlmock.NewRows([]string{"name"}))
			},
			expectNames: nil,
		},
		{
			name: "query error",
			mockSetup: func(m sqlmock.Sqlmock) {
				m.ExpectQuery(`SELECT name FROM users`).
					WillReturnError(errors.New("db error"))
			},
			expectErr:    true,
			expectErrMsg: "db query:",
		},
		{
			name: "scan error",
			mockSetup: func(m sqlmock.Sqlmock) {
				m.ExpectQuery(`SELECT name FROM users`).
					WillReturnRows(sqlmock.NewRows([]string{"name"}).AddRow(nil))
			},
			expectErr:    true,
			expectErrMsg: "rows scanning:",
		},
		{
			name: "rows iteration error",
			mockSetup: func(m sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"name"}).AddRow("Max")
				rows.RowError(0, errors.New("iteration error"))
				m.ExpectQuery(`SELECT name FROM users`).
					WillReturnRows(rows)
			},
			expectErr:    true,
			expectErrMsg: "rows error:",
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			db, mock, err := sqlmock.New()
			require.NoError(t, err)
			defer db.Close()

			tc.mockSetup(mock)

			svc := New(db)
			names, err := svc.GetNames()

			if tc.expectErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tc.expectErrMsg)
				assert.Nil(t, names)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tc.expectNames, names)
			}

			require.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestDBService_GetUniqueNames(t *testing.T) {
	tests := []struct {
		name         string
		mockSetup    func(sqlmock.Sqlmock)
		expectErr    bool
		expectErrMsg string
		expectNames  []string
	}{
		{
			name: "success with unique names",
			mockSetup: func(m sqlmock.Sqlmock) {
				m.ExpectQuery(`SELECT DISTINCT name FROM users`).
					WillReturnRows(sqlmock.NewRows([]string{"name"}).AddRow("Tom").AddRow("Anna"))
			},
			expectNames: []string{"Tom", "Anna"},
		},
		{
			name: "empty result",
			mockSetup: func(m sqlmock.Sqlmock) {
				m.ExpectQuery(`SELECT DISTINCT name FROM users`).
					WillReturnRows(sqlmock.NewRows([]string{"name"}))
			},
			expectNames: nil,
		},
		{
			name: "query error",
			mockSetup: func(m sqlmock.Sqlmock) {
				m.ExpectQuery(`SELECT DISTINCT name FROM users`).
					WillReturnError(errors.New("timeout error"))
			},
			expectErr:    true,
			expectErrMsg: "db query:",
		},
		{
			name: "scan error",
			mockSetup: func(m sqlmock.Sqlmock) {
				m.ExpectQuery(`SELECT DISTINCT name FROM users`).
					WillReturnRows(sqlmock.NewRows([]string{"name"}).AddRow(nil))
			},
			expectErr:    true,
			expectErrMsg: "rows scanning:",
		},
		{
			name: "rows error",
			mockSetup: func(m sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"name"}).AddRow("Leo")
				rows.RowError(0, errors.New("row error"))
				m.ExpectQuery(`SELECT DISTINCT name FROM users`).
					WillReturnRows(rows)
			},
			expectErr:    true,
			expectErrMsg: "rows error:",
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			db, mock, err := sqlmock.New()
			require.NoError(t, err)
			defer db.Close()

			tc.mockSetup(mock)

			svc := New(db)
			names, err := svc.GetUniqueNames()

			if tc.expectErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tc.expectErrMsg)
				assert.Nil(t, names)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tc.expectNames, names)
			}

			require.NoError(t, mock.ExpectationsWereMet())
		})
	}
}
