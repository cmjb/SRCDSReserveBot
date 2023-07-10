package main

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"log"
	"os"
	"os/signal"
	"syscall"
)

var discordGlobal *discordgo.Session

func discordInit() {
	fmt.Println("Discord Initalizing...")

	Token := os.Getenv("DISCORD_TOKEN")
	if Token == "" {
		fmt.Println("Discord token missing. Terminating...")
		return
	}

	discord, err := discordgo.New("Bot " + Token)
	if err != nil {
		fmt.Println(fmt.Sprintf("Error initialising Discord: %s", err))
		return
	}

	discordGlobal = discord

	err = discord.Open()
	if err != nil {
		fmt.Println(fmt.Sprintf("Error initialising Discord: %s", err))
		return
	}
	fmt.Println("Discord Initalized.")
	updateStatus()
	registerCommands()
	loadGuildCommandsViaDB()
	// go routine so we're not hanging on waiting for signal.
	go func() {
		catchSigInt()
	}()

}

func registerCommands() {
	fmt.Println("Instancing commands...")
	discordGlobal.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		if h, ok := commandHandlers[i.ApplicationCommandData().Name]; ok {
			h(s, i)
		}
	})

}

func loadGuildCommandsViaDB() {
	fmt.Println("Loading saved guilds...")
	err, guilds := getGuilds()
	if err != nil {
		log.Panicf("Issue loading guilds %v", err)
	}

	for _, guild := range guilds {
		fmt.Println(guild.GuildId)
		registerCommandsForGuild(guild.GuildId)
	}
}

func registerCommandsForGuild(guildId string) {
	registeredCommands := make([]*discordgo.ApplicationCommand, len(commands))
	for i, v := range commands {
		cmd, err := discordGlobal.ApplicationCommandCreate(discordGlobal.State.User.ID, guildId, v)
		if err != nil {
			log.Panicf("Cannot create '%v' command: %v", v.Name, err)
		}
		registeredCommands[i] = cmd
	}
}

func updateStatus() {
	err := discordGlobal.UpdateGameStatus(0, "0 servers online")
	if err != nil {
		fmt.Println(fmt.Sprintf("Error updating Discord status: %s", err))
		return
	}
}

func catchSigInt() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
	for _ = range c {
		fmt.Println("Cleaning up Discord...")
		err := discordGlobal.Close()
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println("Terminating...")
		os.Exit(0)
	}

}
