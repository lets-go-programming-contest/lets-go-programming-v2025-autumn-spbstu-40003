package db_test

import (
	"errors"
	"testing"

	"github.com/ArtttNik/task-6/internal/db"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	errQuery = errors.New("query error")
	errRows  = errors.New("rows error")
)

func TestNew(t *testing.T) {
	t.Parallel()

	mockDB, mock, err := sqlmock.New()
	require.NoError(t, err)

	mock.ExpectClose()

	s := db.New(mockDB)

	assert.Equal(t, mockDB, s.DB)

	err = mockDB.Close()
	require.NoError(t, err)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestDBService_GetNames(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		setup   func(mock sqlmock.Sqlmock)
		want    []string
		wantErr string
	}{
		{
			name: "success multiple rows",
			setup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"name"}).
					AddRow("Alice").
					AddRow("Bob")
				mock.ExpectQuery("SELECT name FROM users").WillReturnRows(rows)
			},
			want: []string{"Alice", "Bob"},
		},
		{
			name: "success empty",
			setup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"name"})
				mock.ExpectQuery("SELECT name FROM users").WillReturnRows(rows)
			},
			want: nil,
		},
		{
			name: "query error",
			setup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT name FROM users").WillReturnError(errQuery)
			},
			wantErr: "db query: query error",
		},
		{
			name: "scan error",
			setup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"name"}).AddRow(nil)
				mock.ExpectQuery("SELECT name FROM users").WillReturnRows(rows)
			},
			wantErr: "rows scanning",
		},
		{
			name: "rows error after iteration",
			setup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"name"}).
					AddRow("Alice").
					AddRow("Bob").
					RowError(1, errRows)
				mock.ExpectQuery("SELECT name FROM users").WillReturnRows(rows)
			},
			wantErr: "rows error: rows error",
		},
		{
			name: "rows error no rows",
			setup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"name"}).
					AddRow("dummy").
					RowError(0, errRows)
				mock.ExpectQuery("SELECT name FROM users").WillReturnRows(rows)
			},
			wantErr: "rows error: rows error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			mockDB, mock, err := sqlmock.New()
			require.NoError(t, err)

			tt.setup(mock)
			mock.ExpectClose()

			s := db.New(mockDB)
			got, err := s.GetNames()

			closeErr := mockDB.Close()
			require.NoError(t, closeErr)

			if tt.wantErr != "" {
				require.Error(t, err)

				assert.Contains(t, err.Error(), tt.wantErr)
				assert.Nil(t, got)
			} else {
				require.NoError(t, err)

				assert.Equal(t, tt.want, got)
			}

			require.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestDBService_GetUniqueNames(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		setup   func(mock sqlmock.Sqlmock)
		want    []string
		wantErr string
	}{
		{
			name: "success multiple rows",
			setup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"name"}).
					AddRow("Alice").
					AddRow("Bob")
				mock.ExpectQuery("SELECT DISTINCT name FROM users").WillReturnRows(rows)
			},
			want: []string{"Alice", "Bob"},
		},
		{
			name: "success empty",
			setup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"name"})
				mock.ExpectQuery("SELECT DISTINCT name FROM users").WillReturnRows(rows)
			},
			want: nil,
		},
		{
			name: "query error",
			setup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT DISTINCT name FROM users").WillReturnError(errQuery)
			},
			wantErr: "db query: query error",
		},
		{
			name: "scan error",
			setup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"name"}).AddRow(nil)
				mock.ExpectQuery("SELECT DISTINCT name FROM users").WillReturnRows(rows)
			},
			wantErr: "rows scanning",
		},
		{
			name: "rows error after iteration",
			setup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"name"}).
					AddRow("Alice").
					AddRow("Bob").
					RowError(1, errRows)
				mock.ExpectQuery("SELECT DISTINCT name FROM users").WillReturnRows(rows)
			},
			wantErr: "rows error: rows error",
		},
		{
			name: "rows error no rows",
			setup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"name"}).
					AddRow("dummy").
					RowError(0, errRows)
				mock.ExpectQuery("SELECT DISTINCT name FROM users").WillReturnRows(rows)
			},
			wantErr: "rows error: rows error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			mockDB, mock, err := sqlmock.New()
			require.NoError(t, err)

			tt.setup(mock)
			mock.ExpectClose()

			s := db.New(mockDB)
			got, err := s.GetUniqueNames()

			closeErr := mockDB.Close()
			require.NoError(t, closeErr)

			if tt.wantErr != "" {
				require.Error(t, err)

				assert.Contains(t, err.Error(), tt.wantErr)
				assert.Nil(t, got)
			} else {
				require.NoError(t, err)

				assert.Equal(t, tt.want, got)
			}

			require.NoError(t, mock.ExpectationsWereMet())
		})
	}
}
