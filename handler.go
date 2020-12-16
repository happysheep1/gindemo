package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"test.com/studygo/d15/gindemo/ginsession"
)

type UserInfo struct {
	Username string `form:"username"`
	Password string `form:"password"`
}

func loginHandler(c *gin.Context) {
	if c.Request.Method == "POST" {
		toPath := c.DefaultQuery("next", "/index")
		var u UserInfo
		err := c.Bind(&u)
		if err != nil {
			c.HTML(200, "login.html", gin.H{"err": "用户名密码不能为空"})
			return
		}

		if u.Username == "root" && u.Password == "root" {
			//登录成功，保存islogin=true
			//1.先从上下文中取sessiondata
			tmpsd, ok := c.Get(ginsession.SessionContextName)
			if !ok {
				panic("取sessiondata出错")
			}
			sd := tmpsd.(ginsession.SessionData)
			sd.Set("islogin", true)
			fmt.Printf("%#v\n", sd)
			sd.Save()
			c.Redirect(http.StatusMovedPermanently, toPath)
		} else {
			c.HTML(http.StatusOK, "login.html", gin.H{"err": "用户名密码错误"})
		}
		return
	} else {
		c.HTML(200, "login.html", nil)
	}
}
func indexHandler(c *gin.Context) {
	c.HTML(200, "index.html", nil)
}
func homeHandler(c *gin.Context) {
	c.HTML(200, "home.html", nil)
}
func vipHandler(c *gin.Context) {
	c.HTML(200, "vip.html", nil)
}

//判断是否登录的中间件
func IsLoginMiddleware(c *gin.Context) {
	//根据islogin的值判断是否登录
	tmpsd, _ := c.Get(ginsession.SessionContextName)
	sd := tmpsd.(ginsession.SessionData) //断言
	val, err := sd.Get("islogin")
	if err != nil {
		//	取不到就是没有登录
		log.Println("取不到", err)
		c.Redirect(http.StatusFound, "/login")
		return
	}
	islogin, ok := val.(bool)
	if !ok {
		log.Println("不ok")
		c.Redirect(http.StatusFound, "/login")
		return
	}
	if !islogin {
		log.Println("!islogin")
		c.Redirect(http.StatusFound, "/login")
		return
	}
	c.Next()

}
