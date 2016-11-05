package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"time"

	"github.com/jteeuwen/imghash"
)

var env struct {
	auth     string
	url      string
	out      string
	duration time.Duration
}

func main() {

	readParams()

	// ----------------------

	var lastHash, currHash uint64
	var lastImg, currImg []byte
	var err error

	ticker := time.NewTicker(env.duration)

	for t := range ticker.C {
		currHash, currImg, err = currentHash()
		if err != nil {
			fmt.Println("can't retrieve hash:", err)
		}

		dist := imghash.Distance(lastHash, currHash)

		if dist > 10 {
			fmt.Println(time.Now(), "detected a distance:", dist)
			if err = ioutil.WriteFile(env.out+now(t), lastImg, 0644); err != nil {
				fmt.Println("while writing file:", err)
			}
		}

		lastImg, lastHash = currImg, currHash
	}
}

func readParams() {
	env.auth = os.Getenv("AUTH")
	env.url = os.Getenv("URL")
	env.out = os.Getenv("OUTPUT")
	if len(env.out) > 0 && env.out[len(env.out)-1] != '/' {
		env.out += "/"
	}

	var err error
	if env.duration, err = time.ParseDuration(os.Getenv("DURATION")); err != nil {
		env.duration = time.Second
	}
}

func currentHash() (uint64, []byte, error) {
	data, err := get()
	if err != nil {
		return 0, nil, fmt.Errorf("error while hashing the image: %s", err.Error())
	}

	if h, err := hash(data); err != nil {
		return 0, nil, fmt.Errorf("error while hashing the image: %s", err.Error())
	} else {
		return h, data, err
	}
}

func now(t time.Time) string {
	return t.Format("2006-01-02_15-04-05.jpg")
}
