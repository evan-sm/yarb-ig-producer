package main

import (
	"fmt"
	"github.com/wmw9/ig"
	"reflect"
	"testing"
)

func TestApiGetAllUsers(t *testing.T) {
	s := New()
	users := fmt.Sprintf("http://%v/yarb/users", YarbDBApiURL)
	_, _ = s.Client.R().SetResult(&s.Users).Get(users)

	var Users []User
	if reflect.DeepEqual(s.Users, Users) {
		t.Fatalf("List of users is empty!")
	}
	if s.Users[0].Name == "" {
		t.Fatalf("Username not found in userlist")
	}

}

func TestCheckIGStories(t *testing.T) {
	user := 2944757465
	stories := ig.Get(igSessionId).Stories(user)
	fmt.Println("stories:", string(stories))

}
