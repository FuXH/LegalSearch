package logic

import (
	"LegalSearch/constant"
	"fmt"
	"github.com/gin-gonic/gin"
	"time"
)

// GetFullTextReq
type GetFullTextReq struct {
	TextId string `form:"textId"`
}

// GetFullTextRsp
type GetFullTextRsp struct {
	Text string `form:"text"`
}

// GetFullText
func GetFullText(ctx *gin.Context) {
	startTime := time.Now()
	status := "fail"
	defer func() {
		costTime := time.Since(startTime)
		fmt.Println("get_full_text接口耗时(s):", costTime)
		if status != "success" {
			rsp := &GetFullTextRsp{
				Text: "",
			}
			GinReturn(ctx, rsp)
		}
	}()

	param := new(GetFullTextReq)
	if err := ctx.ShouldBind(param); err != nil {
		fmt.Println("get_full_text接口输入参数解析失败, err:", err)
		return
	}
	fmt.Println("get_full_text input:", param)
	if param.TextId == "" {
		fmt.Println("get_full_text输入参数为空")
		return
	}

	data := &struct {
		Text string `json:"wenshu_content"`
	}{}
	if err := GetEsHandler().QueryById(constant.IndexNameInstrument, param.TextId, data); err != nil {
		return
	}

	// rsp
	status = "success"
	rsp := &GetFullTextRsp{
		Text: data.Text,
	}
	GinReturn(ctx, rsp)
}
