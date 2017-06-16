package main

import (
	"database/sql"
	"log"

	"github.com/DavidHuie/gomigrate"
	_ "github.com/go-sql-driver/mysql"
	"github.com/santiclause/goconfig"
)

type Config struct {
	DSN                string `env:"MYSQL_DSN" yaml:"mysql_dsn"`
	MigrationsLocation string `env:"MIGRATIONS_LOCATION" yaml:"migrations_location"`
	goconfig.Config
}

var (
	config = Config{
		MigrationsLocation: "migrations",
	}
)

func main() {
	config.SetFilename("config.yaml")
	goconfig.Load(&config)
	goconfig.ListenForSignals(&config)
	db, err := sql.Open("mysql", config.DSN)
	if err != nil {
		log.Fatal(err)
	}
	migrator, err := gomigrate.NewMigrator(db, gomigrate.Mysql{}, config.MigrationsLocation)
	if err != nil {
		log.Fatal(err)
	}
	err = migrator.Migrate()
	if err != nil {
		log.Fatal(err)
	}
}
