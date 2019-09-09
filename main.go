package main

import (
	"flag"
	"io"
	"log"
	"net/http"
	"net/http/httputil"
)

var (
	url string
)

func main() {
	flag.StringVar(&url, "url", "http://localhost:8080", "setting get url")
	flag.Parse()

	mux := http.NewServeMux()
	mux.HandleFunc("/", index)
	mux.HandleFunc("/dump", dumprequest)
	mux.HandleFunc("/slow", slowrequest)
	log.Fatal(http.ListenAndServe(":8081", mux))

}

func index(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "This is echoecho service\n")
}

func dumprequest(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "This is echoecho service\n")

	dump, _ := httputil.DumpRequest(r, true)
	io.WriteString(w, "===DumpRequest===\n")
	io.WriteString(w, string(dump))
	io.WriteString(w, "\n")

	// Request
	url = url + "/dump"
	req, _ := http.NewRequest("GET", url, nil)

	dumpReq, _ := httputil.DumpRequestOut(req, true)
	io.WriteString(w, "===DumpRequestOut===\n")
	io.WriteString(w, string(dumpReq))

	client := new(http.Client)
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	dumpResp, _ := httputil.DumpResponse(resp, true)
	io.WriteString(w, "===DumpResponse===\n")
	io.WriteString(w, string(dumpResp))

}

func slowrequest(w http.ResponseWriter, r *http.Request) {

	io.WriteString(w, "This is echoecho service\n")
	// Request
	url = url + "/wait"
	req, _ := http.NewRequest("GET", url, nil)
	client := new(http.Client)
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	dumpResp, _ := httputil.DumpResponse(resp, true)
	io.WriteString(w, "===DumpResponse===\n")
	io.WriteString(w, string(dumpResp))
}
