package config

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
)

func LoadDBConfig() *sql.DB {
	db, err := sql.Open("sqlite3", "./database.db")
	if err != nil {
		fmt.Printf("error", err)
	}

	return db
}

func createDb() {

	db := LoadDBConfig()

	fmt.Println("生成数据")

	sql_table := `
		CREATE TABLE IF NOT EXISTS "miners" (
			"mid" INTEGER PRIMARY KEY AUTOINCREMENT,
			"Serverip" VARCHAR(255) NULL,
			"Miner" VARCHAR(255) NULL,
	       "Wallet" VARCHAR(255) NULL,
			"ServerPort" INT(10) NULL,
			"Worker" VARCHAR(255) NULL,
			"MinerIP" VARCHAR(255) NULL,
			"created" TIMESTAMP default (datetime('now','localtime'))
		);
	`
	db.Exec(sql_table)
	db.Close()

}
