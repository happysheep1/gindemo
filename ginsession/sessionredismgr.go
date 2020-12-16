package ginsession

import (
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis"
	uuid "github.com/satori/go.uuid"
	"strconv"
	"sync"
)

type SessionRedisMgr struct {
	client     *redis.Client //redis连接池
	SessionMap map[string]SessionData
	rwlock     sync.RWMutex
}

//初始化redis版管理者
func NewSessionRedisMgr() SessionMgr {
	return &SessionRedisMgr{
		SessionMap: make(map[string]SessionData, 1024),
	}
}
func (s *SessionRedisMgr) Init(addr string, args ...string) {
	var (
		password string
		db       string
	)
	if len(args) == 1 {
		password = args[0]
	} else if len(args) == 2 {
		password = args[0]
		//db = strconv.ParseInt(args[1])
		db = args[1]
	}
	dbVal, err := strconv.Atoi(db)
	if err != nil {
		dbVal = 0
	}
	s.client = redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password, // no password set
		DB:       dbVal,    // use default DB
	})
	_, err = s.client.Ping().Result()
	if err != nil {
		panic(err)
	}

}
func (s *SessionRedisMgr) loadFromRedis(SessionId string) (err error) {
	//	1.连接redis
	val, err := s.client.Get(SessionId).Result()
	if err != nil {
		//	redis中没有这个sessiondata
		fmt.Println("redis中没有这个sessiondata")
		return
	}
	err = json.Unmarshal([]byte(val), &s.SessionMap)
	if err != nil {
		fmt.Println("从redis取出的数据反序列化失败")
		return
	}
	//	2.根据sessionid找到对应的数据
	//	3.吧数据取出反序列化到s.data
	return
}
func (s *SessionRedisMgr) GetSessionData(sessionId string) (sd SessionData, err error) {
	//1.s.session中已经从redis取出数据
	//	2.s.session【sessionid】拿到sessiondata
	if s.SessionMap == nil {
		err := s.loadFromRedis(sessionId)
		if err != nil {
			return nil, err
		}
	}
	s.rwlock.RLock()
	defer s.rwlock.RUnlock()
	sd, ok := s.SessionMap[sessionId]
	if !ok {
		err = fmt.Errorf("错误的sessionid")
		return
	}
	return

}

func (s *SessionRedisMgr) CreateSessionData() (sd SessionData) {
	id := uuid.NewV4()
	sd = NewSessionRedisData(id.String(), s.client)
	s.SessionMap[sd.GenId()] = sd
	return
}
