package util

import (
	"crypto/md5"
	"encoding/hex"
	"github.com/gopherjs/gopherjs/nosync"
	"golang.org/x/crypto/bcrypt"
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

// Crypto 加密工具
type Crypto struct {
}

// Md5 MD5 加密
// @param string data 加密数据
func (u Crypto) Md5(data string) string {
	md5Hash := md5.New()
	md5Hash.Write([]byte(data))
	return hex.EncodeToString(md5Hash.Sum(nil))
}

// PasswordHash 密码加密
// @param password 密码
func (u Crypto) PasswordHash(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

// PasswordVerify 密码验证
// @param password 密码
// @param hash 加密串
func (u Crypto) PasswordVerify(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
