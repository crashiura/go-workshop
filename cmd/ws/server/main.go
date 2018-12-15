package main

import (
	"container/list"
	"flag"
	"log"
	"net"
	"sync"

	"github.com/gobwas/ws"
)

var (
	addr = flag.String("addr", "127.0.0.1:3333", "addr to bind to")
)

type User struct {
	mu      sync.Mutex
	onClose []func()

	once sync.Once
	conn net.Conn
	out  chan ws.Frame
}

func NewUser(conn net.Conn) *User {
	return &User{
		conn: conn,
		out:  make(chan ws.Frame, 100),
	}
}

func (u *User) OnClose(f func()) {
	u.mu.Lock()
	defer u.mu.Unlock()
	u.onClose = append(u.onClose, f)
}

func (u *User) Close() {
	u.once.Do(func() {
		close(u.out)
		u.conn.Close()

		u.mu.Lock()
		fs := u.onClose
		u.onClose = nil
		u.mu.Unlock()

		for _, f := range fs {
			f()
		}
	})
}

func (u *User) WriteMessages(m ws.Frame) {
	for ok := true; ok; m, ok = <-u.out {
		err := ws.WriteFrame(u.conn, m)
		if err != nil {
			log.Printf("write frame error: %v", err)
			u.Close()
			return
		}
	}
}

func (u *User) ReadMessages(broadcast func(ws.Frame)) {
	for {
		f, err := ws.ReadFrame(u.conn)
		if err != nil {
			log.Printf("read frame error: %v", err)
			u.Close()
			return
		}
		broadcast(f)
	}
}

type Room struct {
	mu    sync.Mutex
	users list.List
}

func (r *Room) Register(user *User) (remove func()) {
	r.mu.Lock()
	defer r.mu.Unlock()

	el := r.users.PushBack(user)
	return func() {
		r.mu.Lock()
		defer r.mu.Unlock()
		r.users.Remove(el)
	}
}

func (r *Room) Broadcast(f ws.Frame) {
	r.mu.Lock()
	var victims []*User
	for el := r.users.Front(); el != nil; el = el.Next() {
		// Type assertion is safe here because we put only *User in the list.
		// If it is not a *User – it is mistake and must be fixed asap.
		user := el.Value.(*User)
		select {
		case user.out <- f:
		default:
			victims = append(victims, user)
		}
	}
	r.mu.Unlock()

	for _, user := range victims {
		// FIXME: notify user that his connection is slow.
		user.Close()
	}
}

func main() {
	flag.Parse()

	ln, err := net.Listen("tcp", *addr)
	if err != nil {
		log.Fatalf("listen error: %v", err)
	}
	log.Printf("listening on %s", ln.Addr())

	var chat Room
	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Fatalf("accept error: %v", err)
		}

		user := NewUser(conn)
		cleanup := chat.Register(user)
		user.OnClose(cleanup)

		go user.WriteMessages(ws.NewTextFrame(
			[]byte("greetings from gophers chat server"),
		))
		go user.ReadMessages(func(f ws.Frame) {
			chat.Broadcast(f)
		})
	}
}
