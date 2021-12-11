package bot

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"math/rand"
	"os"
	"os/signal"
	"sneezy-bot/config"
	"strconv"
	"syscall"
	"time"
)

var ownId string

// Extend to multiple servers later but for now cba
var sneezed bool
var sneezeChannel string

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

	loadSneeze(goBot)

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

	if m.Content == config.BotPrefix+"bless" && sneezed {
		if m.ChannelID == sneezeChannel {
			_, err := s.ChannelMessageSend(m.ChannelID, "Thank you "+m.Author.Username+"!")
			sneezed = false
			if err != nil {
				return
			}
		}
	} else if m.Content == "achoo" {
		_, err := s.ChannelMessageSend(m.ChannelID, "Bless you "+m.Author.Username+"!")
		if err != nil {
			return
		}
	}
}

func loadSneeze(goBot *discordgo.Session) {
	// Start a random duration sleeping thread between 30 minutes and 24 hours.
	go func() {
		rand.Seed(time.Now().UnixNano())
		n := 30 + rand.Intn(1410)
		fmt.Println("Sneezing in " + strconv.Itoa(n) + " minutes.")
		time.Sleep(time.Duration(n) * time.Minute)
		sneeze(goBot)
	}()
}

func sneeze(goBot *discordgo.Session) {
	// Send a message in a random public channel.

	// Terrible please delete later.
	var guildId = goBot.State.Guilds[0].ID

	chs, err := goBot.GuildChannels(guildId)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	send, err := goBot.ChannelMessageSend(chs[rand.Intn(len(chs))].ID, "achoo")
	if err != nil {
		return
	}
	sneezeChannel = send.ChannelID
	sneezed = true
	loadSneeze(goBot)
	return
}
