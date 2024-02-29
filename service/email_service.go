package service

import (
	"net/smtp"

	"github.com/jordan-wright/email"
	"github.com/rs/zerolog/log"
)

// 短信通知
// 公众号通知
// 邮件通知

type EmailInfo struct {
	Account            string // 账户
	AuthCode           string // 授权码
	EmailType          string // 邮箱类型
	EmailServerAddress string // 邮件服务地址
}

func SendEmail(from string, toList []string, subject string, content string) {
	em := email.NewEmail()
	// 设置 sender 发送方 的邮箱 ， 此处可以填写自己的邮箱
	em.From = "xx <xxx@qq.com>"

	// 设置 receiver 接收方 的邮箱  此处也可以填写自己的邮箱， 就是自己发邮件给自己
	em.To = []string{"xxx@qq.com"}

	// 设置主题
	em.Subject = "小魔童给你发邮件了"

	// 简单设置文件发送的内容，暂时设置成纯文本
	em.Text = []byte("hello world, 咱们用 golang 发个邮件！！")

	//设置服务器相关的配置
	err := em.Send("smtp.qq.com:25", smtp.PlainAuth("", "自己的邮箱账号", "自己邮箱的授权码", "smtp.qq.com"))
	if err != nil {
		log.Error().Msgf("fail to send email, error is: %v", err)
	}
}
