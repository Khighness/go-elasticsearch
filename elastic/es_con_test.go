package elastic

import (
	"log"
	"os"
	"testing"

	es "github.com/elastic/go-elasticsearch/v7"
)

// @Author Chen Zikang
// @Email  zikang.chen@shopee.com
// @Since  2022-07-16

var (
	client *es.Client
	logger = log.New(os.Stdout, "[ES] ",
		log.Lshortfile|log.Ldate|log.Ltime|log.Lmicroseconds,
	)
)

func init() {
	var err error
	config := es.Config{
		Addresses: []string{"http://0.0.0.0:9200"},
		Username:  "elastic",
		Password:  "KANG1823",
	}
	client, err = es.NewClient(config)
	if err != nil {
		logger.Fatalln(err)
	}
}

func TestNewESClient(t *testing.T) {
	logger.Println(client.Info())
}
