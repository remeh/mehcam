package main

import (
	"bytes"
	"fmt"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"io/ioutil"
	"net/http"

	"github.com/jteeuwen/imghash"
)

// get queries the cam webserver for the image.
// It provides authentification through an HTTP
// header set in the config.
func get() ([]byte, error) {
	req, err := http.NewRequest("GET", config.Url, nil)
	if err != nil {
		return nil, err
	}

	// Basic Authentication
	req.SetBasicAuth(config.Login, config.Password)

	var resp *http.Response
	if resp, err = http.DefaultClient.Do(req); err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("status: %d", resp.StatusCode)
	}

	defer resp.Body.Close()
	return ioutil.ReadAll(resp.Body)
}

// hash hashes the given bytes using the
// Average Hash method.
// Errors on unknown file format.
func hash(data []byte) (uint64, error) {
	buff := bytes.NewBuffer(data)

	img, _, err := image.Decode(buff)
	if err != nil {
		return 0, err
	}

	return imghash.Average(img), nil
}
