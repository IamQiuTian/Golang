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
	directory *string = flag.String("d", "false", "Directory path")
	file      *string = flag.String("f", "false", "file path")
	port      *string = flag.String("p", "8888", "Listening port")
)

func main() {
	flag.Parse()
	if *directory == "false" && *file == "false" {
		flag.Usage()
		return
	}
	if *file != "false" {
		filename := filepath.Base(*file)
		getPublic(*port, filename)
		getPrivate(*port, filename)
		http.HandleFunc(fmt.Sprintf("/%s", filename), func(w http.ResponseWriter, r *http.Request) {
			fmt.Println()
			log.Printf(
				"%s\t%s\t%q\t%s",
				r.RemoteAddr,
				r.Method,
				r.RequestURI,
			)
			http.ServeFile(w, r, *file)
		})

	} else {
		filename := "nil"
		getPublic(*port, filename)
		getPrivate(*port, filename)
		http.Handle("/", http.StripPrefix("/", http.FileServer(http.Dir(*directory))))
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
	if filename == "nil" {
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
				if filename == "nil" {
					fmt.Printf("Tour address: http://%s:%s/\n", ipnet.IP.String(), port)
				} else {
					fmt.Printf("Tour address: http://%s:%s/%s\n", ipnet.IP.String(), port, filename)
				}
			}
		}
	}
}

/*
 Used:
 $ ./easyDown -f FilePath -p 8080
 $ ./easyDown -d DirectoryPath -p 8080
*/
