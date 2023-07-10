package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"regexp"
	"strconv"
	"strings"
	"syscall"
	"time"
	"unsafe"

	steamidconv "github.com/Acidic9/go-steam/steamid"
	log "github.com/Sirupsen/logrus"
	"github.com/kidoman/go-steam"
)

func Monitoring_Connection(ip string) {

	c := make(chan os.Signal, 2)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	pass := config.RconPassword
	steam.SetLog(log.New())

	for {
		opts := &steam.ConnectOptions{RCONPassword: pass}
		rcon, err := steam.Connect(ip, opts)
		defer rcon.Close()
		if err != nil {
			fmt.Println("issue")
		}
		for {
			resp, err := rcon.Send("status")
			if err != nil {
				fmt.Println(err)
				break
			}
			fmt.Println(resp)
			Parse_Response_SteamID3(ip, resp)
			time.Sleep(5 * time.Second)
		}
	}

}

func Kick_User(ip string, id string) {

	pass := config.RconPassword
	opts := &steam.ConnectOptions{RCONPassword: pass}
	rcon, err := steam.Connect(ip, opts)
	defer rcon.Close()
	if err != nil {
		fmt.Println("issue")
	}
	cmd := fmt.Sprintf("kickid \"%s\" \"Sorry, but this server is reserved. Check out our discord to reserve it!\"", id)
	_, err = rcon.Send(cmd)
	if err != nil {
		fmt.Println(err)
	}
}

// STEAM_[0-5]:[01]:\d+ // Steam32
// [0-9]{17}  // Steam64

func Parse_Response_SteamID3(ip string, resp string) {
	if strings.Contains(resp, "[U:") {
		fmt.Println("Player present.")
		err, tempgroup := Get_Temp_Group_By_Server_Ip(ip)
		if err != nil {

		}
		reg := regexp.MustCompile(`\[U:1:[0-9]+\]`).FindAllString(resp, -1)
		fmt.Println(reg)
		fmt.Println(tempgroup)

		if reg != nil {
			fmt.Println(unsafe.Sizeof(tempgroup))
			if len(tempgroup.DiscordId) == 0 {
				Kick_All_Users(ip, reg)
			}
			for _, serverPlayerId := range reg {
				fmt.Println(serverPlayerId)

				playerId := tempgroup.DiscordId

				for _, playerId := range playerId {
					err, user := Get_User(playerId)
					if err != nil {
						panic(err)
					}
					sid := user.SteamId
					sid64s := strconv.FormatUint(sid, 10)
					steamId := steamidconv.NewID(sid64s)
					steamId3 := steamId.To3()
					steamstring := steamId3.String()
					fmt.Println(steamstring)
					if strings.Compare(steamstring, serverPlayerId) != 0 {
						Kick_User(ip, serverPlayerId)
					}
				}
			}
		}

	}
}

func Kick_All_Users(ip string, reg []string) {
	for _, playerId := range reg {
		Kick_User(ip, playerId)
	}
}

func test() {
	debug := flag.Bool("debug", false, "debug")
	flag.Parse()
	if *debug {
		steam.SetLog(log.New())
	}
	addr := "127.0.0.1:27015"
	pass := "testpassword"
	if addr == "" || pass == "" {
		fmt.Println("Please set ADDR & RCON_PASSWORD.")
		return
	}
	for {
		o := &steam.ConnectOptions{RCONPassword: pass}
		rcon, err := steam.Connect(addr, o)
		if err != nil {
			fmt.Println(err)
			time.Sleep(1 * time.Second)
			continue
		}
		defer rcon.Close()
		for {
			resp, err := rcon.Send("status")
			if err != nil {
				fmt.Println(err)
				break
			}
			fmt.Println(resp)
			time.Sleep(5 * time.Second)
		}
	}
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}
