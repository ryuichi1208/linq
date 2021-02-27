package main

import (
	"context"
	"crypto/tls"
	"flag"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httputil"
	"os"
	"runtime/debug"
	"strings"
	"time"

	"gopkg.in/yaml.v3"

	"golang.org/x/sync/errgroup"
)

var (
	fileName    = flag.String("f", "./test.yml", "読み込ませたいyamlファイルのパス")
	worknums    = flag.Int("w", 6, "実行するhttpクライアントの並列数")
	verbosemode = flag.Bool("v", false, "詳細表示モード")
)

type Url struct {
	Url []string `yaml:"url"`
}

type withGoroutineID struct {
	out io.Writer
}

func (w *withGoroutineID) Write(p []byte) (int, error) {
	firstline := []byte(strings.SplitN(string(debug.Stack()), "\n", 2)[0])
	return w.out.Write(append(firstline[:len(firstline)-10], p...))
}

func readYaml(filename string) Url {
	buf, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Fatal(err)
	}

	var u Url
	if err = yaml.Unmarshal(buf, &u); err != nil {
		log.Fatal(err)
	}
	if *verbosemode {
		log.Println(filename)
		log.Println(len(u.Url), u)
	}
	return u

}

func worker(ch <-chan string, id int, cli *http.Client, ctx context.Context) error {
	for {
		select {
		case url, ok := <-ch:
			if ok {
				if *verbosemode {
					req, _ := http.NewRequest("GET", url, nil)
					header, _ := httputil.DumpRequestOut(req, true)
					log.Println(string(header))
				}
				resp, err := cli.Get(url)
				if err != nil {
					return err
				}
				log.Printf("%d: %s", resp.StatusCode, url)
			} else {
				return nil
			}
		case <-ctx.Done():
			log.Println("canseld: ", id)
			return nil
		}
	}
}

func main() {
	flag.Parse()
	urls := readYaml(*fileName)

	if *verbosemode {
		log.SetOutput(&withGoroutineID{out: os.Stderr})
	}

	ch := make(chan string, *worknums)

	cli := &http.Client{
		Timeout: time.Duration(30) * time.Second,
		Transport: &http.Transport{
			MaxIdleConns:    100,
			MaxConnsPerHost: 6,
			TLSClientConfig: &tls.Config{InsecureSkipVerify: false},
		},
	}

	eg, ctx := errgroup.WithContext(context.TODO())
	for i := 0; i < *worknums; i++ {
		i := i
		// 1つでも接続エラーが合った場合は全てのgoroutineをキャンセルする
		eg.Go(func() error {
			return worker(ch, i, cli, ctx)

		})
	}

	for _, url := range urls.Url {
		select {
		case ch <- url:
		case <-time.After(3 * time.Second):
			// チャネルの受信側が存在しなくなる or サーバ側が高負荷なので終了させる
			log.Fatal("Timeout")
		}
	}
	close(ch)

	if err := eg.Wait(); err != nil {
		log.Fatal(err)
	}
}
