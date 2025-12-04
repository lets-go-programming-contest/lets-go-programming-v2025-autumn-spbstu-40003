package db_test

import (
	"fmt"
	"testing"

	"github.com/ArtttNik/task-6/internal/db"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	assert.NoError(t, err)

	mock.ExpectClose()

	s := db.New(mockDB)
	assert.Equal(t, mockDB, s.DB)

	assert.NoError(t, mockDB.Close())
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestDBService_GetNames(t *testing.T) {
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
				mock.ExpectQuery("SELECT name FROM users").WillReturnError(fmt.Errorf("query error"))
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
					RowError(1, fmt.Errorf("rows error"))
				mock.ExpectQuery("SELECT name FROM users").WillReturnRows(rows)
			},
			wantErr: "rows error: rows error",
		},
		{
			name: "rows error no rows",
			setup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"name"}).
					AddRow("dummy").
					RowError(0, fmt.Errorf("rows error"))
				mock.ExpectQuery("SELECT name FROM users").WillReturnRows(rows)
			},
			wantErr: "rows error: rows error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDB, mock, err := sqlmock.New()
			assert.NoError(t, err)

			tt.setup(mock)

			mock.ExpectClose()

			s := db.New(mockDB)
			got, err := s.GetNames()

			if tt.wantErr != "" {
				assert.Error(t, err)
				if err != nil {
					assert.Contains(t, err.Error(), tt.wantErr)
				}
				assert.Nil(t, got)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}

			assert.NoError(t, mockDB.Close())
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestDBService_GetUniqueNames(t *testing.T) {
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
				mock.ExpectQuery("SELECT DISTINCT name FROM users").
					WillReturnError(fmt.Errorf("query error"))
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
					RowError(1, fmt.Errorf("rows error"))
				mock.ExpectQuery("SELECT DISTINCT name FROM users").WillReturnRows(rows)
			},
			wantErr: "rows error: rows error",
		},
		{
			name: "rows error no rows",
			setup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"name"}).
					AddRow("dummy").
					RowError(0, fmt.Errorf("rows error"))
				mock.ExpectQuery("SELECT DISTINCT name FROM users").WillReturnRows(rows)
			},
			wantErr: "rows error: rows error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDB, mock, err := sqlmock.New()
			assert.NoError(t, err)

			tt.setup(mock)

			mock.ExpectClose()

			s := db.New(mockDB)
			got, err := s.GetUniqueNames()

			if tt.wantErr != "" {
				assert.Error(t, err)
				if err != nil {
					assert.Contains(t, err.Error(), tt.wantErr)
				}
				assert.Nil(t, got)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}

			assert.NoError(t, mockDB.Close())
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}
