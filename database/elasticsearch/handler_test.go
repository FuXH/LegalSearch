package elasticsearch

import (
	"fmt"
	es "github.com/olivere/elastic/v7"
	"testing"

	conf "LegalSearch/conf"

	"github.com/stretchr/testify/assert"
)

func TestGetVersion(t *testing.T) {
	esConf := &conf.ElasticSearchConf{
		Addresses: []string{"http://127.0.0.1:9200"},
	}
	esHandler, err := NewEsHandler(esConf)
	assert.Nil(t, nil, err, "es初始化失败")

	info, err := esHandler.GetVersion()
	assert.Equal(t, "7.10.0", info, "es版本信息错误")
}

func TestCreateIndex(t *testing.T) {
	esConf := &conf.ElasticSearchConf{
		Addresses: []string{"http://127.0.0.1:9200"},
	}
	esHandler, err := NewEsHandler(esConf)
	assert.Nil(t, nil, err, "es初始化失败")

	settings := `
{
	"settings": {
        "number_of_shards": 1,
        "number_of_replicas": 0
        // "index.blocks.read_only_allow_delete": null
    },
	"mappings": {
		"properties": {
			"Name": {
				"type": "text",
				"analyzer": "ik_max_word"
			},
			"Title": {
				"type": "text",
				"analyzer": "ik_max_word"
			},
			"Desc": {
				"type": "text",
				"analyzer": "ik_max_word"
			},
			"Age": {
				"type": "integer"
			}
		}
	}
}
`
	err = esHandler.CreateIndex("test_fxh2", settings)
	assert.Nil(t, err, "create index fail")
}

type TestType struct {
	Name  string
	Title string
	Desc  string
	Age   int
}

func TestInsert(t *testing.T) {
	esConf := &conf.ElasticSearchConf{
		Addresses: []string{"http://127.0.0.1:9200"},
	}
	esHandler, err := NewEsHandler(esConf)
	assert.Nil(t, nil, err, "es初始化失败")

	data := TestType{
		Name:  "a2",
		Title: "b2",
		Desc:  "c2",
		Age:   123,
	}
	err = esHandler.Insert(data, "test_fxh")
	assert.Nil(t, err, "插入数据错误")
}

func TestDelete(t *testing.T) {
	esConf := &conf.ElasticSearchConf{
		Addresses: []string{"http://127.0.0.1:9200"},
	}
	esHandler, err := NewEsHandler(esConf)
	assert.Nil(t, nil, err, "es初始化失败")

	err = esHandler.Delete("test_fxh", "idx")
	assert.Nil(t, err, "删除数据失败")
}

func TestQueryById(t *testing.T) {
	esConf := &conf.ElasticSearchConf{
		Addresses: []string{"http://127.0.0.1:9200"},
	}
	esHandler, err := NewEsHandler(esConf)
	assert.Nil(t, nil, err, "es初始化失败")

	data := &TestType{}
	err = esHandler.QueryById("test_fxh", "1", data)
	assert.Nil(t, err, "查询数据失败")
	fmt.Println("data:", data)
	assert.NotNil(t, data, "查询数据为空")
}

func TestGetQueryHits(t *testing.T) {
	esConf := &conf.ElasticSearchConf{
		Addresses: []string{"http://127.0.0.1:9200"},
	}
	esHandler, err := NewEsHandler(esConf)
	assert.Nil(t, nil, err, "es初始化失败")

	data := &TestType{}
	query := es.NewMatchPhrasePrefixQuery("name", "lip")
	searchResult, err := esHandler.BoolQuery("test_fxh", "age", false, 3, nil, query)
	res, err := esHandler.GetQueryHits(data, searchResult)

	assert.Nil(t, err, "get query hits fail")
	for _, v := range res.([]*TestType) {
		fmt.Println("out:", v)
	}
}

func TestGetQueryLen(t *testing.T) {
	esConf := &conf.ElasticSearchConf{
		Addresses: []string{"http://127.0.0.1:9200"},
	}
	esHandler, err := NewEsHandler(esConf)
	assert.Nil(t, nil, err, "es初始化失败")

	searchResult, err := esHandler.BoolQuery("*", "age", false, 3, nil)
	total, err := esHandler.GetQueryLen(searchResult)
	assert.Nil(t, err, "get query len fail")
	fmt.Println("get query len:", total)
}

func TestGetQueryAggs(t *testing.T) {
	esConf := &conf.ElasticSearchConf{
		Addresses: []string{"http://127.0.0.1:9200"},
	}
	esHandler, err := NewEsHandler(esConf)
	assert.Nil(t, nil, err, "es初始化失败")

	aggsMap := make(map[string]es.Aggregation)
	aggsMap["agg_match"] = es.NewFilterAggregation().Filter(es.NewMatchPhraseQuery("obj.a", 123))
	searchResult, _ := esHandler.BoolQuery("test_fxh", "age", false, 3, aggsMap)

	aggsOutput := &struct {
		AggMatch struct {
			DocCount int `json:"doc_count"`
		} `json:"agg_match"`
	}{}
	err = esHandler.GetQueryAggs(aggsOutput, searchResult)
	assert.Nil(t, err, "get query aggs")

	fmt.Println("aggsOutput:", aggsOutput)
}
