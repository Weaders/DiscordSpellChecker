package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/bwmarrin/discordgo"
	"github.com/thoas/go-funk"
)

var (
	Token string
)

func init() {

	flag.StringVar(&Token, "t", "", "Bot Token")
	flag.Parse()
}

type BadWord struct {
	Word string   `json:word`
	S    []string `json:s`
}

func main() {

	Token = "NzM4ODE2NjcyNTUwNDg2MDI3.XyRapQ.6idKBSTTNxDrKRthd7VMuJ6eFRE"

	dg, err := discordgo.New("Bot " + Token)

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

func checkOnErrors(str string) string {

	resp, err := http.Get("https://speller.yandex.net/services/spellservice.json/checkText?text=" + url.PathEscape(str))

	if err == nil {
		body, _ := ioutil.ReadAll(resp.Body)
		resp.Body.Close()
		data := []BadWord{}
		json.Unmarshal(body, &data)

		var result []string

		result = (funk.Map(data, func(word BadWord) string {

			var result strings.Builder

			result.WriteString("Говно слово:" + word.Word + "\r\n")
			result.WriteString("Правильно: " + strings.Join(word.S, ","))

			return result.String()

		})).([]string)

		return strings.Join(result, "\r\n-----\r\n")

	} else {
		return "Ошибка, Яндекс "
	}

}

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {

	// Ignore all messages created by the bot itself
	// This isn't required in this specific example but it's a good practice.
	if m.Author.ID == s.State.User.ID {
		return
	}

	fmt.Println(m.Content)

	result := checkOnErrors(m.Content)

	fmt.Println(result)

	s.ChannelMessageSend(m.ChannelID, result)

	// If the message is "ping" reply with "Pong!"
	if m.Content == "ping" {
		s.ChannelMessageSend(m.ChannelID, "Pong!")
	}

	// If the message is "pong" reply with "Ping!"
	if m.Content == "pong" {
		s.ChannelMessageSend(m.ChannelID, "Ping!")
	}
}
