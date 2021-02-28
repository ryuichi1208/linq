package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httputil"
	"time"
)

func Request(url string, cli *http.Client) {
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("Authorization", "Bearer access-token")
	dump, _ := httputil.DumpRequestOut(req, true)
	fmt.Printf("%s", dump)
	resp, err := cli.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	defer log.Println("close")
	io.Copy(ioutil.Discard, resp.Body)
}

func main() {
	cli := &http.Client{
		Transport: &http.Transport{
			MaxIdleConns:        100,
			MaxIdleConnsPerHost: 3,
			MaxConnsPerHost:     3,
			IdleConnTimeout:     5 * time.Second,
		},
	}

	for i := 0; i < 6; i++ {
		go Request("http://192.168.1.114:30001/abort?exit=sleep", cli)
	}

	time.Sleep(time.Second * 6)
	go Request("http://192.168.1.114:30001/abort?exit=sleep", cli)
	time.Sleep(time.Second * 60)
	Request("http://192.168.1.114:30001/abort?exit=sleep", cli)
}
