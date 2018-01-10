package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"path/filepath"
	"strings"
	"time"
)

var (
	directory *string = flag.String("d", "", "Directory path")
	file      *string = flag.String("f", "", "file path")
	port      *string = flag.String("p", "8888", "Listening port")
)

func main() {
	flag.Parse()
	if len(*directory) == 0 && len(*file) == 0 {
		flag.Usage()
		return
	}
	if len(*file) != 0 {
		filename := filepath.Base(*file)
		getPublic(*port, filename)
		getPrivate(*port, filename)
		http.HandleFunc(fmt.Sprintf("/%s", filename), func(w http.ResponseWriter, r *http.Request) {
			http.ServeFile(w, r, *file)
		})

	} else {
		filename := "false"
		getPublic(*port, filename)
		getPrivate(*port, filename)
		http.Handle("/", http.FileServer(http.Dir(*directory)))
	}
	log.Fatal(http.ListenAndServe(":"+*port, nil))
}

func getPublic(port, filename string) {
	timeout := time.Duration(1 * time.Second)
	client := http.Client{
		Timeout: timeout,
	}
	resp, err := client.Get("http://members.3322.org/dyndns/getip")
	if err != nil {
		return
	}
	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	if filename == "false" {

		fmt.Printf("Tour address: http://%s:%s/\n", strings.Replace(string(b), "\n", "", -1), port)
	} else {
		fmt.Printf("Tour address: http://%s:%s/%s\n", strings.Replace(string(b), "\n", "", -1), port, filename)
	}
}

func getPrivate(port, filename string) {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		log.Fatal(err)
	}
	for _, address := range addrs {
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				if filename == "false" {
					fmt.Printf("Tour address: http://%s:%s/\n", ipnet.IP.String(), port)
				} else {
					fmt.Printf("Tour address: http://%s:%s/%s\n", ipnet.IP.String(), port, filename)
				}
			}
		}
	}
}
