package main

import "time"

type Post struct {
	Id       int       `xml:"id"`
	Created  time.Time `xml:"created"`
	Posted   time.Time `xml:"posted"`
	Title    string    `xml:"title"`
	Contents string    `xml:"contents"`
	Tags     []string  `xml:"tags"`
}

func (p Post) GetPosted() string {
	return p.Posted.Format("Monday, 02.Jan, 2006 15:04")
}
