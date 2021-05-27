package main

import (
	"github.com/go-resty/resty/v2"
	log "github.com/sirupsen/logrus"
	yarb "github.com/wmw9/yarb-struct"
	"time"
)

type SuperAgent struct {
	Client     *resty.Client
	Users      []yarb.User
	UserStruct yarb.User
	Errors     []error
}

const sleepTime = 10 * time.Minute

func main() {
	log.Infof("Yet Another Reposter Bot - Instagram Producer started.\nYarb-DB IP: %v\nYarb-DB Domain: %v\n", YarbDBIp, YarbDBDomain)
	s := New()
	s.checkNewPosts()

}

func (s *SuperAgent) checkNewPosts() {
	for {
		_, err := s.apiGetAllUsers()

		if err != nil {
			log.Errorf("Failed to get list of users from YARB-DB: %v", err)
		}

		if err == nil {
			s.checkIGStories()
		}

		log.Debugf("Sleep for %v\n", sleepTime)
		time.Sleep(sleepTime)
	}
}
