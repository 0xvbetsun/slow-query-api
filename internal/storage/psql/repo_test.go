package psql

import (
	"database/sql"
	"errors"
	"reflect"
	"regexp"
	"testing"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func NewGormMock(t *testing.T) (*sql.DB, *gorm.DB, sqlmock.Sqlmock) {
	dbConn, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected on sqlmock.New", err)
	}

	db, err := gorm.Open(postgres.New(postgres.Config{
		DSN:                  "sqlmock",
		DriverName:           "postgres",
		Conn:                 dbConn,
		PreferSimpleProtocol: true,
	}), &gorm.Config{})

	if err != nil {
		t.Fatalf("Failed to open gorm v2 db, got error: %v", err)
	}

	return dbConn, db, mock
}

func TestPgStatStatementRepo_GetByType(t *testing.T) {
	rows := []string{"queryid", "query", "calls", "min_exec_time",
		"max_exec_time", "mean_exec_time", "total_exec_time", "rows"}

	type args struct {
		minExecTime   uint64
		qType         string
		timeSpentSort string
		limit         int
		offset        int
	}
	tests := map[string]struct {
		args   args
		mock   func(mock sqlmock.Sqlmock)
		exp    []PgStatStatement
		expErr error
	}{
		"all_params": {
			args: args{
				minExecTime:   300,
				qType:         "SELECT",
				timeSpentSort: "ask",
				limit:         10,
			},
			mock: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "pg_stat_statements" WHERE max_exec_time >= $1 AND (query like $2 OR query like $3) ORDER BY max_exec_time ask LIMIT 10`)).
					WithArgs(300, "select%", "SELECT%").
					WillReturnRows(sqlmock.NewRows(rows).
						AddRow(20, "someQUERY", 2, 111.111, 222.222, 333.333, 777777, 5))
			},
			exp: []PgStatStatement{{
				QueryID:       20,
				Query:         "someQUERY",
				Calls:         2,
				MinExecTime:   111.111,
				MaxExecTime:   222.222,
				MeanExecTime:  333.333,
				TotalExecTime: 777777,
				Rows:          5,
			}},
		},
		"no_type_param": {
			args: args{
				minExecTime:   300,
				timeSpentSort: "ask",
				limit:         10,
			},
			mock: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "pg_stat_statements" WHERE max_exec_time >= $1 ORDER BY max_exec_time ask LIMIT 10`)).
					WithArgs(300).WillReturnRows(sqlmock.NewRows(rows))
			},
			exp: []PgStatStatement{},
		},
		"no_params": {
			mock: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "pg_stat_statements" WHERE max_exec_time >= $1`)).
					WithArgs(0).WillReturnRows(sqlmock.NewRows(rows))
			},
			exp: []PgStatStatement{},
		},
		"db_error": {
			mock: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "pg_stat_statements" WHERE max_exec_time >= $1`)).
					WithArgs(0).WillReturnError(errors.New("db_error"))
			},
			expErr: errors.New("db_error"),
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {

			con, db, mock := NewGormMock(t)
			defer con.Close()

			tt.mock(mock)

			got, err := NewPgStatStatementRepo(db).GetByType(
				tt.args.minExecTime,
				tt.args.qType,
				tt.args.timeSpentSort,
				tt.args.limit,
				tt.args.offset,
			)

			if !reflect.DeepEqual(tt.expErr, err) {
				t.Fatalf("expected err: %v, got %v", tt.expErr, err)
			}
			if !reflect.DeepEqual(tt.exp, got) {
				t.Fatalf("expected: %v, got %v", tt.exp, got)
			}
		})
	}
}
