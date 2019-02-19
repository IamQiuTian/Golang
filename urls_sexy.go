package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"
)

var gonumber = 30
var counter chan bool
var wg sync.WaitGroup

func main() {
	counter = make(chan bool, gonumber)
	// https://github.com/alexkimxyz/nsfw_data_scrapper/tree/master/raw_data/sexy
	f, err := os.Open("./urls_sexy.txt")
	if err != nil {
		panic(err)
	}

	reader := bufio.NewScanner(f)
	for reader.Scan() {
		urls := reader.Text()
		counter <- true
		wg.Add(1)

		go func() {
			defer wg.Done()
			status := download(urls)
			fmt.Println(urls, " -> ", status)
		}()
	}

	defer f.Close()
	defer close(counter)
	wg.Wait()
}

func download(urls string) error {
	jpgName := parseUrl(urls)

	client := http.Client{
		Timeout: time.Duration(60 * time.Second),
	}
	resp, err := client.Get(urls)
	if err != nil {
		<-counter
		return err
	}

	out, err := os.Create("sex/" + jpgName)
	if err != nil {
		<-counter
		return err
	}
	io.Copy(out, resp.Body)
	<-counter

	defer resp.Body.Close()
	defer out.Close()
	return errors.New("ok")
}

func parseUrl(url string) string {
	urlList := strings.Split(url, "/")
	urlLen := len(urlList) - 1
	jpgname := urlList[urlLen]

	return jpgname
}
