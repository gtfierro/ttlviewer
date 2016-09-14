package main

import (
	"crypto/md5"
	"encoding/binary"
	"flag"
	"fmt"
	"github.com/gtfierro/ttlviewer/ttl"
	"html/template"
	"log"
	"net/http"
	"time"
)

var port = flag.String("p", "1212", "Port to serve")

func index(w http.ResponseWriter, req *http.Request) {
	defer req.Body.Close()
	template, err := template.ParseFiles("index.template")
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte(err.Error()))
		return
	}
	// get a token
	h := md5.New()
	seed := make([]byte, 16)
	binary.PutVarint(seed, time.Now().UnixNano())
	h.Write(seed)
	token := fmt.Sprintf("%x", h.Sum(nil))
	template.Execute(w, token)
}

func upload(rw http.ResponseWriter, req *http.Request) {
	defer req.Body.Close()
	if req.Method == "GET" {
		index(rw, req)
		return
	}
	if err := req.ParseMultipartForm(5 << 20); err != nil {
		rw.WriteHeader(500)
		rw.Write([]byte(err.Error()))
		return
	}
	file, _, err := req.FormFile("uploadfile")
	if err != nil {
		rw.WriteHeader(500)
		rw.Write([]byte(err.Error()))
		return
	}
	defer file.Close()
	pdf, err := ttl.RunFile(file)
	if err != nil {
		rw.WriteHeader(500)
		rw.Write([]byte(err.Error()))
		return
	}
	rw.Write(pdf)
}

func main() {
	flag.Parse()
	http.HandleFunc("/", index)
	http.HandleFunc("/upload", upload)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))
	log.Printf("Serving on %s...\n", ":"+*port)
	log.Fatal(http.ListenAndServe(":"+*port, nil))
}
