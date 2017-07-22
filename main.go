package main

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/DavidHuie/gomigrate"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"github.com/santiclause/eden/models"
	"github.com/santiclause/goconfig"
)

type Config struct {
	DSN                string        `env:"MYSQL_DSN" yaml:"mysql_dsn"`
	MigrationsLocation string        `env:"MIGRATIONS_LOCATION" yaml:"migrations_location"`
	DiscordAuthToken   string        `env:"DISCORD_AUTH_TOKEN" yaml:"discord_auth_token"`
	DefaultPrefix      string        `env:"DEFAULT_PREFIX" yaml:"default_prefix"`
	Version            string        `env:"VERSION" yaml:"version"`
	IrcServers         []string      `env:"IRC_SERVERS" yaml:"irc_servers"`
	IrcNickname        string        `env:"IRC_NICKNAME" yaml:"irc_nickname"`
	IrcIdent           string        `env:"IRC_IDENT" yaml:"irc_ident"`
	IrcName            string        `env:"IRC_NAME" yaml:"irc_name"`
	IrcNickservPass    string        `env:"IRC_NICKSERV_PASS" yaml:"irc_nickserv_pass"`
	IrcNickservTimeout time.Duration `env:"IRC_NICKSERV_TIMEOUT" yaml:"irc_nickserv_timeout"`
	IrcQuitMessage     string        `env:"IRC_QUIT_MESSAGE" yaml:"irc_quit_message"`
	goconfig.Config
}

var (
	config = Config{
		MigrationsLocation: "migrations",
		IrcNickservTimeout: 15 * time.Second,
	}
	db *gorm.DB
)

func main() {
	config.SetFilename("config.yaml")
	goconfig.Load(&config)
	goconfig.ListenForSignals(&config)
	migrationDB, err := sql.Open("mysql", config.DSN)
	if err != nil {
		log.Fatal(err)
	}
	migrator, err := gomigrate.NewMigrator(migrationDB, gomigrate.Mysql{}, config.MigrationsLocation)
	if err != nil {
		log.Fatal(err)
	}
	err = migrator.Migrate()
	if err != nil {
		log.Fatal(err)
	}
	migrationDB.Close()
	db, err = gorm.Open("mysql", config.DSN)
	if err != nil {
		log.Fatal(err)
	}

	if config.DebugLevel("verbose") {
		db.LogMode(true)
	}

	var servers []*IrcConn
	for server := range config.IrcServers {
		conn := Connect(
			server,
			config.IrcNickname,
			WithIdent(config.IrcIdent),
			WithName(config.IrcName),
			WithVersion(config.Version),
			WithQuitMessage(config.IrcQuitMessage),
		)
		servers = append(servers, conn)
	}

	var user models.User
	err = db.First(&user).Error
	if err != nil {
		log.Fatal(err)
	}
	err = user.GetPermissions(db)
	fmt.Printf("%v\n", user.Permissions)
	// var users []models.User
	// err = db.Find(&users).Error
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// for i, user := range users {
	// 	var roles []models.Role
	// 	// err := db.Model(&user).Related(&roles, "Roles").Error
	// 	err := db.Model(&user).Association("Roles").Find(&roles).Error
	// 	if err != nil {
	// 		log.Fatal(err)
	// 	}
	// 	fmt.Printf("%d: %v\n", i, roles)
	// }
}
