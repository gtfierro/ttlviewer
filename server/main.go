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
	"strings"
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
	var (
		pdf []byte
		dot []byte
		err error
	)
	defer req.Body.Close()
	if req.Method == "GET" {
		index(rw, req)
		return
	}
	// first test for URL
	url := req.FormValue("URL")
	if url != "" { // then we have a url!
		res, err := http.Get(url)
		if err != nil {
			log.Fatal(err)
		}
		defer res.Body.Close()
		pdf, dot, err = ttl.RunFile(res.Body, false)
	} else {
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
		pdf, dot, err = ttl.RunFile(file, false)
	}
	renderOption := req.FormValue("render")
	if err != nil {
		rw.WriteHeader(500)
		rw.Write([]byte(err.Error()))
		return
	}
	if renderOption == "pdf" {
		rw.Write(pdf)
	} else {
		template, err := template.ParseFiles("interact.template")
		if err != nil {
			rw.WriteHeader(500)
			rw.Write([]byte(err.Error()))
			return
		}
		template.Execute(rw, strings.Replace(string(dot), "\n", " ", -1))
	}
}

func main() {
	flag.Parse()
	http.HandleFunc("/", index)
	http.HandleFunc("/upload", upload)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))
	log.Printf("Serving on %s...\n", ":"+*port)
	log.Fatal(http.ListenAndServe(":"+*port, nil))
}
