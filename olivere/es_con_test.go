package olivere

import (
	"context"
	"log"
	"os"
	"testing"

	"github.com/olivere/elastic/v7"
)

// @Author Chen Zikang
// @Email  zikang.chen@shopee.com
// @Since  2022-07-20

var (
	ctx    = context.Background()
	client *elastic.Client
	logger = log.New(os.Stdout, "[ES] ",
		log.Lshortfile|log.Ldate|log.Ltime|log.Lmicroseconds,
	)
)

func init() {
	var err error
	client, err = elastic.NewClient(
		elastic.SetURL("http://127.0.0.1:9200"),
		elastic.SetBasicAuth("elastic", "KANG1823"),
		elastic.SetSniff(false),
		elastic.SetTraceLog(logger),
	)
	if err != nil {
		panic(err)
	}
}

func TestNewESClient(t *testing.T) {
	logger.Println(client.String())
}
