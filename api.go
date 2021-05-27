package main

import (
	"fmt"

	//"context"
	"github.com/go-resty/resty/v2"
	log "github.com/sirupsen/logrus"
	yarb "github.com/wmw9/yarb-struct"
	//	"io"
)

func New() *SuperAgent {
	s := &SuperAgent{
		Users:  make([]yarb.User, 0),
		Client: &resty.Client{},
		Errors: nil,
	}
	s.Client = resty.New()
	s.Client.SetBasicAuth(yarbBasicAuthUser, yarbBasicAuthPass)
	return s
}

func (s *SuperAgent) apiGetAllUsers() (*SuperAgent, error) {
	url := fmt.Sprintf("http://%v/yarb/users", YarbDBApiURL)
	log.Debug("Getting all enabled users from YARB-DB: ", url)

	if _, err := s.Client.R().SetResult(&s.Users).Get(url); err != nil {
		return s, err
	}

	return s, nil
}

func (s *SuperAgent) apiGetIGStoriesTsByID(id string) string {
	url := fmt.Sprintf("http://%v/yarb/user/id/%v/date/instagram_stories", YarbDBApiURL, id)
	log.Debugf("%v\n", url)
	resp, err := s.Client.R().Get(url)
	if err != nil || resp.StatusCode() != 200 {
		log.Errorf("err: %v\nstatus: %v\n resp.String(): %v", err, resp.StatusCode(), resp.String())
	}

	return resp.String()
}

func (s *SuperAgent) apiGetUserByIstagramID(id string) yarb.User {
	log.Printf("req %v/user/ig_id/:id\n", YarbDBApiURL)
	url := fmt.Sprintf("http://%v/yarb/user/ig_id/%v", YarbDBApiURL, id)
	log.Printf("[get]: %v", url)
	resp, err := s.Client.R().SetResult(&s.UserStruct).Get(url)
	if err != nil || resp.StatusCode() != 200 {
		log.Errorf("err: %v\nstatus: %v\n resp.String(): %v", err, resp.StatusCode(), resp.String())
	}

	return s.UserStruct
}

func (s *SuperAgent) apiSendPayloadMakaba(p yarb.Payload) bool {
	url := fmt.Sprintf("http://%v/makaba/post", YarbMakabaApiURL)
	log.Debugf("%v\n", url)
	resp, err := s.Client.R().SetBody(p).Post(url)
	log.Infof("%v:\n%v", resp.String(), err)
	return resp.StatusCode() == 200
}

func (s *SuperAgent) apiSendPayloadTelegram(p yarb.Payload) bool {
	url := fmt.Sprintf("http://%v/yarb/telegram/post", YarbTelegramApiURL)
	log.Debugf("%v\n", url)
	resp, err := s.Client.R().SetBody(p).Post(url)
	log.Infof("%v:\n%v", resp.String(), err)
	return resp.StatusCode() == 200
}

func (s *SuperAgent) apiUpdateIGStories(p yarb.Payload) bool {
	url := fmt.Sprintf("http://%v/yarb/user/name/%v/date/instagram_stories/%v", YarbDBApiURL, p.Person, p.Timestamp)
	log.Debugf("%v\n", url)
	resp, err := s.Client.R().Get(url)
	log.Infof("%v:\n%v", resp.String(), err)
	return resp.StatusCode() == 200
}
