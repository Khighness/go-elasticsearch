package olivere

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

// @Author Chen Zikang
// @Email  zikang.chen@shopee.com
// @Since  2022-07-20

func TestCreateIndex(t *testing.T) {
	a := assert.New(t)
	const mapping = `
	{
		"mappings": {
			"properties": {
				"name": {
					"type": "keyword",
					"normalizer": "lowercase"
				},
				"price": {
					"type": "double"
				},
				"summary": {
					"type": "text"
				},
				"author": {
					"type": "keyword"
				},
				"pub_date": {
					"type": "date"
				},
				"pages": {
					"type": "integer"
				}
			}
		}
	}
	`
	exits, err := client.IndexExists("book").Do(ctx)
	a.Nil(err)
	if exits {
		logger.Fatalln("Index [book] already exists")
	} else {
		response, err := client.
			CreateIndex("book").
			BodyString(mapping).
			Do(ctx)
		a.Nil(err)
		a.EqualValues(true, response.Acknowledged)
		logger.Println(response)
	}
}

func TestGetIndexMapping(t *testing.T) {
	a := assert.New(t)
	mapping, err := client.GetMapping().
		Type(""). // necessarily
		Index("book").
		Do(ctx)
	a.Nil(err)
	logger.Println(mapping)
}

func TestAddMappingField(t *testing.T) {
	a := assert.New(t)
	const mapping = `
	{
		"properties": {
			"isbn": {
				"type": "text"
			}
		}
	}
	`
	response, err := client.
		PutMapping().
		Index("book").
		BodyString(mapping).
		Do(ctx)
	a.Nil(err)
	logger.Println(response)
}
