package util

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/hex"
	"encoding/pem"
	"errors"
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

// aesPKCS7Padding 补码
func (u Crypto) aesPKCS7Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padText := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padText...)
}

// aesPKCS7UnPadding 去码
func (u Crypto) aesPKCS7UnPadding(origData []byte) []byte {
	length := len(origData)
	unPadding := int(origData[length-1])
	return origData[:(length - unPadding)]
}

// AesEncrypt AES加密
// @param data 报文数据
// @param key 加密秘钥(长度必须是16、24或者32字节，分别为AES-128, AES-192, or AES-256模式)
// @return encrypted 加密报文
// @return error
func (u Crypto) AesEncrypt(data, key string) (encrypted string, err error) {
	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		return encrypted, err
	}
	blockSize := block.BlockSize()
	dataByte := u.aesPKCS7Padding([]byte(data), blockSize)
	blockMode := cipher.NewCBCEncrypter(block, []byte(key)[:blockSize])
	encryptedByte := make([]byte, len(dataByte))
	blockMode.CryptBlocks(encryptedByte, dataByte)
	return base64.StdEncoding.EncodeToString(encryptedByte), nil
}

// AesDecrypt AES解密
// @param encrypted 加密报文
// @param key 解密秘钥(长度必须是16、24或者32字节，分别为AES-128, AES-192, or AES-256模式)
// @return decrypted 报文数据
// @return error
func (u Crypto) AesDecrypt(encrypted, key string) (decrypted string, err error) {
	encryptedByte, err := base64.StdEncoding.DecodeString(encrypted)
	if err != nil {
		return decrypted, err
	}
	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		return decrypted, err
	}
	blockSize := block.BlockSize()
	blockMode := cipher.NewCBCDecrypter(block, []byte(key)[:blockSize])
	decryptedByte := make([]byte, len(encryptedByte))
	blockMode.CryptBlocks(decryptedByte, encryptedByte)
	decryptedByte = u.aesPKCS7UnPadding(decryptedByte)
	return string(decryptedByte), nil
}

// RsaEncrypt RSA加密
// @param data 报文数据
// @param publicKey 公钥
// @return encrypted 加密报文
// @return error
func (u Crypto) RsaEncrypt(data, publicKey string) (encrypted string, err error) {
	// 解密pem格式公钥
	block, _ := pem.Decode([]byte(publicKey))
	if block == nil {
		return encrypted, errors.New("public key error")
	}
	// 解析公钥
	pubKey, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return encrypted, err
	}
	// 类型断言
	pub := pubKey.(*rsa.PublicKey)
	//加密
	encryptedByte, err := rsa.EncryptPKCS1v15(rand.Reader, pub, []byte(data))
	if err != nil {
		return encrypted, err
	}

	return base64.StdEncoding.EncodeToString(encryptedByte), nil
}

// RsaDecrypt RSA解密
// @param encrypted 加密报文
// @param privateKey 私钥
// @return decrypted 报文数据
// @return error
func (u Crypto) RsaDecrypt(encrypted, privateKey string) (decrypted string, err error) {
	// base64解码
	encryptedByte, err := base64.StdEncoding.DecodeString(encrypted)
	if err != nil {
		return decrypted, err
	}
	// 解密pem格式私钥
	block, _ := pem.Decode([]byte(privateKey))
	if block == nil {
		return decrypted, errors.New("private key error")
	}
	//解析PKCS1格式的私钥
	priKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return decrypted, err
	}
	// 解密
	decryptedByte, err := rsa.DecryptPKCS1v15(rand.Reader, priKey, encryptedByte)
	if err != nil {
		return decrypted, err
	}
	return string(decryptedByte), nil
}
