package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"net/url"
	"os"
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
		ok, _ := pathExist(*file)
		if !ok {
			fmt.Println("File does not exist")
			return
		}
		filename := filepath.Base(*file)
		getPublic(*port, filename)
		getPrivate(*port, filename)
		print("\n")
		http.HandleFunc(fmt.Sprintf("/%s", filename), func(w http.ResponseWriter, r *http.Request) {
			log.Printf(
				"%s  %s  %s",
				r.RemoteAddr,
				r.Method,
				r.RequestURI,
			)
			http.ServeFile(w, r, *file)
		})

	} else {
		ok, _ := pathExist(*directory)
		if !ok {
			fmt.Println("Directory does not exist")
			return
		}
		filename := "nil"
		getPublic(*port, filename)
		getPrivate(*port, filename)
		print("\n")
		http.Handle("/", func(prefix string, h http.Handler) http.Handler {
			if prefix == "" {
				return h
			}
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				log.Printf(
					"%s  %s  %s",
					r.RemoteAddr,
					r.Method,
					r.RequestURI,
				)
				if p := strings.TrimPrefix(r.URL.Path, prefix); len(p) < len(r.URL.Path) {
					r2 := new(http.Request)
					*r2 = *r
					r2.URL = new(url.URL)
					*r2.URL = *r.URL
					r2.URL.Path = p
					h.ServeHTTP(w, r2)
				} else {
					http.NotFound(w, r)
				}
			})

		}("/", http.FileServer(http.Dir(*directory))))
	}
	log.Fatal(http.ListenAndServe(":"+*port, nil))
}

func pathExist(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
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
$ easyDown -f FilePath -p 8080
$ easyDown -d DirectoryPath -p 8080
*/
