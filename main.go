package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"
)

var cfg = readConfig()
var catalouge Catalouge

func main() {
	fmt.Println("Hello World. This is Blogmaker")
	fmt.Printf("The time is now: %v \n", time.Now())

	defer fmt.Printf("Server ended. The time is now: %v \n", time.Now())

	catalouge = NewCatalouge(cfg.Catalouges.Posts, cfg.Catalouges.Cards)
	resources_fileserver := http.FileServer(http.Dir(cfg.Catalouges.Resources))

	http.HandleFunc("/", handler)
	http.HandleFunc("/blog/", blogHandler)
	http.HandleFunc("/resources/", http.StripPrefix("/resources", resources_fileserver).ServeHTTP)

	log.Fatal(http.ListenAndServe(":3000", nil))
}

func blogHandler(w http.ResponseWriter, r *http.Request) {

	type page struct {
		Title string
		Intro string
		Posts []Post
		Cards []Card
	}

	p := new(page)
	p.Title = cfg.General.Title
	p.Intro = cfg.General.Intro

	p.Posts = catalouge.posts

	// for post_id := range p.Posts {
	// 	fmt.Println(p.Posts[post_id].Title)
	// }

	t, _ := template.ParseFiles(cfg.Catalouges.Templates + "/page_front.html")
	t.Execute(w, p)
}

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hi there, I love %s!", r.URL.Path[1:])
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
