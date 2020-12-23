package logic

// 法律文书的json数据
type LegalDoc struct {
	DefendantInfo []DefendantIndo `json:"defendant_info"` // 被告
	PlaintiffInfo []PlaintiffInfo `json:"plaintiff_info"` // 原告
	TrialJudge    string          `json:"judge"`          // 审理法官
	TrialYear     string          `json:"year"`           // 审理年份
	TrialTime     string          `json:"wenshu_time"`    // 审理具体时间
	TrialArea     string          `json:"area"`           // 审理地区
	TrialCourt    string          `json:"court"`          // 审理法院
	Summarys      []Summary       `json:"case_summary"`   // 文书简介
	InstrumentId  string          `json:"wenshu_id_2"`    // 法律文书ID
	Content       string          `json:"wenshu_content"` // 原文
	CaseId        string          `json:"case_id"`        // 案件id
	CaseType      string          `json:"case_type"`      // 案件类型
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
type DefendantIndo struct {
	Defendant      string `json:"defendant"`
	DefendantAgent string `json:"defendant_agent"`
	LawFirm        string `json:"law_firm"`
}
type PlaintiffInfo struct {
	Plaintiff      string `json:"plaintiff"`
	PlaintiffAgent string `json:"plaintiff_agent"`
	LawFirm        string `json:"law_firm"`
}
type Summary struct {
	DisputeFocus  string   `json:"controversy"` // 争议焦点
	Judgement     string   `json:"judgement"`   // 是否胜诉
	InuseLaw      []string `json:"basis"`       // 常用法条
	JudgeArgument []string `json:"cause"`       // 法官观点
	Evidence      []string `json:"evidence"`    // 证据建议
}
