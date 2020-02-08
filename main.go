package main

import (
	"encoding/json"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"regexp"
	"strconv"
	"net"
)

type Response struct {
	Status string
}

var discordGlobal *discordgo.Session

type conf struct {
	Token       string  `yaml:"token"`
	Channel     string  `yaml:"channel"` 
	LobbyRequestChannel     string  `yaml:"lobby_request_channel"`
	LobbyStatusChannel     string  `yaml:"lobby_status_channel"`
	ModeratorId			   string `yaml:"moderator_id"`
	GuildId			   string `yaml:"guild_id"`
	RconPassword		string `yaml:"rcon_password"`
}


var config conf

func main() {

	configFile, err := ioutil.ReadFile("config.yml")
	if err != nil {
		panic(err)
	}

	err = yaml.Unmarshal(configFile, &config)
	if err != nil {
		panic(err)
	}

	discord, err := discordgo.New("Bot " + config.Token)
	if err != nil {
		panic(err)
	}

	discordGlobal = discord
	
	discord.UpdateStatus(0, ".help for help")
	
	err = discordGlobal.Open()
	if err != nil {
		fmt.Println(err)
	}
	sendStatusMessage(discord, "MahoNet online. Build: Development")
	discord.AddHandler(messageCreate)
	discord.AddHandler(ready)

	err = Start_Up()

	if err != nil {
		fmt.Println("Can't connect to all servers, please do .startup in discord to re-startup")
	}

	fmt.Println("Started.")
	http.HandleFunc("/test/", handleTest)
	log.Fatal(http.ListenAndServe("127.0.0.1:8087", nil))
}

func Start_Up() error {
	err, servers := Get_All_Servers();
	if(err != nil) {
		fmt.Println(err)
	}
	for _, server := range servers {
		Monitoring_Connection(server.ServerIp)
	}
	
	return err
}

func ready(s *discordgo.Session, event *discordgo.Ready) {
	s.UpdateStatus(0, ".help for info")
}

func sendMessage(s *discordgo.Session, msg string) {
	s.ChannelMessageSend(config.Channel, msg)
}

func sendChannelMessage(s *discordgo.Session, c *discordgo.Channel, msg string) {
	s.ChannelMessageSend(c.ID, msg)
}

func sendStatusMessage(s *discordgo.Session, msg string) {
	s.ChannelMessageSend(config.LobbyStatusChannel, msg)
}

func handleTest(w http.ResponseWriter, r *http.Request) {

	sendStatusMessage(discordGlobal, "Test successful!");

	jsonRes := &Response{
		Status: "OK",
	}
	jsonString, _ := json.Marshal(jsonRes)
	fmt.Println(string(jsonString));
	fmt.Fprintf(w, "%s", string(jsonString));
}


