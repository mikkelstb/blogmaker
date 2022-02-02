package main

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"log"
	"regexp"
)

type Catalouge struct {
	path_posts string
	path_cards string
	posts      []Post
	cards      []Card
	//pages catalouge
}

func NewCatalouge(path_posts, path_cards string) Catalouge {
	c := Catalouge{path_posts: path_posts, path_cards: path_cards}

	c.posts = c.readPosts()
	//c.cards = c.readCards()
	//c.pages = c.readPages()

	return c
}

func (c *Catalouge) readPosts() []Post {

	fmt.Println("Reading posts from file")
	fmt.Println(c.path_posts)

	var entries []Post

	files, err := ioutil.ReadDir(fmt.Sprintf("%v", c.path_posts))
	if err != nil {
		log.Fatal(err)
	}

	filepattern, _ := regexp.Compile(`\d{4}_\d{2}_\d{2}_\d{4}.xml`)

	for _, file := range files {

		if !(filepattern.MatchString(file.Name())) {
			continue
		}

		log.Default().Printf("Reading %v \n", file.Name())

		file, err := ioutil.ReadFile(fmt.Sprintf("%v/%v", c.path_posts, file.Name()))
		if err != nil {
			log.Fatal(err)
		}
		e := Post{}
		err = xml.Unmarshal(file, &e)
		if err != nil {
			log.Fatal(err)
		}
		entries = append(entries, e)
	}
	var ea []Post
	for index := len(entries) - 1; index >= 0; index-- {
		ea = append(ea, entries[index])
	}

	return ea
}

// func (c *Catalouge) writePost(e Post) {

// 	year, month, day := e.Created.Date()
// 	hour, min, _ := e.Created.Clock()

// 	filename := fmt.Sprintf("%v/%v_%.2v_%.2v_%.2v%.2v.xml", c.path, year, int(month), day, hour, min)
// 	xml_data, err := xml.Marshal(e)
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	if err == nil {
// 		err = os.WriteFile(filename, xml_data, 0600)

// 		if err != nil {
// 			log.Fatal(err)
// 		}
// 	}
// }
