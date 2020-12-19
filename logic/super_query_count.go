package logic

import (
	"fmt"
	"github.com/gin-gonic/gin"
	es "github.com/olivere/elastic/v7"
	"time"

	"LegalSearch/constant"
)

type QueryCount struct {
	Success bool `form:"success"`
	Count   int  `form:"count"`
}

func SuperQueryCount(ctx *gin.Context) {
	startTime := time.Now()
	status := "success"
	defer func() {
		costTime := time.Since(startTime).Seconds()
		fmt.Println("query_by_condition接口耗时(s):", costTime)
		if status != "success" {
			rsp := &QueryCount{
				Success: false,
				Count:   0,
			}
			GinReturn(ctx, rsp)
		}
	}()

	param := new(QueryByConditionReq)
	if err := ctx.ShouldBind(param); err != nil {
		status = fmt.Sprintf("QueryCount接口输入参数解析失败, err:", err)
		fmt.Println(status)
		return
	}
	fmt.Println("query_count input:", param)

	// 构造查询
	indexName, err := getIndex(param.TrialYear)
	if err != nil {
		status = fmt.Sprintf("QueryCount输入参数错误，审理年份非法:", param.TrialYear, "err:", err)
		fmt.Println(status)
		return
	}
	filters := getFilters(param)
	aggsMap := make(map[string]es.Aggregation)
	aggsMap["win_count"] = es.NewFilterAggregation().Filter(es.NewTermQuery("IsWin", constant.WinString))
	aggsMap["lose_count"] = es.NewFilterAggregation().Filter(es.NewTermQuery("IsWin", constant.LoseString))

	// 查询es
	searchResult, err := GetEsHandler().BoolQuery(indexName,
		constant.SortField, constant.SortOrder, 0,
		aggsMap,
		filters...)
	if err != nil {
		status = fmt.Sprintf("QueryByCondition查询返回错误, err:", err)
		fmt.Println(status)
		return
	}

	// 处理返回数据
	count := aggregateLen(searchResult)

	// rsp
	rsp := &QueryCount{
		Success: true,
		Count:   count,
	}
	GinReturn(ctx, rsp)
}

func aggregateLen(searchResult *es.SearchResult) int {
	// 胜诉概率
	aggsOutput := &struct {
		WinCount struct {
			DocCount int `json:"doc_count"`
		} `json:"win_count"`
		LoseCount struct {
			DocCount int `json:"doc_count"`
		} `json:"lose_count"`
	}{}
	if err := GetEsHandler().GetQueryAggs(aggsOutput, searchResult); err != nil {
		return 0
	}

	return aggsOutput.WinCount.DocCount + aggsOutput.LoseCount.DocCount
}