func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {

	if m.Author.ID == s.State.User.ID {
		return
	}

	c, err := s.State.Channel(m.ChannelID)

	if strings.HasPrefix(m.Content, ".addserver") {

		member := m.Member
		
		err, isModerator := Is_User_Moderator(member, c, s)

		if err != nil {
			return
		}

		if isModerator {

			splitMessage := regexp.MustCompile("\\s").Split(m.Content, 2)
			fmt.Println(splitMessage[1])
			if len(splitMessage) > 0 {
				if Check_Server(splitMessage[1]) {
					sendChannelMessage(discordGlobal, c, "Server exists.")
				} else {
					host, _, err2 := net.SplitHostPort(splitMessage[1])
					err := net.ParseIP(host)
					if err != nil {
						Insert_Server(splitMessage[1])
						sendChannelMessage(discordGlobal, c, "Added Server to pool.")
					} else {
						sendChannelMessage(discordGlobal, c, "Invalid IP Address.")
						fmt.Println(err)
						fmt.Println(err2)
					}

				}
			}
		} else {
			sendChannelMessage(discordGlobal, c, "Sorry. You aren't a Moderator and cannot access this command.")
		}
	}

	if strings.HasPrefix(m.Content, ".steam") {

		splitMessage := regexp.MustCompile("\\s").Split(m.Content, 2)
		
		if len(splitMessage) > 1 {
			if strings.Contains(splitMessage[1], "https://steamcommunity.com/profiles/") {
				sendChannelMessage(discordGlobal, c, "SteamId found.")
				steamUrlSplit := regexp.MustCompile("[/]+").Split(m.Content, 4)
				fmt.Println(steamUrlSplit)
				steam64Id := steamUrlSplit[3]
				fmt.Println(steam64Id)
				steam64idint64, err := strconv.ParseUint(steam64Id, 10, 64)
				if err != nil {
					sendChannelMessage(discordGlobal, c, "No SteamId found, please use your steam community profiles url.")
					fmt.Println(err)
				}
				if Check_User(m.Author.ID) {
					fmt.Println("User exists.")
					sendChannelMessage(discordGlobal, c, "You already have a steamid set! Ask a Moderator to remove it for you.")
				} else {
					Insert_User(m.Author.ID, steam64idint64)
					sendChannelMessage(discordGlobal, c, "Assigned Steam ID to user.")
				}
				
			} else {
				sendChannelMessage(discordGlobal, c, "No SteamId found, please use your steam community profiles url.")
			}
		} else {
			sendChannelMessage(discordGlobal, c, "No SteamId found, please use your steam community profiles url.")
		}

	}

	if strings.HasPrefix(m.Content, ".-build-schema") {
		

		if err != nil {
			sendStatusMessage(discordGlobal, "Panicking! Can't build schema because I can't see the channel!")
			return
		}

		member := m.Member

		err, isModerator := Is_User_Moderator(member, c, s)

		if err != nil {
			return
		}

		if isModerator {
			err := createSchema()
			if err != nil {
				sendStatusMessage(discordGlobal, "Error.")
				sendChannelMessage(discordGlobal, c, "Issue confirmed.")
				fmt.Println(err)
			} else {
				sendStatusMessage(discordGlobal, "Yay it worked.")
				sendChannelMessage(discordGlobal, c, "Schema built.")
			}
		} else {
			sendChannelMessage(discordGlobal, c, "Sorry. You aren't a Moderator and cannot access this command.")
		}

	}

	if strings.HasPrefix(m.Content, ".-restart-up") {
		

		if err != nil {
			sendStatusMessage(discordGlobal, "Panicking! Can't build schema because I can't see the channel!")
			return
		}

		member := m.Member

		err, isModerator := Is_User_Moderator(member, c, s)

		if err != nil {
			return
		}

		if isModerator {
			err := Start_Up()
			if err != nil {
				sendStatusMessage(discordGlobal, "Error.")
				sendChannelMessage(discordGlobal, c, "Issue confirmed.")
				fmt.Println(err)
			} else {
				sendStatusMessage(discordGlobal, "Yay it worked.")
				sendChannelMessage(discordGlobal, c, "Restarted")
			}
		} else {
			sendChannelMessage(discordGlobal, c, "Sorry. You aren't a Moderator and cannot access this command.")
		}

	}

	if strings.HasPrefix(m.Content, ".-build-schema2") {
		

		if err != nil {
			sendStatusMessage(discordGlobal, "Panicking! Can't build schema because I can't see the channel!")
			return
		}

		member := m.Member

		err, isModerator := Is_User_Moderator(member, c, s)

		if err != nil {
			return
		}

		if isModerator {
			err := createSchema()
			if err != nil {
				sendStatusMessage(discordGlobal, "Error.")
				sendChannelMessage(discordGlobal, c, "Issue confirmed.")
				fmt.Println(err)
			} else {
				sendStatusMessage(discordGlobal, "Yay it worked.")
				sendChannelMessage(discordGlobal, c, "Schema built.")
			}
		} else {
			sendChannelMessage(discordGlobal, c, "Sorry. You aren't a Moderator and cannot access this command.")
		}

	}
}

func Is_User_Moderator(member *discordgo.Member, channel *discordgo.Channel, session *discordgo.Session) (error, bool) {

	g, err := session.State.Guild(channel.GuildID)
	if err != nil {
		sendStatusMessage(discordGlobal, "Panicking! Can't build schema because I can't see the server!")
		return err, false
	}

	for _, roleID := range member.Roles {
		role, err := session.State.Role(g.ID, roleID)
		if err != nil {
			sendStatusMessage(discordGlobal, "Panicking! Role issue!")
			return err, false
		}
		if strings.Compare(role.ID, config.ModeratorId) == 0 {
			return nil, true
		} else {
			return nil, false
		}
	}
	return nil, false
}
