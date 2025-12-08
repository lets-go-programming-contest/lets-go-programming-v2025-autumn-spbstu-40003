package db_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gituser549/task-6/internal/db"
	"github.com/gituser549/task-6/internal/util"
)

var errRows = errors.New("rows error")

type mockDBResponse struct {
	names      []string
	namesRows  *sqlmock.Rows
	errFromDB  string
	errFromRow string
}

func TestDbGetNames(t *testing.T) {
	testNamesTable := []mockDBResponse{
		{
			names: []string{"Petya", "Vasya", "Kolya"},
			namesRows: func() *sqlmock.Rows {
				rows := sqlmock.NewRows([]string{"name"})
				for _, name := range []string{"Petya", "Vasya", "Kolya"} {
					rows.AddRow(name)
				}

				return rows
			}(),
		},
		{
			names:     nil,
			namesRows: sqlmock.NewRows([]string{"name"}),
			errFromDB: "sql: no rows in result set",
		},
		{
			names: nil,
			namesRows: func() *sqlmock.Rows {
				rows := sqlmock.NewRows([]string{"name"})
				rows.AddRow(nil)

				return rows
			}(),
			errFromRow: "rows scanning",
		},
		{
			namesRows: func() *sqlmock.Rows {
				rows := sqlmock.NewRows([]string{"name"})
				for _, name := range []string{"Petya", "Vasya", "Kolya"} {
					rows.AddRow(name)
				}
				rows.RowError(0, errRows)

				return rows
			}(),
			errFromRow: "rows error",
		},
	}

	t.Parallel()

	for numTestCase, curRow := range testNamesTable {
		t.Run(fmt.Sprintf("%s #%d", t.Name(), numTestCase), func(t *testing.T) {
			t.Parallel()

			mockDB, mockCfgSetter, err := sqlmock.New()
			if err != nil {
				t.Fatalf("an error '%s' was not expected when creating a mock database", err)
			}

			defer mockDB.Close()

			checkDBService := db.New(mockDB)

			mockCfgSetter.ExpectQuery("SELECT name FROM users").
				WillReturnRows(curRow.namesRows).
				WillReturnError(util.MakeError(curRow.errFromDB))

			names, err := checkDBService.GetNames()

			if !util.IsEmpty(curRow.errFromDB) {
				util.AssertError(t, names, err, curRow.errFromDB)

				return
			}

			if !util.IsEmpty(curRow.errFromRow) {
				util.AssertError(t, names, err, curRow.errFromRow)

				return
			}

			util.AssertNoError(t, curRow.names, names, err)

			if err := mockCfgSetter.ExpectationsWereMet(); err != nil {
				t.Fatalf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}

func TestDbGetUniqueNames(t *testing.T) {
	testDistinctNamesTable := []mockDBResponse{
		{
			names: []string{"Alex", "Victor"},
			namesRows: func() *sqlmock.Rows {
				rows := sqlmock.NewRows([]string{"name"})
				for _, name := range []string{"Alex", "Victor"} {
					rows.AddRow(name)
				}

				return rows
			}(),
		},
		{
			names:     nil,
			namesRows: sqlmock.NewRows([]string{"name"}),
			errFromDB: "sql: no rows in result set",
		},
		{
			names: nil,
			namesRows: func() *sqlmock.Rows {
				rows := sqlmock.NewRows([]string{"name"})
				rows.AddRow(nil)

				return rows
			}(),
			errFromRow: "rows scanning",
		},
		{
			namesRows: func() *sqlmock.Rows {
				rows := sqlmock.NewRows([]string{"name"})
				for _, name := range []string{"Petya", "Vasya", "Kolya"} {
					rows.AddRow(name)
				}
				rows.RowError(0, errRows)

				return rows
			}(),
			errFromRow: "rows error",
		},
	}

	t.Parallel()

	for numTestCase, curRow := range testDistinctNamesTable {
		t.Run(fmt.Sprintf("%s #%d", t.Name(), numTestCase), func(t *testing.T) {
			t.Parallel()

			mockDB, mockCfgSetter, err := sqlmock.New()
			if err != nil {
				t.Fatalf("an error '%s' was not expected when creating a mock database", err)
			}

			defer mockDB.Close()

			checkDBService := db.New(mockDB)

			mockCfgSetter.ExpectQuery("SELECT DISTINCT name FROM users").
				WillReturnRows(curRow.namesRows).
				WillReturnError(util.MakeError(curRow.errFromDB))

			names, err := checkDBService.GetUniqueNames()

			if !util.IsEmpty(curRow.errFromDB) {
				util.AssertError(t, names, err, curRow.errFromDB)

				return
			}

			if !util.IsEmpty(curRow.errFromRow) {
				util.AssertError(t, names, err, curRow.errFromRow)

				return
			}

			util.AssertNoError(t, curRow.names, names, err)

			if err := mockCfgSetter.ExpectationsWereMet(); err != nil {
				t.Fatalf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}
