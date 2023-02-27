// 将pay项目中哈尔封装的邮件包提取到common中

package email

import (
	"context"
	"fmt"
	"gopkg.in/gomail.v2"
)

type Email struct {
	sender   string
	smtpHost string
	smtpPort int
	userName string
	password string
	client   *gomail.Dialer
}

func NewEmail(sender string, smtpHost string, smtpPort int, userName string, password string) *Email {
	return &Email{
		sender:   sender,
		smtpHost: smtpHost,
		smtpPort: smtpPort,
		userName: userName,
		password: password,
	}
}

func (e Email) newClient() *gomail.Dialer {
	return gomail.NewDialer(e.smtpHost, e.smtpPort, e.userName, e.password)
}

func (e Email) getClient() *gomail.Dialer {
	if e.client == nil {
		e.client = e.newClient()
	}
	return e.client
}

func (e Email) SendHtml(ctx context.Context, receivers []string, subject string, html string) error {
	m := gomail.NewMessage()
	m.SetHeader("From", fmt.Sprintf("%s <%s>", e.sender, e.userName))
	m.SetHeader("To", receivers...)
	//m.SetAddressHeader("Cc", "dan@example.com", "Dan")
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", html)
	//m.Attach("/home/Alex/lolcat.jpg")

	if err := e.getClient().DialAndSend(m); err != nil {
		return err
	}
	return nil
}
