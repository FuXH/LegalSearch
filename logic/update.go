package logic

import (
	"LegalSearch/constant"
	"LegalSearch/database/elasticsearch"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// Update接口输入参数
type UpdateReq struct {
	// 导入数据的文件夹路径
	Path string `form:"path"`
}
type UpdateRsp struct {
	InstrumentNumber  int     `form:"instrument_number"`
	ControversyNumber int     `form:"controversy_number"`
	CostTime          float64 `form:"cost_time"`
	Status            string  `form:"status"`
}

func Update(ctx *gin.Context) {
	status := "fail"
	instrumentNumber := 0
	controversyNumber := 0
	startTime := time.Now()
	defer func() {
		if status != "success" {
			rsp := &UpdateRsp{
				CostTime:          0,
				InstrumentNumber:  0,
				ControversyNumber: 0,
				Status:            status,
			}
			GinReturn(ctx, rsp)
		}
	}()

	param := new(UpdateReq)
	if err := ctx.ShouldBind(param); err != nil {
		fmt.Println("Update接口输入参数解析失败, err:", err)
		return
	}
	fmt.Println("update input:", param)
	if param.Path == "" {
		return
	}

	// 读取指定目录下的json数据
	tempBulkData := []elasticsearch.BulkData{}
	tempIndexs := make(map[string]int, 0)
	count := 0
	filepath.Walk(param.Path, func(path string, info os.FileInfo, err error) error {
		// 只读取json文件
		if !strings.HasSuffix(path, ".json") {
			return nil
		}

		// 文件内容tempLegal
		fileContent, err := ioutil.ReadFile(path)
		if err != nil {
			fmt.Println("读取json文件失败, file:", path, "err:", err)
			return nil
		}
		tempLegal := LegalDoc{}
		if err := json.Unmarshal(fileContent, &tempLegal); err != nil {
			fmt.Println("json文件格式错误1, file:", path, "err:", err)
			return nil
		}
		if tempLegal.TrialTime == "" || tempLegal.TrialYear == "" || tempLegal.TrialYear == "Not found" {
			fmt.Println("没有找到'具体审理时间'字段值，该文件不存入es, file:", path)
			return nil
		}

		// 拆分为单个争议焦点的EsData
		tempJsonData := make([]elasticsearch.BulkData, 0)
		for _, val := range legalDocToEsData(&tempLegal) {
			controversyNumber = controversyNumber + 1
			index := fmt.Sprintf("trial_year_%s", val.TrialYear)
			tempIndexs[index] = 1
			tempJsonData = append(tempJsonData, elasticsearch.BulkData{
				Data:  val,
				Index: index,
			})
		}
		tempBulkData = append(tempBulkData, tempJsonData...)

		// 写入原文索引表
		instrumentNumber = instrumentNumber + 1
		tempInstrument := EsDataInstrument{}
		if err := json.Unmarshal(fileContent, &tempInstrument); err != nil {
			fmt.Println("json文件格式错误2, file:", path, "err:", err)
			return nil
		}
		tempBulkData = append(tempBulkData, elasticsearch.BulkData{
			Data:  tempInstrument,
			Index: constant.IndexNameInstrument,
			Id:    tempInstrument.InstrumentId,
		})

		// 分批次写入，避免内存占用太多panic
		count = count + 1
		if (count % 5000) == 0 {
			// 遍历所有年份索引表，创建之前不存在的年份表
			indexs := []string{}
			for key, _ := range tempIndexs {
				indexs = append(indexs, key)
			}
			if err := createAllIndex(indexs); err != nil {
				errMsg := fmt.Errorf("创建索引表失败, indexs:", indexs, "err:", err)
				return errMsg
			}

			// 批量写入es
			if err := GetEsHandler().BulkInsert(tempBulkData); err != nil {
				errMsg := fmt.Errorf("批量写入es错误, err:", err)
				return errMsg
			}

			tempBulkData = []elasticsearch.BulkData{}
			tempIndexs = make(map[string]int, 0)
		}

		return nil
	})

	// 遍历所有年份索引表，创建之前不存在的年份表
	indexs := []string{}
	for key, _ := range tempIndexs {
		indexs = append(indexs, key)
	}
	if err := createAllIndex(indexs); err != nil {
		fmt.Println("创建索引表失败, indexs:", indexs, "err:", err)
		return
	}

	// 批量写入es
	if err := GetEsHandler().BulkInsert(tempBulkData); err != nil {
		fmt.Println("批量写入es错误, err:", err)
		return
	}

	status = "success"
	rsp := &UpdateRsp{
		CostTime:          time.Since(startTime).Seconds(),
		InstrumentNumber:  instrumentNumber,
		ControversyNumber: controversyNumber,
		Status:            status,
	}
	GinReturn(ctx, rsp)
}

// legalDocToEsData
func legalDocToEsData(legalDoc *LegalDoc) []EsDataControversy {
	res := make([]EsDataControversy, 0)

	defendants := []string{}
	for _, val := range legalDoc.DefendantInfo {
		defendants = append(defendants, val.Defendant)
	}
	plaintiffs := []string{}
	for _, val := range legalDoc.PlaintiffInfo {
		plaintiffs = append(plaintiffs, val.Plaintiff)
	}
	for _, val := range legalDoc.Summarys {
		temp := EsDataControversy{
			WenshuId:   legalDoc.InstrumentId, // 原文id
			Defendants: defendants,            // 被告
			Plaintiffs: plaintiffs,            // 原告
			TrialJudge: legalDoc.TrialJudge,   // 法官
			TrialCourt: legalDoc.TrialCourt,   // 审理法院
			TrialYear:  legalDoc.TrialYear,    // 审理年份
			TrialArea:  legalDoc.TrialArea,    // 审理地区
			TrialTime:  legalDoc.TrialTime,    // 审理时间

			DisputeFocus:  val.DisputeFocus,  // 争议焦点
			IsWin:         val.IsWin,         // 是否胜诉
			InuseLaw:      val.InuseLaw,      // 法条
			JudgeArgument: val.JudgeArgument, // 法官观点
			Evidence:      val.Evidence,      // 证据建议
		}
		res = append(res, temp)
	}
	return res
}

// createAllIndex
func createAllIndex(indexs []string) error {
	for _, index := range indexs {
		if GetEsHandler().IsExistIndex([]string{index}) {
			continue
		}

		if err := GetEsHandler().CreateIndex(index, constant.IndexMappingsControversy); err != nil {
			return err
		}
	}
	return nil
}
