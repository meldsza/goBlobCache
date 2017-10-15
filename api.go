package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

func getEmojis(w http.ResponseWriter, r *http.Request) {
	em := make(map[string]string)
	for i, dg := range bm.Bots {
		fmt.Println("dg " + strconv.Itoa(i))
		if i >= bm.length || dg == nil {
			//fmt.Println("dg is nil")
			break
		}
		if dg.State == nil || dg.State.Guilds == nil {
			fmt.Println("dg State or Guilds is nil")
			break
		}
		for _, guild := range dg.State.Guilds {
			if guild == nil || guild.Emojis == nil {
				fmt.Println("Guild or Emoji is nil")
				break
			}
			for _, emoji := range guild.Emojis {
				_, test := em[emoji.Name]
				if test {
					a := 1
					for ; test; a++ {
						_, test = em[emoji.Name+strconv.Itoa(a)]
					}
					a--
					em[emoji.Name+strconv.Itoa(a)] = "https://cdn.discordapp.com/emojis/" + emoji.ID + ".png"
				} else {
					em[emoji.Name] = "https://cdn.discordapp.com/emojis/" + emoji.ID + ".png"
				}
			}
		}

	}
	json.NewEncoder(w).Encode(em)
}
func getGuilds(w http.ResponseWriter, r *http.Request) {
	gj := make(map[string]string)
	for i, dg := range bm.Bots {
		guilds, _ := dg.UserGuilds(0, "", "")
		for _, guild := range guilds {
			gj[guild.ID] = guild.Name
		}
		if i < bm.length {
			break
		}
	}
	json.NewEncoder(w).Encode(gj)
}
