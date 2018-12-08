package main

import (
	"context"
	"flag"
	"html/template"
	"net/http"
	"sync"
)

var (
	addr = flag.String("addr", "127.0.0.1:3030", "addr to bind to")
)

type Message struct {
	User string
	Text string
}

type MessageStorage interface {
	Put(context.Context, Message) error
	List(context.Context) ([]Message, error)
}

type InMemoryStorage struct {
	mu       sync.Mutex
	messages []Message
}

func (s *InMemoryStorage) Put(_ context.Context, msg Message) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.messages = append(s.messages, msg)
	return nil
}

func (s *InMemoryStorage) List(_ context.Context) ([]Message, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.messages, nil
}

type IndexPage struct {
	Title    string
	Messages []Message
}

const indexText = `
<!DOCTYPE html>
<html>
	<head>
		<meta charset="UTF-8">
		<title>{{.Title}}</title>
	</head>
	<body>
		<form method="post" action="/">
			<input type="text" name="user"/>
			<textarea name="message"></textarea>
			<input type="submit" value="post!"/>
		</form>
		{{range .Messages}}
			<p>
				<div><strong>{{ .User }}</strong></div>
				<div>{{ .Text }}</div>
			</p>
		{{else}}
			<div><strong>no posts yet</strong></div>
		{{end}}
	</body>
</html>`

var index = template.Must(template.New("index").Parse(indexText))

func main() {
	flag.Parse()

	var st MessageStorage = &InMemoryStorage{}

	http.HandleFunc("/", func(res http.ResponseWriter, req *http.Request) {
		if req.Method == http.MethodPost {
			if err := req.ParseForm(); err != nil {
				res.WriteHeader(http.StatusBadRequest)
				return
			}
			user := req.Form.Get("user")
			message := req.Form.Get("message")
			if user == "" || message == "" {
				res.WriteHeader(http.StatusBadRequest)
				return
			}
			err := st.Put(context.TODO(), Message{
				User: user,
				Text: message,
			})
			if err != nil {
				res.WriteHeader(http.StatusInternalServerError)
				return
			}
		}

		ms, err := st.List(context.TODO())
		if err != nil {
			res.WriteHeader(http.StatusInternalServerError)
			return
		}
		page := IndexPage{
			Title:    "my-cool-guestbook",
			Messages: ms,
		}
		if err := index.Execute(res, page); err != nil {
			res.WriteHeader(http.StatusInternalServerError)
			return
		}
	})
	http.ListenAndServe(*addr, nil)
}