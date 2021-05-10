package main

import (
	"fmt"
	"github.com/go-resty/resty/v2"
	"log"
)

func New() *SuperAgent {
	s := &SuperAgent{
		Users:  make([]User, 0),
		Client: &resty.Client{},
	}
	s.Client = resty.New()
	s.Client.SetBasicAuth(yarbBasicAuthUser, yarbBasicAuthPass)
	return s
}

func (s *SuperAgent) apiGetAllUsers() *SuperAgent {
	users := fmt.Sprintf("http://%v/yarb/users", YarbDBApiURL)
	log.Printf("[get]: %v", users)
	resp, err := s.Client.R().SetResult(&s.Users).Get(users)
	failOnError(err, "get failed")

	log.Printf("%v", resp.StatusCode())
	return s
}

func (s *SuperAgent) apiGetIGStoriesTsByID(id string) string {
	log.Printf("req %v/user/id/:id/date/instagram_stories\n", YarbDBApiURL)
	url := fmt.Sprintf("http://%v/yarb/user/id/%v/date/instagram_stories", YarbDBApiURL, id)
	log.Printf("[get]: %v", url)
	resp, err := s.Client.R().Get(url)
	failOnError(err, "get failed")

	//log.Println(resp.Status, resp.String())
	return resp.String()
}

func (s *SuperAgent) apiGetUserByIstagramID(id string) User {
	log.Printf("req %v/user/ig_id/:id\n", YarbDBApiURL)
	url := fmt.Sprintf("http://%v/yarb/user/ig_id/%v", YarbDBApiURL, id)
	log.Printf("[get]: %v", url)
	resp, err := s.Client.R().SetResult(&s.UserStruct).Get(url)
	println("\n\nresp.String:", resp.String())
	if err != nil {
		println(resp.Status, resp.Body())
		panic(err)
	}
	failOnError(err, "get failed")

	return s.UserStruct
}

func (s *SuperAgent) apiSendPayloadMakaba(p Payload) bool {
	url := fmt.Sprintf("http://%v/makaba/post", YarbMakabaApiURL)
	println("[post]:", url)
	resp, err := s.Client.R().SetBody(p).Post(url)
	println(resp.String(), err)
	if resp.StatusCode() == 200 {
		return true
	}
	return false
}

func (s *SuperAgent) apiSendPayloadTelegram(p Payload) bool {
	url := fmt.Sprintf("http://%v/yarb/telegram/post", YarbTelegramApiURL)
	println("[post]:", url)
	resp, err := s.Client.R().SetBody(p).Post(url)
	println(resp.String(), err)
	if resp.StatusCode() == 200 {
		return true
	}
	return false
}

func (s *SuperAgent) apiUpdateIGStories(p Payload) bool {
	url := fmt.Sprintf("http://%v/yarb/user/name/%v/date/instagram_stories/%v", YarbDBApiURL, p.Person, p.Timestamp)
	println("[get]:", url)
	resp, err := s.Client.R().Get(url)
	println(resp.String(), err)
	if resp.StatusCode() == 200 {
		return true
	}
	return false
}
