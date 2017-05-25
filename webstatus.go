package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
)

func main() {
	var filepath string
	flag.StringVar(&filepath, "f", "url.txt", "-f <filename>")
	flag.Parse()

	url := openfile(filepath)
	for _, v := range url {
		uL, status, err := urL(v)
		if status != 200 && err != nil {
			fmt.Printf("%v 不可以访问\n", uL)
		}
		fmt.Printf("%v 可以访问\n", uL)
	}
}

func openfile(filename string) []string {
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

func urL(uL string) (string, int, error) {
	u, err := url.Parse(uL)
	if err != nil {
		log.Println(err)
	}
	q := u.Query()
	u.RawQuery = q.Encode()
	res, err := http.Get(u.String())
	if err != nil {
		return uL, 0, errors.New("not open")
	}
	resCode := res.StatusCode
	defer res.Body.Close()

	return uL, resCode, nil
}

/*
http://222.74.213.247 可以访问
http://110.16.70.81 可以访问
http://116.113.100.77 可以访问
http://110.16.70.3 可以访问
http://110.16.70.26 可以访问
http://110.16.70.69 可以访问
http://110.16.70.70 可以访问
http://110.16.70.71 不可以访问
http://110.16.70.71 可以访问
*/
