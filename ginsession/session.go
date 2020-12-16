package ginsession

//SessionData表示一条具体的信息
type SessionData interface {
	Get(key string) (value interface{}, err error)
	Set(key string, value interface{})
	Del(key string)
	Save()
	GenId() string //用于生成uuid作为sessionid
	SetExpire(expired int)
}
