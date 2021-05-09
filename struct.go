package main

import "github.com/go-resty/resty/v2"

type SuperAgent struct {
	Client     *resty.Client
	Users      []User
	UserStruct User
}

type User struct {
	Name   string `bson:"name" json:"name" binding:"required"`
	Social struct {
		InstagramID int    `bson:"instagram_id" json:"instagram_id"`
		Instagram   string `bson:"instagram" json:"instagram"`
	} `bson:"social" json:"social"`
	Date struct {
		InstagramPost    int `bson:"instagram_post" json:"instagram_post"`
		InstagramStories int `bson:"instagram_stories" json:"instagram_stories"`
	} `bson:"date" json:"date"`
	Setting struct {
		Disabled         bool `bson:"disabled" json:"disabled"`
		InstagramPost    bool `bson:"instagram_post" json:"instagram_post"`
		InstagramStories bool `bson:"instagram_stories" json:"instagram_stories"`
		Makaba           bool `bson:"makaba" json:"makaba"`
	} `bson:"setting" json:"setting"`
	Repost struct {
		TgChanID int64  `bson:"tg_chan_id" json:"tg_chan_id"`
		Board    string `bson:"board" json:"board"`
		Thread   string `bson:"thread" json:"thread"`
	} `bson:"repost" json:"repost"`
}

type Payload struct {
	// Body
	Person    string
	Timestamp int
	Caption   string
	From      string
	Type      string
	Source    string
	Files     []string

	// Destination
	TelegramChanID int64
	Board          string
	Thread         string
}
