package proxy

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
)

func copyHeader(dst, src http.Header) {
	for k, vv := range src {
		for _, v := range vv {
			dst.Add(k, v)
		}
	}
}

type Proxy struct {
}

var connectSID string

func (p *Proxy) ServeHTTP(wr http.ResponseWriter, req *http.Request) {
	log.Println(req.RemoteAddr, " ", req.Method, " ", req.URL)

	//    if req.URL.Scheme != "http" && req.URL.Scheme != "https" {
	//        msg := "unsupported protocal scheme "+req.URL.Scheme
	//        http.Error(wr, msg, http.StatusBadRequest)
	//        log.Println(msg)
	//        return
	//    }

	client := &http.Client{}

	//    req.RequestURI = "https://hivemapper.com:10000"
	//    req.URL.Scheme = "https"
	//    req.URL.Host = "hivemapper.com:10000"

	//    if clientIP, _, err := net.SplitHostPort(req.RemoteAddr); err == nil {
	//        appendHostToXForwardHeader(req.Header, clientIP)
	//    }+

	requestURL := req.URL
	requestURL.Scheme = "https"
	requestURL.Host = "hivemapper.com:10000"
	fmt.Println("req url:", requestURL)

	newReq, err := http.NewRequest(req.Method, requestURL.String(), nil)
	if err != nil {
		fmt.Printf("client: could not create request: %s\n", err)
		os.Exit(1)
	}
	//    cookie := req.Header.Get("Cookie")
	//    fmt.Println("Cookie:", cookie)
	//    if cookie != "" {
	//        newReq.Header.Set("Cookie", cookie)
	//    }

	newReq.Header = req.Header
	if connectSID != "" {
		newReq.Header.Set("Cookie", connectSID)
	}
	newReq.Body = req.Body

	resp, err := client.Do(newReq)
	if err != nil {
		http.Error(wr, "Server Error", http.StatusInternalServerError)
		log.Fatal("ServeHTTP:", err)
	}
	defer resp.Body.Close()

	log.Println(req.RemoteAddr, " ", resp.Status)

	for k, _ := range resp.Header {
		v := resp.Header.Get(k)
		if k == "Set-Cookie" {
			v = strings.ReplaceAll(v, "Domain=.hivemapper.com; Path=/;", "")
			v = strings.ReplaceAll(v, "; HttpOnly; Secure", "")
			connectSID = v
		}
		wr.Header().Set(k, v)
	}

	wr.Header().Set("Access-Control-Allow-Origin", "*")
	wr.Header().Set("Access-Control-Allow-Headers", "*")
	wr.Header().Set("Access-Control-Allow-Methods", "GET,PUT,POST,DELETE,PATCH,OPTIONS")
	wr.Header().Set("access-control-expose-headers", "Set-Cookie")

	wr.WriteHeader(resp.StatusCode)
	_, err = io.Copy(wr, resp.Body)
	if err != nil {
		panic(err)
	}
}
