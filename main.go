package main

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"os"
	"os/signal"
	"syscall"
	"encoding/json"
)

// Fember role
var defaultRole = "692808269525286992"

func main() {
	authToken, ok := os.LookupEnv("AUTH_TOKEN")
	if !ok {
		fmt.Println("Environment variable AUTH_TOKEN is not set")
		os.Exit(1)
	}

	// Connect to Discorad
	discord, err := discordgo.New("Bot " + authToken)
	if err != nil {
		panic(err)
	}

	err = discord.Open()
	if err != nil {
		fmt.Println("[ERROR] failed to establish connection to discord")
		panic(err)
	}

	// Add message create/receive handler
	discord.AddHandler(messageCreate)
	fmt.Println("[INFO] add message create handler")

	// Give all members the default role
	//initRoles(discord)
	//fmt.Println("[INFO] completed initRoles")

	// Await termination request from the operating system
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	// Close connection to Discord
	err = discord.Close()
	if err != nil {
		panic(err)
	}
	fmt.Println("Bye")
}

/// Give all members the default role
func initRoles(s *discordgo.Session) {
	// Get list of guilds the bot if part of
	guilds, err := s.UserGuilds(10, "", "")
	if err != nil {
		fmt.Printf("[ERROR] Failed to get list of guilds the bot is part of: %+v\n", err)
		return
	}

	for _, userGuild := range guilds {
		fmt.Printf("[DEBUG][initRoles] start guild: %v(%v)\n", userGuild.Name, userGuild.ID)
		// Get detailed guild information
		guild, err := s.Guild(userGuild.ID)
		if err != nil {
			fmt.Printf("[ERROR] Failed to get guild information for guild: '%v'(%v)\n", guild.Name, guild.ID)
			continue
		}

		// Get list of guild members
		for _, member := range guild.Members {
			user := member.User
			fmt.Printf("[INFO][initRoles] member: %v(%v)\n", user.Username, user.ID)
			if user.Bot {
				// Dont' add default role to Bots
				continue
			}

			// Add default role to the member
			err := s.GuildMemberRoleAdd(userGuild.ID, user.ID, defaultRole)
			if err != nil {
				fmt.Printf("[ERROR] Failed to add role to member (id: %v): %+v\n", user.ID, err)
				continue
			}
		}

		fmt.Printf("[DEBUG][initRoles] end guild: %v(%v)\n", userGuild.Name, userGuild.ID)
	}
}

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Ignore messages created by the bot
	if m.Author.ID == s.State.User.ID {
		return
	}

	fmt.Printf("msg type: %v\n", m.Type)

	// Is this a user joining?
	if m.Type == discordgo.MessageTypeGuildMemberJoin {
		bytes, err := json.Marshal(m)
		if err != nil {
			panic(err)
		}
		fmt.Println(string(bytes))

		fmt.Println("[DEBUG] member joined message")
		// Give the new member the default role
		err = s.GuildMemberRoleAdd(m.GuildID, m.Author.ID, defaultRole)
		if err != nil {
			fmt.Printf("[ERROR] Failed to add role to new member (id: %v): %+v\n", m.Author.ID, err)
			return
		}

		fmt.Printf("[INFO] New user '%v'(%v) was granted default role\n", m.Author.Username, m.Author.ID)

		// Add reaction emoji to join message
		/* BUG: robot is a "Unkown emoji"
		err = s.MessageReactionAdd(m.ChannelID, m.ID, "robot")
		if err != nil {
			fmt.Printf("[ERROR] Failed to add reaction to member join message (member id: %v): %+v\n", m.Author.ID, err)
			return
		}*/
	}
}
