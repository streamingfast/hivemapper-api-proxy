package proxy

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
)

type Proxy struct {
	connectSID string
}

func (p *Proxy) ServeHTTP(wr http.ResponseWriter, req *http.Request) {
	log.Println(req.RemoteAddr, " ", req.Method, " ", req.URL)

	if req.URL.Path == "/reset" {
		p.connectSID = ""
		wr.WriteHeader(200)
		wr.Write([]byte("RESET COMPLETED"))
		addAccessControlHeaders(wr)
		return
	}

	client := &http.Client{}

	requestURL := req.URL
	requestURL.Scheme = "https"
	requestURL.Host = "hivemapper.com:10000"
	fmt.Println("req url:", requestURL)

	newReq, err := http.NewRequest(req.Method, requestURL.String(), nil)
	if err != nil {
		fmt.Printf("client: could not create request: %s\n", err)
		os.Exit(1)
	}

	newReq.Header = req.Header
	if p.connectSID != "" {
		newReq.Header.Set("Cookie", p.connectSID)
		fmt.Println("cookie", p.connectSID)
	}
	newReq.Body = req.Body

	resp, err := client.Do(newReq)
	if err != nil {
		http.Error(wr, "Server Error", http.StatusInternalServerError)
		log.Fatal("ServeHTTP:", err)
	}
	defer resp.Body.Close()

	log.Println(req.RemoteAddr, " ", resp.Status)

	for k := range resp.Header {
		v := resp.Header.Get(k)
		if k == "Set-Cookie" {
			v = strings.ReplaceAll(v, "Domain=.hivemapper.com; Path=/;", "")
			v = strings.ReplaceAll(v, "; HttpOnly; Secure", "")
			p.connectSID = v
		}
		wr.Header().Set(k, v)
	}

	addAccessControlHeaders(wr)

	wr.WriteHeader(resp.StatusCode)
	_, err = io.Copy(wr, resp.Body)
	if err != nil {
		panic(err)
	}
}

func addAccessControlHeaders(wr http.ResponseWriter) {
	wr.Header().Set("Access-Control-Allow-Origin", "*")
	wr.Header().Set("Access-Control-Allow-Headers", "*")
	wr.Header().Set("Access-Control-Allow-Methods", "GET,PUT,POST,DELETE,PATCH,OPTIONS")
	wr.Header().Set("Access-Control-Expose-Headers", "Set-Cookie")
}
