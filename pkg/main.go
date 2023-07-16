package main

import (
	"fmt"
	"os"
)

func main() {
	fmt.Println("Loading SRCDSBot...")

	err := createSchema()
	if err != nil {
		fmt.Println(err)
		fmt.Println("Error creating schema, maybe it already exists?")
	}
	discordInit()

	err = startUp()
	if err != nil {
		fmt.Println("Issue starting up.")
		fmt.Println(err)
		dead()
	}
	httpInit()
}

func startUp() error {
	err, servers := getAllServers()

	for _, server := range servers {
		heartbeat(server.ServerIp, server.Password)
	}

	return err
}

func dead() {
	discordGlobal.Close()
	os.Exit(1)
}
