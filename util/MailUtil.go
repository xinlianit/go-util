package util

import (
	"fmt"
	"gopkg.in/gomail.v2"
	"sync"
)

var (
	mailUtilInstance *Mail
	mailUtilOnce     sync.Once
)

func MailUtil() *Mail {
	mailUtilOnce.Do(func() {
		mailUtilInstance = new(Mail)
	})
	return mailUtilInstance
}

// 邮件服务配置
type MailServerConfig struct {
	// 主机
	Host string
	// 端口
	Port int
	// 用户
	User string
	// 认证; 密码或授权码: 如 阿里云邮箱(密码), QQ邮箱(授权码)...
	Auth string
}

// 邮件工具
type Mail struct {
	serverConfig *MailServerConfig
}

// 配置邮件服务器
func (u Mail) Config(config *MailServerConfig) Mail {
	u.serverConfig = config
	return u
}

// 邮件发送
// @param to 收件人邮箱
// @param subject 邮件主题
// @param body 邮件内容
func (u Mail) Send(to string, subject string, body string) (bool, error) {
	toList := []string{to}
	if err := u.sendMail(toList, subject, body); err != nil {
		return false, err
	}

	return true, nil
}

// 批量邮件发送
// @param toList 收件人邮箱列表
// @param subject 邮件主题
// @param body 邮件内容
func (u Mail) SendMulti(toList []string, subject string, body string) (bool, error) {
	if err := u.sendMail(toList, subject, body); err != nil {
		return false, err
	}

	return true, nil
}

// 发送邮件
// @param toList 收件人邮箱列表
// @param subject 邮件主题
// @param body 邮件内容
func (u Mail) sendMail(toList []string, subject string, body string) error {
	mail := gomail.NewMessage()

	//设置发件人
	mail.SetHeader("From", u.serverConfig.User)

	//设置发送给多个用户
	mail.SetHeader("To", toList...)

	//设置邮件主题
	mail.SetHeader("Subject", subject)

	//设置邮件正文
	mail.SetBody("text/html", body)

	d := gomail.NewDialer(u.serverConfig.Host, u.serverConfig.Port, u.serverConfig.User, u.serverConfig.Auth)

	err := d.DialAndSend(mail)
	if err != nil {
		fmt.Println(err)
	}
	return err
}
