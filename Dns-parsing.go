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
		fmt.Println(v, ": ", result)
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
www.nmgjxw.gov.cn :  [222.74.213.236]
www.nmggzw.gov.cn :  [110.16.70.8]
nmgfy.chinacourt.org :  [111.202.173.230]
ssmzggbm.nmgov.edu.cn :  [101.7.0.40]
nmdswx.nmds.gov.cn :  [218.21.128.235]
www.nmg.gov.cn :  [60.31.197.41]
www.ordos.gov.cn :  [58.18.251.58]
www.nmschoolfootball.com :  [101.7.0.75]
zypt.nmgov.edu.cn :  [101.7.0.42]
nmgmwzy.nmgov.edu.cn :  [101.7.0.43]
nmgs.gov.cn :  [116.112.10.219 222.74.203.130]
neimenggu.mca.gov.cn :  [116.193.41.42]
www.nmds.gov.cn :  [218.21.128.234]
www.nmagri.gov.cn :  [110.16.70.147]
www.nmgat.gov.cn :  [58.18.164.36]
www.nmfda.gov.cn :  [110.16.70.10]
www.nmgcz.gov.cn :  [202.99.230.232]
www.nmgepb.gov.cn :  [202.99.246.243]
www.nmgfgw.gov.cn :  [61.138.111.230]
www.nmggtt.gov.cn :  [116.114.18.210]
www.nmgjrb.gov.cn :  [110.16.70.32]
*/
