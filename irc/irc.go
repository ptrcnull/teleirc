package irc

import (
	"crypto/tls"
	"fmt"
	"log"
	"os"
	"strings"

	irc "github.com/fluffle/goirc/client"
)

const badChars = "_*[]()~`>#+-=|{}.!\\"

func cleanMsg(msg string) string {
	for _, char := range badChars {
		ch := string(char)
		msg = strings.Replace(msg, ch, "\\"+ch, -1)
	}
	return msg
}

func ConnectIRC(incoming chan string, outgoing chan string, mainquit chan bool) {
	cfg := irc.NewConfig("tg")
	cfg.SSL = true
	cfg.SSLConfig = &tls.Config{ServerName: os.Getenv("IRC_HOST"), InsecureSkipVerify: true}
	cfg.Server = os.Getenv("IRC_HOST") + ":6697"
	cfg.Me.Ident = "teleirc"
	cfg.Me.Name = "bridge with Telegram"
	c := irc.Client(cfg)

	// Add handlers to do things here!
	// e.g. join a channel on connect.
	c.HandleFunc(irc.CONNECTED,
		func(conn *irc.Conn, line *irc.Line) {
			passwd := os.Getenv("IRC_PASSWD")
			if passwd != "" {
				log.Println("Trying to get oper")
				conn.Oper("tg", passwd)
			}
			conn.Join("#telegram")
		})
	// And a signal on disconnect
	quit := make(chan bool)
	c.HandleFunc(irc.DISCONNECTED,
		func(conn *irc.Conn, line *irc.Line) { quit <- true })

	c.HandleFunc(irc.PRIVMSG, func(conn *irc.Conn, line *irc.Line) {
		if line.Target() != "#telegram" {
			return
		}
		outgoing <- fmt.Sprintf("*%s*: %s", line.Nick, cleanMsg(line.Text()))
	})

	// Tell client to connect.
	if err := c.Connect(); err != nil {
		fmt.Printf("Connection error: %s\n", err.Error())
	}

	go func() {
		for m := range incoming {
			if strings.Contains(m, ": > ") {
				m = strings.Join(strings.Split(m, "> "), "\x0303> ")
			}
			c.Privmsg("#telegram", m)
		}
	}()

	// Wait for disconnect
	<-quit
	mainquit <- true
}
