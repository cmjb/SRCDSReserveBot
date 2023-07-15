package main

import (
	"fmt"
	steamidconv "github.com/Acidic9/go-steam/steamid"
	"regexp"
	"strconv"
	"strings"
)

func convertSteamID3(steamId uint64) string {
	steamId64 := strconv.FormatUint(steamId, 10)
	steamIds := steamidconv.NewID(steamId64)
	steamId3 := steamIds.To3()
	return steamId3.String()
}

func isSteamLobby(status string) bool {
	return strings.Contains(status, "[U:")
}

func isSourceLobby(status string) bool {
	return strings.Contains(status, "STEAM_")
}

func parseSteamProfile(discordId string, url string) (error, string) {
	splitMessage := regexp.MustCompile("\\s").Split(url, 2)
	fmt.Println(splitMessage)
	steamUrlSplit := regexp.MustCompile("[/]+").Split(splitMessage[0], 4)
	fmt.Println(steamUrlSplit)
	steam64Id := steamUrlSplit[3]
	steam64idint64, err := strconv.ParseUint(steam64Id, 10, 64)
	if err != nil {
		return err, "Error"
	}
	if checkUser(discordId) {
		return err, "User exists"
	} else {
		insertUser(discordId, steam64idint64)
		return err, "Added user ID: " + steam64Id
	}
}
