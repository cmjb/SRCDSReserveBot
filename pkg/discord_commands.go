package main

import (
	"github.com/bwmarrin/discordgo"
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

			session.InteractionRespond(interaction.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "This isn't implemented yet.",
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
