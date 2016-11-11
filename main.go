package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"time"

	"github.com/jteeuwen/imghash"
)

// env contains the configuration for the
// execution of the app.
var env struct {
	// camera auth
	login    string
	password string
	// camera url
	url string
	// output directory
	out string
	// fetch frequency
	duration time.Duration
	// minimum distance to consider motion
	dist uint64
	// yo api key
	yo_api_key string
	// yo username to push
	yo string
	// addr to listen on
	addr string
}

func main() {

	if err := readParams(); err != nil {
		fmt.Println("error: can't read mandatory params:", err)
		os.Exit(1)
	}

	// ----------------------

	var lastHash, currHash uint64
	var lastImg, currImg []byte
	var err error

	ticker := time.NewTicker(env.duration)

	// start the web server
	if len(env.addr) > 0 {
		go serve()
	}

	// timed infinite loop
	for t := range ticker.C {

		// retrieve the current img and its hash
		currHash, currImg, err = current()
		if err != nil {
			fmt.Println("can't retrieve hash:", err)
		}

		// compute the distance between previous image and the current one
		dist := imghash.Distance(lastHash, currHash)
		if dist > env.dist && lastImg != nil {
			// write the file
			fmt.Println(time.Now(), "detected a distance:", dist)
			if err = ioutil.WriteFile(env.out+filename(t), currImg, 0644); err != nil {
				fmt.Println("while writing file:", err)
			}

			// send notification
			if len(env.yo_api_key) != 0 && len(env.addr) != 0 {
				if err = send(env.out + filename(t)); err != nil {
					fmt.Println("during push notification:", err)
				}
			}
		}

		// store for next iteration
		lastImg, lastHash = currImg, currHash
	}
}

// readParams reads in the execution environment for some configuration
func readParams() error {
	var err error

	// authorization and webserver url

	env.login = os.Getenv("LOGIN")
	env.password = os.Getenv("PASSWORD")
	env.url = os.Getenv("URL")
	env.yo_api_key = os.Getenv("YO_API_KEY")
	env.yo = os.Getenv("YO")

	if len(env.login) == 0 || len(env.url) == 0 {
		return fmt.Errorf("no url or no authorization info provided.")
	}

	if len(env.yo_api_key) == 0 || len(env.yo) == 0 || len(env.addr) {
		fmt.Println("no Yo API Token, Yo username or addr to listen to, notification disabled.")
	}

	// output directory

	env.out = os.Getenv("OUTPUT")
	if len(env.out) > 0 && env.out[len(env.out)-1] != '/' {
		env.out += "/"
	}

	// ----------------------

	if env.duration, err = time.ParseDuration(os.Getenv("DURATION")); err != nil {
		fmt.Println("warning: can't read DURATION env var. Fallback on 1s")
		env.duration = time.Second
	}

	if env.dist, err = strconv.ParseUint(os.Getenv("DIST"), 10, 64); err != nil {
		fmt.Println("warning: can't read DIST env var. Fallback on 10")
		env.dist = 10
	}

	return nil
}

// current snapshots the image currently seen by the cam.
// It returns the hash and the image (in the original format).
func current() (uint64, []byte, error) {
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

// filename returns a convenient filename using the given time
// as input.
func filename(t time.Time) string {
	return t.Format("2006-01-02_15-04-05.jpg")
}
