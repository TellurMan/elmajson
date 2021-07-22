package main

import (
	"elmajson/myjson"
	"fmt"
	"os"
	"sync"
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

	client := myjson.NewClient(url, keyName, key)

	err := client.GetJSON(params, &news)
	if err != nil {
		fmt.Println(err)
		return
	}
}

func showTitleDesc(news interface{}) string {
	answer := ""
	switch newsType := news.(type) {
	case NewsInfoNewsapi:
		for _, item := range newsType.Articles {
			answer += fmt.Sprint(item.Title, "\n", item.Description, "\n\n")
		}
		return answer
	case NewsInfoMediastack:
		for _, item := range newsType.Data {
			answer += fmt.Sprint(item.Title, "\n", item.Description, "\n\n")
		}
		return answer
	default:
		answer += "error showTitleDesc"
	}
	return answer
}

func main() {
	query := "tesla"
	if len(os.Args) > 1 {
		query = os.Args[1]
	}

	fmt.Println("start work")

	var wg sync.WaitGroup       // для ожидания завершения всех горутин
	out := make(chan string, 1) // канал передачи данных

	for i := 1; i <= 10; i++ { // запрашиваю 10 страниц
		wg.Add(1)        // увеличение количества горутин, которые надо подождать
		go func(i int) { // горутина запроса
			newsInfoNewsapi := NewsInfoNewsapi{}
			getOneQuery(
				"https://newsapi.org/v2/everything",
				"apiKey",
				"95f7dff332fc496e945a2d707fe50730",
				myjson.Params{
					"q":        query,
					"from":     "2021-06-24",
					"pageSize": "10",
					"page":     fmt.Sprint(i),
					"language": "ru",
				},
				&newsInfoNewsapi)
			// отправка в канал человекочитаемой инфы по новостям
			out <- fmt.Sprint("-------------------страница ", i, "---------------------\n", showTitleDesc(newsInfoNewsapi))
			defer wg.Done() // уменьшение количества горутин, которых надо подождать
		}(i)
	}

	go func() { // приём из канала инфы по новостям
		for s := range out {
			fmt.Println(s)
		}
	}()
	wg.Wait()  // ожидание завершения всех запросов к сайту
	close(out) // закрытие канала для завершения приёма из канала

	/*
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
	*/
}
