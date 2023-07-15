package main

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"log"
	"regexp"
)

var (
	dmPermission                   = false
	defaultMemberPermissions int64 = discordgo.PermissionManageServer
	commands                       = []*discordgo.ApplicationCommand{
		{
			Name:                     "register-server",
			Description:              "Register a server to be used for reservations.",
			DefaultMemberPermissions: &defaultMemberPermissions,
			DMPermission:             &dmPermission,
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "server-ip",
					Description: "Server IP",
					Required:    true,
				},
			},
		},
		{
			Name:                     "register-commands-in-guild",
			Description:              "Register where the bot should be present in",
			DefaultMemberPermissions: &defaultMemberPermissions,
			DMPermission:             &dmPermission,
			Options: []*discordgo.ApplicationCommandOption{

				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "guild-id",
					Description: "Given Server ID",
					Required:    true,
				},
			},
		},
		{
			Name:                     "register-steamid",
			Description:              "Register user steam id",
			DefaultMemberPermissions: &defaultMemberPermissions,
			DMPermission:             &dmPermission,
			Options: []*discordgo.ApplicationCommandOption{

				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "steamid",
					Description: "Given Steam ID",
					Required:    true,
				},
			},
		},
		{
			Name:                     "reserve-server",
			Description:              "Reserve server",
			DefaultMemberPermissions: &defaultMemberPermissions,
			DMPermission:             &dmPermission,
		},
	}

	rconCommand = discordgo.ApplicationCommand{

		Name:                     "set-rcon-password",
		Description:              "Set a rcon password",
		DefaultMemberPermissions: &defaultMemberPermissions,
		DMPermission:             &dmPermission,
		Options: []*discordgo.ApplicationCommandOption{

			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "rcon-password",
				Description: "RCON Password",
				Required:    true,
			},
		},
	}

	commandHandlers = map[string]func(session *discordgo.Session, interaction *discordgo.InteractionCreate){
		"register-server": func(session *discordgo.Session, interaction *discordgo.InteractionCreate) {
			options := interaction.ApplicationCommandData().Options

			optionMap := make(map[string]*discordgo.ApplicationCommandInteractionDataOption, len(options))

			for _, opt := range options {
				optionMap[opt.Name] = opt
			}
			serverIp := ""
			if option, ok := optionMap["server-ip"]; ok {
				serverIp = option.StringValue()
				ipRegexCompile := regexp.MustCompile(`\b(?:(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\.){3}(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?):\d{1,5}\b`).FindAllString(serverIp, -1)
				if ipRegexCompile != nil {
					err, _ := addServer(serverIp)
					if err != nil {
						session.InteractionRespond(interaction.Interaction, &discordgo.InteractionResponse{
							Type: discordgo.InteractionResponseChannelMessageWithSource,
							Data: &discordgo.InteractionResponseData{
								Content: "Unable to save ip address: " + err.Error(),
							},
						})
						return
					}
					session.InteractionRespond(interaction.Interaction, &discordgo.InteractionResponse{
						Type: discordgo.InteractionResponseChannelMessageWithSource,
						Data: &discordgo.InteractionResponseData{
							Content: "Server saved. ip: " + serverIp,
						},
					})
					channel, err := session.UserChannelCreate(interaction.Interaction.Member.User.ID)
					if err != nil {
						fmt.Println(err)
					}

					_, err = session.ChannelMessageSend(channel.ID, "Please enter a RCON password.")

					_, err = discordGlobal.ApplicationCommandCreate(discordGlobal.State.User.ID, channel.GuildID, &rconCommand)
					if err != nil {
						log.Panicf("Cannot create '%v' command: %v", rconCommand.Name, err)
					}

				} else {
					session.InteractionRespond(interaction.Interaction, &discordgo.InteractionResponse{
						Type: discordgo.InteractionResponseChannelMessageWithSource,
						Data: &discordgo.InteractionResponseData{
							Content: "Ip address and port format not recognized.",
						},
					})
					return
				}
			}
			session.InteractionRespond(interaction.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "The IP Address you chose is invalid. " + serverIp,
				},
			})

		},
		"register-commands-in-guild": func(session *discordgo.Session, interaction *discordgo.InteractionCreate) {
			session.InteractionRespond(interaction.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "Ximf is the only one that can config this at the moment.",
				},
			})
		},
		"register-steamid": func(session *discordgo.Session, interaction *discordgo.InteractionCreate) {

			options := interaction.ApplicationCommandData().Options

			optionMap := make(map[string]*discordgo.ApplicationCommandInteractionDataOption, len(options))

			for _, opt := range options {
				optionMap[opt.Name] = opt
			}
			steamId := ""
			if option, ok := optionMap["steamid"]; ok {
				steamId = option.StringValue()
			}
			fmt.Println(steamId)
			err, result := parseSteamProfile(interaction.Member.User.ID, steamId)
			if err != nil {
				fmt.Println(err)
			}
			err = session.InteractionRespond(interaction.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: result,
				},
			})
			if err != nil {
				fmt.Println(err)
				return
			}
		},
		"reserve-server": func(session *discordgo.Session, interaction *discordgo.InteractionCreate) {

			session.InteractionRespond(interaction.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "Ximf is the only one that can config this at the moment.",
				},
			})
		},
	}
)
