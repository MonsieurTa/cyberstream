package tracker

type Decoder interface {
	GetString(key string) string
	GetList(key string) []interface{}
	GetInt(key string) int
	GetDict(key string) map[string]interface{}
}
