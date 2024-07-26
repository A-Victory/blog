package conn

import "github.com/A-Victory/blog/database"

type DB struct {
	Conn *database.DBconn
}

func NewConn(conn *database.DBconn) *DB {
	return &DB{
		Conn: conn,
	}
}
