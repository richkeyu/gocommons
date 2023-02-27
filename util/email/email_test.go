package email

import "testing"

func TestEmail_SendHtml(t *testing.T) {
	var (
		userName  = "ope@im30.net"
		sender    = "KOP"
		password  = "cqD63fDdR3EsrKX6"
		host      = "smtp.exmail.qq.com"
		port      = 465
		receivers = []string{"xxx@example.com"}
	)
	email := NewEmail(sender, host, port, userName, password)
	subject := "测试邮件"
	html := `<html><body><div>sdjfkldsfjkldsfjksdlfjdsklfjsdklf啦啦啦</div></body></html>`
	email.SendHtml(nil, receivers, subject, html)
	html = `<html><body><div>哈哈哈哈kfjdlskfjskldfslk</div></body></html>`
	email.SendHtml(nil, receivers, subject, html)
}
