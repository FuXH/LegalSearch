package elasticsearch

import (
	"context"
	"encoding/json"
	"fmt"
	"reflect"

	"LegalSearch/conf"

	es "github.com/olivere/elastic/v7"
)

type EsHandler struct {
	client    *es.Client
	addresses []string
}

func NewEsHandler(conf *conf.ElasticSearchConf) (*EsHandler, error) {
	esClient, err := es.NewClient(
		es.SetURL(conf.Addresses...),
	)
	if err != nil {
		return nil, err
	}
	handler := &EsHandler{
		client:    esClient,
		addresses: conf.Addresses,
	}
	return handler, nil
}

// 获取es版本
func (p *EsHandler) GetVersion() (string, error) {
	client := p.client
	if len(p.addresses) == 0 {
		return "", fmt.Errorf("没有配置es地址")
	}
	address := p.addresses[0]
	version, err := client.ElasticsearchVersion(address)
	if err != nil {
		return "", err
	}

	return version, nil
}

// 创建索引
func (p *EsHandler) CreateIndex(index string, mappings string) error {
	client := p.client
	isExist, err := client.IndexExists(index).Do(context.Background())
	if err != nil {
		return err
	}
	if isExist {
		return fmt.Errorf("创建索引失败，索引已存在")
	}
	createIndex, err := client.CreateIndex(index).Body(mappings).Header("include_type_name", "true").Do(context.Background())
	if err != nil {
		return err
	}
	if !createIndex.Acknowledged {
		return fmt.Errorf("创建索引未应答")
	}

	return nil
}

// 索引是否存在
func (p *EsHandler) IsExistIndex(indexs []string) bool {
	client := p.client
	exist, err := client.IndexExists(indexs...).Do(context.Background())
	if err != nil {
		return false
	}
	return exist
}

type BulkData struct {
	Data  interface{}
	Index string
	Id    string
}

// 批量写入数据
func (p *EsHandler) BulkInsert(data interface{}) error {
	client := p.client
	bulkRequest := client.Bulk()

	for _, value := range data.([]BulkData) {
		req := es.NewBulkIndexRequest().Index(value.Index).Doc(value.Data)
		if value.Id != "" {
			req = req.Id(value.Id)
		}
		bulkRequest.Add(req)
	}
	bulkResponse, err := bulkRequest.Do(context.Background())
	if err != nil {
		return err
	}

	// 统计写入的耗时
	fmt.Println("导入数据数量:", len(bulkResponse.Indexed()), " 耗时:", bulkResponse.Took)

	return nil
}

// 插入数据
func (p *EsHandler) Insert(data interface{}, index string) error {
	client := p.client
	put1, err := client.Index().
		Index(index).
		BodyJson(data).
		Do(context.Background())
	if err != nil {
		return err
	}
	fmt.Printf("Indexed tweet %s to index %s, type %s\n", put1.Id, put1.Index, put1.Type)

	return nil
}

// 指定ID插入数据
func (p *EsHandler) InsertWithId(data interface{}, index string, id string) error {
	client := p.client
	put1, err := client.Index().
		Index(index).
		Id(id).
		BodyJson(data).
		Do(context.Background())
	if err != nil {
		return err
	}
	fmt.Printf("Indexed tweet %s to index %s, type %s\n", put1.Id, put1.Index, put1.Type)

	return nil
}

// 删除数据
func (p *EsHandler) Delete(index string, id string) error {
	client := p.client
	res, err := client.Delete().
		Index(index).
		Id(id).
		Do(context.Background())
	if err != nil {
		return err
	}
	fmt.Println("res:", res)
	//if res.Result
	return nil
}

// 检索 - 通过id
func (p *EsHandler) QueryById(index string, id string, out interface{}) error {
	client := p.client

	get, err := client.Get().Index(index).Id(id).Do(context.Background())
	if err != nil {
		return err
	}
	if !get.Found {
		return fmt.Errorf("未找到该id数据, index: %s, id: %s", index, id)
	}

	if err := json.Unmarshal(get.Source, out); err != nil {
		return err
	}

	return nil
}

// BoolQuery bool查询
// outType: 查询数据的结构体指针
// sortField: 排序字段（数字类型）
// order: true-升序排序，false-降序排序
// size: 查询结果返回的最大数量
// aggs: 聚合统计的条件
func (p *EsHandler) BoolQuery(index []string,
	sortField string, order bool, size int,
	aggs map[string]es.Aggregation,
	filters ...es.Query) (*es.SearchResult, error) {
	client := p.client

	// match
	searchService := client.Search().
		Index(index...).
		Pretty(true)
	if sortField != "" {
		searchService = searchService.Sort(sortField, order)
	}
	if size != 0 {
		searchService = searchService.Size(size)
	}
	// aggs
	for key, val := range aggs {
		searchService = searchService.Aggregation(key, val)
	}
	// filters
	boolQuery := es.NewBoolQuery().Filter(filters...)
	searchService = searchService.Query(boolQuery)

	searchResult, err := searchService.Do(context.Background())
	if err != nil {
		return nil, err
	}

	return searchResult, nil
}

// GetQueryHits 获取查询的数据
// outType: 查询数据的结构体指针
// 返回值: 查询数据的数组
func (p *EsHandler) GetQueryHits(outType interface{}, searchResult *es.SearchResult) (interface{}, error) {
	// 根据传入的outType创建返回数组
	refVal := reflect.ValueOf(outType)
	res := reflect.MakeSlice(reflect.SliceOf(refVal.Type()),
		0,
		int(searchResult.Hits.TotalHits.Value))
	for index := 0; index < len(searchResult.Hits.Hits); index++ {
		temp := reflect.New(refVal.Elem().Type())
		if err := json.Unmarshal(searchResult.Hits.Hits[index].Source, temp.Interface()); err != nil {
			return nil, err
		}
		res = reflect.Append(res, temp)
	}

	return res.Interface(), nil
}

// GetQueryLen 获取匹配的总数量
func (p *EsHandler) GetQueryLen(searchResult *es.SearchResult) (int64, error) {
	return searchResult.Hits.TotalHits.Value, nil
}

// GetQueryAggs 在匹配查询的基础上做聚合统计
// aggsOutput: 聚合统计的输出结构体
func (p *EsHandler) GetQueryAggs(aggsOutput interface{},
	searchResult *es.SearchResult) error {
	aggByte, err := json.Marshal(searchResult.Aggregations)
	if err != nil {
		return err
	}
	if err := json.Unmarshal(aggByte, aggsOutput); err != nil {
		return err
	}

	return nil
}

func (p *EsHandler) CleanEsIndex() {
	client := p.client

	indexs, _ := client.IndexNames()
	client.DeleteIndex(indexs...).Do(context.Background())
}

func (p *EsHandler) GetIndexs() string {
	client := p.client
	res := ""

	indexs, _ := client.IndexNames()
	for _, index := range indexs {
		res = res + index + "; "
	}

	return res
}
