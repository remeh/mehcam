package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/jteeuwen/imghash"
)

// config contains the configuration for the
// execution of the app.
// It is read from the "mehcam.conf" file.
var config struct {
	// camera auth
	Login    string
	Password string
	// camera url
	Url string
	// output directory
	Out string
	// fetch frequency (in seconds)
	Frequency int
	// minimum distance to consider motion
	Dist uint64
	// Pushover api key
	PoApiKey string
	// Pushover push to push
	PoUser string
	// addr to listen on. E.g. ':8080'
	Addr string
	// base link is the base of link for notification.
	// E.g. http://10.0.0.5:8000
	// This way, the link sent in the push will be:
	// http://10.0.0.5:8000/?f=AbCdEfGhIjKlMnOpqRTstVuWxZaBcDEf
	BaseLink string
}

func main() {

	if err := readConfig(); err != nil {
		log.Println("error: can't read mandatory params:", err)
		os.Exit(1)
	}

	// ----------------------

	var lastHash, currHash uint64
	var lastImg, currImg []byte
	var err error

	duration, _ := time.ParseDuration(fmt.Sprintf("%ds", config.Frequency))
	ticker := time.NewTicker(duration)

	// start the web server
	if len(config.Addr) > 0 {
		go serve()
	}

	// timed infinite loop
	for t := range ticker.C {

		// retrieve the current img and its hash
		currHash, currImg, err = current()
		if err != nil {
			log.Println("can't retrieve hash:", err)
		}

		// compute the distance between previous image and the current one
		dist := imghash.Distance(lastHash, currHash)
		if dist > config.Dist && lastImg != nil {
			filepath := config.Out + filename(t)

			// write the file
			log.Println("detected a distance:", dist)
			if err = ioutil.WriteFile(filepath, currImg, 0644); err != nil {
				log.Println("while writing file:", err)
			}

			// send notification in a goroutine
			if len(config.PoApiKey) != 0 && len(config.PoUser) != 0 && len(config.Addr) != 0 {
				id := addPic(filepath)
				go send(t, id)
			}
		}

		// store for next iteration
		lastImg, lastHash = currImg, currHash
	}
}

// readConfig reads in the configuration file to set the config var.
func readConfig() error {

	// read the configuration file
	if _, err := toml.DecodeFile("mehcam.conf", &config); err != nil {
		return fmt.Errorf("while reading mehcam.conf: %v", err)
	}

	// test some read values
	if len(config.Login) == 0 || len(config.Url) == 0 {
		return fmt.Errorf("no url or no authorization info provided.")
	}

	if len(config.PoApiKey) == 0 || len(config.PoUser) == 0 ||
		len(config.Addr) == 0 || len(config.BaseLink) == 0 {
		log.Println("no Po configuration or addr to listen to, notification disabled")
	}

	// output directory must end with a /

	if len(config.Out) > 0 && config.Out[len(config.Out)-1] != '/' {
		config.Out += "/"
	}

	if len(config.BaseLink) > 0 && config.BaseLink[len(config.BaseLink)-1] != '/' {
		config.BaseLink += "/"
	}

	if config.Frequency == 0 {
		config.Frequency = 5
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
