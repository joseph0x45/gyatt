package db

import "github.com/jmoiron/sqlx"

type Conn struct {
	db *sqlx.DB
}

func GetConn() {

}
