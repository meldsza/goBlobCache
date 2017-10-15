package main

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
)

var (
	bm *BotManager
)

func init() {
	bm = new(BotManager)

	i := 1
	/*
		for _, env := range os.Environ() {
			fmt.Println(env)
		}
	*/
	for os.Getenv("T"+strconv.Itoa(i)) != "" {
		fmt.Println("Got token " + strconv.Itoa(i))
		NewBot(bm, os.Getenv("T"+strconv.Itoa(i)))
		i++
	}
}

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/", getEmojis).Methods("GET")
	router.HandleFunc("/guilds", getGuilds).Methods("GET")
	port := os.Getenv("PORT")
	if port == "" {
		port = "7070"
	}
	log.Fatal(http.ListenAndServe(":"+port, router))

	// Wait here until CTRL-C or other term signal is received.
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc
	for i, dg := range bm.Bots {
		// Cleanly close down the Discord session.
		if i < bm.length {
			dg.Close()
		}
	}
}

// message is created on any channel that the autenticated bot has access to.
func MessageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	if bm.active == bm.length {
		return //cant add more
	}
	if strings.Contains(m.Content, "discordapp.com/invite/") {
		index := strings.Index(m.Content, "discordapp.com/invite/")
		invite := strings.Split(m.Content[index:len(m.Content)], " ")[0]
		Join(strings.Replace(invite, "discordapp.com/invite/", "", len(invite)), bm.Bots[bm.length-1].Token)
	}
	if strings.Contains(m.Content, "discord.gg/") {
		index := strings.Index(m.Content, "discord.gg/")
		invite := strings.Split(m.Content[index:len(m.Content)], " ")[0]
		Join(strings.Replace(invite, "discord.gg/", "", len(invite)), bm.Bots[bm.length-1].Token)
	}
}

func GuildAdd(s *discordgo.Session, m *discordgo.MessageCreate) {
	guilds, _ := s.UserGuilds(0, "", "")
	if len(guilds) == 100 && bm.active != bm.length {
		bm.active++
		fmt.Println("Active changed to " + strconv.Itoa(bm.active))
	}
}
