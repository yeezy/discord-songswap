package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/bwmarrin/discordgo"
)

type SongLinks struct {
	Entityuniqueid  string `json:"entityUniqueId"`
	Usercountry     string `json:"userCountry"`
	Pageurl         string `json:"pageUrl"`
	Linksbyplatform struct {
		Amazonmusic struct {
			URL            string `json:"url"`
			Entityuniqueid string `json:"entityUniqueId"`
		} `json:"amazonMusic"`
		Deezer struct {
			URL            string `json:"url"`
			Entityuniqueid string `json:"entityUniqueId"`
		} `json:"deezer"`
		Applemusic struct {
			URL                 string `json:"url"`
			Nativeappurimobile  string `json:"nativeAppUriMobile"`
			Nativeappuridesktop string `json:"nativeAppUriDesktop"`
			Entityuniqueid      string `json:"entityUniqueId"`
		} `json:"appleMusic"`
		Itunes struct {
			URL                 string `json:"url"`
			Nativeappurimobile  string `json:"nativeAppUriMobile"`
			Nativeappuridesktop string `json:"nativeAppUriDesktop"`
			Entityuniqueid      string `json:"entityUniqueId"`
		} `json:"itunes"`
		Pandora struct {
			URL            string `json:"url"`
			Entityuniqueid string `json:"entityUniqueId"`
		} `json:"pandora"`
		Spotify struct {
			URL                 string `json:"url"`
			Nativeappuridesktop string `json:"nativeAppUriDesktop"`
			Entityuniqueid      string `json:"entityUniqueId"`
		} `json:"spotify"`
		Tidal struct {
			URL            string `json:"url"`
			Entityuniqueid string `json:"entityUniqueId"`
		} `json:"tidal"`
		Youtube struct {
			URL            string `json:"url"`
			Entityuniqueid string `json:"entityUniqueId"`
		} `json:"youtube"`
		Youtubemusic struct {
			URL            string `json:"url"`
			Entityuniqueid string `json:"entityUniqueId"`
		} `json:"youtubeMusic"`
	} `json:"linksByPlatform"`
}

func odesliCall(link string) *discordgo.MessageEmbed {

	url := "https://api.song.link/v1-alpha.1/links?url=" + link + "=US"
	method := "GET"

	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)

	if err != nil {
		fmt.Println(err)

	}
	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)

	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)

	}

	var allLinks SongLinks
	if err := json.Unmarshal(body, &allLinks); err != nil {
		fmt.Println("error: ", err)
	}

	embed := &discordgo.MessageEmbed{
		Author: &discordgo.MessageEmbedAuthor{},
		Color:  0x00ff00, // Green
		Fields: []*discordgo.MessageEmbedField{
			&discordgo.MessageEmbedField{
				Name:   "Amazon",
				Value:  "[Here](" + allLinks.Linksbyplatform.Amazonmusic.URL + ")",
				Inline: true,
			},
			&discordgo.MessageEmbedField{
				Name:   "Deezer",
				Value:  "[Here](" + allLinks.Linksbyplatform.Deezer.URL + ")",
				Inline: true,
			},
			&discordgo.MessageEmbedField{
				Name:   "Applemusic",
				Value:  "[Here](" + allLinks.Linksbyplatform.Applemusic.URL + ")",
				Inline: true,
			},
			&discordgo.MessageEmbedField{
				Name:   "Itunes",
				Value:  "[Here](" + allLinks.Linksbyplatform.Itunes.URL + ")",
				Inline: true,
			},
			&discordgo.MessageEmbedField{
				Name:   "Pandora",
				Value:  "[Here](" + allLinks.Linksbyplatform.Pandora.URL + ")",
				Inline: true,
			},
			&discordgo.MessageEmbedField{
				Name:   "Spotify",
				Value:  "[Here](" + allLinks.Linksbyplatform.Spotify.URL + ")",
				Inline: true,
			},
			&discordgo.MessageEmbedField{
				Name:   "Tidal",
				Value:  "[Here](" + allLinks.Linksbyplatform.Tidal.URL + ")",
				Inline: true,
			},
			&discordgo.MessageEmbedField{
				Name:   "Youtube",
				Value:  "[Here](" + allLinks.Linksbyplatform.Youtube.URL + ")",
				Inline: true,
			},
			&discordgo.MessageEmbedField{
				Name:   "Youtube Music",
				Value:  "[Here](" + allLinks.Linksbyplatform.Youtubemusic.URL + ")",
				Inline: true,
			},
		},
		Timestamp: time.Now().Format(time.RFC3339), // Discord wants ISO8601; RFC3339 is an extension of ISO8601 and should be completely compatible.
		Title:     "Smart Links For All Platforms",
	}
	return embed
}

var (
	Token string
)

func init() {
	flag.StringVar(&Token, "t", "", "Bot Token")
	flag.Parse()
}

func main() {
	dg, err := discordgo.New("Bot " + Token)
	if err != nil {
		fmt.Println("error creating Discord session,", err)
		return
	}

	// Register the messageCreate func as a callback for MessageCreate events.
	dg.AddHandler(messageCreate)

	// In this example, we only care about receiving message events.
	dg.Identify.Intents = discordgo.IntentsGuildMessages

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

// This function will be called (due to AddHandler above) every time a new
// message is created on any channel that the authenticated bot has access to.
func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {

	// Ignore all messages created by the bot itself
	if m.Author.ID == s.State.User.ID {
		return
	}
	// If the message begins with "-swap" reply with embed
	if strings.Contains(m.Content, "-swap") {
		messageLink := strings.Replace(m.Content, "-swap ", "", -1)
		s.ChannelMessageSendEmbed(m.ChannelID, odesliCall(messageLink))
	}
}
