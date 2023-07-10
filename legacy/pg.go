package main

import (
	"fmt"
	"time"

	"github.com/go-pg/pg"
	"github.com/go-pg/pg/orm"
)

type User struct {
	Id        int64
	DiscordId string
	SteamId   uint64
	SteamId32 string
	Admin     bool
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

func (u User) String() string {
	return fmt.Sprintf("User<%d %s %s>", u.Id, u.DiscordId, u.SteamId32)
}

type Server struct {
	Id       int64
	ServerIp string `pg:",unique"`
	Game     string
	Owner    *User
}

type EntityError struct {
	Issue string
}

func (s Server) String() string {
	return fmt.Sprintf("Server<%d %s %s>", s.Id, s.ServerIp, s.Game)
}

func Get_DB() *pg.DB {
	db := pg.Connect(&pg.Options{
		Addr:     "10.0.0.10:5432",
		User:     "postgres",
		Password: "example",
		Database: "schedulebot",
	})
	return db
}

func Check_User(id string) bool {
	db := Get_DB()
	defer db.Close()

	var user2 User
	err := db.Model(&user2).WhereIn("discord_id IN (?)", []string{id}).Select()
	if err != nil {
		if err == pg.ErrNoRows {
			fmt.Println(id)
			fmt.Println(err)
			return false
		}
		fmt.Println(err)
		return true
	} else {
		fmt.Println(id)
		return true
	}
}

func Insert_User(id string, steamid uint64) {
	db := Get_DB()
	defer db.Close()
	user := &User{
		DiscordId: id,
		SteamId:   steamid,
		SteamId32: "",
		Admin:     false,
	}
	err := db.Insert(user)
	if err != nil {
		panic(err)
	}
}

func Get_User(id string) (error, User) {
	db := Get_DB()
	defer db.Close()

	var user User
	err := db.Model(&user).WhereIn("discord_id IN (?)", []string{id}).Limit(1).Select()

	return err, user
}

func Check_Server(ip string) bool {
	db := Get_DB()
	defer db.Close()

	var server Server
	err := db.Model(&server).WhereIn("server_ip IN (?)", []string{ip}).Select()
	if err != nil {
		if err == pg.ErrNoRows {
			fmt.Println(ip)
			fmt.Println(err)
			return false
		}
		fmt.Println(err)
		return true
	} else {
		fmt.Println(ip)
		return true
	}
}

func Get_Server(ip string) (error, Server) {
	db := Get_DB()
	defer db.Close()

	var server Server
	err := db.Model(&server).WhereIn("server_ip IN (?)", []string{ip}).Limit(1).Select()

	return err, server
}

func Get_Temp_Group_By_Server_Ip(ip string) (error, TempGroup) {
	var tempgroup TempGroup
	err, server := Get_Server(ip)
	if err != nil {
		return err, tempgroup
	}

	db := Get_DB()
	defer db.Close()

	err = db.Model(&tempgroup).Relation("Server").Where("server.id = ? AND active = TRUE", server.Id).Limit(1).Select()
	return err, tempgroup
}

func Create_Temp_Group(discordId string) (error, *TempGroup) {
	db := Get_DB()
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

func Get_All_Servers() (error, []Server) {
	db := Get_DB()
	defer db.Close()

	var servers []Server
	err := db.Model(&servers).Select()

	return err, servers
}

func Insert_Server(ip string) bool {
	db := Get_DB()
	defer db.Close()
	server := &Server{
		ServerIp: ip,
	}
	err := db.Insert(server)
	if err != nil {
		panic(err)
	} else {
		return true
	}
}

func createSchema() error {
	db := Get_DB()
	for _, model := range []interface{}{(*User)(nil), (*TempGroup)(nil), (*Server)(nil)} {
		err := db.CreateTable(model, &orm.CreateTableOptions{
			Temp: false,
		})
		if err != nil {
			return err
		}
	}
	return nil
}
