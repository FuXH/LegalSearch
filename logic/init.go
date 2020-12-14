// 负责初始化es连接，并创建索引表结构
package logic

import (
	"LegalSearch/conf"
	"LegalSearch/constant"
	"LegalSearch/database/elasticsearch"
	"fmt"
)

var (
	EsHandler *elasticsearch.EsHandler
)

func GetEsHandler() *elasticsearch.EsHandler {
	return EsHandler
}

func InitEsHandler() error {
	config := conf.GetConfig()

	temp, err := elasticsearch.NewEsHandler(config.EsConfig)
	if err != nil {
		return err
	}
	EsHandler = temp

	if !EsHandler.IsExistIndex(constant.IndexNameInstrument) {
		if err := EsHandler.CreateIndex(constant.IndexNameInstrument, constant.IndexMappingsInstrument); err != nil {
			return fmt.Errorf("创建instrument索引表失败, err:", err)
		}
	}

	return nil
}
