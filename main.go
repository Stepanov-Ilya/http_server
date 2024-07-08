package main

import (
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"sync"
	"time"
)

const (
	baseUrl       = "localhost:8080"
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

type SyncMap struct {
	data  map[int64]*News
	mutex sync.RWMutex
}

var newsMap = &SyncMap{
	data: make(map[int64]*News),
}

func createNewsHandler(writer http.ResponseWriter, request *http.Request) {
	info := &NewsInfo{}
	if err := json.NewDecoder(request.Body).Decode(info); err != nil {
		http.Error(writer, "Failed to decode news data", http.StatusBadRequest)
		return
	}

	rand.Seed(time.Now().UnixNano())
	now := time.Now()

	news := &News{
		ID:        rand.Int63(),
		Info:      *info,
		CreatedAt: now,
		UpdatedAt: now,
	}

	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(writer).Encode(news); err != nil {
		http.Error(writer, "Failed to encode news data", http.StatusInternalServerError)
		return
	}

	newsMap.mutex.Lock()
	defer newsMap.mutex.Unlock()

	newsMap.data[news.ID] = news
}

func getNewsHandler(writer http.ResponseWriter, request *http.Request) {
	newsId := chi.URLParam(request, "id")
	id, err := parseNewsID(newsId)
	if err != nil {
		http.Error(writer, "Invalid news ID", http.StatusBadRequest)
		return
	}

	newsMap.mutex.RLock()
	defer newsMap.mutex.RUnlock()

	news, ok := newsMap.data[id]
	if !ok {
		http.Error(writer, "News not found", http.StatusNotFound)
		return
	}

	writer.Header().Set("Content-Type", "application/json")
	if err = json.NewEncoder(writer).Encode(news); err != nil {
		http.Error(writer, "Failed to encode news data", http.StatusInternalServerError)
		return
	}
}

func parseNewsID(idStr string) (int64, error) {
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func main() {
	router := chi.NewRouter()

	router.Post(createPostfix, createNewsHandler)
	router.Get(getPostfix, getNewsHandler)

	err := http.ListenAndServe(baseUrl, router)
	if err != nil {
		log.Fatal(err)
	}
}
