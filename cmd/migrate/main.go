package main

import (
	"flag"
	"fmt"
	"log"

	. "ruuviLogger"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {
	down := flag.Bool("down", false, "apply all migrations down; otherwise up")

	flag.Parse()

	m, err := migrate.New(
		"file://db/migrations",
		fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable",
			Cfg.DbUser, Cfg.DbPwd, Cfg.DbHost, Cfg.DbPort, Cfg.DbName))
	if err != nil {
		log.Fatal(err)
	}

	if *down {
		if err := m.Down(); err != nil {
			log.Fatal(err)
		}
	} else {
		if err := m.Up(); err != nil {
			log.Fatal(err)
		}
	}
}
