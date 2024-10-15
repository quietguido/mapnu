package httpserver

import (
	"context"
	"net/http"
	"time"
)

const (
	ReadHeaderTimeout = 10 * time.Second
	ReadTimeout       = 2 * time.Minute
	MaxHeaderBytes    = 300 * 1024
)

type St struct {
	addr   string
	server *http.Server
	eChan  chan error
}

func New(addr string, handler http.Handler) *St {
	return &St{
		addr: addr,
		server: &http.Server{
			Addr:              addr,
			Handler:           handler,
			ReadHeaderTimeout: ReadHeaderTimeout,
			ReadTimeout:       ReadHeaderTimeout,
			MaxHeaderBytes:    MaxHeaderBytes,
		},
		eChan: make(chan error, 1),
	}
}

func (st *St) Start() {
	go func() {
		err := st.server.ListenAndServe()
		if err != nil && err == http.ErrServerClosed {
			st.eChan <- err
		}
	}()
}

func (st *St) Wait() <-chan error {
	return st.eChan
}

func (st *St) Shutdown(t time.Duration) error {
	defer close(st.eChan)

	ctx, ctxcancel := context.WithTimeout(context.Background(), t)
	defer ctxcancel()

	err := st.server.Shutdown(ctx)
	if err != nil {
		return err
	}

	return nil
}
