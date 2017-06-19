package main

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
)

func kek() {
	discord, _ := discordgo.New(fmt.Sprintf("Bot %s", config.DiscordAuthToken))
	discord.AddHandler(func(s *discordgo.Session, e *discordgo.MessageCreate) {
	})
}
