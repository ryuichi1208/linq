package main

import (
	"context"
	"crypto/tls"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httputil"
	"os"
	"os/signal"
	"runtime/debug"
	"strings"
	"syscall"
	"time"

	"linq/cmd"

	. "github.com/logrusorgru/aurora/v3"
	"golang.org/x/sync/errgroup"
	"gopkg.in/yaml.v3"
)

var (
	fileName    = flag.String("f", "./test.yml", "読み込ませたいyamlファイルのパス")
	worknums    = flag.Int("w", 6, "実行するhttpクライアントの並列数")
	single      = flag.Bool("s", false, "シングルフライトを用いた実装")
	verbosemode = flag.Bool("v", false, "詳細表示モード")
)

type Url struct {
	Url []string `yaml:"url"`
}

type withGoroutineID struct {
	out io.Writer
}

func (u *Url) UrlListInit(filename string) error {
	_, err := os.Stat(filename)
	return err
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
		log.Println(Bold(Cyan(fmt.Sprintf("[DEBUG]: %s", filename))))
		log.Println(Bold(Cyan(len(u.Url))))
		log.Println(Bold(Cyan(u)))
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
					// httpのヘッダー情報をdump
					log.Println(string(header))
				}
				resp, err := cli.Get(url)
				if err != nil {
					return err
				}
				defer func() {
					if _, err := io.Copy(ioutil.Discard, resp.Body); err != nil {
						log.Fatal("io error", err)
					}
					resp.Body.Close()
				}()
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

func workerSingleFlyght(ch <-chan string, id int, cli *http.Client) {
	// return
}

func main() {
	flag.Parse()
	flag.Usage = func() {
		usageTxt := `Usage example [option]
An example of customizing usage output

    -s, --s  STRING argument, default: String help message
    -i, --i  INTEGER argument, default: Int help message
    -b, --b  BOOLEAN argument, default: Bool help message
`
		fmt.Fprintf(os.Stderr, "%s\n", usageTxt)
	}

	sig := make(chan os.Signal, 2)
	signal.Notify(sig, syscall.SIGTERM, os.Interrupt)

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
	if !*single {
		for i := 0; i < *worknums; i++ {
			i := i
			// 1つでも接続エラーが合った場合は全てのgoroutineをキャンセルする
			eg.Go(func() error {
				return worker(ch, i, cli, ctx)

			})
		}
	} else {
		log.Println("single flyght")
		for i := 0; i < *worknums; i++ {
			i := i
			go workerSingleFlyght(ch, i, cli)
		}
	}

	for _, url := range urls.Url {
		select {
		case ch <- url:
		case <-time.After(3 * time.Second):
			// チャネルの受信側が存在しなくなる or サーバ側が高負荷なので終了させる
			log.Fatal(Bold(Red("[ERROR] channel sending Timeout")))
		case <-sig:
			// キーボード割り込みとその他ネットワーク関連のシグナルは全てエラーとする
			log.Fatal(Bold(Red("[ERROR] signal recved")))
		}
	}
	close(ch)

	if err := eg.Wait(); err != nil {
		log.Fatal(err)
	}
	fmt.Println(cmd.CommandNameUsage())
}
