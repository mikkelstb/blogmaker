package main

type config struct {
	General    general    `json:"general"`
	Catalouges catalouges `json:"catalouges"`
}

type general struct {
	Title string `json:"title"`
	Intro string `json:"intro"`
}

type catalouges struct {
	Resources string `json:"resources"`
	Images    string `json:"images"`
	Templates string `json:"templates"`
	Posts     string `json:"posts"`
	Cards     string `json:"cards"`
}
