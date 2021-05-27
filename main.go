package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

type Post struct {
	ID int `json:"id"`
	//Title string `json:"title"`
	//Url   string `json:"url"`
}
type Posts []Post

func (posts Posts) IDList() []int {
	var list []int
	for _, post := range posts {
		list = append(list, post.ID)
	}
	return list
}

var ids []int

func getList() ([]int, error) {
	posts := Posts{}

	url := "https://www.v2ex.com/api/topics/latest.json"
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(body, &posts)
	if err != nil {
		return nil, err
	}

	return posts.IDList(), nil
}

func difference(a, b []int) (diff []int) {
	m := make(map[int]bool)
	for _, item := range b {
		m[item] = true
	}
	for _, item := range a {
		if _, ok := m[item]; !ok {
			diff = append(diff, item)
		}
	}
	return
}

func init() {
	ids, _ = getList()
}

func push(id int) {
	token := os.Getenv("HedwigToken")
	channel := "@V2EXChannel"
	url := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage?chat_id=%s&text=https://www.v2ex.com/t/%d", token, channel, id)
	_, err := http.Get(url)
	if err != nil {
		return
	}
}

func main() {
	for {
		fetchIds, err := getList()
		if err != nil {
			time.Sleep(60 * time.Second)
			continue
		}
		newIds := difference(fetchIds, ids)
		for _, id := range newIds {
			log.Println(id)
			go push(id)
		}
		ids = fetchIds
		time.Sleep(30 * time.Second)
	}

}
