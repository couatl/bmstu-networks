package main

import (
	"crypto/tls"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/smtp"
	"strings"
)

type Mail struct {
	senderId string
	toIds    []string
	subject  string
	body     string
}

type SmtpServer struct {
	host string
	port string
}

func (s *SmtpServer) ServerName() string {
	return s.host + ":" + s.port
}

func (mail *Mail) BuildMessage() string {
	message := ""
	message += fmt.Sprintf("From: %s\r\n", mail.senderId)
	if len(mail.toIds) > 0 {
		message += fmt.Sprintf("To: %s\r\n", strings.Join(mail.toIds, ";"))
	}

	message += fmt.Sprintf("Subject: %s\r\n", mail.subject)
	message += "\r\n" + mail.body

	return message
}

var (
	serverFile = flag.String("s", "server.json", "File of server config")
	mailFile = flag.String("m", "mail.json", "File of mail connection config")
)

func main() {
	flag.Parse()

	var mail Mail
	data, err := ioutil.ReadFile(*serverFile)
	if err != nil {
		panic(err.Error())
	}

	err = json.Unmarshal(data, &mail)
	if err != nil {
		panic(err.Error())
	}

	messageBody := mail.BuildMessage()

	smtpServer := SmtpServer{}
	data, err = ioutil.ReadFile(*mailFile)
	if err != nil {
		panic(err.Error())
	}

	err = json.Unmarshal(data, &smtpServer)
	if err != nil {
		panic(err.Error())
	}

	log.Println(smtpServer.host, messageBody)

	auth := smtp.PlainAuth("", mail.senderId, "IU8-NETWORKS", smtpServer.host)

	tlsConfig := &tls.Config{
		InsecureSkipVerify: true,
		ServerName:         smtpServer.host,
	}

	conn, err := tls.Dial("tcp", smtpServer.ServerName(), tlsConfig)
	if err != nil {
		log.Panic(err)
	}

	client, err := smtp.NewClient(conn, smtpServer.host)
	if err != nil {
		log.Panic(err)
	}

	if err = client.Auth(auth); err != nil {
		log.Panic(err)
	}

	if err = client.Mail(mail.senderId); err != nil {
		log.Panic(err)
	}
	for _, k := range mail.toIds {
		if err = client.Rcpt(k); err != nil {
			log.Panic(err)
		}
	}

	w, err := client.Data()
	if err != nil {
		log.Panic(err)
	}

	_, err = w.Write([]byte(messageBody))
	if err != nil {
		log.Panic(err)
	}

	err = w.Close()
	if err != nil {
		log.Panic(err)
	}

	client.Quit()

	log.Println("Mail sent successfully")
}
