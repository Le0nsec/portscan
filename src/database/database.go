package database

import (
	"database/sql"
	"log"
	"portscan/database/sqlite3"
)

var DB *sql.DB

func Conn() {
	DB = sqlite3.SQLiteConn()
	log.Println("Connected to SQLite3 successfully!")
}
