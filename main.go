package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

type sendMessageReqBody struct {
	ChatID string `json:"chat_id"`
	Text   string `json:"text"`
}

type Post struct {
	ID    int    `json:"id"`
	Title string `json:"title"`
	Url   string `json:"url"`
	Node  struct {
		Title string `json:"title"`
	} `json:"node"`
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
var posts = Posts{}

func getList() ([]int, error) {

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

	for _, post := range posts {
		if post.ID == id {
			node := post.Node.Title
			title := post.Title
			link := post.Url
			reqBody := &sendMessageReqBody{
				ChatID: channel,
				Text:   fmt.Sprintf("#%s %s %s", node, title, link),
			}
			reqBytes, err := json.Marshal(reqBody)
			if err != nil {
				return
			}
			url := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", token)
			_, err = http.Post(url, "application/json", bytes.NewBuffer(reqBytes))
			if err != nil {
				return
			}

		}
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
