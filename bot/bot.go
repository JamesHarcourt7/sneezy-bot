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
var leaderBoard map[string]int

func Start(data map[string]int) {
	// Create bot session.
	goBot, err := discordgo.New("Bot " + config.Token)
	// Specify intents
	goBot.Identify.Intents = discordgo.MakeIntent(discordgo.IntentsAll)

	if err != nil {
		fmt.Println(err.Error())
		return
	}
	leaderBoard = data

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
	goBot.AddHandler(guildCreateHandler)

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
			_, ok := leaderBoard[m.Author.ID]
			if ok {
				leaderBoard[m.Author.ID] = leaderBoard[m.Author.ID] + 10
			} else {
				leaderBoard[m.Author.ID] = 10
			}
			err = config.WriteData(leaderBoard)
			if err != nil {
				fmt.Println(err.Error())
			}
			fmt.Println("Updated leaderboard.")
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
		n := 1//0 + rand.Intn(1)
		fmt.Println("Sneezing in " + strconv.Itoa(n) + " minutes.")
		time.Sleep(time.Duration(n) * time.Minute)
		sneeze(goBot)
	}()
}

func sneeze(goBot *discordgo.Session) {
	// Send a message in a random public channel.

	fmt.Println("Achoo")
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

func guildCreateHandler(s *discordgo.Session, gu *discordgo.GuildCreate) {
	time.Sleep(5*time.Second)
	guild, err := s.State.Guild(gu.ID)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	edits := false
	for _, member := range guild.Members {
		if member.User.Bot {
			continue
		}
		if member.User.ID == s.State.User.ID {
			fmt.Println("found self")
			continue
		}
		fmt.Println(member.User.Username)
		_, ok := leaderBoard[member.User.ID]
		if !ok {
			edits = true
			fmt.Println("initialised member " + member.User.Username)
			leaderBoard[member.User.ID] = 0
		}
	}
	if edits {
		err = config.WriteData(leaderBoard)
		if err != nil {
			fmt.Println(err.Error())
			return
		}
	}
}
