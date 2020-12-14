package logic

import (
	"LegalSearch/constant"
	"fmt"
	es "github.com/olivere/elastic/v7"
	"time"

	"github.com/gin-gonic/gin"
)

// 模糊搜索输入参数
type QueryReq struct {
	SearchParam string `form:"searchParam"`
}

// 模糊搜索
func Query(ctx *gin.Context) {
	startTime := time.Now()
	status := "success"
	defer func() {
		costTime := time.Since(startTime)
		fmt.Println("query接口耗时(s):", costTime)
		if status != "success" {
			rsp := &ApiQueryRsp{
				Success:   false,
				ErrorMsg:  status,
				ErrorCode: 1,
			}
			GinReturn(ctx, rsp)
		}
	}()

	param := new(QueryReq)
	if err := ctx.ShouldBind(param); err != nil {
		status = fmt.Sprintf("Query接口输入参数解析失败, err:", err)
		fmt.Println(status)
		return
	}
	fmt.Println("query input:", param)
	if param.SearchParam == "" {
		status := fmt.Sprintf("query输入参数为空")
		fmt.Println(status)
		return
	}

	// 构造查询
	indexName := constant.IndexNameInstrument
	filters := getFuzzFilter(param)
	aggsMap := make(map[string]es.Aggregation)
	aggsMap["win_count"] = es.NewFilterAggregation().Filter(es.NewTermQuery("IsWin", 1))
	aggsMap["lose_count"] = es.NewFilterAggregation().Filter(es.NewTermQuery("IsWin", 2))

	// 查询es
	searchResult, err := GetEsHandler().BoolQuery(indexName,
		constant.SortField, constant.SortOrder, constant.SortSize,
		aggsMap,
		filters...)
	if err != nil {
		status = fmt.Sprintf("模糊查询失败, err:", err)
		fmt.Println(status)
		return
	}

	// 处理返回数据
	queryRes := aggregateFuzzQueryData(searchResult)

	rsp := &ApiQueryRsp{
		Success:   true,
		ErrorMsg:  "",
		ErrorCode: 0,
		Content:   queryRes,
	}
	GinReturn(ctx, rsp)
}

func getFuzzFilter(req *QueryReq) []es.Query {
	query := make([]es.Query, 0)

	if req.SearchParam != "" {
		query = append(query,
			es.NewMatchQuery("Content", req.SearchParam))
	}

	return query
}

func aggregateFuzzQueryData(searchResult *es.SearchResult) *QueryRes {
	queryRes := &QueryRes{
		WinRate: 0,
	}

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
		return queryRes
	}
	winRate := float64(aggsOutput.WinCount.DocCount) / float64(aggsOutput.WinCount.DocCount+aggsOutput.LoseCount.DocCount)
	fmt.Println("number:", aggsOutput.WinCount.DocCount, aggsOutput.WinCount.DocCount+aggsOutput.LoseCount.DocCount)

	// 证据建议、常用法条、法官意见
	// 根据匹配查询后，展示时间最近的10条数据
	outType := &EsDataInstrument{}
	hits, err := GetEsHandler().GetQueryHits(outType, searchResult)
	if err != nil {
		return queryRes
	}
	tempEvidences := []string{}
	tempInuseLaws := []string{}
	tempJudgeArguments := []JudgeArgument{}
	fmt.Println("hits:", hits, hits.([]*EsDataInstrument)[0])
	for _, val := range hits.([]*EsDataInstrument) {
		for _, caseSummary := range val.Summarys {
			tempEvidences = append(tempEvidences, caseSummary.Evidence...)
			tempInuseLaws = append(tempInuseLaws, caseSummary.InuseLaw...)
			tempJudgeArguments = append(tempJudgeArguments,
				JudgeArgument{
					//Data:   val.JudgeArgument,
					TextId: val.InstrumentId,
				})
		}
	}

	queryRes = &QueryRes{
		WinRate:       winRate,
		Evidence:      tempEvidences,
		InuseLaw:      tempInuseLaws,
		JudgeArgument: tempJudgeArguments,
	}
	fmt.Println("res:", *queryRes)
	return queryRes
}
