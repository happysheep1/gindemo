package ginsession

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-redis/redis"
	"sync"
	"time"
)

type SessionRedisData struct {
	ID     string //这个技术sessionid
	Data   map[string]interface{}
	rwlock sync.RWMutex //读写锁,锁的是上面的Data
	//	过期时间
	Expire   int
	Client   *redis.Client
	IsModify bool
}

func NewSessionRedisData(id string, client *redis.Client) SessionData {
	return &SessionRedisData{
		ID:     id,
		Data:   make(map[string]interface{}, 16),
		Expire: 30,
		Client: client,
	}
}

func (s *SessionRedisData) Get(key string) (value interface{}, err error) {
	//获取读锁
	s.rwlock.RLock()
	defer s.rwlock.RUnlock()
	value, ok := s.Data[key]
	if !ok {

		err = errors.New("key不存在")

		return
	}
	return
}
func (s *SessionRedisData) Set(key string, value interface{}) {
	//获取写锁
	s.rwlock.Lock()
	defer s.rwlock.Unlock()
	s.Data[key] = value
	s.IsModify = true

}
func (s *SessionRedisData) Del(key string) {
	//获取写锁
	s.rwlock.Lock()
	defer s.rwlock.Unlock()
	//go doc builtin.delete
	delete(s.Data, key)
	s.IsModify = true
}
func (s *SessionRedisData) Save() {
	if !s.IsModify {
		return
	}
	val, err := json.Marshal(s.Data)
	if err != nil {
		fmt.Println("序列化失败")
		return
	}
	s.Client.Set(s.ID, val, time.Second*time.Duration(s.Expire))
	s.IsModify = false
}

func (s *SessionRedisData) GenId() string {
	return s.ID
}

//设置过期时间
func (s *SessionRedisData) SetExpire(expired int) {
	s.Expire = expired
}
