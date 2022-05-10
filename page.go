package main

type page struct {
	Title string
	Intro string
	Url   string
	Posts []Post
	Post  Post
	Cards []Card
	Tags  map[string]int
}

func newPage() *page {

	p := new(page)

	p.Title = cfg.General.Title
	p.Intro = cfg.General.Intro
	p.Url = cfg.General.Url

	return p
}
