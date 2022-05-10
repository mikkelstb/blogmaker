package main

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"time"
)

type Catalouge struct {
	path_posts string
	path_cards string
	posts      []Post
	cards      []Card

	// post_id -> position on slice
	post_index map[string]int

	// tag_name -> []post_id
	tag_index map[string][]string
}

func NewCatalouge(path_posts, path_cards string) Catalouge {
	c := Catalouge{path_posts: path_posts, path_cards: path_cards}

	c.post_index = make(map[string]int)
	c.tag_index = map[string][]string{}

	c.readPosts()
	c.makePostIndex()
	c.makeTagIndex()

	//fmt.Println(c.tag_index)

	c.cards = c.readCards()

	//c.pages = c.readPages()

	return c
}

func (c *Catalouge) savePost(id, title, contents string) error {

	post := c.posts[c.post_index[id]]

	post.Title = title
	post.Contents = contents
	post.Edited = time.Now()

	xmldata, err := xml.Marshal(post)

	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	err = os.WriteFile(fmt.Sprintf("%v/%v", c.path_posts, post.filename), xmldata, 0666)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	return nil
}

func (c *Catalouge) readCards() []Card {
	var cards []Card

	files, err := ioutil.ReadDir(fmt.Sprintf("%v", c.path_cards))
	if err != nil {
		log.Fatal(err)
	}

	filepattern, _ := regexp.Compile(`.*?\.xml`)

	for _, file := range files {

		if !(filepattern.MatchString(file.Name())) {
			continue
		}

		//log.Default().Printf("Reading %v \n", file.Name())

		file, err := ioutil.ReadFile(fmt.Sprintf("%v/%v", c.path_cards, file.Name()))
		if err != nil {
			log.Fatal(err)
		}
		e := Card{}
		err = xml.Unmarshal(file, &e)
		if err != nil {
			log.Fatal(err)
		}
		cards = append(cards, e)
	}
	return cards
}

func (c *Catalouge) readPosts() {

	//fmt.Println("Reading posts from file")
	//fmt.Println(c.path_posts)

	var entries []Post

	files, err := ioutil.ReadDir(fmt.Sprintf("%v", c.path_posts))
	if err != nil {
		log.Fatal(err)
	}

	filepattern, _ := regexp.Compile(`\d{4}_\d{2}_\d{2}_\d{4}\.xml`)

	for _, file := range files {

		if !(filepattern.MatchString(file.Name())) {
			continue
		}

		xml_file, err := ioutil.ReadFile(fmt.Sprintf("%v/%v", c.path_posts, file.Name()))
		if err != nil {
			log.Fatal(err)
		}
		e := Post{}
		e.filename = file.Name()
		err = xml.Unmarshal(xml_file, &e)
		if err != nil {
			log.Fatal(err)
		}

		//Append next to beginning of slice
		entries = append([]Post{e}, entries...)
	}
	c.posts = entries
}

// Populate the map post_index. This makes an easy hashtable that converts post id's into slice-index for c.posts.

func (c *Catalouge) makePostIndex() {

	for p := range c.posts {
		c.post_index[c.posts[p].Id] = p
	}
}

func (c *Catalouge) makeTagIndex() {
	for p := range c.posts {
		for _, t := range c.posts[p].GetTags() {
			c.tag_index[t] = append(c.tag_index[t], c.posts[p].Id)
		}
	}
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
