package ginsession

import (
	"fmt"
	uuid "github.com/satori/go.uuid"
	"sync"
)

//SessionMgr 全局session管理
type SessionMemoryMgr struct {
	SessionMap map[string]SessionData
	rwlock     sync.RWMutex
}

//根据sessionid找到对应的数据
func (s *SessionMemoryMgr) GetSessionData(sessionId string) (sd SessionData, err error) {
	s.rwlock.RLock()
	defer s.rwlock.RUnlock()
	sd, ok := s.SessionMap[sessionId]
	if !ok {
		err = fmt.Errorf("错误的sessionid")
		return
	}
	return
}

//初始化内存版管理者
//*SessionMemoryMgr会出错，直接返回接口就行
func NewSessionMemoryMgr() SessionMgr {
	return &SessionMemoryMgr{
		SessionMap: make(map[string]SessionData, 1024),
	}
}
func (s *SessionMemoryMgr) Init(addr string, args ...string) {

}

//创建一个session记录
func (s *SessionMemoryMgr) CreateSessionData() (sd SessionData) {
	//造一个sessionid
	//ts := time.Now().UnixNano()
	id := uuid.NewV4()
	sd = NewSessionMemoryData(id.String())

	s.SessionMap[sd.GenId()] = sd //创建之后需要把这个sessionid保存起来，不然下次还是新创建一个mgr，
	//	造一个sessiondata
	//	返回sessiondata
	return
}
