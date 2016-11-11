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

// TODO(remy): document this.
var mapping map[string]string = make(map[string]string)

// serve creates an HTTP server listening
// for GET call to gets some pictures captured
// by the camera in case of motion.
func serve() {
	// init the seed because we'll need some
	// random later.
	rand.Seed(time.Now().UnixNano())

	// create the http server
	s := &http.Server{
		Addr:    env.addr,
		Handler: picHandler{},
	}
	log.Fatal(s.ListenAndServe())
}

func addPic(filename string) {
	// generate an random ID
	id := make([]byte, 32)
	for i := range id {
		id[i] = idsChar[rand.Intn(len(idsChar))]
	}

	// we'll now serve this file using this newly generated ID
	mapping[string(id)] = filename
}

// ----------------------

type picHandler struct{}

func (p picHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	filename := r.Header.Get("f")
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