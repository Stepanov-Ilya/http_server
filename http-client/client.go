package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/brianvoe/gofakeit"
	"github.com/fatih/color"
	"io"
	"log"
	"net/http"
	"time"
)

const (
	baseUrl       = "http://localhost:8080"
	createPostfix = "/news"
	getPostfix    = "/news/%d"
)

type NewsInfo struct {
	Title    string    `json:"title"`
	Context  string    `json:"context"`
	Reporter string    `json:"reporter"`
	Country  string    `json:"country"`
	Time     time.Time `json:"time"`
}

type News struct {
	ID        int64     `json:"id"`
	Info      NewsInfo  `json:"info"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

func createNews() (News, error) {
	news := NewsInfo{
		Title:    gofakeit.Sentence(3),
		Context:  gofakeit.Quote(),
		Reporter: gofakeit.Name(),
		Country:  gofakeit.Country(),
		Time:     gofakeit.Date(),
	}

	data, err := json.Marshal(news)
	if err != nil {
		return News{}, err
	}

	resp, err := http.Post(baseUrl+createPostfix, "application/json", bytes.NewBuffer(data))
	if err != nil {
		return News{}, err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {

		}
	}(resp.Body)

	if resp.StatusCode != http.StatusCreated {
		return News{}, err
	}

	var createdNews News
	if err = json.NewDecoder(resp.Body).Decode(&createdNews); err != nil {
		return News{}, err
	}

	return createdNews, nil
}

func getNews(id int64) (News, error) {
	resp, err := http.Get(fmt.Sprintf(baseUrl+getPostfix, id))
	if err != nil {
		log.Fatal("Failed to get news: ", err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(resp.Body)

	if resp.StatusCode == http.StatusNotFound {
		return News{}, err
	}

	if resp.StatusCode != http.StatusOK {
		return News{}, errors.New(fmt.Sprintf("Failed to get news: %d", resp.StatusCode))
	}

	var news News
	if err = json.NewDecoder(resp.Body).Decode(&news); err != nil {
		return News{}, err
	}

	return news, nil
}

func main() {
	news, err := createNews()
	if err != nil {
		log.Fatal("Failed to create news: ", err)
	}
	log.Printf(color.RedString("News created:\n"), color.GreenString("&v", news))

	news, err = getNews(news.ID)
	if err != nil {
		log.Fatal("Failed to get news: ", err)
	}
	log.Printf(color.RedString("News created:\n"), color.GreenString("&v", news))
}
