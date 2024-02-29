package es

import (
	"context"
	"testing"

	"github.com/olivere/elastic/v7"
	"github.com/rs/zerolog/log"
)

func ConnectElasticSearch() (client *elastic.Client, err error) {
	client, err = elastic.NewClient(
		elastic.SetURL("http://es-cn-8ed2l11hp004g3orv.elasticsearch.aliyuncs.com:9200"),
		elastic.SetBasicAuth("elastic", "PP47bTzS#ktoAvEL"),
		elastic.SetSniff(false))
	if err != nil {
		log.Info().Msgf("fail to connect to server, error is: %v", err)
		return
	}
	return
}

func TestElasticSearch(t *testing.T) {
	client, err := ConnectElasticSearch()
	if err != nil {
		log.Info().Msgf("fail to connect to server, error is: %v", err)
		return
	}
	log.Info().Msgf("client is: %v, error is: %v", client, err)
}

func TestCreateIndex(t *testing.T) {
	client, err := ConnectElasticSearch()
	if err != nil {
		log.Info().Msgf("fail to connect to server, error is: %v", err)
		return
	}
	rs, err := client.CreateIndex("internal_test_808").Do(context.Background())
	if err != nil {
		log.Error().Msgf("error is: %v", err)
		return
	}
	log.Info().Msgf("res is: %#v", rs)
}

func TestDeleteIndex(t *testing.T) {
	client, err := ConnectElasticSearch()
	if err != nil {
		log.Info().Msgf("fail to connect to server, error is: %v", err)
		return
	}
	rs, err := client.DeleteIndex("internal_test_808").Do(context.Background())
	if err != nil {
		log.Error().Msgf("error is: %v", err)
		return
	}
	log.Info().Msgf("res is: %#v", rs)
}

func TestFindKeyword(t *testing.T) {
	client, err := ConnectElasticSearch()
	if err != nil {
		log.Info().Msgf("fail to connect to server, error is: %v", err)
		return
	}
	keys := "Receive"
	// 参考: https://cloud.tencent.com/developer/article/1911255
	// termsQuery := elastic.NewTermsQuery("username", "cat", "bob") // NewTermQuery 单字段多值精确匹配
	// rangeQuery := elastic.NewRangeQuery("age").Gte(18).Lte(35) // 范围查询
	// From(0) 从第0行记录开始
	// Size 单页大小

	// boolQuery := elastic.NewBoolQuery()
	// boolQuery.Must(termQuery, rangeQuery)
	// boolQuery.Must(termQuery, rangeQuery) // Must/Filter 类似于sql中的and

	res, err := client.Search().Index("2jmc-zhcd-808-gateway*").Size(100).Sort("@timestamp", false).Query(elastic.NewMatchQuery("message", keys)).Do(context.Background())
	//从搜索结果中取数据的方法
	//for _, item := range res.Each(reflect.TypeOf(typ)) {
	//	if t, ok := item.(Task); ok {
	//		fmt.Println(t)
	//	}
	//}
	if err != nil {
		log.Info().Msgf("fail to search, error is: %v", err)
		return
	}
	//log.Info().Msgf("res is: %#v", res.Hits)
	log.Info().Msgf("res is: %s", res.Hits.Hits[0].Source)
}
