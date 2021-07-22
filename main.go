package main

import (
	"elmajson/myjson"
	"fmt"
	"time"
)

type NewsInfoNewsapi struct {
	Status       string `json:"status"`
	TotalResults int    `json:"totalResults"`
	Articles     []struct {
		Source struct {
			ID   interface{} `json:"id"`
			Name string      `json:"name"`
		} `json:"source"`
		Author      string `json:"author"`
		Title       string `json:"title"`
		Description string `json:"description"`
		URL         string `json:"url"`
		URLToImage  string `json:"urlToImage"`
		PublishedAt string `json:"publishedAt"`
		Content     string `json:"content"`
	} `json:"articles"`
}

type NewsInfoMediastack struct {
	Pagination struct {
		Limit  int `json:"limit"`
		Offset int `json:"offset"`
		Count  int `json:"count"`
		Total  int `json:"total"`
	} `json:"pagination"`
	Data []struct {
		Author      string    `json:"author"`
		Title       string    `json:"title"`
		Description string    `json:"description"`
		URL         string    `json:"url"`
		Source      string    `json:"source"`
		Image       string    `json:"image"`
		Category    string    `json:"category"`
		Language    string    `json:"language"`
		Country     string    `json:"country"`
		PublishedAt time.Time `json:"published_at"`
	} `json:"data"`
}

func getOneQuery(url, keyName, key string, params myjson.Params, news interface{}) {
	fmt.Println("Получаем новости из", url)
	client := myjson.NewClient(url, keyName, key)

	err := myjson.GetJSON(*client, params, &news)
	if err != nil {
		fmt.Println(err)
		return
	}
}

func showTitleDesc(news interface{}) {
	switch newsType := news.(type) {
	case NewsInfoNewsapi:
		for _, item := range newsType.Articles {
			fmt.Println(item.Title)
			fmt.Println(item.Description)
			fmt.Println()
		}
	case NewsInfoMediastack:
		for _, item := range newsType.Data {
			fmt.Println(item.Title)
			fmt.Println(item.Description)
			fmt.Println()
		}
	default:
		fmt.Println("error showTitleDesc")
	}
}

func main() {
	fmt.Println("start work")

	newsInfoNewsapi := NewsInfoNewsapi{}
	getOneQuery(
		"https://newsapi.org/v2/everything",
		"apiKey",
		"95f7dff332fc496e945a2d707fe50730",
		myjson.Params{
			"q":        "tesla",
			"from":     "2021-06-21",
			"sortBy":   "publishedAt",
			"pageSize": "5",
			"language": "ru",
		},
		&newsInfoNewsapi)

	showTitleDesc(newsInfoNewsapi)
	fmt.Println()

	newsInfoMediastack := NewsInfoMediastack{}
	getOneQuery(
		"http://api.mediastack.com/v1/news",
		"access_key",
		"30b6cdc840f7697d8db0a3ef74384183",
		myjson.Params{
			"date":    "2021-06-21,2021-07-21",
			"sort":    "popularity",
			"limit":   "5",
			"sources": "en",
		},
		&newsInfoMediastack)

	showTitleDesc(newsInfoMediastack)
}
