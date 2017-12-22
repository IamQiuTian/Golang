package main

import (
	"flag"
	"io"
	"log"
	"net"
)

type route struct {
	src  string
	dest string
}

func main() {
	s := flag.String("s", "nil", "Listened ip")
	d := flag.String("d", "nil", "Mapping ip")
	flag.Parse()

	if *s == "nil" || *d == "nil" {
		flag.Usage()
		return
	}

	r := &route{
		src:  *s,
		dest: *d,
	}
	go r.listen()
	<-chan bool(nil)
}

func (r *route) handle(src net.Conn) {
	dst, err := net.Dial("tcp", r.dest)
	if err != nil {
		log.Printf("Error connecting to destination: %s\n", r.dest)
		return
	}
func main() {
	s := flag.String("s", "nil", "Listened ip")
	d := flag.String("d", "nil", "Mapping ip")
	flag.Parse()

	if *s == "nil" || *d == "nil" {
		flag.Usage()
		return
	}

	r := &route{
		src:  *s,
		dest: *d,
	}
	go r.listen()
	<-chan bool(nil)
}

func (r *route) handle(src net.Conn) {
	dst, err := net.Dial("tcp", r.dest)
	if err != nil {

	log.Printf("Routing %s -> %s", r.src, r.dest)

	go io.Copy(src, dst)
	io.Copy(dst, src)

	log.Printf("Connection closed: %s -> %s", r.src, r.dest)
}

func (r *route) listen() {
	log.Printf("Creating listener: %s -> %s\n", r.src, r.dest)

	l, err := net.Listen("tcp", r.src)
	if err != nil {
		log.Fatalf("Error setting up listener: %s -> %s: %s", r.src, r.dest, err)
	}

	for {
		conn, err := l.Accept()
		if err != nil {
			log.Printf("Error accepting connection: %s", err)
		}
		go r.handle(conn)
	}
}
