package slack

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// Message slack message
type Message struct {
	Channel     string        `json:"channel"`
	Username    string        `json:"username"`
	IconEmoji   string        `json:"icon_emoji"`
	Text        string        `json:"text"`
	Attachments []*Attachment `json:"attachments"`
}

// AddAttachment add a attachment
func (sm *Message) AddAttachment(sa *Attachment) {
	sm.Attachments = append(sm.Attachments, sa)
}

// Post post slack message
func (sm *Message) Post(url string, timeout time.Duration) error {
	bs, err := json.Marshal(sm)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(bs))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: timeout}
	res, err := client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("%s - POST %s", res.Status, url)
	}
	return nil
}

// Attachment slack attachment
type Attachment struct {
	Callback   string   `json:"callback"`
	Color      string   `json:"color"`
	Pretext    string   `json:"pretext"`
	AuthorName string   `json:"author_name"`
	AuthorLink string   `json:"author_link"`
	AuthorIcon string   `json:"author_icon"`
	Title      string   `json:"title"`
	TitleLink  string   `json:"title_link"`
	Text       string   `json:"text"`
	ImageURL   string   `json:"image_url"`
	ThumbURL   string   `json:"thumb_url"`
	Fields     []*Field `json:"fields"`
}

// AddField add a field
func (sm *Attachment) AddField(sf *Field) {
	sm.Fields = append(sm.Fields, sf)
}

// Field slack field
type Field struct {
	Title string `json:"title"`
	Value string `json:"value"`
	Short string `json:"short"`
}
