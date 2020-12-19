package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"strconv"

	"LegalSearch/conf"
	_ "LegalSearch/database/elasticsearch"
	"LegalSearch/logic"

	"github.com/gin-gonic/gin"
)

var (
	ConfPath = flag.String("conf", "conf.yaml", "yaml配置文件的相对路径")
)

func LoginPage(ctx *gin.Context) {
	ctx.HTML(http.StatusOK, "index.html", gin.H{})
}

func HomePage(ctx *gin.Context) {
	ctx.HTML(http.StatusOK, "index.html", gin.H{})
}

func ContentPage(ctx *gin.Context) {
	ctx.HTML(http.StatusOK, "content.html", gin.H{})
}

func main() {
	// init
	if err := InitConnect(); err != nil {
		return
	}

	// 启动监听服务
	config := conf.GetConfig()
	r := gin.Default()

	// 加载页面内容
	r.LoadHTMLGlob("./html/template/*")
	r.StaticFS("/assets", http.Dir("./html/assets"))
	r.GET("/login", LoginPage)
	r.GET("/home", HomePage)
	r.GET("/content", ContentPage)

	// api接口
	r.GET("/api/queryByCondition", logic.QueryByCondition)
	r.GET("/api/queryCount", logic.SuperQueryCount)
	r.GET("/api/query", logic.Query)
	r.GET("/api/update", logic.Update)
	r.GET("/api/getFullText", logic.GetFullText)
	r.GET("/api/operation", logic.ElasticOperation)

	//r.Use(func(ctx *gin.Context) {
	//	ctx.Header("Access-Control-Allow-Origin", "*")
	//	ctx.Next()
	//})

	//r.Run(":" + strconv.Itoa(config.Server.Port))
	server := &http.Server{Handler: r}
	l, err := net.Listen("tcp4", ":"+strconv.Itoa(config.Server.Port))
	if err != nil {
		log.Fatal(err)
	}
	err = server.Serve(l)

	return
}

func InitConnect() error {
	flag.Parse()
	// 读取配置文件
	if err := conf.InitYmlFile(*ConfPath); err != nil {
		fmt.Println("初始化yaml配置文件失败, err:", err)
		return err
	}
	config := conf.GetConfig()
	fmt.Println("config:", config)

	// 初始化es
	if err := logic.InitEsHandler(); err != nil {
		fmt.Println("初始化es连接失败, err:", err)
		return err
	}

	return nil
}
