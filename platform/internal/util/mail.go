package util

import (
	"crypto/tls"
	"github.com/StellrisJAY/cloud-emu/platform/internal/conf"
	"net/smtp"
)

type EmailHelper struct {
	c *conf.Smtp
}

func NewEmailHelper(c *conf.Smtp) *EmailHelper {
	return &EmailHelper{c: c}
}

func (e *EmailHelper) Send(to string, subject, body string) error {
	conn, err := tls.Dial("tcp", e.c.Addr, &tls.Config{InsecureSkipVerify: true, ServerName: e.c.Host})
	if err != nil {
		return err
	}
	defer conn.Close()
	c, err := smtp.NewClient(conn, e.c.Host)
	if err != nil {
		return err
	}
	defer c.Quit()
	if err := c.Auth(smtp.PlainAuth("", e.c.UserName, e.c.Password, e.c.Host)); err != nil {
		return err
	}

	if err = c.Mail(e.c.UserName); err != nil {
		return err
	}
	if err = c.Rcpt(to); err != nil {
		return err
	}

	msg := "From:" + e.c.From + "\r\n" +
		"To:" + to + "\r\n" +
		"Subject:" + subject + "\r\n\r\n" +
		body + "\r\n"

	// 写入邮件内容
	wc, err := c.Data()
	if err != nil {
		return err
	}
	defer wc.Close()
	_, err = wc.Write([]byte(msg))
	if err != nil {
		return err
	}
	return nil
}
