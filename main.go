package main

import (
	"errors"
	"fmt"
	//"github.com/k0kubun/pp"
	"github.com/tidwall/gjson"
	"github.com/wmw9/ig"
	"log"
	"strconv"
	"time"
)

func main() {
	log.Printf("Yet Another Reposter Bot - Instagram Producer started.\nYarb-DB IP: %v\nYarb-DB Domain: %v\n", YarbDBIp, YarbDBDomain)
	s := New()
	s.checkNewPosts()

}

func (s *SuperAgent) checkNewPosts() {
	for {
		log.Printf("Trying to get users list")
		s.apiGetAllUsers()
		s.checkIGStories()
		//		for k, v := range s.User {
		//			log.Printf("k: %v; v:%v\n", k, v)
		//		}
		time.Sleep(10 * time.Minute)
	}
}

func (s *SuperAgent) checkIGStories() {
	userlist := make([]int, 0)
	for k, v := range s.Users {
		log.Printf("k: %v; v:%v\n", k, v)
		userlist = append(userlist, v.Social.InstagramID)
	}
	stories := ig.Get(igSessionId).Stories(userlist)
	//	println(string(stories))
	s.getNewIGStories(stories)
}

func (s *SuperAgent) getNewIGStories(stories []byte) {
	result := gjson.GetBytes(stories, "@valid")
	result.ForEach(func(key, value gjson.Result) bool {
		json := value.String()
		//log.Println("\n\njson: ", json)

		id := gjson.Get(json, `id`).String()
		user := s.apiGetUserByIstagramID(id)

		latest_reel_media := s.apiGetIGStoriesTsByID(id) // Request latest_reel_media by instagram id
		//println("id:", id, "latest_reel_media:", latest_reel_media)

		// Drop old posts
		items, err := dropOldStories(json, latest_reel_media)
		if err != nil {
			log.Printf("%v", err) // All posts are old
			return true
		}
		if items == "" {
			return true
		}
		//println("\n\nitems:", items[0:100])

		payload := preparePayloadFromIGStories(user, items)
		//pp.Println(payload)
		println("Trying to send payload to consumers")

		// Send payload
		if user.Setting.InstagramStories {
			println("Stories reposts enabled")
			if user.Setting.Makaba {
				println("Makaba reposts enabled")
				ok := s.apiSendPayloadMakaba(payload)
				if ok {
					_ = s.apiSendPayloadTelegram(payload)
					println("Trying to update Instagram Stories timestamp")
					ok := s.apiUpdateIGStories(payload)
					if !ok {
						panic("Couldn't update IG Stories timestamp")
					}
				}
			}
		}

		time.Sleep(20 * time.Second)
		return true
	})
}

func dropOldStories(json string, ts string) (string, error) {
	path := fmt.Sprintf(`items.#(taken_at>%v)#`, ts)
	items := gjson.Get(json, path).String()
	if items == "[]" {
		return "", errors.New("All IG stories are all old.")
	}
	return items, nil
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

func preparePayloadFromIGStories(u User, items string) Payload {
	var files []string
	var count int
	p := Payload{}

	result := gjson.Get(items, "@valid")
	result.ForEach(func(key, value gjson.Result) bool {
		//log.Printf("\n\n%v", value.String()[0:20])

		switch gjson.Get(value.String(), "media_type").String() {
		case "2": // video
			if count == 4 {
				return false
			}
			count++
			//log.Printf("count: %v", count)
			url := gjson.Get(value.String(), "video_versions.0.url").String()
			//log.Printf("Got video story: %v", url)
			files = append(files, url) // add .mp4
			p.Timestamp = getTimestampIGStory(value.String())
		case "1": // image
			if count == 4 {
				return false
			}
			count++
			//log.Printf("count: %v", count)
			url := gjson.Get(value.String(), "image_versions2.candidates.0.url").String()
			//log.Printf("Got image story: %v", url)
			files = append(files, url) // add .jpg
			p.Timestamp = getTimestampIGStory(value.String())
		}
		return true // keep iterating
	})

	p.Person = u.Name
	p.From = "instagram"
	p.Type = "story"
	p.Source = fmt.Sprintf("https://instagram.com/%v/", u.Social.Instagram)
	p.TelegramChanID = u.Repost.TgChanID
	p.Board = u.Repost.Board
	p.Thread = u.Repost.Thread
	p.Files = files

	return p
}

func getTimestampIGStory(item string) int {
	time64 := gjson.Get(item, "taken_at").String()
	time, err := strconv.Atoi(time64)
	if err != nil {
		log.Printf("%v", err)
	}
	return time
}
