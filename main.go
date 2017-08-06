package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"os/signal"
	"sync"
	"time"

	"github.com/DavidHuie/gomigrate"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"github.com/santiclause/goconfig"
)

type Config struct {
	DSN                string        `env:"MYSQL_DSN" yaml:"mysql_dsn"`
	MigrationsLocation string        `env:"MIGRATIONS_LOCATION" yaml:"migrations_location"`
	DiscordAuthToken   string        `env:"DISCORD_AUTH_TOKEN" yaml:"discord_auth_token"`
	Version            string        `env:"VERSION" yaml:"version"`
	IrcServers         []string      `env:"IRC_SERVERS" yaml:"irc_servers"`
	IrcChannels        []string      `env:"IRC_CHANNELS" yaml:"irc_channels"`
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

	fmt.Println("Hello!")
	fmt.Printf("List of servers: %v\n", config.IrcServers)
	fmt.Printf("List of channels: %v\n", config.IrcChannels)

	var servers []*IrcConn
	for _, server := range config.IrcServers {
		conn, err := Connect(
			server,
			config.IrcNickname,
			WithAutojoinChannels(config.IrcChannels),
			WithIdent(config.IrcIdent),
			WithName(config.IrcName),
			WithNickservPassword(config.IrcNickservPass),
			WithNickservTimeout(config.IrcNickservTimeout),
			WithVersion(config.Version),
			WithQuitMessage(config.IrcQuitMessage),
		)
		if err == nil {
			servers = append(servers, conn)
		} else {
			log.Printf("Shit's fucked, failed to connect. %s\n", err)
		}
	}

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt)

	<-sig
	wait := sync.WaitGroup{}
	fmt.Println("Closing...")
	for _, server := range servers {
		if !server.conn.Connected() {
			continue
		}
		wait.Add(1)
		go (func() {
			timeout := time.After(15 * time.Second)
			select {
			case <-server.Close():
				wait.Done()
			case <-timeout:
				wait.Done()
			}
		})()
	}
	wait.Wait()
	fmt.Println("Goodbye!")

	// var user models.User
	// err = db.First(&user).Error
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// err = user.GetPermissions(db)
	// fmt.Printf("%v\n", user.Permissions)
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
