package http2

import (
	"context"
	"crypto/tls"
	"github.com/julienschmidt/httprouter"
	"golang.org/x/sys/unix"
	"net"
	"net/http"
	"sync"
	"syscall"
	"time"
)

type (
	SocketOptions struct {
		ReuseAddr bool
		ReusePort bool
	}

	TLSConfig struct {
		Enabled         bool
		CertificateFile string
		KeyFile         string
		MinVersion      uint16
		CipherSuites    []uint16
	}

	Config struct {
		ListenAddr      string
		SockOptions     SocketOptions
		TLSConfig       TLSConfig
		ReadTimeout     time.Duration
		WriteTimeout    time.Duration
		IdleTimeout     time.Duration
		PreFunc         func(http.ResponseWriter, *http.Request, httprouter.Params) error
		PostFunc        func(http.ResponseWriter, *http.Request, httprouter.Params)
		NotFoundHandler http.Handler
	}

	Server struct {
		sync.RWMutex

		server   *http.Server
		config   Config
		preFunc  func(http.ResponseWriter, *http.Request, httprouter.Params) error
		postFunc func(http.ResponseWriter, *http.Request, httprouter.Params)
		router   *httprouter.Router
	}
)

func NewServer(cfg Config) *Server {
	s := &Server{}
	s.config = cfg

	if s.config.ListenAddr == "" {
		s.config.ListenAddr = "0.0.0.0:8080"
		if s.config.TLSConfig.Enabled {
			s.config.ListenAddr = "0.0.0.0:8443"
		}
	}

	s.router = httprouter.New()
	if cfg.NotFoundHandler != nil {
		s.router.NotFound = cfg.NotFoundHandler
	}

	s.server = &http.Server{
		Addr:         s.config.ListenAddr,
		Handler:      s.router,
		ReadTimeout:  s.config.ReadTimeout,
		WriteTimeout: s.config.WriteTimeout,
		IdleTimeout:  s.config.IdleTimeout,
	}

	if s.config.TLSConfig.Enabled {
		if s.config.TLSConfig.MinVersion == 0 {
			s.config.TLSConfig.MinVersion = tls.VersionTLS12
		}

		tc := tls.Config{}

		tc.MinVersion = s.config.TLSConfig.MinVersion
		tc.NextProtos = []string{"h2"}

		s.server.TLSConfig = &tc
	}

	s.preFunc = cfg.PreFunc
	s.postFunc = cfg.PostFunc

	return s
}

func (s *Server) setSocketOpts(_, _ string, rawConn syscall.RawConn) (err error) {
	err = rawConn.Control(func(fd uintptr) {
		var er error
		if s.config.SockOptions.ReuseAddr {
			if er = unix.SetsockoptInt(int(fd), unix.SOL_SOCKET, unix.SO_REUSEADDR, 1); er != nil {
				return
			}
		} else if s.config.SockOptions.ReusePort {
			if er = unix.SetsockoptInt(int(fd), unix.SOL_SOCKET, unix.SO_REUSEPORT, 1); err != nil {
				return
			}
		}
	})
	return
}

func (s *Server) Listen() (err error) {
	var listener net.Listener
	lc := net.ListenConfig{Control: s.setSocketOpts}
	if listener, err = lc.Listen(context.Background(), "tcp", s.config.ListenAddr); err == nil {
		if s.config.TLSConfig.Enabled {
			if err = s.server.ServeTLS(listener, s.config.TLSConfig.CertificateFile, s.config.TLSConfig.KeyFile); err != nil && err != http.ErrServerClosed {
				return
			}
		}
		if err = s.server.Serve(listener); err != nil && err != http.ErrServerClosed {
			return
		}
	}
	return
}

func (s *Server) Shutdown() (err error) {
	return s.server.Shutdown(context.Background())
}

func (s *Server) RegisterHandler(method, path string, handler httprouter.Handle) {
	s.router.Handle(method, path, handler)
}
