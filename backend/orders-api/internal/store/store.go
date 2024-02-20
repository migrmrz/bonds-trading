package store

import (
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/nleof/goyesql"
)

const (
	statementGetUser = "getUser"
	// statementGetActiveOrders       = "getActiveOrders"
	// statementGetActiveOrdersByUser = "getActiveOrdersByUser"
	// statementInsertOrder           = "insertOrder"
	// statementUpdateOrder           = "updateOrder"
	// statementCancelOrder           = "cancelOrder"
)

type DBConfig struct {
	Host            string        `mapstructure:"host"`
	Port            string        `mapstructure:"port"`
	DBName          string        `mapstructure:"db"`
	User            string        `mapstructure:"user"`
	MaxOpenConn     int           `mapstructure:"max-open-conn"`
	MaxIdleConn     int           `mapstructure:"max-idle-conn"`
	ConnMaxLifetime time.Duration `mapstructure:"conn-max-lifetime"`
}

func (dbc DBConfig) GetConnectionInfo() string {
	return fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=disable", dbc.Host, dbc.Port, dbc.User, dbc.DBName)
}

type DatabaseStorage struct {
	db               *sqlx.DB
	statementGetUser *sqlx.Stmt
}

type Statement struct {
	queryName goyesql.Tag
	stmt      *sqlx.Stmt
}

// GetUser response message struct
// ! probably not needed. Check this
type User struct {
	UserID         int     `db:"user_id" json:"user_id"`
	Username       string  `db:"username" json:"username"`
	HashedPassword string  `db:"hashed_password"`
	Email          string  `db:"email"`
	CreatedAt      string  `db:"created_at"`
	LastLogin      *string `db:"last_login"`
}

// type Order struct {
// 	OrderID   int       `db:"order_id" json:"order_id"`
// 	BondID    int       `db:"bond_id" json:"bond_id"`
// 	Quantity  int       `db:"quantity" json:"quantity"`
// 	Action    string    `db:"action" json:"action"`
// 	Price     float32   `db:"price" json:"price"`
// 	Status    string    `db:"status" json:"status"`
// 	Username  string    `db:"username" json:"username"`
// 	CreatedAt time.Time `db:"created_at" json:"created_at"`
// }

func NewStorage(db *sqlx.DB, queries goyesql.Queries, dbconfig *DBConfig) (*DatabaseStorage, error) {
	db.SetMaxOpenConns(dbconfig.MaxOpenConn)
	db.SetMaxIdleConns(dbconfig.MaxIdleConn)
	db.SetConnMaxLifetime(dbconfig.ConnMaxLifetime)

	statements := []Statement{
		{
			queryName: statementGetUser,
		},
	}

	for i, statement := range statements {
		stmt, err := preparex(db, queries, statement.queryName)
		if err != nil {
			return nil, err
		}

		statements[i].stmt = stmt
	}

	storage := DatabaseStorage{
		db:               db,
		statementGetUser: getStmt(statements, statementGetUser),
	}

	return &storage, nil
}

func preparex(db *sqlx.DB, queries goyesql.Queries, queryName goyesql.Tag) (*sqlx.Stmt, error) {
	stmt, err := db.Preparex(queries[queryName])
	if err != nil {
		return nil, fmt.Errorf("unable to preprare statements: %s", err.Error())
	}

	return stmt, nil
}

func getStmt(stmts []Statement, queryName goyesql.Tag) *sqlx.Stmt {
	for _, statement := range stmts {
		if queryName == statement.queryName {
			return statement.stmt
		}
	}

	return nil
}
