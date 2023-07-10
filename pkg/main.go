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
	httpInit()
}

func dead() {
	discordGlobal.Close()
	os.Exit(1)
}
