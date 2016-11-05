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
	auth     string
	url      string
	out      string
	duration time.Duration
	dist     uint64
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
			fmt.Println(time.Now(), "detected a distance:", dist)
			if err = ioutil.WriteFile(env.out+filename(t), currImg, 0644); err != nil {
				fmt.Println("while writing file:", err)
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

	env.auth = os.Getenv("AUTH")
	env.url = os.Getenv("URL")

	if len(env.auth) == 0 || len(env.url) == 0 {
		return fmt.Errorf("no url or no authorization info provided.")
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
