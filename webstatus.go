package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"runtime"
	"strings"
	"sync"
	"time"
)

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())
}

func main() {
	filepath := flag.String("f", "url.txt", "-f <filename>")
	flag.Parse()

	wg := sync.WaitGroup{}
	url := Openfile(*filepath)
	wg.Add(len(url))
	for _, v := range url {
		go UrL(&wg, v)
	}
	wg.Wait()
}

func Openfile(filename string) []string {
	var lo []string
	f, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	reader := bufio.NewScanner(f)
	for reader.Scan() {
		if strings.HasPrefix(reader.Text(), "http://") {
			lo = append(lo, reader.Text())
		} else {
			lo = append(lo, "http://"+reader.Text())
		}
	}
	return lo
}

func UrL(wg *sync.WaitGroup, uL string) {
	u, err := url.Parse(uL)
	if err != nil {
		log.Println(err)
	}
	q := u.Query()
	u.RawQuery = q.Encode()

	timeout := time.Duration(1 * time.Second)
	client := http.Client{
		Timeout: timeout,
	}
	res, err := client.Get(u.String())
	if err != nil {
		fmt.Printf("%v Connect failed\n", uL)
		wg.Done()
		return
	}
	resCode := res.StatusCode
	defer res.Body.Close()

	if resCode != 200 {
		fmt.Printf("%v Connect failed\n", uL)
	} else {
		fmt.Printf("%v Connect Success\n", uL)
	}
	wg.Done() 
}
