package main

import (
	"log"
	"net/http"
	"net/url"
)

const (
	YoApiUrl = "https://api.justyo.co/yo/"
)

// send uses the JustYo API to send a push notification to my phone
// with a link to the serve image using its unique ID.
func send(id string) {
	if len(id) == 0 {
		return
	}

	values := url.Values{}
	values.Add("api_token", config.YoApiKey)
	values.Add("username", config.Yo)
	values.Add("link", config.BaseLink+"?f="+id)

	resp, err := http.PostForm(YoApiUrl, values)
	if err != nil {
		log.Println("error while sending the notification: %v", err)
		return
	}

	if resp.StatusCode != 200 {
		log.Println("error while sending the notification: %v", resp.Status)
	}
}
