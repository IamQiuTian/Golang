package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"runtime"
	"strings"
)

var url []string

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())
}

func main() {
	var filepath string
	c := make(chan bool)

	flag.StringVar(&filepath, "f", "url.txt", "-f <filename>")
	flag.Parse()

	go func() {
		url = Openfile(filepath)
		c <- true
	}()
	<-c

	for _, v := range url {
		result, _ := dnsurL(v)
		if result != nil {
			fmt.Println("[+] ", v, ": ", result)
		} else {
			fmt.Println("[-] ", v, ": ", "Did not find dns parsing records")

		}
	}
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
			ustr := strings.Split(reader.Text(), "//")[1]
			lo = append(lo, ustr)
		} else {
			lo = append(lo, reader.Text())
		}
	}
	return lo
}

func dnsurL(uL string) ([]string, error) {
	ns, err := net.LookupHost(uL)
	if err != nil {
		return nil, err
	}
	return ns, err
}

/*
[+]  www.btjxw.gov.cn :  [1.24.191.100]
[+]  www.bynrjxw.gov.cn :  [221.199.203.239]
[-]  www.als.gov.cn/jw :  Did not find dns parsing records
[+]  jw.xlgl.gov.cn :  [110.16.70.58]
[+]  jjj.elht.gov.cn :  [222.74.213.233]
[+]  www.wlcbjxw.gov.cn :  [124.67.110.69]
[+]  jw.xam.gov.cn :  [58.18.7.197 58.18.7.202]
[+]  www.hlbrjxw.gov.cn :  [58.18.185.134]
[+]  www.nmgwh.gov.cn :  [110.16.70.8]
[+]  www.nmgxwcbgdj.gov.cn :  [222.74.213.250 110.16.70.105]
[+]  www.bjnmghotel.com :  [115.29.108.91]
[+]  www.xincheng-hotel.com.cn :  [121.41.231.83]
[+]  www.hhhtbtjc.com :  [42.123.76.89]
[+]  www.hctrust.cn :  [162.251.95.52 162.251.95.54]
*/
