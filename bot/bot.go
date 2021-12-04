package bot

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"os"
	"os/signal"
	"sneezy-bot/config"
	"syscall"
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

	// Wait here until CTRL-C or other term signal is received.
	fmt.Println("Bot is now running.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	// Cleanly close down the Discord session.
	err = goBot.Close()
	if err != nil {
		fmt.Println(err.Error())
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
	} else if m.Content == "achoo" {
		_, err := s.ChannelMessageSend(m.ChannelID, "Bless you!")
		if err != nil {
			return
		}
	}
}
