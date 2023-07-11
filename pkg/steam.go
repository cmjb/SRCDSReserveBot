package main

import (
	steamidconv "github.com/Acidic9/go-steam/steamid"
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
