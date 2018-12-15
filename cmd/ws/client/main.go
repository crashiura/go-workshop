package main

import (
	"bufio"
	"flag"
	"log"
	"net"
	"os"

	"github.com/gobwas/ws"
)

var (
	addr = flag.String("addr", "127.0.0.1:3333", "addr of the server")
)

func main() {
	flag.Parse()
	log.SetFlags(0)

	conn, err := net.Dial("tcp", *addr)
	if err != nil {
		log.Fatalf("dial error: %v", err)
	}
	go func() {
		s := bufio.NewScanner(os.Stdin)
		for s.Scan() {
			err := ws.WriteFrame(conn, ws.NewTextFrame(s.Bytes()))
			if err != nil {
				log.Fatalf("write frame error: %v", err)
			}
		}
	}()
	for {
		f, err := ws.ReadFrame(conn)
		if err != nil {
			log.Fatalf("read frame error: %v", err)
		}
		switch f.Header.OpCode {
		case ws.OpText:
			log.Printf("> %s", f.Payload)
		default:
			log.Printf("unexpected frame type: %s", f.Header.OpCode)
		}
	}
}
