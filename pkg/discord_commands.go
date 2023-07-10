package main

import (
	"github.com/bwmarrin/discordgo"
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
	}
)
