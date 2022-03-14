package main

import (
	"fmt"
	uuid "github.com/satori/go.uuid"
	"net/http"
)

var count = 0

func getUser(w http.ResponseWriter, r *http.Request) user {
	var u = user{}
	c, err := r.Cookie("session")
	if err != nil {
		id := uuid.NewV4()
		c = &http.Cookie{
			Name:  "session",
			Value: id.String(),
		}
	}
	http.SetCookie(w, c)
	// if the user exists already, get user
	//var u user
	if un, ok := dbSessions[c.Value]; ok {
		u = dbUsers[un]
	}
	fmt.Printf("#%d: %v\n", count, u)
	count++
	return u
}

func alreadyLoggedIn(w http.ResponseWriter, r *http.Request) bool {
	c, err := r.Cookie("session")
	if err != nil {
		return false
	}
	un := dbSessions[c.Value]
	_, ok := dbUsers[un]
	return ok
}

func logout(w http.ResponseWriter, r *http.Request) {
	c, err := r.Cookie("session")
	if err != nil {
		http.Error(w, "you have to login first!", http.StatusForbidden)
		return
	}
	c = &http.Cookie{
		Name:   "session",
		MaxAge: -1,
	}
	http.SetCookie(w, c)
	http.Redirect(w, r, "/login", http.StatusSeeOther)
	fmt.Println("user has been logged out")
}
