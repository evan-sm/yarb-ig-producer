package main

import (
	"errors"
	"fmt"

	//"github.com/k0kubun/pp"
	"strconv"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/tidwall/gjson"
	"github.com/wmw9/ig"
	yarb "github.com/wmw9/yarb-struct"
)

func (s *SuperAgent) checkIGStories() {
	userlist := make([]int, 0)
	for _, v := range s.Users {
		log.Debugf("%v\n", v)
		userlist = append(userlist, v.Social.InstagramID)
	}

	log.Infof("%v users to check", len(userlist))
	stories := ig.Get(igSessionId).Stories(userlist)
	log.Trace(string(stories))

	s.getNewIGStories(stories)
}

func (s *SuperAgent) getNewIGStories(stories []byte) {
	result := gjson.GetBytes(stories, "@valid")
	result.ForEach(func(key, value gjson.Result) bool {
		json := value.String()
		id := gjson.Get(json, `id`).String()
		user := s.apiGetUserByIstagramID(id)

		latest_reel_media := s.apiGetIGStoriesTsByID(id) // Request latest_reel_media by instagram id

		// Drop old posts
		items, err := dropOldStories(json, latest_reel_media)
		if err != nil {
			log.Debugf("%v", err) // All posts are old
			return true           // Next user
		}

		payload := preparePayloadFromIGStories(user, items)
		log.Trace(payload)
		log.Info("Trying to send payload to consumers")

		// Send payload
		if user.Setting.InstagramStories {

			log.Debugf("%v instagram stories reposts are enabled", user.Name)
			if err := SendToPubSub("yarb-313112", "yarb-telegram", payload); err != nil {
				log.Errorf("SendToPubSub telegram failed: %v", err)
			}

			if user.Setting.Makaba {

				log.Debugf("%v makaba reposts are enabled", user.Name)
				//if ok := s.apiSendPayloadMakaba(payload); ok {
				//	_ = s.apiSendPayloadTelegram(payload)
				//	log.Info("Trying to update Instagram Stories timestamp")
				//	if ok := s.apiUpdateIGStories(payload); ok {
				//		log.Errorf("Couldn't update IG Stories timestamp\n")
				//	}
				//}
				if err := SendToPubSub("yarb-313112", "yarb-makaba", payload); err != nil {
					log.Errorf("SendToPubSub makaba failed: %v", err)
				}
				//	log.Info("Trying to update Instagram Stories timestamp")
				//	if ok := s.apiUpdateIGStories(payload); ok {
				//		log.Errorf("Couldn't update IG Stories timestamp\n")
				//	}
			}
		}

		log.Trace("Sleep for %v", 20*time.Second)
		time.Sleep(20 * time.Second)

		return true
	})
}

func dropOldStories(json string, ts string) (string, error) {
	path := fmt.Sprintf(`items.#(taken_at>%v)#`, ts)
	items := gjson.Get(json, path).String()
	if items == "[]" {
		return "", errors.New("all IG stories are old")
	}
	return items, nil
}

func preparePayloadFromIGStories(u yarb.User, items string) yarb.Payload {
	var files []string
	var count int
	p := yarb.Payload{}

	result := gjson.Get(items, "@valid")
	result.ForEach(func(key, value gjson.Result) bool {
		//log.Printf("\n\n%v", value.String()[0:20])

		switch gjson.Get(value.String(), "media_type").String() {
		case "2": // video
			if count == 4 {
				return false // exit
			}
			count++
			url := gjson.Get(value.String(), "video_versions.0.url").String()
			log.Tracef("Got video story: %v", url)
			files = append(files, url) // add .mp4
			p.Timestamp = getTimestampIGStory(value.String())
		case "1": // image
			if count == 4 {
				return false // exit
			}
			count++
			url := gjson.Get(value.String(), "image_versions2.candidates.0.url").String()
			log.Tracef("Got image story: %v", url)
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
		log.Printf("%v\n", err)
	}
	return time
}
