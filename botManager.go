package main

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"io/ioutil"
	"net/http"
)

/**
Manages Bots
*/
type BotManager struct {
	Bots   [100]*discordgo.Session
	active int
	length int
}

func PushBot(b *BotManager, dg *discordgo.Session) {
	fmt.Println("Adding bot: ", b.length)
	b.Bots[b.length] = dg
	b.length++
	guilds, _ := dg.UserGuilds(0, "", "")
	if len(guilds) == 100 && b.active == b.length-1 {
		b.active++
	}
}

func NewBot(b *BotManager, Token string) (*discordgo.Session, error) {
	dg, err := discordgo.New(Token)
	if err != nil {
		fmt.Println("error creating Discord session,", err)
		return nil, err
	}
	err = dg.Open()
	if err != nil {
		fmt.Println("error opening connection,", err)
		return nil, err
	}

	dg.AddHandler(MessageCreate)
	PushBot(b, dg)
	return dg, err
}

func Join(invite string, token string) {
	fmt.Println("Recived invite " + invite)
	url := "https://discordapp.com/api/invite/" + invite
	req, err := http.NewRequest("POST", url, nil)
	req.Header.Set("Accept", "*/*")
	req.Header.Set("Accept-Encoding", "gzip, deflate")
	req.Header.Set("Authorization", token)
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/58.0.3029.110 Safari/537.36")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	fmt.Println("response Status:", resp.Status)
	fmt.Println("response Headers:", resp.Header)
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println("response Body:", string(body))
}
