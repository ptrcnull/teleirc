package telegram

type Response struct {
	Ok     bool     `json:"ok"`
	Result []Update `json:"result"`
}

type Author struct {
	ID        int64  `json:"id"`
	IsBot     bool   `json:"is_bot"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Username  string `json:"username"`
}

type Update struct {
	UpdateID      int64    `json:"update_id"`
	Message       *Message `json:"message"`
	EditedMessage *Message `json:"edited_message"`
}

type Chat struct {
	ID    int64  `json:"id"`
	Title string `json:"title"`
	Type  string `json:"type"`
}

type Sticker struct {
	Width        int    `json:"width"`
	Height       int    `json:"height"`
	SetName      string `json:"set_name"`
	IsAnimated   bool   `json:"is_animated"`
	FileId       string `json:"file_id"`
	FileUniqueId string `json:"file_unique_id"`
	FileSize     int64  `json:"file_size"`
}

type File struct {
	FileID       string `json:"file_id"`
	FileUniqueId string `json:"file_unique_id"`
	FileSize     int64  `json:"file_size"`
	Width        int    `json:"width"`
	Height       int    `json:"height"`
}

type PollOption struct {
	Text       string `json:"text"`
	VoterCount int    `json:"voter_count"`
}

type Poll struct {
	ID       string       `json:"id"`
	Question string       `json:"question"`
	Options  []PollOption `json:"options"`
	IsClosed bool         `json:"is_closed"`
}

type Message struct {
	MessageID           int64    `json:"message_id"`
	From                Author   `json:"from"`
	Text                string   `json:"text"`
	Caption             *string  `json:"caption"`
	Chat                *Chat    `json:"chat"`
	Sticker             *Sticker `json:"sticker"`
	Photo               []File   `json:"photo"`
	Document            *File    `json:"document"`
	NewChatTitle        *string  `json:"new_chat_title"`
	LeftChatParticipant *Author  `json:"left_chat_participant"`
	NewChatParticipant  *Author  `json:"new_chat_participant"`
	ReplyToMessage      *Message `json:"reply_to_message"`
}
