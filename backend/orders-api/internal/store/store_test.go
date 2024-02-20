package store

import (
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/nleof/goyesql"
	"github.com/stretchr/testify/assert"
)

func TestNewStore(t *testing.T) {
	mockDB, mock, _ := sqlmock.New()
	sqlxDB := sqlx.NewDb(mockDB, "sqlmock")

	mock.ExpectPrepare("SELECT * ")

	defer sqlxDB.Close()

	queries, err := goyesql.ParseFile("./testdata/statements.sql")
	if err != nil {
		t.Error("error getting queries")

		return
	}

	cases := []struct {
		name          string
		db            *sqlx.DB
		config        *DBConfig
		queries       goyesql.Queries
		expectedStore *DatabaseStorage
		expectedError error
	}{
		{
			name: "success",
			db:   sqlxDB,
			config: &DBConfig{
				MaxOpenConn:     10,
				MaxIdleConn:     10,
				ConnMaxLifetime: time.Hour,
			},
			queries:       queries,
			expectedStore: &DatabaseStorage{},
			expectedError: nil,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			actualStore, actualError := NewStorage(tc.db, tc.queries, tc.config)
			assert.IsType(t, tc.expectedStore, actualStore)
			assert.Equal(t, tc.expectedError, actualError)
		})
	}
}
