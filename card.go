package main

import "html/template"

type Card struct {
	Id       int    `xml:"id"`
	Title    string `xml:"title"`
	Contents string `xml:"contents"`
}

func (c Card) GetContentsAsHtml() template.HTML {
	return template.HTML(c.Contents)
}
