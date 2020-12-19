package constant

// 排序相关
const (
	// 排序的字段名称, 审理具体时间
	SortField = "TrialTime"
	// 排序方式，降序排序
	SortOrder = false
	// 排序后返回的查询数量
	SortSize = 500
)

// es索引设置
const (
	IndexMappingsControversy = `
{
    "settings": {
        "number_of_shards": 1,
        "number_of_replicas": 0
        // "index.blocks.read_only_allow_delete": null
    },
    "mappings": {
		"dynamic": "false",
		"properties": {
			"InstrumentId": {
				"type": "keyword"
			},
			"Defendants": {
				"type": "text",
				"analyzer": "ik_smart"
			},
			"Plaintiffs": {
				"type": "text",
				"analyzer": "ik_smart"
			},
			"TrialJudge": {
				"type": "text",
				"analyzer": "ik_smart"
			},
			"TrialTime": {
				"type": "date",
				"format": "yyyy-MM-dd"
			},
			"TrialArea": {
				"type": "keyword"
			},
			"TrialCourt": {
				"type": "text",
				"analyzer": "ik_smart"
			},
			"DisputeFocus": {
				"type": "text",
				"analyzer": "ik_max_word"
			},
			"IsWin": {
				"type": "keyword"
			},
			"InuseLaw": {
				"type": "keyword"
			},
			"JudgeArgument": {
				"type": "keyword"
			},
			"Evidence": {
				"type": "keyword"
			}
        }
    }
}
`
	IndexNameInstrument     = "instrument"
	IndexMappingsInstrument = `
{
    "settings": {
        "number_of_shards": 1,
        "number_of_replicas": 0
	},
	"mappings": {
		"dynamic": "false",
		"properties": {
			"InstrumentId": {
				"type": "keyword"
			},
			"wenshu_content": {
				"type": "text",
				"analyzer": "ik_max_word"
			},
			"CaseId": {
				"type": "keyword"
			},
			"Cause": {
				"type": "keyword"
			},
			"CaseType": {
				"type": "keyword"
			},
			"Summary": {
				"properties": {
					"DisputeFocus": {
						"type": "keyword"
					},
					"IsWin": {
						"type": "keyword"
					},
					"InuseLaw": {
						"type": "keyword"
					},
					"JudgeArgument": {
						"type": "keyword"
					},
					"Evidence": {
						"type": "keyword"
					}
				}
			},
			"FeeMedical": {
				"properties": {
					"Criterion": {
						"type": "keyword"
					},
					"Days": {
						"type": "keyword"
					},
					"SumOfMoney": {
						"type": "keyword"
					}
				}
			},
			"FeeMess": {
				"properties": {
					"Criterion": {
						"type": "keyword"
					},
					"Days": {
						"type": "keyword"
					},
					"SumOfMoney": {
						"type": "keyword"
					}
				}
			},
			"FeeNurse": {
				"properties": {
					"Criterion": {
						"type": "keyword"
					},
					"Days": {
						"type": "keyword"
					},
					"SumOfMoney": {
						"type": "keyword"
					}
				}
			},
			"FeeNutrition": {
				"properties": {
					"Criterion": {
						"type": "keyword"
					},
					"Days": {
						"type": "keyword"
					},
					"SumOfMoney": {
						"type": "keyword"
					}
				}
			},
			"FeePostCure": {
				"properties": {
					"Criterion": {
						"type": "keyword"
					},
					"Days": {
						"type": "keyword"
					},
					"SumOfMoney": {
						"type": "keyword"
					}
				}
			},
			"FeeLossWorking": {
				"properties": {
					"Criterion": {
						"type": "keyword"
					},
					"Days": {
						"type": "keyword"
					},
					"SumOfMoney": {
						"type": "keyword"
					}
				}
			},
			"FeeTrffic": {
				"properties": {
					"Criterion": {
						"type": "keyword"
					},
					"Days": {
						"type": "keyword"
					},
					"SumOfMoney": {
						"type": "keyword"
					}
				}
			}
		}
	}
}
`
)

const (
	WinString  = "1"
	LoseString = "0"
)
