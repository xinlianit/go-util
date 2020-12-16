package util

import (
	jsoniter "github.com/json-iterator/go"
	"sync"
)

var (
	serializeUtilInstance *Serialize
	serializeUtilOnce     sync.Once
	json                  = jsoniter.ConfigCompatibleWithStandardLibrary
)

func SerializeUtil() *Serialize {
	serializeUtilOnce.Do(func() {
		serializeUtilInstance = new(Serialize)
	})

	return serializeUtilInstance
}

type Serialize struct {
}

// Json 编码
func (u Serialize) JsonEncode(v interface{}) (string, error) {
	data, err := json.Marshal(v)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

// Json 解码
func (u Serialize) JsonDecode(data string, v interface{}) error {
	return json.Unmarshal([]byte(data), v)
}
