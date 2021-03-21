package main

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"sync"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/saintfish/chardet"
	"golang.org/x/net/html/charset"
)

const BASE_URL string = "https://news.livedoor.com/topics/category"

func initHttp() *http.Client {
	return &http.Client{
		Transport: &http.Transport{
			DialContext: (&net.Dialer{
				Timeout:   30 * time.Second,
				KeepAlive: 30 * time.Second,
				DualStack: true,
			}).DialContext,
			ForceAttemptHTTP2:     true,
			MaxIdleConns:          500,
			MaxIdleConnsPerHost:   100,
			IdleConnTimeout:       90 * time.Second,
			TLSHandshakeTimeout:   10 * time.Second,
			ResponseHeaderTimeout: 10 * time.Second,
			ExpectContinueTimeout: 1 * time.Second,
		},
		Timeout: 60 * time.Second,
	}
}

func scrapeWeb(wg *sync.WaitGroup, cli *http.Client, url string, ch chan<- string) {
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("user-agent", "google/1.0")
	res, err := cli.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		io.Copy(ioutil.Discard, res.Body)
		res.Body.Close()
	}()
	buf, _ := ioutil.ReadAll(res.Body)

	// 文字コード判定
	det := chardet.NewTextDetector()
	detRslt, _ := det.DetectBest(buf)

	// 文字コード変換
	bReader := bytes.NewReader(buf)
	reader, _ := charset.NewReaderLabel(detRslt.Charset, bReader)

	// HTMLパース
	doc, _ := goquery.NewDocumentFromReader(reader)
	doc.Find("ul.articleList").Find("a").Each(func(_ int, s *goquery.Selection) {
		url, _ := s.Attr("href")
		ch <- url
	})
	wg.Done()
}

func requestsWorker(wg *sync.WaitGroup, cli *http.Client, ch <-chan string) {
	for {
		select {
		case url := <-ch:
			func() {
				req, _ := http.NewRequest("GET", url, nil)
				req.Header.Set("user-agent", "google/1.0")
				res, err := cli.Do(req)
				fmt.Printf("%s: %d\n", url, res.StatusCode)
				if err != nil {
					log.Fatal(err)
				}
				defer func() {
					io.Copy(ioutil.Discard, res.Body)
					res.Body.Close()
				}()
			}()
		case <-time.After(1 * time.Second):
			fmt.Println("Time Out!")
			wg.Done()
			return
		}
	}
}

func main() {
	wg := sync.WaitGroup{}

	// CATEGORY_LIST := []string{"main", "dom", "world", "eco", "ent", "sports", "gourmet", "love"}
	CATEGORY_LIST := []string{"main"}

	client := initHttp()

	ch := make(chan string, 10)

	for _, cat := range CATEGORY_LIST {
		wg.Add(1)
		go scrapeWeb(&wg, client, fmt.Sprintf("%s/%s/", BASE_URL, cat), ch)
	}

	for i := 0; i < 16; i++ {
		wg.Add(1)
		go requestsWorker(&wg, client, ch)
	}
	wg.Wait()
}
