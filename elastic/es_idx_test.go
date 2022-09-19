package elastic

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

// @Author Chen Zikang
// @Email  zikang.chen@shopee.com
// @Since  2022-07-16

func TestCreateIndex(t *testing.T) {
	a := assert.New(t)
	response, err := client.Indices.Create("book-0.1.0", client.Indices.Create.WithBody(strings.NewReader(`
	{
		"aliases": {
			"book":{}
		},
		"settings": {
			"analysis": {
				"normalizer": {
					"lowercase": {
						"type": "custom",
						"char_filter": [],
						"filter": ["lowercase"]
					}
				}
			}
		},
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
				"pubDate": {
					"type": "date"
				},
				"pages": {
					"type": "integer"
				}
			}
		}
	}
	`)))
	a.Nil(err)
	logger.Println(response)
}

func TestGetIndex(t *testing.T) {
	a := assert.New(t)
	response, err := client.Indices.Get([]string{"book-0.1.0"})
	a.Nil(err)
	logger.Println(response)
}

func TestDeleteIndex(t *testing.T) {
	a := assert.New(t)
	response, err := client.Indices.Delete([]string{"book-0.1.0"})
	a.Nil(err)
	logger.Println(response)
}
