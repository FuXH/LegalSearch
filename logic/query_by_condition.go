package logic

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"LegalSearch/constant"

	"github.com/gin-gonic/gin"
	es "github.com/olivere/elastic/v7"
)

// 高级搜索输入参数
type QueryByConditionReq struct {
	Defendant    string `form:"defendant"`
	Plaintiff    string `form:"plaintiff"`
	TrialJudge   string `form:"trialJudge"`
	TrialYear    string `form:"trialYear"`
	TrialArea    string `form:"trialArea"`
	TrialCourt   string `form:"trialCourt"`
	DisputeFocus string `form:"disputeFocus"`
}

// 搜索的返回值
type ApiQueryRsp struct {
	Success   bool      `form:"success"`
	ErrorMsg  string    `form:"errorMsg"`
	ErrorCode int       `form:"errorCode"`
	Content   *QueryRes `form:"content"`
}

// 搜索的数据查询结果
type QueryRes struct {
	WinRate       float64         `form:"winRate"`
	Evidence      []string        `form:"evidence"`
	InuseLaw      []string        `form:"inuseLaw"`
	JudgeArgument []JudgeArgument `form:"judgeArgument"`
}

// 法官意见
type JudgeArgument struct {
	Data   string `form:"data"`
	TextId string `form:"textId"`
}

// 高级搜索
func QueryByCondition(ctx *gin.Context) {
	startTime := time.Now()
	status := "success"
	defer func() {
		costTime := time.Since(startTime).Seconds()
		fmt.Println("query_by_condition接口耗时(s):", costTime)
		if status != "success" {
			rsp := &ApiQueryRsp{
				Success:   false,
				ErrorMsg:  status,
				ErrorCode: 1,
			}
			GinReturn(ctx, rsp)
		}
	}()

	param := new(QueryByConditionReq)
	if err := ctx.ShouldBind(param); err != nil {
		status = fmt.Sprintf("QueryByCondition接口输入参数解析失败, err:", err)
		fmt.Println(status)
		return
	}
	fmt.Println("query input:", param)

	// 构造查询
	indexName, err := getIndex(param.TrialYear)
	if err != nil {
		status := fmt.Sprintf("QueryByCondition输入参数错误，审理年份非法:", param.TrialYear, "err:", err)
		fmt.Println(status)
		return
	}
	filters := getFilters(param)
	aggsMap := make(map[string]es.Aggregation)
	aggsMap["win_count"] = es.NewFilterAggregation().Filter(es.NewTermQuery("IsWin", "1.0"))
	aggsMap["lose_count"] = es.NewFilterAggregation().Filter(es.NewTermQuery("IsWin", "2.0"))

	// 查询es
	searchResult, err := GetEsHandler().BoolQuery(indexName,
		constant.SortField, constant.SortOrder, constant.SortSize,
		aggsMap,
		filters...)
	if err != nil {
		status := fmt.Sprintf("QueryByCondition查询返回错误, err:", err)
		fmt.Println(status)
		return
	}

	// 处理返回数据
	queryRes := aggregateData(searchResult)

	// rsp
	rsp := &ApiQueryRsp{
		Success:   true,
		ErrorMsg:  "",
		ErrorCode: 0,
		Content:   queryRes,
	}
	GinReturn(ctx, rsp)
}

// 根据年份获取对应索引表
func getIndex(trailYear string) (string, error) {
	// 未指定年份，则搜索所有索引表
	if trailYear == "" {
		return "*", nil
	}

	// 输入年份不为数字
	year, err := strconv.Atoi(trailYear)
	if err != nil {
		return "", err
	}
	indexName := fmt.Sprintf("trial_year_%d", year)

	// 索引不存在
	if exist := GetEsHandler().IsExistIndex(indexName); !exist {
		return "", fmt.Errorf("索引不存在, index:", indexName)
	}

	fmt.Println("index_name:", indexName)
	return indexName, nil
}

// 创建筛选
func getFilters(param *QueryByConditionReq) []es.Query {
	query := make([]es.Query, 0)

	if param.Defendant != "" {
		//被告
		query = append(query,
			es.NewMatchQuery("Defendants", param.Defendant))
	}
	if param.Plaintiff != "" {
		// 原告
		query = append(query,
			es.NewMatchQuery("Plaintiffs", param.Plaintiff))
	}
	if param.TrialJudge != "" {
		//审理法官
		query = append(query,
			es.NewMatchQuery("TrialJudge", param.TrialJudge))
	}
	if param.TrialArea != "" {
		// 审理地区
		query = append(query,
			es.NewMatchQuery("TrialArea", param.TrialArea))
	}
	if param.TrialCourt != "" {
		// 审理法院，精确值
		query = append(query,
			es.NewMatchQuery("TrialCourt", param.TrialCourt))
	}
	if param.DisputeFocus != "" {
		// 争议焦点，模糊匹配
		query = append(query,
			es.NewMatchQuery("DisputeFocus", param.DisputeFocus))
	}
	fmt.Println("filters:", query)

	return query
}

// 处理查询数据
func aggregateData(searchResult *es.SearchResult) *QueryRes {
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
	outType := &EsDataControversy{}
	hits, err := GetEsHandler().GetQueryHits(outType, searchResult)
	if err != nil {
		return queryRes
	}
	tempEvidences := []string{}
	tempInuseLaws := []string{}
	tempJudgeArguments := []JudgeArgument{}
	fmt.Println("hits:", hits, hits.([]*EsDataControversy))
	for _, val := range hits.([]*EsDataControversy) {
		fmt.Println("es query_controversy data:", *val)
		tempEvidences = append(tempEvidences, val.Evidence...)
		tempInuseLaws = append(tempInuseLaws, val.InuseLaw...)
		tempJudgeArguments = append(tempJudgeArguments,
			JudgeArgument{
				//Data:   val.JudgeArgument,
				TextId: val.InstrumentId,
			})
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

// 数组去重
func RemoveRepeat(src []string) []string {
	res := make([]string, 0)
	tempMap := make(map[string]int, 0)

	for _, val := range src {
		mapLen := len(tempMap)
		tempMap[val] = 1
		if len(tempMap) != mapLen {
			res = append(res, val)
		}
	}

	return res
}

// 框架返回rsp参数
func GinReturn(ctx *gin.Context, rsp interface{}) {
	tempByte, _ := json.Marshal(rsp)
	rspMap := gin.H{}
	json.Unmarshal(tempByte, &rspMap)

	ctx.JSON(200, rspMap)
}
