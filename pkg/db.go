package main

import (
	"errors"
	"fmt"
	"github.com/go-pg/pg"
	"github.com/go-pg/pg/orm"
	"os"
	"time"
)

type User struct {
	Id          int64
	DiscordId   string
	SteamId     uint64
	SteamId32   string
	Admin       bool
	CreatedTime time.Time
}

type Guild struct {
	Id          int64 `pg:",unique"`
	GuildId     string
	Active      bool
	Name        string
	CreatedTime time.Time
}

type Server struct {
	Id       int64
	ServerIp string `pg:",unique"`
	Game     string
	Owner    *User
	Password string
}

type TempGroup struct {
	Id           int64
	DiscordId    []string
	Server       *Server
	Owner        string
	Active       bool
	ReservedTime time.Time
	ExpiryTime   time.Time
}

func db() *pg.DB {

	Host := os.Getenv("DB_HOST")
	if Host == "" {
		fmt.Println("Database host missing. Terminating...")
		dead()
	}

	Name := os.Getenv("DB_NAME")
	if Name == "" {
		fmt.Println("Database name missing. Terminating...")
		dead()
	}

	User := os.Getenv("DB_USER")
	if User == "" {
		fmt.Println("Database user missing. Terminating...")
		dead()
	}

	Pass := os.Getenv("DB_PASS")
	if Pass == "" {
		fmt.Println("Database pass missing. Terminating...")
		dead()
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

func getServer(ip string) (error, *Server) {
	db := db()
	defer db.Close()

	var server Server
	err := db.Model(&server).WhereIn("server_ip IN (?)", []string{ip}).Limit(1).Select()

	return err, &server

}

func getAllServers() (error, []Server) {
	db := db()
	defer db.Close()

	var servers []Server
	err := db.Model(&servers).Select()

	return err, servers
}

func addServer(ip string) (error, *Server) {
	db := db()
	defer db.Close()

	_, existingServer := getServer(ip)

	if existingServer.ServerIp == ip {
		return errors.New("duplicate server exists"), existingServer
	}

	server := &Server{
		ServerIp: ip,
	}

	err := db.Insert(server)

	return err, server

}

func getTempGroupByServerIp(ip string) (error, TempGroup) {
	var tempgroup TempGroup
	err, server := getServer(ip)
	if err != nil {
		return err, tempgroup
	}

	db := db()
	defer db.Close()

	err = db.Model(&tempgroup).Relation("Server").Where("server.id = ? AND active = TRUE", server.Id).Limit(1).Select()
	return err, tempgroup
}

func getUser(id string) (error, User) {
	db := db()
	defer db.Close()

	var user User
	err := db.Model(&user).WhereIn("discord_id IN (?)", []string{id}).Limit(1).Select()

	return err, user
}

func createTempGroup(discordId string) (error, *TempGroup) {
	db := db()
	defer db.Close()
	nowtime := time.Now()
	expirytime := nowtime.Add(time.Hour * 3)
	tempgroup := &TempGroup{
		Owner:        discordId,
		Active:       true,
		ReservedTime: nowtime,
		ExpiryTime:   expirytime,
	}
	err := db.Insert(tempgroup)
	return err, tempgroup
}

func getNumberOfOnlineServers() (error, int) {
	db := db()
	defer db.Close()

	var servers []Server
	err := db.Model(&servers).Where("active = TRUE").Select()

	return err, len(servers)
}

func createSchema() error {
	db := db()
	for _, model := range []interface{}{(*User)(nil), (*Guild)(nil), (*TempGroup)(nil), (*Server)(nil)} {
		err := db.CreateTable(model, &orm.CreateTableOptions{
			Temp: false,
		})
		if err != nil {
			return err
		}
	}
	return nil
}
