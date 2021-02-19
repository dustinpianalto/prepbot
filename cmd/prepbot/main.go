package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
	"github.com/dustinpianalto/disgoman"
	"github.com/dustinpianalto/prepbot/internal/exts/members"
	"github.com/dustinpianalto/prepbot/internal/exts/messages"
)

func main() {
	Token := os.Getenv("DISCORD_TOKEN")
	dg, err := discordgo.New("Bot " + Token)
	if err != nil {
		fmt.Println("There was an error when creating the Discord Session, ", err)
		return
	}
	dg.State.MaxMessageCount = 100
	dg.StateEnabled = true

	dg.Identify = discordgo.Identify{
		Intents: discordgo.MakeIntent(discordgo.IntentsAll),
	}

	owners := []string{
		"351794468870946827",
	}

	manager := disgoman.CommandManager{
		Prefixes:         getPrefixes,
		Owners:           owners,
		StatusManager:    disgoman.GetDefaultStatusManager(),
		ErrorChannel:     make(chan disgoman.CommandError, 10),
		Commands:         make(map[string]*disgoman.Command),
		IgnoreBots:       true,
		CheckPermissions: false,
	}

	dg.AddHandler(messages.CleanAmazonURLs)
	dg.AddHandler(members.OnGuildMemberAddLogging)
	dg.AddHandler(members.OnGuildMemberRemoveLogging)

	err = dg.Open()
	if err != nil {
		fmt.Println("There was an error opening the connection, ", err)
		return
	}

	// Start the Error handler in a goroutine
	go ErrorHandler(manager.ErrorChannel)

	fmt.Println("The Bot is now running.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	fmt.Println("Shutting Down...")
	err = dg.Close()
	if err != nil {
		fmt.Println(err)
	}

}

func getPrefixes(guildID string) []string {
	return []string{"!"}
}

func ErrorHandler(ErrorChan chan disgoman.CommandError) {
	for ce := range ErrorChan {
		msg := ce.Message
		if msg == "" {
			msg = ce.Error.Error()
		}
		_, _ = ce.Context.Send(msg)
		fmt.Println(ce.Error)
	}
}
