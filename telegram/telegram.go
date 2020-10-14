package telegram

import (
	"fmt"
	"git.ddd.rip/ptrcnull/telegram"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

func parse(message *telegram.Message) string {
	text := message.Text
	if message.Caption != nil {
		text = *message.Caption
	}
	if len(message.Photo) != 0 {
		text = fmt.Sprintf("<%dx photo> %s", len(message.Photo), text)
	}
	return text
}

func format(message *telegram.Message, isEdited bool) []string {
	var lines []string
	text := parse(message)
	if isEdited {
		text += " (edited)"
	}
	if message.ReplyToMessage != nil {
		reply := parse(message.ReplyToMessage)
		lines = append(lines, "> In reply to @"+message.ReplyToMessage.From.Username)
		for _, line := range strings.Split(reply, "\n") {
			lines = append(lines, "> "+line)
		}
	}

	lines = append(lines, strings.Split(text, "\n")...)

	user := message.From.Username
	for i := range lines {
		lines[i] = user + ": " + lines[i]
	}
	return lines
}

var offset int64

func ConnectTelegram(incoming chan string, outgoing chan string) {
	if offset == 0 {
		offsetBytes, err := ioutil.ReadFile("offset")
		if err != nil {
			panic(err)
		}
		offset, err = strconv.ParseInt(string(offsetBytes), 10, 64)
		if err != nil {
			panic(err)
		}
	}

	client := telegram.Client{
		Key:    os.Getenv("TG_TOKEN"),
		Offset: offset,
	}

	go func() {
		for m := range incoming {
			_, err := client.SendMarkdownMessage(os.Getenv("TG_CHAT_ID"), m)
			if err != nil {
				outgoing <- "error: " + err.Error()
			}
		}
	}()

	go func() {
		for {
			res, err := client.GetUpdates()
			if err != nil {
				log.Println(err)
			} else {
				for _, m := range res.Result {
					log.Println(m)

					var message *telegram.Message
					edited := false

					if m.Message != nil {
						message = m.Message
					}
					if m.EditedMessage != nil {
						message = m.EditedMessage
						edited = true
					}

					lines := format(message, edited)
					for _, line := range lines {
						outgoing <- line
					}

					log.Println("Update: ", m.UpdateID)
					log.Println(m.Message)
					log.Println(m.EditedMessage)
				}
				offset = res.Result[len(res.Result)-1].UpdateID + 1
				err = ioutil.WriteFile("offset", []byte(strconv.FormatInt(offset, 10)), 0644)
				if err != nil {
					log.Println(err)
				}
			}
			time.Sleep(time.Millisecond * 100)
		}
	}()
}
