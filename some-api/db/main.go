package db

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
	migrate "github.com/rubenv/sql-migrate"
)

var Connection *sql.DB

// Connect -
func Connect() *sql.DB {
	var useMigrationsGo = true
	var err error
	if Connection != nil {
		return Connection
	}

	Connection, err = sql.Open("mysql", "root:admin@tcp(localhost:3307)/pep?parseTime=true")

	err = Connection.Ping()
	if err != nil {
		log.Fatal(err)
	}

	if useMigrationsGo {

		migrations := &migrate.FileMigrationSource{
			Dir: "db/migrations",
		}

		//Executa os migraions carregados
		n, err := migrate.Exec(Connection, "mysql", migrations, migrate.Up)
		if err != nil {
			_, _ = migrate.Exec(Connection, "mysql", migrations, migrate.Down)
			log.Fatal(err)
		}
		fmt.Printf("Applied %d migrations!\n", n)
	}

	return Connection
}

func CloseConn() {
	Connection.Close()
}
