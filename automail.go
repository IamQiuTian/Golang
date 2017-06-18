package main

import (
	"encoding/json"
	"flag"
	"log"
	"net/smtp"
	"os"
	"strings"
	//使用了https://github.com/scorredoira/email中的email.go
	"github.com/email"
)

type Email struct {
	Title       string   `json:"title"`
	Message     string   `json:"message"`
	FromName    string   `json:"fromName"`
	FromAddress string   `json:"fromAddress"`
	To          []string `json:"to"`
	Filepath    string   `json:"filepath"`
	EmailServer string   `json:"emailserver"`
	EmailDome   string   `json:"Emaildome"`
	Username    string   `json:"username"`
	Password    string   `json:"password"`
}

var conf Email

func Error(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func txt_email() {
	auth := smtp.PlainAuth("", conf.Username, conf.Password, conf.EmailDome)
	to := conf.To
	nickname := conf.FromName
	user := conf.FromAddress
	subject := conf.Title
	content_type := "Content-Type: text/plain; charset=UTF-8"
	body := conf.Message
	msg := []byte("To: " + strings.Join(to, ",") + "\r\nFrom: " + nickname +
		"<" + user + ">\r\nSubject: " + subject + "\r\n" + content_type + "\r\n\r\n" + body)
	err := smtp.SendMail(conf.EmailServer, auth, user, to, msg)
	Error(err)
}

func annex_email() {
	m := email.NewMessage(conf.Title, conf.Message)
	m.From.Name = conf.FromName
	m.From.Address = conf.FromAddress
	m.To = conf.To
	err := m.Attach(conf.Filepath)
	Error(err)
	err = email.Send(conf.EmailServer, smtp.PlainAuth("", conf.Username, conf.Password, conf.EmailDome), m)
	Error(err)
}

func main() {
	confile := flag.String("config", "email.conf", "--config filename")
	flag.Parse()

	file, err := os.Open(*confile)
	defer file.Close()
	if err != nil {
		flag.Usage()
	}

	err = json.NewDecoder(file).Decode(&conf)
	Error(err)

	if len(conf.Filepath) > 0 {
		annex_email()
	} else {
		txt_email()
	}
}

/*
./automail -config email.conf

cat email.conf
{
"Title":"邮件标题",
"Message":"邮件内容",
"FromName":"发件人名称",
"FromAddress":"发件人邮箱",
"To":[收件人邮箱],
"Filepath":"附件路径",
"EmailServer":"邮件服务器",
"EmailDome":"邮件域名",
"Username":"发件邮箱",
"Password":"邮箱鉴权码"
}
*/
