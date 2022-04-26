package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"path"
	"time"

	"github.com/gorilla/sessions"
)

var cfg = readConfig()
var catalouge Catalouge
var key = []byte("test")
var store = sessions.NewCookieStore(key)

type page struct {
	Title string
	Intro string
	Url   string
	Posts []Post
	Post  Post
	Cards []Card
	Tags  map[string]int
}

func main() {

	fmt.Println("Hello World. This is Blogmaker")
	fmt.Printf("The time is now: %v \n", time.Now())

	defer fmt.Printf("Server ended. The time is now: %v \n", time.Now())

	resources_fileserver := http.FileServer(http.Dir(cfg.Catalouges.Resources))

	http.HandleFunc("/", blogHandler)
	http.HandleFunc("/post/", postHandler)
	http.HandleFunc("/tag/", searchHandler)
	http.HandleFunc("/resources/", http.StripPrefix("/resources", resources_fileserver).ServeHTTP)
	http.HandleFunc("/login/", loginHandler)
	http.HandleFunc("/logout/", logouthandler)

	log.Fatal(http.ListenAndServe(":3000", nil))
}

func blogHandler(w http.ResponseWriter, r *http.Request) {

	fmt.Println("Incoming request")

	catalouge = NewCatalouge(cfg.Catalouges.Posts, cfg.Catalouges.Cards)

	p := new(page)
	p.Title = cfg.General.Title
	p.Intro = cfg.General.Intro

	p.Posts = catalouge.posts
	p.Cards = catalouge.cards

	p.Tags = make(map[string]int)

	for k, v := range catalouge.tag_index {
		p.Tags[k] = len(v)
	}

	for c := range p.Cards {
		fmt.Println(p.Cards[c].Contents)
	}

	t, _ := template.ParseFiles(
		cfg.Catalouges.Templates+"/page_front.html",
		cfg.Catalouges.Templates+"/head.html",
		cfg.Catalouges.Templates+"/posts.html",
	)
	t.Execute(w, p)
}

func loginHandler(w http.ResponseWriter, r *http.Request) {

	session, _ := store.Get(r, "session_data")

	session.Values["authenticated"] = true
	session.Save(r, w)

}

func logouthandler(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "session_data")

	// Revoke users authentication
	session.Values["authenticated"] = false
	session.Save(r, w)
}

func postHandler(w http.ResponseWriter, r *http.Request) {

	catalouge = NewCatalouge(cfg.Catalouges.Posts, cfg.Catalouges.Cards)

	p := new(page)
	p.Title = cfg.General.Title
	p.Intro = cfg.General.Intro
	p.Cards = catalouge.cards
	p.Url = cfg.General.Url

	post_request_id := path.Base(r.URL.Path)

	p.Post = catalouge.posts[catalouge.post_index[post_request_id]]

	t, _ := template.ParseFiles(
		cfg.Catalouges.Templates+"/post_single.html",
		cfg.Catalouges.Templates+"/head.html",
	)
	t.Execute(w, p)
}

func searchHandler(w http.ResponseWriter, r *http.Request) {

	catalouge = NewCatalouge(cfg.Catalouges.Posts, cfg.Catalouges.Cards)

	p := new(page)
	p.Title = cfg.General.Title
	p.Intro = cfg.General.Intro
	p.Cards = catalouge.cards
	p.Url = cfg.General.Url

	tag_request := path.Base(r.URL.Path)
	post_ids := catalouge.tag_index[tag_request]

	for _, pid := range post_ids {
		p.Posts = append(p.Posts, catalouge.posts[catalouge.post_index[pid]])
	}

	t, _ := template.ParseFiles(
		cfg.Catalouges.Templates+"/post_search.html",
		cfg.Catalouges.Templates+"/head.html",
		cfg.Catalouges.Templates+"/posts.html",
	)
	t.Execute(w, p)
}
