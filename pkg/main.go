package main

import (
	"fmt"
	"os"
)

func main() {
	fmt.Println("Loading SRCDSBot...")

	discordInit()
	httpInit()
}

func close() {
	discordGlobal.Close()
	os.Exit(1)
}
