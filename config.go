package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
)

type config struct {
	General    general    `json:"general"`
	Catalouges catalouges `json:"catalouges"`
}

type general struct {
	Title string `json:"title"`
	Intro string `json:"intro"`
	Url   string `json:"url"`
}

type catalouges struct {
	Resources string `json:"resources"`
	Images    string `json:"images"`
	Templates string `json:"templates"`
	Posts     string `json:"posts"`
	Cards     string `json:"cards"`
}

func readConfig() *config {

	json_file, err := os.Open("config.json")
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

	fmt.Println(config_file.General.Title)
	fmt.Println(config_file.General.Intro)

	return config_file
}
