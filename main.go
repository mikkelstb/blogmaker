package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"path"
	"time"

	"github.com/gorilla/sessions"
)

var cfg *config
var catalouge Catalouge
var key []byte
var store *sessions.CookieStore
var config_file_name string = "config.json"

func main() {

	fmt.Println("Hello World. This is Blogmaker")
	fmt.Printf("The time is now: %v \n", time.Now())
	fmt.Printf("Reading config file: %s \n", config_file_name)

	// Read config file
	//
	cfg = readConfig(config_file_name)
	catalouge = NewCatalouge(cfg.Catalouges.Posts, cfg.Catalouges.Cards)
	var err error

	// Read key file used for encrypting cookies
	//
	data, err := os.ReadFile(cfg.General.Login_key_file)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	fmt.Println(string(data))
	store = sessions.NewCookieStore(data)

	resources_fileserver := http.FileServer(http.Dir(cfg.Catalouges.Resources))

	http.HandleFunc("/", blogHandler)
	http.HandleFunc("/post/", postHandler)
	http.HandleFunc("/tag/", searchHandler)
	http.HandleFunc("/resources/", http.StripPrefix("/resources", resources_fileserver).ServeHTTP)
	http.HandleFunc("/login/", loginHandler)
	http.HandleFunc("/logout/", logouthandler)
	http.HandleFunc("/edit/", postEditHandler)

	log.Fatal(http.ListenAndServe(":3000", nil))
}

func blogHandler(w http.ResponseWriter, r *http.Request) {

	session, _ := store.Get(r, "session_data")

	fmt.Println(session.Values["authenticated"])
	fmt.Println(session.Values["logintime"])

	p := newPage()

	p.Posts = catalouge.posts
	p.Cards = catalouge.cards

	p.Tags = make(map[string]int)

	for k, v := range catalouge.tag_index {
		p.Tags[k] = len(v)
	}

	t, _ := template.ParseFiles(
		cfg.Catalouges.Templates+"/page_front.html",
		cfg.Catalouges.Templates+"/head.html",
		cfg.Catalouges.Templates+"/posts.html",
	)
	t.Execute(w, p)
}

func loginHandler(w http.ResponseWriter, r *http.Request) {

	session, err := store.Get(r, "session_data")
	fmt.Println("Login: " + session.Name())

	if err != nil {

		http.Error(w, err.Error(), http.StatusForbidden)
	}

	session.Values["authenticated"] = true
	session.Values["logintime"] = time.Now().Format(time.Kitchen)

	err = session.Save(r, w)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

}

func logouthandler(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "session_data")
	fmt.Println("Logout: " + session.Name())

	// Revoke users authentication
	session.Values["authenticated"] = false
	session.Save(r, w)
}

func postHandler(w http.ResponseWriter, r *http.Request) {

	p := newPage()
	p.Cards = catalouge.cards

	post_request_id := path.Base(r.URL.Path)

	p.Post = catalouge.posts[catalouge.post_index[post_request_id]]

	t, _ := template.ParseFiles(
		cfg.Catalouges.Templates+"/post_single.html",
		cfg.Catalouges.Templates+"/head.html",
	)
	t.Execute(w, p)
}

func postEditHandler(w http.ResponseWriter, r *http.Request) {

	var err error
	post_request_id := path.Base(r.URL.Path)

	p := newPage()
	p.Cards = catalouge.cards

	if r.Method == "POST" {
		err = r.ParseForm()

		if err != nil {
			w.Write([]byte(err.Error()))
		}
		catalouge.savePost(r.FormValue("id"), r.FormValue("title"), r.FormValue("contents"))
		catalouge.readPosts()
	}

	p.Post = catalouge.posts[catalouge.post_index[post_request_id]]

	t, _ := template.ParseFiles(
		cfg.Catalouges.Templates+"/post_edit.html",
		cfg.Catalouges.Templates+"/head.html",
	)
	t.Execute(w, p)
}

func searchHandler(w http.ResponseWriter, r *http.Request) {

	p := newPage()
	p.Cards = catalouge.cards

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
