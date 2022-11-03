package sqlite3

import (
	"database/sql"
	"log"
	"portscan/config"

	_ "github.com/mattn/go-sqlite3"
)

// SQLiteConn 连接 SQLite3 数据库，返回(*sql.DB)。
func SQLiteConn() *sql.DB {
	db, err := sql.Open("sqlite3", config.SqliteDBFile)
	if err != nil {
		log.Fatalln("sqlite connect error", err)
	}
	return db
}
