package telegram

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"
)

var apiUrl = "https://api.telegram.org/" + os.Getenv("TG_TOKEN") + "/"

func parse(message *Message) string {
	text := message.Text
	if message.Caption != nil {
		text = *message.Caption
	}
	if len(message.Photo) != 0 {
		text = fmt.Sprintf("<%dx photo> %s", len(message.Photo), text)
	}
	return text
}

func format(message *Message, isEdited bool) []string {
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

func ConnectTelegram(incoming chan string, outgoing chan string) {
	go func() {
		for m := range incoming {
			SendUpdate(m)
		}
	}()
	go func() {
		for {
			body, err := getUpdates()
			if err != nil {
				log.Println(err)
			} else {
				for _, m := range body.Result {
					log.Println(m)

					var message *Message
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
			}
			time.Sleep(time.Millisecond * 100)
		}
	}()
}

// var offset = "263196900"
var offset int64

func getUpdates() (Response, error) {
	var res Response

	if offset == 0 {
		offsetBytes, err := ioutil.ReadFile("offset")
		if err != nil {
			panic(err)
		}
		offset, err = strconv.ParseInt(string(offsetBytes), 10, 64)
	}

	resp, err := http.Get(apiUrl + "getUpdates?offset=" + strconv.FormatInt(offset, 10))
	if err != nil {
		return res, err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return res, err
	}

	if string(body) != `{"ok":true,"result":[]}` {
		log.Println(string(body))
	}

	err = json.Unmarshal(body, &res)

	if len(res.Result) > 0 {
		log.Println("Last update: ", res.Result[len(res.Result)-1].UpdateID)
		offset = res.Result[len(res.Result)-1].UpdateID + 1
		err1 := ioutil.WriteFile("offset", []byte(strconv.FormatInt(offset, 10)), 0644)
		if err1 != nil {
			log.Println(err1)
		}
	}

	return res, err
}

func SendUpdate(update string) {
	response, err := http.PostForm(apiUrl+"sendMessage", url.Values{
		"chat_id":    {os.Getenv("TG_CHAT_ID")},
		"text":       {update},
		"parse_mode": {"MarkdownV2"},
	})

	if err != nil {
		log.Println(err)
		return
	}

	defer func() {
		err := response.Body.Close()
		if err != nil {
			log.Println(err)
		}
	}()

	body, err := ioutil.ReadAll(response.Body)

	if err != nil {
		log.Println(err)
		return
	}

	fmt.Printf("%s\n", string(body))
}
