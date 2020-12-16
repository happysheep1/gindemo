package ginsession

import (
	"errors"
	"sync"
)

//内存版session
//SessionData表示一条具体的信息
type SessionMemoryData struct {
	ID     string
	Data   map[string]interface{}
	rwlock sync.RWMutex //读写锁,锁的是上面的Data
	//	过期时间
}

//创建一个构造函数
func NewSessionMemoryData(id string) SessionData {
	return &SessionMemoryData{
		ID:   id,
		Data: make(map[string]interface{}, 16),
	}
}

func (s *SessionMemoryData) Get(key string) (value interface{}, err error) {
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
func (s *SessionMemoryData) Set(key string, value interface{}) {
	//获取写锁
	s.rwlock.Lock()
	defer s.rwlock.Unlock()
	s.Data[key] = value

}
func (s *SessionMemoryData) Del(key string) {
	//获取写锁
	s.rwlock.Lock()
	defer s.rwlock.Unlock()
	//go doc builtin.delete
	delete(s.Data, key)
}
func (s *SessionMemoryData) Save() {

}
func (s *SessionMemoryData) GenId() string {
	return s.ID
}

//设置过期时间
func (s *SessionMemoryData) SetExpire(expired int) {

}
