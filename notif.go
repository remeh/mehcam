package main

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"time"
)

const (
	PoApiUrl = "https://api.pushover.net/1/messages.json"
)

// send uses the Pushover API to send a push notification to my phone
// with a link to the serve image using its unique ID.
func send(t time.Time, id string) {
	if len(id) == 0 {
		return
	}

	values := url.Values{}
	values.Add("token", config.PoToken)
	values.Add("user", config.PoUser)
	values.Add("title", "Motion detection")
	values.Add("message", fmt.Sprintf("Motion detected at %s", t.Format(time.ANSIC)))
	values.Add("url", config.BaseLink+"?f="+id)
	values.Add("url_title", "Open picture")
	values.Add("sound", "gamelan")

	resp, err := http.PostForm(PoApiUrl, values)
	if err != nil {
		log.Printf("error while sending the notification: %v", err)
		return
	}

	if resp.StatusCode != 200 {
		log.Printf("error while sending the notification: %v", resp.Status)
	}
}
