package eventsender

import (
	"bufio"
	"net"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestEventSender(t *testing.T) {
	server := NewStoppableHTTPServer(":8080", func(w http.ResponseWriter, r *http.Request) {
		reader := bufio.NewReader(r.Body)
		line, err := reader.ReadString('\n')
		assert.NoError(t, err)
		assert.Equal(t, "123\n", line)
		line, err = reader.ReadString('\n')
		assert.NoError(t, err)
		assert.Equal(t, "\"abc\"\n", line)
		_, err = reader.ReadString('\n')
		assert.Error(t, err)
		assert.Equal(t, "EOF", err.Error())
	})
	go server.ListenAndServe()
	client := New("POST", "http://localhost:8080")
	err := client.SendEvent(123)
	assert.NoError(t, err)
	err = client.SendEvent("abc")
	assert.NoError(t, err)
	err = client.Close()
	assert.NoError(t, err)
	time.Sleep(100 * time.Millisecond)
	server.Stop()
	assert.NoError(t, client.Error())

}

type StoppableHTTPServer struct {
	srv *http.Server
	ln  net.Listener
}

// New creates a new webserver
func NewStoppableHTTPServer(addr string, handler http.HandlerFunc) *StoppableHTTPServer {
	mux := http.NewServeMux()
	mux.HandleFunc("/", handler)
	srv := &http.Server{
		Addr:           addr,
		Handler:        mux,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	server := &StoppableHTTPServer{srv, nil}
	return server
}

// ListenAndServe starts the webserver
func (srv *StoppableHTTPServer) ListenAndServe() error {
	ln, err := net.Listen("tcp", srv.srv.Addr)
	if err != nil {
		return err
	}
	srv.ln = ln
	return srv.srv.Serve(ln)
}

// Stop stops the webserver
func (srv *StoppableHTTPServer) Stop() error {
	if srv.ln != nil {
		return srv.ln.Close()
	}
	return nil
}
