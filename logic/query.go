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
	fmt.Println("query")
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
		status = fmt.Sprintf("query输入参数为空")
		fmt.Println(status)
		return
	}

	// 构造查询
	indexName := []string{constant.IndexNameInstrument}
	filters := getFuzzFilter(param)
	aggsMap := make(map[string]es.Aggregation)
	aggsMap["terms_id"] = es.NewTermsAggregation().Field("wenshu_id_2")
	//aggsMap["win_count"] = es.NewFilterAggregation().Filter(es.NewTermQuery("case_summary.judgement", constant.WinString))
	//aggsMap["lose_count"] = es.NewFilterAggregation().Filter(es.NewTermQuery("case_summary.judgement", constant.LoseString))

	// 查询es
	searchResult, err := GetEsHandler().BoolQuery(indexName,
		"", constant.SortOrder, constant.SortSize,
		aggsMap,
		filters...)
	if err != nil {
		status = fmt.Sprintf("模糊查询失败, err:", err)
		fmt.Println(status)
		return
	}

	// 获取匹配的案件id
	matchIds, err := getMatchInstrumentIds(searchResult)
	if err != nil {
		status = fmt.Sprintf("获取匹配的案件id失败, err:", err)
		fmt.Println(status)
		return
	}

	// 获取案件id对应的返回数据
	queryRes := getMatchControversyInfo(matchIds)

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
			es.NewMatchQuery("wenshu_content", req.SearchParam))
	}

	return query
}

func getMatchInstrumentIds(searchResult *es.SearchResult) ([]string, error) {
	resMatchIds := []string{}

	aggsOutput := &struct {
		TermsId struct {
			Buckets []struct {
				Key      string `json:"key"`
				DocCount int    `json:"doc_count"`
			} `json:"buckets"`
		} `json:"terms_id"`
	}{}
	if err := GetEsHandler().GetQueryAggs(aggsOutput, searchResult); err != nil {
		return resMatchIds, err
	}

	for _, val := range aggsOutput.TermsId.Buckets {
		resMatchIds = append(resMatchIds, val.Key)
	}

	return resMatchIds, nil
}

func getMatchControversyInfo(matchIds []string) *QueryRes {
	indexName := []string{"trial_year_*"}
	filters := []es.Query{}
	for _, id := range matchIds {
		filters = append(filters,
			es.NewTermQuery("WenshuId", id))
	}
	aggsMap := make(map[string]es.Aggregation)
	aggsMap["win_count"] = es.NewFilterAggregation().Filter(es.NewTermQuery("IsWin", constant.WinString))
	aggsMap["lose_count"] = es.NewFilterAggregation().Filter(es.NewTermQuery("IsWin", constant.LoseString))
	aggsMap["inuselaw_count"] = es.NewTermsAggregation().Field("InuseLaw")

	// 查询es
	searchResult, err := GetEsHandler().BoolShouldQuery(indexName,
		"", constant.SortOrder, constant.SortSize,
		aggsMap,
		1, filters...)
	if err != nil {
		fmt.Println("模糊查询第二次失败, err:", err)
		return nil
	}

	return aggregateData(searchResult)
}

//func aggregateFuzzQueryData(searchResult *es.SearchResult) *QueryRes {
//	var winRate float64 = 0
//	queryRes := &QueryRes{
//		WinRate: winRate,
//	}
//
//	// 胜诉概率
//	aggsOutput := &struct {
//		WinCount struct {
//			DocCount int `json:"doc_count"`
//		} `json:"win_count"`
//		LoseCount struct {
//			DocCount int `json:"doc_count"`
//		} `json:"lose_count"`
//	}{}
//	if err := GetEsHandler().GetQueryAggs(aggsOutput, searchResult); err != nil {
//		return queryRes
//	}
//	if aggsOutput.WinCount.DocCount+aggsOutput.LoseCount.DocCount != 0 {
//		winRate = float64(aggsOutput.WinCount.DocCount) / float64(aggsOutput.WinCount.DocCount+aggsOutput.LoseCount.DocCount)
//	}
//	fmt.Println("number:", aggsOutput.WinCount.DocCount, aggsOutput.WinCount.DocCount+aggsOutput.LoseCount.DocCount)
//
//	// 证据建议、常用法条、法官意见
//	// 根据匹配查询后，展示时间最近的10条数据
//	outType := &EsDataInstrument{}
//	hits, err := GetEsHandler().GetQueryHits(outType, searchResult)
//	if err != nil {
//		return queryRes
//	}
//	tempEvidences := []string{}
//	tempInuseLaws := []InuseLawInfo{}
//	//tempInuseLaws := []string{}
//	tempJudgeArguments := []JudgeArgument{}
//	count := 0
//
//	for _, val := range hits.([]*EsDataInstrument) {
//		for _, caseSummary := range val.CaseSummary {
//			if len(caseSummary.JudgeArgument) == 0 {
//				continue
//			}
//
//			tempEvidences = append(tempEvidences, caseSummary.Evidence...)
//			//tempInuseLaws = append(tempInuseLaws, caseSummary.InuseLaw...)
//			for _, lawInfo := range caseSummary.InuseLaw {
//				//tempInuseLaws = append(tempInuseLaws, lawInfo)
//				tempInuseLaws = append(tempInuseLaws, InuseLawInfo{
//					Data:    lawInfo,
//					UseRate: 0,
//				})
//			}
//			for _, judgeContent := range caseSummary.JudgeArgument {
//				// 仅展示法官意见为胜诉的案件
//				//if caseSummary.IsWin != constant.WinString {
//				//	continue
//				//}
//
//				temp := JudgeArgument{
//					Data:   judgeContent,
//					TextId: val.InstrumentId,
//				}
//				tempJudgeArguments = append(tempJudgeArguments, temp)
//			}
//
//			count = count + 1
//			if count >= 10 {
//				break
//			}
//		}
//		if count >= 10 {
//			break
//		}
//	}
//
//	queryRes = &QueryRes{
//		WinRate:       winRate,
//		Evidence:      tempEvidences,
//		InuseLaw:      tempInuseLaws,
//		JudgeArgument: tempJudgeArguments,
//	}
//	fmt.Println("res:", *queryRes)
//	return queryRes
//}
