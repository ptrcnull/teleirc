package main

import (
	"github.com/ptrcnull/teleirc/irc"
	"github.com/ptrcnull/teleirc/telegram"
)

func main() {
	toIRC := make(chan string)
	toTelegram := make(chan string)
	quit := make(chan bool)

	go irc.ConnectIRC(toIRC, toTelegram, quit)
	go telegram.ConnectTelegram(toTelegram, toIRC)
	<-quit
}
