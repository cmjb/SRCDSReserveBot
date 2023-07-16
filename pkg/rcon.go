package main

import (
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/kidoman/go-steam"
	"regexp"
	"strings"
	"sync"
	"time"
)

type ServerInstance struct {
	IP   string
	sync sync.Mutex
}

func heartbeat(ip string, pass string) {
	var server ServerInstance
	server.IP = ip
	server.sync.Lock()
	defer server.sync.Unlock()
	steam.SetLog(log.New())
	for {

		opts := &steam.ConnectOptions{RCONPassword: pass}
		rcon, err := steam.Connect(ip, opts)
		if err != nil {
			fmt.Println(err)
		}
		defer rcon.Close()

		for {
			response, err := rcon.Send("status")
			if err != nil {
				fmt.Println(err)
			}
			scanPlayers(ip, response, pass)
			time.Sleep(5 * time.Second)
		}
	}
}

func scanPlayers(ip string, response string, pass string) {
	if isSteamLobby(response) {
		regEx := regexp.MustCompile(`\[U:1:[0-9]+`).FindAllString(response, -1)
		err, tempgroup := getTempGroupByServerIp(ip)
		if err != nil {
			fmt.Println(err)
		}

		if regEx != nil {
			if len(tempgroup.DiscordId) == 0 {
				kickAllUsers(ip, regEx, pass)
			}

			for _, serverPlayerId := range regEx {
				playerIds := tempgroup.DiscordId

				for _, playerId := range playerIds {
					err, user := getUser(playerId)
					if err != nil {
						fmt.Println(err)
					}
					steamId := convertSteamID3(user.SteamId)
					if strings.Compare(steamId, serverPlayerId) != 0 {
						kickUser(ip, serverPlayerId, pass)
					}
				}
			}
		}
	}
}

func kickAllUsers(ip string, regEx []string, pass string) {
	for _, playerId := range regEx {
		kickUser(ip, playerId, pass)
	}
}

func kickUser(ip string, playerId string, pass string) {
	opts := &steam.ConnectOptions{RCONPassword: pass}
	rcon, err := steam.Connect(ip, opts)
	defer rcon.Close()
	if err != nil {
		fmt.Println(err)
	}

	cmd := fmt.Sprintf("kickid \"%s\" \"Sorry, but this server is reserved. Check out our discord to reserve it!\"", playerId)
	_, err = rcon.Send(cmd)
	if err != nil {
		fmt.Println(err)
	}
}

func runHeartBeatOnAllServers() {
	err, servers := getAllServers()
	if err != nil {
		return
	}
	for _, server := range servers {
		heartbeat(server.ServerIp, server.Password)
	}
}

func runHeartBeatOnServer(ip string) {

}
