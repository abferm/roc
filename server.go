package roc

import (
	"io"
	"net"
	"sync"
	"time"
)

type Handler interface {
	ServeROC(ReponseWriter io.Writer, Request *Message)
}

// SyncServer should be used for single connection mediums, like RS-232/485
func NewSyncServer(l net.Listener, handler Handler) (server *SyncServer) {
	server = new(SyncServer)
	server.l, server.handler = l, handler
	return
}

type SyncServer struct {
	l       net.Listener
	handler Handler
	stop    bool
	mu      sync.Mutex
	done    chan struct{}
}

// Responds to a single response at a time
func (srv *SyncServer) Serve() (err error) {
	logger.Debugf("Listening on %s", srv.l.Addr())
	defer srv.l.Close()
	srv.mu.Lock()
	srv.stop = false
	srv.done = make(chan struct{})
	srv.mu.Unlock()
	for {
		srv.mu.Lock()
		if srv.stop {
			close(srv.done)
			srv.mu.Unlock()
			return
		}
		srv.mu.Unlock()

		var conn net.Conn
		conn, err = srv.l.Accept()
		if err != nil {
			return
		}
		request := new(Message)
		readErr := request.read(conn)
		if readErr != nil {
			logger.Errorf("Invalid request: %s", readErr.Error())
			// TODO: Error response should go here
			continue
		}
		srv.handler.ServeROC(conn, request)
	}
	return
}

func (srv *SyncServer) Stop() {
	srv.mu.Lock()
	srv.stop = true
	srv.mu.Unlock()
}

func (srv *SyncServer) Wait(timeout time.Duration) (ok bool) {
	select {
	case <-time.After(timeout):
		ok = false
	case <-srv.done:
		ok = true
	}
	return
}
