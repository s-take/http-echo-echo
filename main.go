package main

import (
	"flag"
	"io"
	"log"
	"net/http"
	"net/http/httputil"
	"time"
)

var (
	url string
)

var client = http.Client{
	Timeout: time.Millisecond * 30000,
}

var headers = []string{
	"x-request-id",
	"x-b3-traceid",
	"x-b3-spanid",
	"x-b3-parentspanid",
	"x-b3-sampled",
	"x-b3-flags",
	"x-ot-span-context",
}

func main() {
	flag.StringVar(&url, "url", "http://localhost:8080", "setting get url")
	flag.Parse()

	http.DefaultTransport.(*http.Transport).MaxIdleConns = 0
	http.DefaultTransport.(*http.Transport).MaxIdleConnsPerHost = 1000

	mux := http.NewServeMux()
	mux.HandleFunc("/", dump)
	mux.HandleFunc("/slow", slow)
	mux.HandleFunc("/error", err)
	log.Fatal(http.ListenAndServe(":8081", mux))

}

func dump(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "This is echoecho service\n")

	dump, _ := httputil.DumpRequest(r, true)
	io.WriteString(w, "===DumpRequest===\n")
	io.WriteString(w, string(dump))
	io.WriteString(w, "\n")

	// Request
	u := url
	req, _ := http.NewRequest("GET", u, nil)
	// add trace header
	for _, header := range headers {
		if r.Header.Get(header) != "" {
			req.Header.Add(header, r.Header.Get(header))
		}
	}

	dumpReq, _ := httputil.DumpRequestOut(req, true)
	io.WriteString(w, "===DumpRequestOut===\n")
	io.WriteString(w, string(dumpReq))

	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	// _, err = ioutil.ReadAll(resp.Body)

	dumpResp, _ := httputil.DumpResponse(resp, true)
	io.WriteString(w, "===DumpResponse===\n")
	io.WriteString(w, string(dumpResp))

}

func slow(w http.ResponseWriter, r *http.Request) {

	io.WriteString(w, "This is echoecho service\n")
	// Request
	u := url + "/slow"
	req, _ := http.NewRequest("GET", u, nil)
	client := new(http.Client)

	resp, err := client.Do(req)
	if err != nil {
		log.Println(err)
	}
	defer resp.Body.Close()
	// _, err = ioutil.ReadAll(resp.Body)

	dumpResp, _ := httputil.DumpResponse(resp, true)
	io.WriteString(w, "===DumpResponse===\n")
	io.WriteString(w, string(dumpResp))
}

func err(w http.ResponseWriter, r *http.Request) {

	io.WriteString(w, "This is echoecho service\n")
	// Request
	u := url + "/error"
	req, _ := http.NewRequest("GET", u, nil)

	resp, err := client.Do(req)
	if err != nil {
		log.Println(err)
	}
	defer resp.Body.Close()
	// _, err = ioutil.ReadAll(resp.Body)

	dumpResp, _ := httputil.DumpResponse(resp, true)
	io.WriteString(w, "===DumpResponse===\n")
	io.WriteString(w, string(dumpResp))
}
