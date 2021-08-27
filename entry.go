package main

import "time"

type Entry struct {
	Id       int       `xml:"id"`
	Created  time.Time `xml:"created"`
	Posted   time.Time `xml:"posted"`
	Title    string    `xml:"title"`
	Contents string    `xml:"contents"`
}
