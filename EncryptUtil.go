package util

import (
	"crypto/md5"
	"encoding/hex"
	"github.com/gopherjs/gopherjs/nosync"
)

var (
	cryptoUtilInstance *Crypto
	cryptoUtilOnce     nosync.Once
)

func CryptoUtil() *Crypto {
	cryptoUtilOnce.Do(func() {
		cryptoUtilInstance = new(Crypto)
	})
	return cryptoUtilInstance
}

// 加密工具
type Crypto struct {
}

// MD5 加密
// @param string data 加密数据
func (u Crypto) Md5(data string) string {
	md5Hash := md5.New()
	md5Hash.Write([]byte(data))
	return hex.EncodeToString(md5Hash.Sum(nil))
}
