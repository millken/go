package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"

	"golang.org/x/net/html"
	"golang.org/x/net/html/charset"
	"golang.org/x/text/transform"
)

func isTitleElement(n *html.Node) bool {
	return n.Type == html.ElementNode && n.Data == "title"
}

func traverse(n *html.Node) (string, bool) {
	if isTitleElement(n) {
		return n.FirstChild.Data, true
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		result, ok := traverse(c)
		if ok {
			return result, ok
		}
	}
	return "", false
}

func getHtmlTitle(body io.Reader) (string, bool) {

	doc, err := html.Parse(body)
	if err != nil {
		fmt.Println(err)
		return "", false
	}
	return traverse(doc)
}

func splitLines(s string) []string {
	var lines []string
	sc := bufio.NewScanner(strings.NewReader(s))
	for sc.Scan() {
		lines = append(lines, sc.Text())
	}
	return lines
}

var timeoutFlag *int

func main() {
	urlFileFlag := flag.String("urlFile", "", "包含了一组URL的文件路径")
	timeoutFlag = flag.Int("timeout", 7, "超时时间")
	flag.Parse()

	ch := make(chan string)

	var urlArr []string
	if *urlFileFlag != "" {
		b, err := ioutil.ReadFile(*urlFileFlag) // just pass the file name
		if err != nil {
			panic(err)
		}

		urlArr = splitLines(string(b))
		for _, url := range urlArr {
			go fetch(url, ch) // start a goroutine
		}
	} else {
		urlArr = os.Args[1:]
		for _, url := range urlArr {
			go fetch(url, ch) // start a goroutine
		}
	}

	for range urlArr {
		fmt.Println(<-ch) // receive from channel ch
	}
}

func fetch(url string, ch chan<- string) {
	start := time.Now()
	client := &http.Client{
		Timeout: time.Duration(*timeoutFlag) * time.Second,
		// CheckRedirect: func(req *http.Request, via []*http.Request) error {
		// 	return http.ErrUseLastResponse
		// },
	}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		ch <- fmt.Sprint(err) // send to channel ch
		return
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/83.0.4103.116 Safari/537.36")

	resp, err := client.Do(req)
	if err != nil {
		ch <- fmt.Sprintf(`0s %s %d, %s, "%s"`, url, 000, "", "")
		return
	}
	defer resp.Body.Close() // don't leak resources
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		secs := time.Since(start).Seconds()
		ch <- fmt.Sprintf(`%.2fs %s %d, %s, "%s"`, secs, url, 000, "", "")
		return
	}
	//fmt.Printf("%s", body)
	e, charset, _ := charset.DetermineEncoding(body, "")
	utf8Reader := transform.NewReader(bytes.NewReader(body), e.NewDecoder())
	title, _ := getHtmlTitle(utf8Reader)

	secs := time.Since(start).Seconds()
	ch <- fmt.Sprintf(`%.2fs %s %d, %s, "%s"`, secs, url, resp.StatusCode, charset, title)
}
