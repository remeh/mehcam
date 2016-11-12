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

// send uses the JustYo API to send a push notification to my phone
// with a link to the serve image using its unique ID.
func send(id string) {
	if len(id) == 0 {
		return
	}

	values := url.Values{}
	values.Add("token", config.PoApiKey)
	values.Add("user", config.PoUser)
	values.Add("title", "Motion detection")
	// TODO(remy): use the time from the main
	values.Add("message", fmt.Sprintf("Motion detected at %s", time.Now().String()))
	values.Add("url", config.BaseLink+"?f="+id)
	values.Add("url_title", "Open picture")

	resp, err := http.PostForm(PoApiUrl, values)
	if err != nil {
		log.Printf("error while sending the notification: %v", err)
		return
	}

	if resp.StatusCode != 200 {
		log.Printf("error while sending the notification: %v", resp.Status)
	}
}
