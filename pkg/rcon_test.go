package main

import (
	"github.com/cmjb/go-steam"
	log "github.com/sirupsen/logrus"
	"os"
	"testing"
)

func TestRCON(t *testing.T) {

	serverIp := os.Getenv("TEST_SERVER_IP")
	if serverIp == "" {
		t.Error("No test server ip selected.")
		return
	}

	serverPass := os.Getenv("TEST_SERVER_PASS")
	if serverPass == "" {
		t.Error("No test server password selected.")
		return
	}
	steam.SetLog(log.New())

	opts := &steam.ConnectOptions{RCONPassword: serverPass}
	rcon, err := steam.Connect(serverIp, opts)
	if err != nil {
		t.Error(err)
	}

	response, err := rcon.Send("status")
	if err != nil {
		t.Error(err)
	}

	t.Log(response)
	if err != nil {
		t.Error(err)
	}

	t.Cleanup(func() {
		rcon.Close()
	})
}
