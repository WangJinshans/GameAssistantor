package service

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"game_assistantor/model"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/olivere/elastic/v7"
	"github.com/rs/zerolog/log"
	"github.com/tidwall/gjson"
)

var indexName string = "summary_info"
var esClient *elastic.Client

func init() {
	var err error
	esClient, err = ConnectElasticSearch()
	if err != nil {
		log.Info().Msgf("fail to connect to elastic, error is: %v", err)
		return
	}
}

func ConnectElasticSearch() (client *elastic.Client, err error) {
	client, err = elastic.NewClient(
		elastic.SetURL("http://127.0.0.1:9200"),
		// elastic.SetBasicAuth("elastic", "PP47bTzS#ktoAvEL"),
		elastic.SetSniff(false))
	if err != nil {
		log.Info().Msgf("fail to connect to server, error is: %v", err)
		return
	}
	return
}

func CreateIndex() {
	mapping := `"mappings":{
        "properties":{
            "code":{
                "type":"keyword"
            },
            "name":{
                "type":"keyword"
            },
            "title":{
                "type":"text"
            },
			"content": {
				"type": "text",
			},
			"rate": {
				"type": "float"
			}
        }
    }`
	esClient.CreateIndex(indexName).BodyString(mapping)
}

func FetchTradeSummary(tradeDate string) {
	client := http.Client{}

	data := map[string]string{
		"date": tradeDate,
		"pc":   "1",
	}
	bs, err := json.Marshal(data)
	if err != nil {
		log.Info().Msgf("fail to decode body, error is: %v", err)
		return
	}
	req, err := http.NewRequest("POST", "https://app.jiucaigongshe.com/jystock-app/api/v1/action/field", bytes.NewReader(bs))
	if err != nil {
		log.Info().Msgf("fail to crate request, error is: %v", err)
		return
	}
	t := time.Now().Unix() * 1000
	req.Header.Set("authority", "app.jiucaigongshe.com")
	req.Header.Set("method", "POST")
	req.Header.Set("path", "/jystock-app/api/v2/article/community")
	req.Header.Set("scheme", "https")
	req.Header.Set("accept", "application/json, text/plain, */*")
	req.Header.Set("accept-encoding", "gzip, deflate, br")
	req.Header.Set("accept-language", "en,zh-CN;q=0.9,zh;q=0.8")
	req.Header.Set("content-type", "application/json")
	req.Header.Set("platform", "3")
	req.Header.Set("origin", "https://www.jiucaigongshe.com")
	req.Header.Set("referer", "https://www.jiucaigongshe.com/")
	req.Header.Set("timestamp", fmt.Sprintf("%d", t))
	req.Header.Set("token", "bbb3f25d0ab1078e367857b05c39ed36")
	req.Header.Set("user-agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/92.0.4515.107 Safari/537.36")

	res, err := client.Do(req)
	if err != nil {
		log.Info().Msgf("fail to get response, error is: %v", err)
		return
	}

	content, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Info().Msgf("fail to get content, error is: %v", err)
		return
	}

	// log.Info().Msgf("content is: %v", string(content))
	value := gjson.Get(string(content), "data")

	value.ForEach(func(key, v gjson.Result) bool {
		// log.Info().Msgf("key is: %v, value is: %v", key.String(), v.String())
		filedId := gjson.Get(v.String(), "action_field_id")
		tradeDate := gjson.Get(v.String(), "date")
		bankName := gjson.Get(v.String(), "name")
		dataList := gjson.Get(v.String(), "list")
		if filedId.String() == "" {
			return true
		}
		if strings.Contains(bankName.String(), "ST") {
			return true
		}
		log.Info().Msgf("trade date is: %v, bank name is: %s", tradeDate.String(), bankName.String())
		for _, item := range dataList.Array() {
			code := gjson.Get(item.String(), "code")
			name := gjson.Get(item.String(), "name")
			if strings.Contains(name.String(), "ST") {
				continue
			}
			log.Info().Msgf("item is: %v", item.String())

			// article := gjson.Get(item.String(), "article")
			rate := gjson.Get(item.String(), "article.action_info.shares_range").Float() / 1000
			expound := gjson.Get(item.String(), "article.action_info.expound").String()
			fmt.Printf("name is: %s\n", name)
			// fmt.Printf("article is: %s\n", article)
			fmt.Printf("rate is: %.2f\n", rate)
			fmt.Printf("expound is: %s\n", expound)
			index := strings.Index(expound, "\n")
			var title string
			if index > 0 {
				title = expound[:index]
				expound = expound[index:]
			}

			info := &model.TradeSummary{
				Name:    name.String(),
				Code:    code.String(),
				Rate:    rate,
				Title:   title,
				Content: expound,
			}
			_, err = esClient.Index().Index(indexName).BodyJson(info).Do(context.Background())
			if err != nil {
				log.Info().Msgf("Failed to index tweet: %s", err)
			}

			esClient.Stop()
		}
		return true
	})
}

func SearchContent() {
	query := elastic.NewQueryStringQuery("+\"机器人\" +\"特斯拉\"")
	// query := elastic.NewMatchPhrasePrefixQuery("content", "+\"华为\" +\"算力\"")

	result, err := esClient.Search().Index(indexName).Size(10).Query(query).Do(context.Background())
	// result, err := esClient.Search().Index(indexName).From(10).Size(1).Query(query).Do(context.Background())
	// 处理返回的结果
	log.Info().Msgf("hit is: %v", len(result.Hits.Hits))
	for _, hit := range result.Hits.Hits {
		log.Info().Msgf("hit is: %#v", string(hit.Source))
		source := make(map[string]interface{})
		err = json.Unmarshal(hit.Source, &source)
		if err != nil {
			log.Info().Msgf("Failed to parse document source: %v", err)
		} else {
			fieldName, ok := source["content"]
			if ok {
				log.Info().Msgf("content is: %v", fieldName.(string))
			}
		}
	}
}
