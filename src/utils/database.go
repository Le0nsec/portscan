package utils

import (
	"database/sql"
	"log"
	"os"
	"portscan/config"
)

// DBInit 实现初始化数据库等功能。
func DBInit() {
	if _, err := os.Stat(config.DataPath); os.IsNotExist(err) {
		err := os.MkdirAll(config.DataPath, 0755)
		if err != nil {
			log.Fatalf("mkdir error: %s", err.Error())
		}
	}
	// 判断是否存在 sqlite3 数据库文件，不存在则执行初始化，存在则不执行
	if _, err := os.Stat(config.SqliteDBFile); os.IsNotExist(err) {
		log.Println("sqlite3 initialization...")

		db, err := sql.Open("sqlite3", config.SqliteDBFile)
		if err != nil {
			log.Fatalln("sqlite connect error: ", err)
		}
		log.Println("create db file success!")

		_, err = db.Exec(sqlite3_sql)
		if err != nil {
			log.Fatalln("db init error:", err)
		}
		log.Println("db init success!")
		return
	} else {
		log.Println("sqlite3 db file already exist!")
		return
	}

}

var sqlite3_sql = `
CREATE TABLE "scan" (
	"id"	INTEGER NOT NULL UNIQUE,
	"host"	TEXT NOT NULL,
	"port"	TEXT NOT NULL,
	"threads"	INTEGER NOT NULL,
	"timeout"	INTEGER NOT NULL,
	"created"	INTEGER NOT NULL,
	"status"	INTEGER NOT NULL,
	PRIMARY KEY("id" AUTOINCREMENT)
);
CREATE TABLE "detail" (
	"id"	INTEGER NOT NULL UNIQUE,
	"scan_id"	INTEGER NOT NULL UNIQUE,
	"json_str"	TEXT,
	PRIMARY KEY("id" AUTOINCREMENT)
);
`
