package olivere

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"testing"
	"time"

	"github.com/olivere/elastic/v7"
	"github.com/stretchr/testify/assert"
)

// @Author Chen Zikang
// @Email  zikang.chen@shopee.com
// @Since  2022-07-20

type Book struct {
	ID      string     `json:"-"`
	Author  string     `json:"author"`
	Name    string     `json:"name"`
	Pages   int        `json:"pages"`
	Price   float64    `json:"price"`
	PubDate *time.Time `json:"pubDate"`
	Summary string     `json:"summary"`
	ISBN    string     `json:"isbn"`
}

func TestCreateDoc(t *testing.T) {
	a := assert.New(t)
	now := time.Now()
	book := Book{
		ID:      "1",
		Author:  "金庸",
		Name:    "《笑傲江湖》",
		Pages:   2999,
		Price:   133.9,
		PubDate: &now,
		Summary: "令狐冲和东方不败",
		ISBN:    "1-1-1-1",
	}
	response, err := client.Index().
		Index("book").
		BodyJson(book).
		Id(book.ID).
		Do(ctx)
	a.Nil(err)
	logger.Println(response)
}

func TestUpdateDoc(t *testing.T) {
	a := assert.New(t)
	response, err := client.Update().Index("book").Id("1").
		Script(elastic.NewScript(
			`ctx._source.price += params.delta_price; 
					ctx._source.pub_date = params.pub_date`).
			Param("delta_price", 3.1).
			Param("pub_date", time.Now())).
		Do(ctx)
	a.Nil(err)
	logger.Println(response)
}

func TestUpsertDoc(t *testing.T) {
	a := assert.New(t)
	response, err := client.Update().Index("book").Id("1").
		Script(elastic.NewScript(`ctx._source.is_deleted = params.deleted`).
			Param("deleted", true)).
		Upsert(map[string]interface{}{"is_deleted": false}).
		Do(ctx)
	a.Nil(err)
	logger.Println(response)
}

func TestSearch(t *testing.T) {
	a := assert.New(t)
	query := elastic.NewBoolQuery()
	query.Should(
		elastic.NewMatchQuery("name", "《笑傲江湖》"),
		elastic.NewMatchQuery("name", "《神雕侠侣》"),
	)
	query.Filter(
		elastic.NewRangeQuery("price").Gte(100).Lte(200),
	)
	result, err := client.Search("book").
		Query(query).
		Pretty(true).
		Do(ctx)
	a.Nil(err)
	for _, hit := range result.Hits.Hits {
		var book Book
		err := json.Unmarshal(hit.Source, &book)
		a.Nil(err)
		logger.Printf("%+v", book)
	}
}

func TestBulk(t *testing.T) {
	a := assert.New(t)
	books := make([]Book, 0)
	for i := 0; i < 100; i++ {
		now := time.Now()
		books = append(books, Book{
			ID:      fmt.Sprintf("%d", i+2),
			Author:  "KHighness",
			Name:    fmt.Sprintf("BOOK-%d", i+2),
			Pages:   int(rand.Int63n(10000)),
			Price:   rand.Float64() * float64(rand.Int63n(1000)),
			PubDate: &now,
			Summary: fmt.Sprintf("SUMMARY-%d", i+2),
			ISBN:    fmt.Sprintf("ISBN-%d", i+2),
		})
	}
	bulkRequest := client.Bulk()
	for _, book := range books {
		req := elastic.NewBulkCreateRequest().Index("book").
			Id(book.ID).Doc(book)
		bulkRequest.Add(req)
	}
	bulkResponse, err := bulkRequest.Do(ctx)
	a.Nil(err)
	logger.Println(bulkResponse)
}
