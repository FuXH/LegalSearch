package logic

// es存储的结构体，争议焦点索引表
type EsDataControversy struct {
	WenshuId      string   `json:"WenshuId"`      // 法律文书的id
	Defendants    []string `json:"Defendants"`    // 被告
	Plaintiffs    []string `json:"Plaintiffs"`    // 原告
	TrialJudge    string   `json:"TrialJudge"`    // 审理法官
	TrialYear     string   `json:"TrialYear"`     // 审理年份
	TrialTime     string   `json:"TrialTime"`     // 审理具体时间
	TrialArea     string   `json:"TrialArea"`     // 审理地区
	TrialCourt    string   `json:"TrialCourt"`    // 审理法院
	DisputeFocus  string   `json:"DisputeFocus"`  // 争议焦点
	IsWin         string   `json:"IsWin"`         // 是否胜诉, else-未知，'1.0'-胜诉，'2.0'-败诉
	InuseLaw      []string `json:"InuseLaw"`      // 常用法条
	JudgeArgument []string `json:"JudgeArgument"` // 法官观点
	Evidence      []string `json:"Evidence"`      // 证据建议
}

// es存储的结构体，法律文书的索引表
type EsDataInstrument struct {
	InstrumentId string    `json:"wenshu_id_2"`    // 法律文书的ID
	Content      string    `json:"wenshu_content"` // 原文
	CaseId       string    `json:"case_id"`        // 案件号
	Cause        string    `json:"cause"`          // 纠纷原因
	CaseType     string    `json:"case_type"`      // 案件类型
	CaseSummary  []Summary `json:"case_summary"`   // 概要
	// 费用相关
	FeeMedical     FeeInfo `json:"fee_medical"`
	FeeMess        FeeInfo `json:"fee_mess"`
	FeeNurse       FeeInfo `json:"fee_nurse"`
	FeeNutrition   FeeInfo `json:"fee_nutrition"`
	FeePostCure    FeeInfo `json:"fee_post_cure"`
	FeeLossWorking FeeInfo `json:"fee_loss_working"`
	FeeTrffic      FeeInfo `json:"fee_trffic"`
	// 结构体不明确，预留
	//FeeDisable
	//FeeDeath
	//FeeBury
	//FeeLife
	//FeeTrafficForProcessBury
	//FeeLossWorkingForProcessBury
	//FeeMind
	//FeeAppraise
}
type FeeInfo struct {
	Criterion  string `json:"criterion"`
	Days       string `json:"days"`
	SumOfMoney string `json:"sum_of_money"`
}
