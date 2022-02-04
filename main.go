package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"
	"time"
)

var cfg = readConfig()
var catalouge Catalouge

func main() {
	fmt.Println("Hello World. This is Blogmaker")
	fmt.Printf("The time is now: %v \n", time.Now())

	defer fmt.Printf("Server ended. The time is now: %v \n", time.Now())

	resources_fileserver := http.FileServer(http.Dir(cfg.Catalouges.Resources))

	http.HandleFunc("/", handler)
	http.HandleFunc("/blog/", blogHandler)
	http.HandleFunc("/blog/post/", postHandler)
	//http.handlefunc("blog/tag/", searchHandler)
	http.HandleFunc("/resources/", http.StripPrefix("/resources", resources_fileserver).ServeHTTP)

	log.Fatal(http.ListenAndServe(":3000", nil))
}

func handler(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, cfg.General.Url, http.StatusPermanentRedirect)
}

func blogHandler(w http.ResponseWriter, r *http.Request) {

	type page struct {
		Title string
		Intro string
		Posts []Post
		Cards []Card
	}

	catalouge = NewCatalouge(cfg.Catalouges.Posts, cfg.Catalouges.Cards)

	p := new(page)
	p.Title = cfg.General.Title
	p.Intro = cfg.General.Intro

	p.Posts = catalouge.posts
	p.Cards = catalouge.cards

	for c := range p.Cards {
		fmt.Println(p.Cards[c].Contents)
	}

	t, _ := template.ParseFiles(cfg.Catalouges.Templates + "/page_front.html")
	t.Execute(w, p)
}

func postHandler(w http.ResponseWriter, r *http.Request) {

	type page struct {
		Title string
		Intro string
		Url   string
		Post  Post
		Cards []Card
	}

	catalouge = NewCatalouge(cfg.Catalouges.Posts, cfg.Catalouges.Cards)

	p := new(page)
	p.Title = cfg.General.Title
	p.Intro = cfg.General.Intro
	p.Cards = catalouge.cards
	p.Url = cfg.General.Url

	post_request_id := path.Base(r.URL.Path)

	p.Post = catalouge.posts[catalouge.post_index[post_request_id]]

	t, _ := template.ParseFiles(cfg.Catalouges.Templates + "/post_single.html")
	t.Execute(w, p)
}

func readConfig() *config {

	json_file, err := os.Open("config.json")
	if err != nil {
		log.Fatal(err)
	}
	defer json_file.Close()

	json_data, err := ioutil.ReadAll(json_file)
	if err != nil {
		log.Fatal(err)
	}

	config_file := new(config)
	if err := json.Unmarshal(json_data, &config_file); err != nil {
		log.Fatal(err)
	}

	fmt.Println(config_file.General.Title)
	fmt.Println(config_file.General.Intro)

	return config_file
}
