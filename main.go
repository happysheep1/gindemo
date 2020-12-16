package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"test.com/studygo/d15/gindemo/ginsession"
)

//测试session中间件
func main() {
	r := gin.Default()
	ginsession.InitSessionMgr("redis", "127.0.0.1:6379")
	option := &ginsession.Options{MaxAge: 3600, Path: "/", Domain: "127.0.0.1", Secure: false, HttpOnly: true}
	r.Use(ginsession.SessionMiddleware(ginsession.SMgr, option))
	r.LoadHTMLGlob("templates/*")
	r.Any("/login", loginHandler)
	r.GET("/index", indexHandler)
	r.GET("/home", homeHandler)
	r.GET("/vip", IsLoginMiddleware, vipHandler)
	//没有匹配的路由都走这个
	r.NoRoute(func(c *gin.Context) {
		c.HTML(http.StatusNotFound, "404.html", nil)
	})
	r.Run()

	//目前问题
	//1.save不管修改没都会保存，定义一个标志位
	//2.从redis中加载数据，并发查询redis ，sync.once
	//3.cookie设置写死了

}
