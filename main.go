package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
)

func main() {
	key := os.Getenv("DISCORD_KEY")
	discord, err := discordgo.New("Bot " + key)
	if err != nil {
		fmt.Println("Error creating Discord session,", err)
		return
	}

	discord.AddHandler(messageCreate)
	discord.Identify.Intents = discordgo.MakeIntent(discordgo.IntentsGuildMessages)

	err = discord.Open()
	if err != nil {
		fmt.Println("Error opening connection,", err)
		return
	}

	fmt.Println("Bot is now running. Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	discord.Close()
}

func messageCreate(session *discordgo.Session, msg *discordgo.MessageCreate) {
	if msg.Author.ID == session.State.User.ID {
		return
	}

	var reply string

	// TODO: Move this or the entire function to a separate file?
	switch msg.Content {
	case "good morning":
		reply = "<:kys:620483919774744596>"
	case "gsgo":
		reply = "<@!199174365408002049> <@!199462593264484352> <@!357042093199196162> <@!199917417953099778> <@!379870473481355264>"
	}

	session.ChannelMessageSend(msg.ChannelID, reply)
}
