package main

import (
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"time"
)

// idsChar are the chars used to generate an ID.
const idsChar = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

// mapping from an unique ID to the actual filepath on the disk.
var mapping map[string]string = make(map[string]string)

// serve creates an HTTP server listening
// for GET call to gets some pictures captured
// by the camera in case of motion.
func serve() {
	// init the seed because we'll need some
	// random later.
	rand.Seed(time.Now().UnixNano())

	// create the http server
	http.Handle("/", picHandler{})
	log.Fatal(http.ListenAndServe(config.Addr, nil))
}

// addPic adds a filename indexed by a generated
// unique ID to the mapping map.
func addPic(filename string) string {
	// generate a random ID
	id := make([]byte, 64)
	for i := range id {
		id[i] = idsChar[rand.Intn(len(idsChar))]
	}

	// we'll now serve this file using this newly generated ID
	sid := string(id)
	mapping[sid] = filename
	return sid
}

// ----------------------

// picHandler resolves the given file ID using the map mapping
// to read the actual image from fs and to serve it.
type picHandler struct{}

func (p picHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	r.ParseForm()
	filename := r.Form.Get("f")
	filepath, exists := mapping[filename]

	// unknown image
	if !exists {
		w.WriteHeader(404)
		return
	}

	// read the file on disk
	data, err := ioutil.ReadFile(filepath)
	if err != nil {
		log.Printf("error: while serving file %s during read: %v", err)
		w.WriteHeader(404)
		return
	}

	// render the file
	w.Header().Set("Content-Type", http.DetectContentType(data))
	w.Write(data)
}
