package main

import (
	"strings"
	"time"
)

type Post struct {
	Id       string    `xml:"id"`
	Created  time.Time `xml:"created"`
	Posted   time.Time `xml:"posted"`
	Title    string    `xml:"title"`
	Contents string    `xml:"contents"`
	Tags     string    `xml:"tags"`
	Language string    `xml:"language"`
}

func (p Post) GetPosted() string {
	return p.Posted.Format("Monday, 02.Jan, 2006 15:04")
}

func (p Post) GetTags() []string {

	return strings.Split(p.Tags, " ")
}
