package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"

	"github.com/weaders/DiscordSpellChecker/speller"

	"github.com/bwmarrin/discordgo"
)

var (
	token string
)

func init() {

	flag.StringVar(&token, "t", "", "Bot Token")
	flag.Parse()

}

func main() {

	dg, err := discordgo.New("Bot " + token)

	if err != nil {
		fmt.Println("error creating Discord session,", err)
		return
	}

	dg.AddHandler(messageCreate)

	// In this example, we only care about receiving message events.
	dg.Identify.Intents = discordgo.MakeIntent(discordgo.IntentsGuildMessages)

	// Open a websocket connection to Discord and begin listening.
	err = dg.Open()
	if err != nil {
		fmt.Println("error opening connection,", err)
		return
	}

	// Wait here until CTRL-C or other term signal is received.
	fmt.Println("Bot is now running.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	// Cleanly close down the Discord session.
	dg.Close()

}

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {

	if m.Author.ID == s.State.User.ID {
		return
	}

	if m.Content == "!help" {
		sendMsg(s, m.ChannelID, "!git - ссылка на гит исходники")
	} else if m.Content == "!count" {

		members, _ := s.GuildMembers(m.GuildID, "0", 20)

		var ids []string

		for _, v := range members {
			ids = append(ids, v.User.ID)
		}

		result := speller.CounterForUsers(ids)

		var resultStr strings.Builder

		for k, v := range result {

			user := findUser(members, k)

			resultStr.WriteString(user.Username + ": " + strconv.Itoa(v))

			sendMsg(s, m.ChannelID, resultStr.String())

		}

	} else if m.Content == "!git" {
		sendMsg(s, m.ChannelID, "https://github.com/Weaders/DiscordSpellChecker")
	} else {

		result := speller.CheckString(m.Content)

		speller.AddCountForUser(m.Author.ID, len(result))

		sendMsg(s, m.ChannelID, strings.Join(result, "\r\n-------\r\n"))
	}

}

func sendMsg(s *discordgo.Session, channelID string, msg string) {

	_, err := s.ChannelMessageSend(channelID, msg)

	if err != nil {
		println(err)
	}

}

func findUser(users []*discordgo.Member, id string) *discordgo.User {
	for _, v := range users {
		if v.User.ID == id {
			return v.User
		}
	}

	return nil
}
