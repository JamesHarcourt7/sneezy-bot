package main

import (
	"fmt"
	"sneezy-bot/bot"
	"sneezy-bot/config"
)

func main() {
	err := config.ReadConfig()
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	bot.Start()
	<-make(chan struct{})

	return
}
