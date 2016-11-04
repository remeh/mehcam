package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"time"

	"github.com/jteeuwen/imghash"
)

var env struct {
	auth string
	url  string
	out  string
}

func main() {
	// read environment var

	env.auth = os.Getenv("AUTH")
	env.url = os.Getenv("URL")
	env.out = os.Getenv("OUTPUT")

	// ----------------------

	var lastHash, currHash uint64
	var lastImg, currImg []byte
	var err error

	for {
		currHash, currImg, err = currentHash()
		if err != nil {
			fmt.Println("can't retrieve hash:", err)
		}

		dist := imghash.Distance(lastHash, currHash)

		if dist > 10 {
			fmt.Println(time.Now(), "detected a distance:", dist)
			if err = ioutil.WriteFile(env.out+now(), lastImg, 0644); err != nil {
				fmt.Println("while writing file:", err)
			}
		}

		lastImg, lastHash = currImg, currHash
		time.Sleep(time.Second * 1)
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

func now() string {
	return time.Now().Format("2006-01-02-15-04-05.jpg")
}
