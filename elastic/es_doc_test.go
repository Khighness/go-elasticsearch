package elastic

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/elastic/go-elasticsearch/v7/esapi"
	"github.com/stretchr/testify/assert"
	"log"
	"net/http"
	"testing"
	"time"
)

// @Author Chen Zikang
// @Email  zikang.chen@shopee.com
// @Since  2022-07-17

type DocWrapper struct {
	Doc interface{} `json:"doc"`
}

type Book struct {
	ID      string     `json:"-"`
	Author  string     `json:"author"`
	Name    string     `json:"name"`
	Pages   int        `json:"pages"`
	Price   float64    `json:"price"`
	PubDate *time.Time `json:"pubDate"`
	Summary string     `json:"summary"`
}

type SearchResponse struct {
	Took     int64 `json:"took"`
	TimeOut  int64 `json:"time_out"`
	MaxScore int64 `json:"max_score"`
}

type SearchShards struct {
	Total      int64 `json:"total"`
	Successful int64 `json:"successful"`
	Skipped    int64 `json:"skipped"`
	Failed     int64 `json:"failed"`
}

type SearchHits struct {
}

func TestCreateDocument(t *testing.T) {
	a := assert.New(t)
	body := &bytes.Buffer{}
	pubDate := time.Now()
	err := json.NewEncoder(body).Encode(&Book{
		Author:  "金庸",
		Name:    "笑傲江湖",
		Pages:   1978,
		Price:   99.9,
		PubDate: &pubDate,
		Summary: "...",
	})
	if err != nil {
		return
	}
	a.Nil(err)
	// 创建文档
	response, err := client.Index("book-0.1.0", body, client.Index.WithDocumentID("10001"))
	a.Nil(err)
	logger.Println(response)
}

func TestUpdateDocument(t *testing.T) {
	a := assert.New(t)
	body := &bytes.Buffer{}
	now := time.Now()
	err := json.NewEncoder(body).Encode(&Book{
		Author:  "金庸",
		Name:    "《神雕侠侣》！！！",
		Pages:   2020,
		Price:   119.9,
		PubDate: &now,
		Summary: "杨过...",
	},
	)
	request := esapi.IndexRequest{
		Index:      "book-0.1.0",
		DocumentID: "10002",
		Body:       body,
	}
	response, err := request.Do(context.Background(), client)
	a.Nil(err)
	logger.Println(response)
}

func TestUpdateDocumentPartly(t *testing.T) {
	a := assert.New(t)
	body := map[string]interface{}{
		"doc": map[string]interface{}{
			"name": "《天龙八部》",
		},
	}
	buf, err := json.Marshal(body)
	a.Nil(err)
	request, err := http.NewRequest(http.MethodPost, "http://0.0.0.0:9200/book-0.1.0/_update/10002", bytes.NewReader(buf))
	a.Nil(err)
	//request := esapi.UpdateRequest{
	//	Index:      "book-0.1.0",
	//	DocumentID: "10002",
	//	Body:       bytes.NewReader(buf),
	//}
	response, err := client.Perform(request)
	//response, err := client.Update("book-0.1.0", "10002", bytes.NewReader(buf), client.Update.WithSourceIncludes("name"))
	a.Nil(err)
	logger.Println(response)
}

func TestGetDocument(t *testing.T) {
	a := assert.New(t)
	response, err := client.Get("book-0.1.0", "10001")
	a.Nil(err)
	logger.Println(response)
}

func TestSearch(t *testing.T) {
	a := assert.New(t)
	body := &bytes.Buffer{}
	body.WriteString(`
	{
		"_source":{
		  "excludes": ["author"]
		}, 
		"query": {
		  "match_phrase": {
			"author": "金庸"
		  }
		},
		"sort": [
		  {
			"pages": {
			  "order": "desc"
			}
		  }
		], 
		"from": 0,
		"size": 5
	}
	`)
	response, err := client.Search(client.Search.WithIndex("book-0.1.0"), client.Search.WithBody(body))
	a.Nil(err)
	var result map[string]interface{}
	err = json.NewDecoder(response.Body).Decode(&result)
	a.Nil(err)

	for _, hit := range result["hits"].(map[string]interface{})["hits"].([]interface{}) {
		h := hit.(map[string]interface{})
		var book Book
		buf, err := json.Marshal(h["_source"])
		a.Nil(err)
		err = json.Unmarshal(buf, &book)
		a.Nil(err)
		log.Printf(" ID=%v, %v", h["_id"], book)
	}
	logger.Println(response)
}

func TestDeleteDocument(t *testing.T) {
	a := assert.New(t)
	response, err := client.Delete("book-0.1.0", "10003")
	a.Nil(err)
	logger.Println(response)
}

func TestBulk(t *testing.T) {
	books := []*Book{
		{
			ID:     "10002",
			Author: "金庸",
			Name:   "神雕侠侣",
		},
		{
			ID:     "10003",
			Author: "金庸",
			Name:   "连城诀",
		},
	}
	a := assert.New(t)
	body := &bytes.Buffer{}
	for _, book := range books {
		meta := []byte(fmt.Sprintf(`{"index": {"_id": "%s"} }%s`, book.ID, "\n"))
		data, err := json.Marshal(book)
		a.Nil(err)
		data = append(data, "\n"...)
		body.Grow(len(meta) + len(data))
		body.Write(meta)
		body.Write(data)
	}
	response, err := client.Bulk(body, client.Bulk.WithIndex("book-0.1.0"))
	a.Nil(err)
	logger.Println(response)
}
