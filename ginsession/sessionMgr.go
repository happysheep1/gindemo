package ginsession

import "github.com/gin-gonic/gin"

var (
	SMgr SessionMgr
	sd   SessionData
)

const (
	SessionCookieName  = "sessionId"
	SessionContextName = "session"
)

type Options struct {
	MaxAge   int
	Domain   string
	Path     string
	Secure   bool
	HttpOnly bool
}

//SessionMgr 全局session管理
type SessionMgr interface {
	Init(addr string, args ...string) //用来执行具体的连接
	GetSessionData(sessionId string) (sd SessionData, err error)
	CreateSessionData() (sd SessionData)
}

func InitSessionMgr(name, addr string, args ...string) {
	switch name {
	case "memory":
		SMgr = NewSessionMemoryMgr()
	case "redis":
		SMgr = NewSessionRedisMgr()
	}
	SMgr.Init(addr, args...)

}

//实现一个gin框架的中间件
func SessionMiddleware(s SessionMgr, o *Options) gin.HandlerFunc {
	//	1.从请求的cookie中获取sessionid
	//	1.1取不到，说明第一次，创建sessiondata和sessionid
	//	1.2取到sessionid
	//	2.根据sessionid取到data
	//3.思考如何实现后续的请求方法都能使用sessiondata
	//	利用c.set("session",sessiondata)
	//	c.next()
	if s == nil {
		panic("必须先初始化")
	}
	return func(c *gin.Context) {
		sessionId, err := c.Cookie(SessionCookieName)
		if err != nil {
			//	1.1取不到，说明第一次，创建sessiondata和sessionid
			sd = s.CreateSessionData()
			sessionId = sd.GenId()
		} else {
			sd, err = s.GetSessionData(sessionId)
			if err != nil {
				//	1.1取不到，说明第一次，创建sessiondata和sessionid
				sd = s.CreateSessionData()
				sessionId = sd.GenId()
				//更新用户cookie的id
			}
		}
		c.Set(SessionContextName, sd)
		//gin框架中药回写cookie，需要在处理请求返回之前
		//就是需要在next之前
		c.SetCookie(SessionCookieName, sessionId, o.MaxAge, o.Path, o.Domain, o.Secure, o.HttpOnly)
		c.Next()
	}
}
