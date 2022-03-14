package main

import (
	"fmt"
	uuid "github.com/satori/go.uuid"
	"log"
	"net/http"
	"text/template"
)

type user struct {
	UserName string
	Password string
	First    string
	Last     string
}

var dbSessions = map[string]string{} //cookie's value -> user name
var dbUsers = map[string]user{}      // User name  ->  user information
var tpl *template.Template

func init() {
	tpl = template.Must(template.ParseGlob("templates/*"))
	u := user{
		UserName: "test@test",
		Password: "password",
		First:    "test",
		Last:     "test",
	}
	dbUsers["test@test"] = u
}
func main() {
	http.HandleFunc("/", index)
	http.HandleFunc("/login", login)
	http.HandleFunc("/bar", bar)
	http.HandleFunc("/signup", signUp)
	http.HandleFunc("/logout", logout)
	http.Handle("/favicon.ico", http.NotFoundHandler())
	err := http.ListenAndServe(":80", nil)
	if err != nil {
		log.Fatalln(err)
	}
}
func index(w http.ResponseWriter, r *http.Request) {
	var u = user{}
	u = getUser(w, r)
	err := tpl.ExecuteTemplate(w, "index.gohtml", u)
	if err != nil {
		log.Fatalln(err)
	}
}

func signUp(w http.ResponseWriter, r *http.Request) {
	var u = user{}
	if alreadyLoggedIn(w, r) {
		http.Error(w, "you are already login!", http.StatusForbidden)
		return
	}
	//process the form
	if r.Method == http.MethodPost {
		un := r.FormValue("username")
		p := r.FormValue("password")
		f := r.FormValue("firstname")
		l := r.FormValue("lastname")

		//is user name already taken?
		if _, ok := dbUsers[un]; ok {
			http.Error(w, "username is already taken", http.StatusForbidden)
			return
		}
		//create session
		sID := uuid.NewV4()
		c := &http.Cookie{
			Name:  "session",
			Value: sID.String(),
		}
		http.SetCookie(w, c)
		dbSessions[c.Value] = un

		//encrypt password
		/*bs, err := bcrypt.GenerateFromPassword([]byte(p), bcrypt.MinCost)
		if err != nil {
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}
		*/
		u = user{un, p, f, l}
		dbSessions[c.Value] = un
		dbUsers[un] = u
		//getUser(w, r)
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
	fmt.Printf("#%d: %v\n", count, u)
	count++
	err := tpl.ExecuteTemplate(w, "signup.gohtml", u)
	if err != nil {
		log.Fatalln(err)
	}

}

func login(w http.ResponseWriter, r *http.Request) {
	if alreadyLoggedIn(w, r) {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	//process the form
	if r.Method == http.MethodPost {
		un := r.FormValue("username")
		p := r.FormValue("password")

		//is there any username?
		u, ok := dbUsers[un]
		if !ok {
			http.Error(w, "username or password deos not match!!!!!", http.StatusForbidden)
			return
		}
		if p != u.Password {
			http.Error(w, "username or password deos not match!", http.StatusForbidden)
			return
		}
		//create session
		sID := uuid.NewV4()
		c := &http.Cookie{
			Name:  "session",
			Value: sID.String(),
		}
		dbSessions[c.Value] = un
		dbUsers[un] = u
		http.SetCookie(w, c)
		http.Redirect(w, r, "/bar", http.StatusSeeOther)
	}
	err := tpl.ExecuteTemplate(w, "login.gohtml", nil)
	if err != nil {
		log.Fatalln(err)
	}
}
func bar(w http.ResponseWriter, r *http.Request) {
	if !alreadyLoggedIn(w, r) {
		http.Error(w, "you have to login first to come here!", http.StatusForbidden)
		return
	}
	c, _ := r.Cookie("session")
	un := dbSessions[c.Value]
	u := dbUsers[un]
	fmt.Println("user has been logged")
	err := tpl.ExecuteTemplate(w, "bar.gohtml", u)
	if err != nil {
		http.Error(w, "internal server error!", http.StatusInternalServerError)
	}
}
