package main

import (
	"bytes"
	"flag"
	"fmt"
	"log"
	"net"
)

var (
	addr = flag.String("addr", "127.0.0.1:3030", "addr to bind to")
)

func main() {
	flag.Parse()

	ln, err := net.Listen("tcp", *addr)
	if err != nil {
		log.Fatalf("listen error: %v", err)
	}

	log.Printf("listening on %s", ln.Addr())

	for {
		conn, err := ln.Accept()
		if err != nil {
			// FIXME: maybe repeat.
			log.Fatalf("accept() error: %v", err)
		}

		log.Printf(
			"accepted new conn: %s -> %s",
			conn.LocalAddr(), conn.RemoteAddr(),
		)

		var (
			bts = make([]byte, 4096)

			lines [][]byte
			line  []byte
		)
		for {
			n, err := conn.Read(bts)
			if err != nil {
				log.Printf("read error: %v", err)
				break
			}
			log.Printf("got bytes: %s", bts[:n])

			data := bts[:n]
			for {
				i := bytes.IndexByte(data, '\n')
				if i == -1 {
					line = append(line, data...)
					break
				}
				line = append(line, data[:i]...)
				lines = append(lines, line)
				log.Printf("request line: %s", line)
				line = nil
				data = data[i+1:]
			}

			if !bytes.HasPrefix(lines[0], []byte("GET")) {
				log.Printf("not a GET")
				conn.Close()
				continue
			}
			proto := bytes.Index(lines[0], []byte("HTTP/1.1"))
			if proto == -1 {
				log.Printf("no proto")
				conn.Close()
				continue
			}

			resource := bytes.TrimSpace(lines[0][3:proto])
			log.Printf("client wants: %s", resource)

			// May send response.
			fmt.Fprintf(conn, "HTTP/1.1 200 OK\nContent-Length: 0\n\n")
			conn.Close()
		}
	}
}
