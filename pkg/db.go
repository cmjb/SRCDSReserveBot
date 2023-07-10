package main

import (
	"fmt"
	"github.com/go-pg/pg"
	"os"
)

type User struct {
	Id        int64
	DiscordId string
	SteamId   uint64
	SteamId32 string
	Admin     bool
	CreatedAt int64
}

type Guild struct {
	Id        int64 `pg:",unique"`
	Active    bool
	Name      string
	CreatedAt int64
}

type Server struct {
	Id       int64
	ServerIp string `pg:",unique"`
	Game     string
	Owner    *User
}

func db() *pg.DB {

	Host := os.Getenv("DB_HOST")
	if Host == "" {
		fmt.Println("Database host missing. Terminating...")
		close()
	}

	Name := os.Getenv("DB_NAME")
	if Name == "" {
		fmt.Println("Database name missing. Terminating...")
		close()
	}

	User := os.Getenv("DB_USER")
	if User == "" {
		fmt.Println("Database user missing. Terminating...")
		close()
	}

	Pass := os.Getenv("DB_PASS")
	if Pass == "" {
		fmt.Println("Database pass missing. Terminating...")
		close()
	}

	return pg.Connect(&pg.Options{
		Addr:     Host,
		User:     User,
		Password: Pass,
		Database: Name,
	})
}

func getGuilds() (error, []Guild) {
	db := db()
	defer db.Close()

	var guilds []Guild
	err := db.Model(&guilds).Select()

	return err, guilds
}

func addServer(ip string) *Server {
	db := db()
	defer db.Close()
	server := &Server{
		ServerIp: ip,
	}
	err := db.Insert(server)
	if err != nil {
		panic(err)
	}

	return server

}
