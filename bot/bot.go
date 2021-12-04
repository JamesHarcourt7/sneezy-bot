package bot

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"sneezy-bot/config"
)

var ownId string

func Start() {
	// Create bot session.
	goBot, err := discordgo.New("Bot " + config.Token)

	if err != nil {
		fmt.Println(err.Error())
		return
	}

	// Make bot a user using User function
	u, err := goBot.User("@me")
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	// Store ID in ownID
	ownId = u.ID

	// Add message handler function
	goBot.AddHandler(messageHandler)

	err = goBot.Open()
	if err != nil {
		fmt.Println(err.Error())
		return
	}
}

func messageHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Don't reply to own messages.
	if m.Author.ID == ownId {
		return
	}

	if m.Content == config.BotPrefix+"bless" {
		_, err := s.ChannelMessageSend(m.ChannelID, "Thank you.")
		if err != nil {
			return
		}
	}
}
