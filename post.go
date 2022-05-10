package main

import (
	"strings"
	"time"
)

type Post struct {
	Id       string    `xml:"id"`
	Created  time.Time `xml:"created"`
	Posted   time.Time `xml:"posted"`
	Edited   time.Time `xml:"edited"`
	Title    string    `xml:"title"`
	Contents string    `xml:"contents"`
	Tags     string    `xml:"tags"`
	Language string    `xml:"language"`
	filename string
}

func (p Post) GetPosted() string {
	return p.Posted.Format("Monday, 02.Jan, 2006")
}

func (p Post) GetEdited() string {
	return p.Edited.Format("Monday, 02.Jan, 2006")
}

func (p Post) GetTags() []string {
	return strings.Split(p.Tags, " ")
}
