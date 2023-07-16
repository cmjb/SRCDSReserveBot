package main

import "testing"

func TestParseSteamProfile(t *testing.T) {
	err, result := parseSteamProfile("tempDiscord", "https://steamcommunity.com/profiles/76561197993684328")
	t.Log(result)
	if err != nil {
		t.Error(err)
	}

	t.Cleanup(func() {
		err, _ := deleteUser("tempDiscord")
		if err != nil {
			t.Error(err)
		}
	})
}
