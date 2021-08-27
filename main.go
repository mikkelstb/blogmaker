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
var blogfolder = blogFolder{path: cfg.BlogCatalouge}

func main() {
	fmt.Println("Hello World. This is Blogmaker")
	fmt.Printf("The time is now: %v \n", time.Now())

	http.HandleFunc("/blog/", blogHandler)
	http.HandleFunc("/test/", handler)

	log.Fatal(http.ListenAndServe(":3000", nil))

	// for _, entry := range entries {
	// 	fmt.Printf("Found: %v @ %v \n", entry.Title, entry.Posted.String())
	// }

	// fmt.Println("-----")

	// entries = blogfolder.getEntries(12, 0)

	// for _, entry := range entries {
	// 	fmt.Printf("Found: %v @ %v \n", entry.Title, entry.Posted.String())
	// }

}

func blogHandler(w http.ResponseWriter, r *http.Request) {

	type page struct {
		Title   string
		Intro   string
		Entries []Entry
	}

	p := new(page)
	p.Title = cfg.Title
	p.Intro = cfg.Intro

	p.Entries = blogfolder.getEntries(10, 0)
	fmt.Println(p)

	t, _ := template.ParseFiles("./templates/blog.html")
	t.Execute(w, p)
}

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hi there, I love %s!", r.URL.Path[1:])
}

func readConfig() *config {

	json_file, err := os.Open("setup.json")
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

	fmt.Println(config_file.Title)
	fmt.Println(config_file.Intro)
	fmt.Println(config_file.BlogCatalouge)

	return config_file
}

type config struct {
	Title         string `json:"title"`
	Intro         string `json:"intro"`
	BlogCatalouge string `json:"blogCatalouge"`
}
